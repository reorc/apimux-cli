package output

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestWriteCapabilityResponseJSONModeWritesDataOnlyWithNewline(t *testing.T) {
	var stdout bytes.Buffer
	renderer := Renderer{Stdout: &stdout}

	if err := renderer.WriteCapabilityResponse([]byte(`{"ok":true,"data":{"asin":"B001","title":"Desk Lamp","brand_store":{"text":"Acme"},"buybox":{"availability":"In Stock"},"images":[{"link":"a"},{"link":"b"}],"variants":[{"asin":"A"},{"asin":"B"}],"feature_bullets":["one","two","three","four"]},"meta":{"capability":"amazon.get_product","source":"upstream"}}`), BodyOutputCompact, false); err != nil {
		t.Fatalf("WriteCapabilityResponse() error = %v", err)
	}

	got := stdout.String()
	if !bytes.Contains([]byte(got), []byte(`"images":{"columns":["variant","link"]`)) || !bytes.Contains([]byte(got), []byte(`"variants":{"columns":["asin","title","dimensions"]`)) {
		t.Fatalf("expected compact output with embedded tables, got %q", got)
	}
}

func TestWriteCapabilityResponsePrettyModeIndentsDataOnly(t *testing.T) {
	var stdout bytes.Buffer
	renderer := Renderer{Stdout: &stdout}

	if err := renderer.WriteCapabilityResponse([]byte(`{"ok":true,"data":{"id":"1","nested":{"ok":true}},"meta":{"capability":"amazon.get_product"}}`), BodyOutputDataPretty, false); err != nil {
		t.Fatalf("WriteCapabilityResponse() error = %v", err)
	}

	want := "{\n  \"id\": \"1\",\n  \"nested\": {\n    \"ok\": true\n  }\n}\n"
	if got := stdout.String(); got != want {
		t.Fatalf("unexpected pretty output:\n got: %q\nwant: %q", got, want)
	}
}

func TestWriteCapabilityResponseAutoFallsBackToDataWhenProjectionUnsupported(t *testing.T) {
	var stdout bytes.Buffer
	renderer := Renderer{Stdout: &stdout}

	if err := renderer.WriteCapabilityResponse([]byte(`{"ok":true,"data":{"id":"1","nested":{"ok":true}},"meta":{"capability":"unsupported.capability"}}`), BodyOutputCompact, false); err != nil {
		t.Fatalf("WriteCapabilityResponse() error = %v", err)
	}

	if got, want := stdout.String(), "{\"id\":\"1\",\"nested\":{\"ok\":true}}\n"; got != want {
		t.Fatalf("unexpected fallback output:\n got: %q\nwant: %q", got, want)
	}
}

func TestWriteCapabilityResponseDebugModePreservesEnvelopeWithoutProviderSource(t *testing.T) {
	var stdout bytes.Buffer
	renderer := Renderer{Stdout: &stdout}

	if err := renderer.WriteCapabilityResponse([]byte(`{"ok":true,"data":{"id":"1"},"meta":{"capability":"amazon.get_product","source":"upstream","cache_hit":false}}`), BodyOutputCompact, true); err != nil {
		t.Fatalf("WriteCapabilityResponse() error = %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal debug output: %v", err)
	}
	meta, _ := got["meta"].(map[string]any)
	if got["ok"] != true {
		t.Fatalf("expected ok=true in debug output, got %#v", got)
	}
	if _, exists := meta["source"]; exists {
		t.Fatalf("expected provider source to be stripped, got %#v", got)
	}
	if meta["capability"] != "amazon.get_product" {
		t.Fatalf("expected capability metadata, got %#v", got)
	}
}

func TestWriteCapabilityResponseErrorModeWritesErrorOnly(t *testing.T) {
	var stdout bytes.Buffer
	renderer := Renderer{Stdout: &stdout}

	if err := renderer.WriteCapabilityResponse([]byte(`{"ok":false,"error":{"code":"invalid_request","message":"bad input"},"meta":{"capability":"amazon.get_product","source":"upstream"}}`), BodyOutputCompact, false); err != nil {
		t.Fatalf("WriteCapabilityResponse() error = %v", err)
	}

	want := "{\"error\":{\"code\":\"invalid_request\",\"message\":\"bad input\"}}\n"
	if got := stdout.String(); got != want {
		t.Fatalf("unexpected error output:\n got: %q\nwant: %q", got, want)
	}
}

func TestWriteLocalErrorWritesJSONErrorObject(t *testing.T) {
	var stdout bytes.Buffer
	renderer := Renderer{Stdout: &stdout}

	if err := renderer.WriteLocalError("request failed", "transport_error"); err != nil {
		t.Fatalf("WriteLocalError() error = %v", err)
	}

	want := "{\"error\":{\"code\":\"transport_error\",\"message\":\"request failed\",\"type\":\"internal\"}}\n"
	if got := stdout.String(); got != want {
		t.Fatalf("unexpected local error output:\n got: %q\nwant: %q", got, want)
	}
}
