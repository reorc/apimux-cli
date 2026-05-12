package update

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Manifest struct {
	LatestVersion string `json:"latest_version"`
}

type CheckResult struct {
	CurrentVersion  string `json:"current_version"`
	LatestVersion   string `json:"latest_version,omitempty"`
	UpdateAvailable bool   `json:"update_available"`
	ManifestURL     string `json:"manifest_url,omitempty"`
	Status          string `json:"status"`
}

type UpgradeResult struct {
	CheckResult
	ExecutablePath string `json:"executable_path,omitempty"`
	ArchiveURL     string `json:"archive_url,omitempty"`
	StatusMessage  string `json:"message,omitempty"`
}

func Check(ctx context.Context, client *http.Client, currentVersion string, manifestURL string) (CheckResult, error) {
	result := CheckResult{
		CurrentVersion: strings.TrimSpace(currentVersion),
		ManifestURL:    strings.TrimSpace(manifestURL),
		Status:         "unknown",
	}

	if result.CurrentVersion == "" {
		result.CurrentVersion = "dev"
	}
	if result.ManifestURL == "" {
		return result, fmt.Errorf("release manifest URL is not configured")
	}
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, result.ManifestURL, nil)
	if err != nil {
		return result, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("release manifest returned status %d", resp.StatusCode)
	}

	var manifest Manifest
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return result, err
	}

	latest := strings.TrimSpace(manifest.LatestVersion)
	if latest == "" {
		return result, fmt.Errorf("release manifest did not include latest_version")
	}

	result.LatestVersion = latest
	switch compareVersions(result.CurrentVersion, latest) {
	case -1:
		result.UpdateAvailable = true
		result.Status = "update_available"
	case 0:
		result.Status = "up_to_date"
	default:
		result.Status = "ahead_of_manifest"
	}

	return result, nil
}

func Upgrade(ctx context.Context, client *http.Client, currentVersion string, manifestURL string, executablePath string) (UpgradeResult, error) {
	check, err := Check(ctx, client, currentVersion, manifestURL)
	result := UpgradeResult{CheckResult: check}
	if err != nil {
		return result, err
	}
	if !check.UpdateAvailable {
		result.StatusMessage = "apimux is already up to date"
		return result, nil
	}
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}

	releaseBaseURL, err := releaseBaseURLFromManifest(manifestURL)
	if err != nil {
		return result, err
	}
	releaseVersion, releaseTag, err := normalizeReleaseVersion(check.LatestVersion)
	if err != nil {
		return result, err
	}
	archiveName, err := archiveName(releaseVersion)
	if err != nil {
		return result, err
	}
	archiveURL := strings.TrimRight(releaseBaseURL, "/") + "/" + releaseTag + "/" + archiveName
	checksumURL := strings.TrimRight(releaseBaseURL, "/") + "/" + releaseTag + "/apimux_" + releaseVersion + "_checksums.txt"
	result.ArchiveURL = archiveURL

	targetPath := strings.TrimSpace(executablePath)
	if targetPath == "" {
		targetPath, err = os.Executable()
		if err != nil {
			return result, err
		}
	}
	targetPath, err = filepath.Abs(targetPath)
	if err != nil {
		return result, err
	}
	result.ExecutablePath = targetPath
	if err := validateReplaceTarget(targetPath); err != nil {
		return result, err
	}

	tmpDir, err := os.MkdirTemp("", "apimux-upgrade-*")
	if err != nil {
		return result, err
	}
	defer os.RemoveAll(tmpDir)

	archivePath := filepath.Join(tmpDir, archiveName)
	checksumPath := filepath.Join(tmpDir, "checksums.txt")
	if err := downloadFile(ctx, client, archiveURL, archivePath); err != nil {
		return result, err
	}
	if err := downloadFile(ctx, client, checksumURL, checksumPath); err != nil {
		return result, err
	}
	if err := verifyChecksum(archivePath, checksumPath, archiveName); err != nil {
		return result, err
	}
	extractedPath, err := extractBinary(archivePath, tmpDir)
	if err != nil {
		return result, err
	}
	if err := replaceExecutable(targetPath, extractedPath); err != nil {
		return result, err
	}

	result.Status = "upgraded"
	result.UpdateAvailable = false
	result.StatusMessage = "apimux upgraded to " + check.LatestVersion
	return result, nil
}

