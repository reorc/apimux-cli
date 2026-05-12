package update

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestUpgradeDownloadsVerifiesAndReplacesBinary(t *testing.T) {
	archiveName, err := archiveName("1.2.0")
	if err != nil {
		t.Fatalf("archive name: %v", err)
	}
	newBinary := []byte("#!/bin/sh\necho upgraded\n")
	archiveBody := buildTarGz(t, "apimux", newBinary)
	archiveSum := fmt.Sprintf("%x", sha256.Sum256(archiveBody))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/releases/latest/download/latest.json":
			_, _ = w.Write([]byte(`{"latest_version":"1.2.0"}`))
		case "/releases/download/v1.2.0/" + archiveName:
			_, _ = w.Write(archiveBody)
		case "/releases/download/v1.2.0/apimux_1.2.0_checksums.txt":
			_, _ = fmt.Fprintf(w, "%s  %s\n", archiveSum, archiveName)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	targetPath := filepath.Join(t.TempDir(), "apimux")
	if err := os.WriteFile(targetPath, []byte("old binary"), 0o755); err != nil {
		t.Fatalf("write target: %v", err)
	}

	result, err := Upgrade(context.Background(), server.Client(), "1.1.0", server.URL+"/releases/latest/download/latest.json", targetPath)
	if err != nil {
		t.Fatalf("upgrade: %v", err)
	}
	if result.Status != "upgraded" || result.LatestVersion != "1.2.0" || result.UpdateAvailable {
		t.Fatalf("unexpected result: %#v", result)
	}
	got, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("read target: %v", err)
	}
	if string(got) != string(newBinary) {
		t.Fatalf("binary was not replaced: %q", string(got))
	}
	info, err := os.Stat(targetPath)
	if err != nil {
		t.Fatalf("stat target: %v", err)
	}
	if info.Mode().Perm()&0o111 == 0 {
		t.Fatalf("expected executable permissions, got %s", info.Mode().Perm())
	}
}

func TestUpgradeRejectsChecksumMismatch(t *testing.T) {
	archiveName, err := archiveName("1.2.0")
	if err != nil {
		t.Fatalf("archive name: %v", err)
	}
	archiveBody := buildTarGz(t, "apimux", []byte("new binary"))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/releases/latest/download/latest.json":
			_, _ = w.Write([]byte(`{"latest_version":"1.2.0"}`))
		case "/releases/download/v1.2.0/" + archiveName:
			_, _ = w.Write(archiveBody)
		case "/releases/download/v1.2.0/apimux_1.2.0_checksums.txt":
			_, _ = fmt.Fprintf(w, "%064x  %s\n", 0, archiveName)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	targetPath := filepath.Join(t.TempDir(), "apimux")
	if err := os.WriteFile(targetPath, []byte("old binary"), 0o755); err != nil {
		t.Fatalf("write target: %v", err)
	}

	_, err = Upgrade(context.Background(), server.Client(), "1.1.0", server.URL+"/releases/latest/download/latest.json", targetPath)
	if err == nil || !strings.Contains(err.Error(), "checksum mismatch") {
		t.Fatalf("expected checksum mismatch, got %v", err)
	}
	got, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("read target: %v", err)
	}
	if string(got) != "old binary" {
		t.Fatalf("target changed after failed upgrade: %q", string(got))
	}
}

func TestUpgradeRejectsSymlinkTarget(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permissions vary on Windows")
	}
	dir := t.TempDir()
	realPath := filepath.Join(dir, "apimux-real")
	linkPath := filepath.Join(dir, "apimux")
	if err := os.WriteFile(realPath, []byte("old binary"), 0o755); err != nil {
		t.Fatalf("write target: %v", err)
	}
	if err := os.Symlink(realPath, linkPath); err != nil {
		t.Fatalf("symlink target: %v", err)
	}
	err := validateReplaceTarget(linkPath)
	if err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("expected symlink error, got %v", err)
	}
}

func buildTarGz(t *testing.T, name string, body []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	if err := tw.WriteHeader(&tar.Header{
		Name: name,
		Mode: 0o755,
		Size: int64(len(body)),
	}); err != nil {
		t.Fatalf("write tar header: %v", err)
	}
	if _, err := tw.Write(body); err != nil {
		t.Fatalf("write tar body: %v", err)
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("close tar: %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("close gzip: %v", err)
	}
	return buf.Bytes()
}
