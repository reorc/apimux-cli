package client

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestExecuteCapabilityUsesServiceBoundary(t *testing.T) {
	var gotPath string
	var gotMethod string
	var gotAuth string
	var gotParams map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		gotAuth = r.Header.Get("Authorization")
		var body struct {
			Params map[string]any `json:"params"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		gotParams = body.Params
		_, _ = w.Write([]byte(`{"ok":true,"data":{"asin":"B0CM5JV26D"}}`))
	}))
	defer server.Close()

	client := NewWithHTTPClient(Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	}, server.Client())

	resp, err := client.ExecuteCapability(context.Background(), "amazon.get_product", map[string]any{
		"asin":   "B0CM5JV26D",
		"market": "US",
	})
	if err != nil {
		t.Fatalf("execute capability: %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Fatalf("expected POST, got %s", gotMethod)
	}
	if gotPath != "/v1/capabilities/amazon.get_product" {
		t.Fatalf("unexpected path: %s", gotPath)
	}
	if gotAuth != "Bearer test-key" {
		t.Fatalf("unexpected auth header: %q", gotAuth)
	}
	if gotParams["asin"] != "B0CM5JV26D" || gotParams["market"] != "US" {
		t.Fatalf("unexpected params: %#v", gotParams)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}
}

func TestGetSchemaCoversSuccessNotFoundAndTimeout(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var gotPath string
		var gotAuth string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			gotAuth = r.Header.Get("Authorization")
			_, _ = w.Write([]byte(`{"name":"amazon.get_product"}`))
		}))
		defer server.Close()

		client := NewWithHTTPClient(Config{BaseURL: server.URL, APIKey: "test-key"}, server.Client())
		resp, err := client.GetSchema(context.Background(), "amazon.get_product")
		if err != nil {
			t.Fatalf("GetSchema() error = %v", err)
		}
		if gotPath != "/v1/schema/amazon.get_product" {
			t.Fatalf("unexpected path: %s", gotPath)
		}
		if gotAuth != "Bearer test-key" {
			t.Fatalf("unexpected auth header: %q", gotAuth)
		}
		if resp.StatusCode != http.StatusOK || string(resp.Body) != `{"name":"amazon.get_product"}` {
			t.Fatalf("unexpected response: %#v", resp)
		}
	})

	t.Run("not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"ok":false,"error":{"code":"unknown_capability"}}`))
		}))
		defer server.Close()

		client := NewWithHTTPClient(Config{BaseURL: server.URL}, server.Client())
		resp, err := client.GetSchema(context.Background(), "missing.capability")
		if err != nil {
			t.Fatalf("GetSchema() error = %v", err)
		}
		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("unexpected status: %d", resp.StatusCode)
		}
		if string(resp.Body) != `{"ok":false,"error":{"code":"unknown_capability"}}` {
			t.Fatalf("unexpected body: %s", string(resp.Body))
		}
	})

	t.Run("timeout", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond)
			_, _ = w.Write([]byte(`{"name":"late"}`))
		}))
		defer server.Close()

		httpClient := server.Client()
		httpClient.Timeout = 20 * time.Millisecond
		client := NewWithHTTPClient(Config{BaseURL: server.URL}, httpClient)

		_, err := client.GetSchema(context.Background(), "amazon.get_product")
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("expected deadline exceeded, got %v", err)
		}
	})
}

func TestExecuteCapabilityCoversSuccessAndErrorResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
	}{
		{name: "success", statusCode: http.StatusOK, body: `{"ok":true,"data":{"asin":"B001"}}`},
		{name: "error envelope", statusCode: http.StatusBadRequest, body: `{"ok":false,"error":{"code":"invalid_market","message":"market is invalid"}}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer server.Close()

			client := NewWithHTTPClient(Config{BaseURL: server.URL}, server.Client())
			resp, err := client.ExecuteCapability(context.Background(), "amazon.get_product", map[string]any{"asin": "B001"})
			if err != nil {
				t.Fatalf("ExecuteCapability() error = %v", err)
			}
			if resp.StatusCode != tc.statusCode {
				t.Fatalf("unexpected status: %d", resp.StatusCode)
			}
			if string(resp.Body) != tc.body {
				t.Fatalf("unexpected body: %s", string(resp.Body))
			}
		})
	}
}

func TestListSchemasUsesSchemaCollectionEndpoint(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"data":[]}`))
	}))
	defer server.Close()

	client := NewWithHTTPClient(Config{BaseURL: server.URL}, server.Client())
	resp, err := client.ListSchemas(context.Background())
	if err != nil {
		t.Fatalf("ListSchemas() error = %v", err)
	}
	if gotPath != "/v1/schema" {
		t.Fatalf("unexpected path: %s", gotPath)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}
}

func TestClientWrapsTransportErrors(t *testing.T) {
	client := NewWithHTTPClient(Config{BaseURL: "http://127.0.0.1:1"}, &http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, io.ErrUnexpectedEOF
		}),
	})

	_, err := client.ExecuteCapability(context.Background(), "amazon.get_product", map[string]any{"asin": "B001"})
	if err == nil {
		t.Fatalf("expected transport error")
	}
	if !strings.Contains(err.Error(), "request failed:") {
		t.Fatalf("unexpected error: %v", err)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
