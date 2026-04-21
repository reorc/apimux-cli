package command

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInvalidOutputModeReturnsNonZero(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--output", "jsno",
		"schema", "list",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode == 0 {
		t.Fatalf("expected non-zero exit code")
	}
	if !strings.Contains(stdout.String(), `"cli_invalid_output_mode"`) {
		t.Fatalf("expected canonical local error, got %s", stdout.String())
	}
}

func TestInvalidFormatReturnsNonZero(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--format", "wide",
		"schema", "list",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode == 0 {
		t.Fatalf("expected non-zero exit code")
	}
	if !strings.Contains(stdout.String(), `"cli_invalid_format"`) {
		t.Fatalf("expected canonical local error, got %s", stdout.String())
	}
}

func TestSchemaListCallsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/schema" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":{"capabilities":[{"name":"amazon.get_product"}]}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"schema", "list",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"amazon.get_product"`) {
		t.Fatalf("expected service-backed schema payload, got %s", stdout.String())
	}
}
