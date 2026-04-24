package command

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"

	"github.com/reorc/apimux-cli/internal/config"
)

func TestAuthLoginNoBrowserSavesKey(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("APIMUX_CONFIG_DIR", tempDir)

	pollCalls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/cli-auth/start":
			_, _ = w.Write([]byte(`{"device_code":"dev-1","user_code":"ABCD-EFGH","verification_uri":"https://apimux.test/auth/cli?user_code=ABCD-EFGH","verification_uri_complete":"https://apimux.test/auth/cli?user_code=ABCD-EFGH","expires_in":900,"interval":0}`))
		case "/api/cli-auth/poll":
			pollCalls++
			if pollCalls == 1 {
				w.WriteHeader(http.StatusAccepted)
				_, _ = w.Write([]byte(`{"status":"authorization_pending","interval":0}`))
				return
			}
			_, _ = w.Write([]byte(`{"api_key":"saved-key","api_key_id":"key-1"}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root := NewRoot(&stdout, &stderr)

	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"auth", "login", "--no-browser", "--device-name", "test-box",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.APIKey != "saved-key" {
		t.Fatalf("expected saved API key, got %#v", cfg)
	}
	if !strings.Contains(stdout.String(), "Visit: https://apimux.test/auth/cli?user_code=ABCD-EFGH") {
		t.Fatalf("expected verification URL in stdout, got %s", stdout.String())
	}
	if !strings.Contains(stdout.String(), "Authorized. API key saved") {
		t.Fatalf("expected success message, got %s", stdout.String())
	}
}

func TestAuthLoginSSHEnvironmentUsesNoBrowserBehavior(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("APIMUX_CONFIG_DIR", tempDir)
	t.Setenv("SSH_CONNECTION", "client server")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/cli-auth/start":
			_, _ = w.Write([]byte(`{"device_code":"dev-1","user_code":"ABCD-EFGH","verification_uri":"https://apimux.test/auth/cli?user_code=ABCD-EFGH","verification_uri_complete":"https://apimux.test/auth/cli?user_code=ABCD-EFGH","expires_in":900,"interval":0}`))
		case "/api/cli-auth/poll":
			_, _ = w.Write([]byte(`{"api_key":"saved-key","api_key_id":"key-1"}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root := NewRoot(&stdout, &stderr)

	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"auth", "login",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Open the URL above in your browser to approve access.") {
		t.Fatalf("expected no-browser guidance, got %s", stdout.String())
	}
	if strings.Contains(stderr.String(), "open browser failed") {
		t.Fatalf("expected browser open to be skipped, got stderr %s", stderr.String())
	}
}

func TestHeadlessAuthEnvDetection(t *testing.T) {
	t.Run("ssh", func(t *testing.T) {
		t.Setenv("SSH_CONNECTION", "client server")
		if !isHeadlessAuthEnv() {
			t.Fatal("expected SSH_CONNECTION to be headless")
		}
	})

	t.Run("ci", func(t *testing.T) {
		t.Setenv("CI", "true")
		if !isHeadlessAuthEnv() {
			t.Fatal("expected CI=true to be headless")
		}
	})

	t.Run("no browser", func(t *testing.T) {
		t.Setenv("NO_BROWSER", "1")
		if !isHeadlessAuthEnv() {
			t.Fatal("expected NO_BROWSER=1 to be headless")
		}
	})

	t.Run("linux without display", func(t *testing.T) {
		if runtime.GOOS != "linux" {
			t.Skip("linux-specific detection")
		}
		t.Setenv("DISPLAY", "")
		t.Setenv("WSL_DISTRO_NAME", "")
		t.Setenv("WSL_INTEROP", "")
		if !isHeadlessAuthEnv() {
			t.Fatal("expected linux without DISPLAY to be headless")
		}
	})
}

func TestAuthLoginWebURLDoesNotPersistBaseURL(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("APIMUX_CONFIG_DIR", tempDir)

	if err := config.Save(config.Config{BaseURL: "http://service.example"}); err != nil {
		t.Fatalf("seed config: %v", err)
	}

	pollCalls := 0
	webServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/cli-auth/start":
			_, _ = w.Write([]byte(`{"device_code":"dev-1","user_code":"ABCD-EFGH","verification_uri":"https://apimux.test/auth/cli?user_code=ABCD-EFGH","verification_uri_complete":"https://apimux.test/auth/cli?user_code=ABCD-EFGH","expires_in":900,"interval":0}`))
		case "/api/cli-auth/poll":
			pollCalls++
			_, _ = w.Write([]byte(`{"api_key":"saved-key","api_key_id":"key-1"}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer webServer.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root := NewRoot(&stdout, &stderr)

	exitCode, err := root.Execute(context.Background(), []string{
		"auth", "login", "--no-browser", "--web-url", webServer.URL,
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	if pollCalls == 0 {
		t.Fatalf("expected polling against web URL")
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.BaseURL != "http://service.example" {
		t.Fatalf("expected service base URL to remain unchanged, got %#v", cfg)
	}
	if cfg.APIKey != "saved-key" {
		t.Fatalf("expected saved API key, got %#v", cfg)
	}
}