func compareVersions(current string, latest string) int {
	currentParts := parseVersion(current)
	latestParts := parseVersion(latest)
	if len(currentParts) == 0 || len(latestParts) == 0 {
		return strings.Compare(current, latest)
	}
	limit := len(currentParts)
	if len(latestParts) > limit {
		limit = len(latestParts)
	}
	for i := 0; i < limit; i++ {
		currentValue := 0
		if i < len(currentParts) {
			currentValue = currentParts[i]
		}
		latestValue := 0
		if i < len(latestParts) {
			latestValue = latestParts[i]
		}
		switch {
		case currentValue < latestValue:
			return -1
		case currentValue > latestValue:
			return 1
		}
	}
	return 0
}

func normalizeReleaseVersion(value string) (string, string, error) {
	version := strings.TrimSpace(value)
	if version == "" {
		return "", "", fmt.Errorf("release version is empty")
	}
	version = strings.TrimPrefix(version, "v")
	return version, "v" + version, nil
}

func releaseBaseURLFromManifest(manifestURL string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(manifestURL))
	if err != nil {
		return "", err
	}
	const suffix = "/latest/download/latest.json"
	if !strings.HasSuffix(parsed.Path, suffix) {
		return "", fmt.Errorf("cannot derive release download URL from manifest URL %s", manifestURL)
	}
	parsed.Path = strings.TrimSuffix(parsed.Path, suffix) + "/download"
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String(), nil
}

func archiveName(version string) (string, error) {
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	switch goos {
	case "darwin", "linux":
	default:
		return "", fmt.Errorf("unsupported platform: %s/%s", goos, goarch)
	}
	switch goarch {
	case "amd64", "arm64":
	default:
		return "", fmt.Errorf("unsupported platform: %s/%s", goos, goarch)
	}
	return fmt.Sprintf("apimux_%s_%s_%s.tar.gz", version, goos, goarch), nil
}

func validateReplaceTarget(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("current executable is a symlink; upgrade with your package manager or rerun the install script")
	}
	if strings.Contains(path, "/Cellar/") || strings.Contains(path, "/Homebrew/") {
		return fmt.Errorf("current executable appears to be package-manager managed; upgrade with your package manager or rerun the install script")
	}
	dir := filepath.Dir(path)
	probe, err := os.CreateTemp(dir, ".apimux-upgrade-probe-*")
	if err != nil {
		return fmt.Errorf("cannot write to install directory %s; upgrade with your package manager or rerun the install script: %w", dir, err)
	}
	probePath := probe.Name()
	_ = probe.Close()
	_ = os.Remove(probePath)
	return nil
}

func downloadFile(ctx context.Context, client *http.Client, fileURL string, dst string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fileURL, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download %s returned status %d", fileURL, resp.StatusCode)
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func verifyChecksum(archivePath string, checksumPath string, archiveName string) error {
	body, err := os.ReadFile(checksumPath)
	if err != nil {
		return err
	}
	var expected string
	for _, line := range strings.Split(string(body), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == archiveName {
			expected = fields[0]
			break
		}
	}
	if expected == "" {
		return fmt.Errorf("checksum entry not found for %s", archiveName)
	}
	archiveBody, err := os.ReadFile(archivePath)
	if err != nil {
		return err
	}
	actual := fmt.Sprintf("%x", sha256.Sum256(archiveBody))
	if !strings.EqualFold(expected, actual) {
		return fmt.Errorf("checksum mismatch for %s", archiveName)
	}
	return nil
}

func extractBinary(archivePath string, tmpDir string) (string, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer gzipReader.Close()
	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if filepath.Base(header.Name) != "apimux" || header.FileInfo().IsDir() {
			continue
		}
		outPath := filepath.Join(tmpDir, "apimux")
		out, err := os.OpenFile(outPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o755)
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(out, tarReader); err != nil {
			_ = out.Close()
			return "", err
		}
		if err := out.Close(); err != nil {
			return "", err
		}
		return outPath, nil
	}
	return "", fmt.Errorf("archive did not contain apimux binary")
}

func replaceExecutable(targetPath string, newBinaryPath string) error {
	dir := filepath.Dir(targetPath)
	tmp, err := os.CreateTemp(dir, ".apimux-upgrade-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	in, err := os.Open(newBinaryPath)
	if err != nil {
		_ = tmp.Close()
		return err
	}
	defer in.Close()
	if _, err := io.Copy(tmp, in); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmpPath, 0o755); err != nil {
		return err
	}
	return os.Rename(tmpPath, targetPath)
}

func parseVersion(value string) []int {
	value = strings.TrimSpace(strings.TrimPrefix(value, "v"))
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ".")
	out := make([]int, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			return nil
		}
		number, err := strconv.Atoi(part)
		if err != nil {
			return nil
		}
		out = append(out, number)
	}
	return out
}
