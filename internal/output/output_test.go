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
	if !bytes.Contains([]byte(got), []byte(`"images":[{"link":"a"},{"link":"b"}]`)) || !bytes.Contains([]byte(got), []byte(`"variants":[{"asin":"A"},{"asin":"B"}]`)) {
		t.Fatalf("expected compact output with nested arrays intact, got %q", got)
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

func TestWriteCapabilityResponseCompactModeExposesPaginationMetadata(t *testing.T) {
	var stdout bytes.Buffer
	renderer := Renderer{Stdout: &stdout}

	input := `{
		"ok": true,
		"data": [
			{"post_id": "abc123", "title": "First Post", "score": 100},
			{"post_id": "def456", "title": "Second Post", "score": 50}
		],
		"meta": {
			"capability": "reddit.search",
			"cursor": "t3_xyz789",
			"has_more": true
		}
	}`

	if err := renderer.WriteCapabilityResponse([]byte(input), BodyOutputCompact, false); err != nil {
		t.Fatalf("WriteCapabilityResponse() error = %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal compact output: %v", err)
	}

	// Verify data is present
	if _, ok := got["data"]; !ok {
		t.Fatalf("expected data field in compact output, got %#v", got)
	}

	// Verify pagination metadata is exposed
	meta, ok := got["meta"].(map[string]any)
	if !ok {
		t.Fatalf("expected meta field in compact output, got %#v", got)
	}
	if meta["cursor"] != "t3_xyz789" {
		t.Fatalf("expected cursor in meta, got %#v", meta)
	}
	if meta["has_more"] != true {
		t.Fatalf("expected has_more=true in meta, got %#v", meta)
	}
}

func TestWriteCapabilityResponseCompactModeExposesPartialFailureMetadata(t *testing.T) {
	var stdout bytes.Buffer
	renderer := Renderer{Stdout: &stdout}

	input := `{
		"ok": true,
		"data": [
			{"month":"2026-01","sales_volume":1000,"brand_count":50,"avg_price":null}
		],
		"meta": {
			"capability": "amazon.get_category_trend",
			"partial": true,
			"subrequest_count": 3,
			"subrequests": [
				{"name": "sales_volume", "ok": true},
				{"name": "brand_count", "ok": true},
				{"name": "avg_price", "ok": false, "error": {"code": "upstream_timeout"}}
			]
		}
	}`

	if err := renderer.WriteCapabilityResponse([]byte(input), BodyOutputCompact, false); err != nil {
		t.Fatalf("WriteCapabilityResponse() error = %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal compact output: %v", err)
	}

	// Verify projected data is present
	data, ok := got["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected wrapped data field in compact output, got %#v", got)
	}
	if _, ok := data["items"]; !ok {
		t.Fatalf("expected projected items field in compact output, got %#v", got)
	}

	// Verify partial-failure metadata is exposed
	meta, ok := got["meta"].(map[string]any)
	if !ok {
		t.Fatalf("expected meta field in compact output, got %#v", got)
	}
	if meta["partial"] != true {
		t.Fatalf("expected partial=true in meta, got %#v", meta)
	}
	if meta["subrequest_count"] != float64(3) {
		t.Fatalf("expected subrequest_count=3 in meta, got %#v", meta)
	}
	if _, ok := meta["subrequests"]; !ok {
		t.Fatalf("expected subrequests in meta, got %#v", meta)
	}
}

func TestWriteCapabilityResponseCompactModeWithoutMetadataOmitsMeta(t *testing.T) {
	var stdout bytes.Buffer
	renderer := Renderer{Stdout: &stdout}

	input := `{
		"ok": true,
		"data": {"asin": "B001", "title": "Product"},
		"meta": {
			"capability": "amazon.get_product"
		}
	}`

	if err := renderer.WriteCapabilityResponse([]byte(input), BodyOutputCompact, false); err != nil {
		t.Fatalf("WriteCapabilityResponse() error = %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal compact output: %v", err)
	}

	// When no critical metadata exists, output should be data only (no meta wrapper)
	if _, ok := got["meta"]; ok {
		t.Fatalf("expected no meta field when no critical metadata present, got %#v", got)
	}
	if _, ok := got["data"]; ok {
		t.Fatalf("expected unwrapped data when no critical metadata present, got %#v", got)
	}
}

func TestWriteCapabilityResponseCompactAmazonGetCategoryTrendUsesRequestedMetrics(t *testing.T) {
	var stdout bytes.Buffer
	renderer := Renderer{Stdout: &stdout}

	input := `{
		"ok": true,
		"data": [{"month":"2026-01","sales_volume":1200,"brand_count":15,"seller_count":8}],
		"meta": {
			"capability": "amazon.get_category_trend",
			"metrics": ["sales_volume", "brand_count"]
		}
	}`

	if err := renderer.WriteCapabilityResponse([]byte(input), BodyOutputCompact, false); err != nil {
		t.Fatalf("WriteCapabilityResponse() error = %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal compact output: %v", err)
	}

	items, ok := got["items"].(map[string]any)
	if !ok {
		t.Fatalf("expected items table, got %#v", got)
	}
	columns, _ := items["columns"].([]any)
	if len(columns) != 3 || columns[0] != "month" || columns[1] != "sales_volume" || columns[2] != "brand_count" {
		t.Fatalf("unexpected compact category trend columns: %#v", columns)
	}
}
