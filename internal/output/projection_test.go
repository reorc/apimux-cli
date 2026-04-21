package output

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestParseFormat(t *testing.T) {
	tests := []struct {
		value string
		want  Format
		ok    bool
	}{
		{"", FormatAuto, true},
		{"auto", FormatAuto, true},
		{"data", FormatData, true},
		{"compact", FormatCompact, true},
		{"columnar", FormatCompact, true},
		{"weird", "", false},
	}

	for _, tt := range tests {
		got, ok := ParseFormat(tt.value)
		if got != tt.want || ok != tt.ok {
			t.Fatalf("ParseFormat(%q) = (%q,%v), want (%q,%v)", tt.value, got, ok, tt.want, tt.ok)
		}
	}
}

func TestProjectCapabilityCompactAmazonGetProduct(t *testing.T) {
	payload := json.RawMessage(`{
		"asin":"B001",
		"title":"Desk Lamp",
		"brand_store":{"text":"Acme"},
		"price":{"display":"$19.99"},
		"buybox":{"availability":"In Stock"},
		"rating":4.5,
		"review_count":123,
		"main_image":"https://example.com/main.jpg",
		"images":[{"variant":"MAIN","link":"a"},{"variant":"PT01","link":"b"}],
		"variants":[{"asin":"B001-A"},{"asin":"B001-B"}],
		"feature_bullets":["one","two","three","four"]
	}`)

	body, err := projectCapability("amazon.get_product", payload, FormatCompact)
	if err != nil {
		t.Fatalf("projectCapability() error = %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("unmarshal compact projection: %v", err)
	}
	if got["brand"] != "Acme" || got["image_count"] != float64(2) || got["variant_count"] != float64(2) {
		t.Fatalf("unexpected compact projection: %#v", got)
	}
	bullets, _ := got["feature_bullets"].([]any)
	if len(bullets) != 3 {
		t.Fatalf("expected first 3 bullets, got %#v", got["feature_bullets"])
	}
	if _, exists := got["images"]; !exists {
		t.Fatalf("expected compact projection to include embedded images table, got %#v", got)
	}
}

func TestProjectCapabilityColumnarSearchProducts(t *testing.T) {
	payload := json.RawMessage(`[
		{"position":1,"asin":"A1","title":"Desk Lamp","price":{"display":"$19.99"},"rating":4.5,"review_count":10},
		{"position":2,"asin":"A2","title":"Floor Lamp","price":{"display":"$29.99"},"rating":4.7,"review_count":20}
	]`)

	body, err := projectCapability("amazon.search_products", payload, FormatColumnar)
	if err != nil {
		t.Fatalf("projectCapability() error = %v", err)
	}

	var got struct {
		Columns []string `json:"columns"`
		Rows    [][]any  `json:"rows"`
	}
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("unmarshal columnar projection: %v", err)
	}
	if len(got.Columns) != 6 || len(got.Rows) != 2 {
		t.Fatalf("unexpected table projection: %#v", got)
	}
	if got.Columns[0] != "position" || got.Columns[3] != "price" {
		t.Fatalf("unexpected columns: %#v", got.Columns)
	}
}

func TestProjectCapabilityColumnarGoogleTrends(t *testing.T) {
	payload := json.RawMessage(`{
		"search_parameters":{"q":"AI","geo":"US","time":"today 12-m"},
		"timeline_data":[
			{"date":"Jan","timestamp":"1","values":[{"query":"AI","value":50}]},
			{"date":"Feb","timestamp":"2","values":[{"query":"AI","value":60}]}
		]
	}`)

	body, err := projectCapability("google_trends.get_interest_over_time", payload, FormatColumnar)
	if err != nil {
		t.Fatalf("projectCapability() error = %v", err)
	}

	text := string(body)
	if !strings.Contains(text, `"timeline":{"columns":["date","timestamp","value"]`) {
		t.Fatalf("expected columnar timeline projection, got %s", text)
	}
	if !strings.Contains(text, `"query":"AI"`) {
		t.Fatalf("expected scalar context preserved, got %s", text)
	}
}

func TestProjectCapabilityUnsupported(t *testing.T) {
	_, err := projectCapability("unsupported.capability", json.RawMessage(`[]`), FormatCompact)
	var unsupported *UnsupportedProjectionError
	if err == nil || !strings.Contains(err.Error(), "not supported") {
		t.Fatalf("expected unsupported projection error, got %v", err)
	}
	if !errors.As(err, &unsupported) {
		t.Fatalf("expected UnsupportedProjectionError, got %T", err)
	}
}

func TestProjectCapabilityAutoFallsBackToDataWhenUnsupported(t *testing.T) {
	payload := json.RawMessage(`[{"id":"p1","title":"hello"}]`)
	body, err := projectCapability("unsupported.capability", payload, FormatAuto)
	if err != nil {
		t.Fatalf("projectCapability() error = %v", err)
	}
	if string(body) != string(payload) {
		t.Fatalf("expected raw payload fallback, got %s", string(body))
	}
}

func TestProjectionRulesCoverAllAgentTestCapabilities(t *testing.T) {
	required := []string{
		"amazon.get_product",
		"amazon.get_keyword_overview",
		"amazon.get_asins_sales_history",
		"douyin.search_videos",
		"douyin.get_video_detail",
		"douyin.get_video_comments",
		"google_ads.search_advertisers",
		"google_ads.list_ad_creatives",
		"google_ads.get_ad_details",
		"meta_ads.search_ads",
		"meta_ads.get_ad_detail",
		"reddit.search",
		"reddit.get_subreddit_feed",
		"tiktok.search_videos",
		"tiktok.list_comments",
		"xiaohongshu.search_notes",
		"xiaohongshu.get_note_detail",
		"amazon.search_category",
		"amazon.search_products",
		"amazon.expand_keywords",
		"amazon.get_product_reviews",
		"amazon.get_keyword_trends",
		"google_trends.get_interest_over_time",
		"reddit.get_post_detail",
		"reddit.get_post_comments",
		"tiktok.shop_products",
		"tiktok.shop_product_info",
		"xiaohongshu.get_note_comments",
		"douyin.get_comment_replies",
	}
	if len(projectionRules) != len(required) {
		t.Fatalf("projectionRules size = %d, want %d", len(projectionRules), len(required))
	}
	for _, capability := range required {
		if _, ok := projectionRules[capability]; !ok {
			t.Fatalf("missing projection rule for %s", capability)
		}
	}
}
