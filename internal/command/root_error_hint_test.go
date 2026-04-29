package command

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServiceErrorHintRemainsVisibleInCLIOutput(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/amazon.get_product" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"ok":false,"error":{"type":"validation","code":"invalid_asin","message":"asin must be a 10-character string","hint":"Check if the ASIN format is valid and retry."},"meta":{"capability":"amazon.get_product"}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"amazon", "get_product",
		"--asin", "B08N5WRWNW",
		"--market", "US",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 2 {
		t.Fatalf("expected exit code 2, got %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"hint":"Check if the ASIN format is valid and retry."`) {
		t.Fatalf("expected hint in stdout, got %s", stdout.String())
	}
}

func TestServiceNonJSONHTTPErrorReturnsStructuredCLIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/trendcloud.get_market_trend" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`error loading upstream trendcloud data`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"trendcloud", "get_market_trend",
		"--metrics", "sales,volume",
		"--filters-json", `{"platforms":["douyin"],"brands":["yili"]}`,
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"code":"cli_http_error"`) || !strings.Contains(stdout.String(), `"status_code":502`) {
		t.Fatalf("expected structured HTTP error, got %s", stdout.String())
	}
	if strings.Contains(stdout.String(), "invalid character") {
		t.Fatalf("expected HTTP error instead of raw JSON parse error, got %s", stdout.String())
	}
}

func TestServiceJSONHTTPErrorWithoutEnvelopeReturnsStructuredCLIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/amazon.get_product" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"message":"upstream temporarily unavailable"}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"amazon", "get_product",
		"--asin", "B08N5WRWNW",
		"--market", "US",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"code":"cli_http_error"`) || !strings.Contains(stdout.String(), `"body":{"message":"upstream temporarily unavailable"}`) {
		t.Fatalf("expected structured HTTP error with JSON body, got %s", stdout.String())
	}
}
