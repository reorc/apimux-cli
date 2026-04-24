package command

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/reorc/apimux-cli/internal/buildinfo"
)

func TestVersionCommandPrintsBuildMetadata(t *testing.T) {
	originalVersion := buildinfo.Version
	originalCommit := buildinfo.Commit
	originalBuildDate := buildinfo.BuildDate
	t.Cleanup(func() {
		buildinfo.Version = originalVersion
		buildinfo.Commit = originalCommit
		buildinfo.BuildDate = originalBuildDate
	})

	buildinfo.Version = "1.2.3"
	buildinfo.Commit = "abc1234"
	buildinfo.BuildDate = "2026-04-17T19:00:00Z"

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{"version"})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"version": "1.2.3"`) {
		t.Fatalf("expected version metadata, got %s", stdout.String())
	}
	if !strings.Contains(stdout.String(), `"commit": "abc1234"`) {
		t.Fatalf("expected commit metadata, got %s", stdout.String())
	}
}

func TestConfigInitIsDeprecatedOnboardingAlias(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("APIMUX_CONFIG_DIR", tempDir)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{"config", "init"})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "apimux auth login") {
		t.Fatalf("expected auth login guidance, got %s", stdout.String())
	}
	if !strings.Contains(stdout.String(), "apimux config set --api-key") {
		t.Fatalf("expected manual config fallback guidance, got %s", stdout.String())
	}
}

func TestUpgradeCheckReportsVersionStatus(t *testing.T) {
	originalVersion := buildinfo.Version
	t.Cleanup(func() {
		buildinfo.Version = originalVersion
	})
	buildinfo.Version = "1.0.0"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"latest_version":"1.2.0"}`))
	}))
	defer server.Close()

	t.Setenv("APIMUX_CLI_RELEASE_MANIFEST_URL", server.URL)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{"upgrade", "--check"})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}

	var payload map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	if payload["status"] != "update_available" {
		t.Fatalf("unexpected status: %#v", payload)
	}
	if payload["latest_version"] != "1.2.0" {
		t.Fatalf("unexpected latest version: %#v", payload)
	}
}

func TestConfigInitFollowedByGoogleTrendsUsesSavedConfig(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("APIMUX_CONFIG_DIR", tempDir)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/schema/google_trends.get_interest_over_time":
			_, _ = w.Write([]byte(`{"ok":true,"data":{"name":"google_trends.get_interest_over_time","parameters":[{"name":"q","type":"string","required":true,"description":"Comma-separated query terms, up to 5 values"},{"name":"time","type":"string","description":"Approved Google Trends time range","default":"today 12-m"},{"name":"geo","type":"string","description":"ISO country or region code"},{"name":"cat","type":"string","description":"Google Trends category id"},{"name":"gprop","type":"string","description":"Google Trends property filter","enum":["images","news","froogle","youtube"]},{"name":"tz","type":"integer","description":"Timezone offset in minutes"}]}}`))
			return
		case "/v1/capabilities/google_trends.get_interest_over_time":
			if got := r.Header.Get("Authorization"); got != "Bearer saved-key" {
				t.Fatalf("unexpected authorization header: %s", got)
			}
			_, _ = w.Write([]byte(`{"ok":true,"data":{"search_parameters":{"q":"AI","geo":"US","time":"today 12-m"},"timeline_data":[{"date":"Jan","timestamp":"1","values":[{"query":"AI","value":50}]}]},"meta":{"capability":"google_trends.get_interest_over_time"}}`))
			return
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	var initStdout bytes.Buffer
	var initStderr bytes.Buffer

	initRoot := NewRoot(&initStdout, &initStderr)
	exitCode, err := initRoot.Execute(context.Background(), []string{"config", "set", "--base-url", server.URL, "--api-key", "saved-key"})
	if err != nil {
		t.Fatalf("execute init: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected init exit code 0, got %d, stderr=%s", exitCode, initStderr.String())
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root := NewRoot(&stdout, &stderr)
	exitCode, err = root.Execute(context.Background(), []string{
		"google_trends", "get_interest_over_time",
		"--q", "AI",
	})
	if err != nil {
		t.Fatalf("execute google trends: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"query":"AI"`) || !strings.Contains(stdout.String(), `"timeline":{"columns":[`) {
		t.Fatalf("expected default compact google trends output, got %s", stdout.String())
	}
	if strings.Contains(stdout.String(), `"meta"`) || strings.Contains(stdout.String(), `"ok"`) {
		t.Fatalf("expected data-only output without envelope fields, got %s", stdout.String())
	}
}
