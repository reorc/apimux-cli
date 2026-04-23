package command

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
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
