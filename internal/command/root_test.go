package command

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInvalidOutputReturnsNonZero(t *testing.T) {
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
	if !strings.Contains(stdout.String(), `"cli_invalid_output"`) {
		t.Fatalf("expected canonical local error, got %s", stdout.String())
	}
}

func TestUnknownFormatFlagReturnsNonZero(t *testing.T) {
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
	if !strings.Contains(stdout.String(), `"cli_invalid_flags"`) {
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

func TestSchemaCapabilitiesOutputsNames(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/schema" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":[{"name":"amazon.get_product"},{"name":"reddit.search"}]}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"schema", "capabilities",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	output := stdout.String()
	if !strings.Contains(output, "amazon.get_product\n") {
		t.Fatalf("expected one name per line, got %s", output)
	}
	if !strings.Contains(output, "reddit.search\n") {
		t.Fatalf("expected one name per line, got %s", output)
	}
	if strings.Contains(output, "{") || strings.Contains(output, "[") {
		t.Fatalf("expected plain text output, not JSON, got %s", output)
	}
}

func TestSchemaCapabilitiesJSONFlag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/schema" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":[{"name":"amazon.get_product"},{"name":"reddit.search"}]}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"schema", "capabilities", "--json",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	output := stdout.String()
	if !strings.Contains(output, `["amazon.get_product","reddit.search"]`) {
		t.Fatalf("expected JSON array output, got %s", output)
	}
}

func TestCompletionZshGeneratesScript(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"completion", "zsh",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "#compdef apimux") {
		t.Fatalf("expected zsh completion script, got %s", stdout.String())
	}
}

func TestCompletionBashGeneratesScript(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"completion", "bash",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "bash completion for apimux") {
		t.Fatalf("expected bash completion script, got %s", stdout.String())
	}
}

func TestCompletionListsSourceSubcommands(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"__complete", "amazon", "",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "query_aba_keywords") {
		t.Fatalf("expected source subcommands in completion output, got %s", stdout.String())
	}
}

func TestCompletionRequiresSupportedShell(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"completion",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode == 0 {
		t.Fatalf("expected non-zero exit code")
	}
	if !strings.Contains(stdout.String(), `"cli_invalid_command"`) {
		t.Fatalf("expected canonical local error, got %s", stdout.String())
	}
}
