package update

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
