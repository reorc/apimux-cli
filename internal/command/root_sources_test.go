package command

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func assertDataOnlyOutputContains(t *testing.T, stdout string, want string) {
	t.Helper()
	if !strings.Contains(stdout, want) {
		t.Fatalf("expected data-only output to contain %s, got %s", want, stdout)
	}
	if strings.Contains(stdout, `"ok"`) {
		t.Fatalf("expected data-only output without ok field, got %s", stdout)
	}
}

func assertCompactTableOutputContains(t *testing.T, stdout string, columns string, want string) {
	t.Helper()
	if !strings.Contains(stdout, columns) || !strings.Contains(stdout, want) {
		t.Fatalf("expected compact table output to contain %s and %s, got %s", columns, want, stdout)
	}
	if strings.Contains(stdout, `"ok"`) {
		t.Fatalf("expected compact output without ok field, got %s", stdout)
	}
}

func TestAmazonGetProductCallsServiceAndPrintsDataOnly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/amazon.get_product" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":{"asin":"B0CM5JV26D","title":"Desk Lamp","brand_store":{"text":"Acme"},"buybox":{"availability":"In Stock"},"images":[{"link":"a"},{"link":"b"}],"variants":[{"asin":"A"},{"asin":"B"}],"feature_bullets":["one","two","three","four"]},"meta":{"capability":"amazon.get_product"}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"amazon", "get_product",
		"--asin", "B0CM5JV26D",
		"--market", "US",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"images":[{"link":"a"},{"link":"b"}]`) || !strings.Contains(stdout.String(), `"variants":[{"asin":"A"},{"asin":"B"}]`) {
		t.Fatalf("expected default compact projection, got %s", stdout.String())
	}
	if strings.Contains(stdout.String(), `"meta"`) || strings.Contains(stdout.String(), `"ok"`) {
		t.Fatalf("expected compact output without envelope fields, got %s", stdout.String())
	}
}

func TestAmazonGetProductDebugPrintsSanitizedEnvelope(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/amazon.get_product" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":{"asin":"B0CM5JV26D","market":"US"},"meta":{"capability":"amazon.get_product","source":"upstream","cache_hit":false}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"--debug",
		"amazon", "get_product",
		"--asin", "B0CM5JV26D",
		"--market", "US",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"ok":true`) || !strings.Contains(stdout.String(), `"capability":"amazon.get_product"`) {
		t.Fatalf("expected sanitized debug envelope, got %s", stdout.String())
	}
	if strings.Contains(stdout.String(), `"source":"upstream"`) {
		t.Fatalf("expected provider source to be stripped in debug output, got %s", stdout.String())
	}
}

func TestAmazonGetProductCompactProjection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/amazon.get_product" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":{"asin":"B0CM5JV26D","title":"Desk Lamp","brand_store":{"text":"Acme"},"price":{"display":"$19.99"},"buybox":{"availability":"In Stock"},"rating":4.6,"review_count":100,"main_image":"https://example.com/main.jpg","images":[{"variant":"MAIN","link":"a"},{"variant":"PT01","link":"b"}],"variants":[{"asin":"A"},{"asin":"B"}],"feature_bullets":["one","two","three","four"]},"meta":{"capability":"amazon.get_product"}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"amazon", "get_product",
		"--asin", "B0CM5JV26D",
		"--market", "US",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"feature_bullets":["one","two","three","four"]`) {
		t.Fatalf("expected compact projection, got %s", stdout.String())
	}
	if !strings.Contains(stdout.String(), `"images":[{"link":"a","variant":"MAIN"},{"link":"b","variant":"PT01"}]`) {
		t.Fatalf("expected compact projection to keep image list, got %s", stdout.String())
	}
}

func TestGoogleTrendsOutputPrettyUsesCompactProjection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/google_trends.get_interest_over_time" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":{"search_parameters":{"q":"AI","geo":"US","time":"today 12-m"},"timeline_data":[{"date":"Jan","timestamp":"1","values":[{"query":"AI","value":50}]}]},"meta":{"capability":"google_trends.get_interest_over_time"}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"--output", "pretty",
		"google_trends", "get_interest_over_time",
		"--q", "AI",
		"--geo", "US",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"columns": [`) || !strings.Contains(stdout.String(), `"value"`) {
		t.Fatalf("expected pretty compact table output, got %s", stdout.String())
	}
	if !strings.Contains(stdout.String(), "\n  \"timeline\"") {
		t.Fatalf("expected pretty output indentation, got %s", stdout.String())
	}
}

func TestAmazonExpandKeywordsOutputCompact(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/amazon.expand_keywords" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":[{"keyword":"desk lamp"}],"meta":{"capability":"amazon.expand_keywords"}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"--output", "compact",
		"amazon", "expand_keywords",
		"--keyword", "desk lamp",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"columns":["keyword","est_searches_num","searches_rank","match_types"]`) {
		t.Fatalf("expected compact table output, got %s", stdout.String())
	}
}

func TestAmazonExpandKeywordsCallsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/amazon.expand_keywords" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":[{"keyword":"yoga mat","est_searches_num":12000}],"meta":{"capability":"amazon.expand_keywords","total":1}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"amazon", "expand_keywords",
		"--keyword", "yoga mat",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"columns":["keyword","est_searches_num","searches_rank","match_types"]`) || !strings.Contains(stdout.String(), `"yoga mat"`) {
		t.Fatalf("expected compact keyword table output, got %s", stdout.String())
	}
}

func TestAmazonGetProductFormatDataExplicitlyDisablesDefaultCompact(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/amazon.get_product" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":{"asin":"B0CM5JV26D","market":"US","images":[{"link":"a"},{"link":"b"}],"variants":[{"asin":"A"},{"asin":"B"}]},"meta":{"capability":"amazon.get_product"}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"--output", "data",
		"amazon", "get_product",
		"--asin", "B0CM5JV26D",
		"--market", "US",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"images":[`) || !strings.Contains(stdout.String(), `"variants":[`) {
		t.Fatalf("expected explicit data-only output, got %s", stdout.String())
	}
	if strings.Contains(stdout.String(), `"image_count"`) || strings.Contains(stdout.String(), `"variant_count"`) {
		t.Fatalf("expected compact projection to be disabled, got %s", stdout.String())
	}
}

func TestAmazonGetProductOutputDataPrettyPrettyPrintsRawData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/amazon.get_product" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":{"asin":"B0CM5JV26D","market":"US","images":[{"link":"a"}]},"meta":{"capability":"amazon.get_product"}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"--output", "data-pretty",
		"amazon", "get_product",
		"--asin", "B0CM5JV26D",
		"--market", "US",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "\n  \"asin\": \"B0CM5JV26D\"") {
		t.Fatalf("expected pretty raw data output, got %s", stdout.String())
	}
	if strings.Contains(stdout.String(), `"image_count"`) {
		t.Fatalf("expected raw data output, got %s", stdout.String())
	}
}

func TestAmazonGetKeywordOverviewCallsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/amazon.get_keyword_overview" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":{"keyword":"yoga mat","est_searches_num":12000},"meta":{"capability":"amazon.get_keyword_overview"}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"amazon", "get_keyword_overview",
		"--keyword", "yoga mat",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	assertDataOnlyOutputContains(t, stdout.String(), `"est_searches_num":12000`)
}

func TestAmazonGetKeywordTrendsCallsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/amazon.get_keyword_trends" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":[{"keyword":"yoga mat"}],"meta":{"capability":"amazon.get_keyword_trends"}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"amazon", "get_keyword_trends",
		"--keywords", "yoga mat,pilates ring",
		"--granularity", "week",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"keyword":"yoga mat"`) {
		t.Fatalf("expected compact list output, got %s", stdout.String())
	}
}

func TestAmazonListASINKeywordsCallsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/amazon.list_asin_keywords" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":[{"keyword":"yoga mat"}],"meta":{"capability":"amazon.list_asin_keywords","total":1}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"amazon", "list_asin_keywords",
		"--asin", "B0CM5JV26D",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	assertCompactTableOutputContains(t, stdout.String(), `"columns":["keyword","kw_characters","conversion_characters","exposure_type","last_rank_str","ad_last_rank_str","est_searches_num","searches_rank","ratio_score"]`, `"yoga mat"`)
}

func TestAmazonQueryABAKeywordsCallsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/amazon.query_aba_keywords" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":[{"keyword":"yoga mat"}],"meta":{"capability":"amazon.query_aba_keywords","current_page":1,"has_more":false}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"amazon", "query_aba_keywords",
		"--keyword", "yoga mat",
		"--node-ids", "12345",
		"--page", "1",
		"--page-size", "40",
		"--market", "US",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	assertCompactTableOutputContains(t, stdout.String(), `"columns":["keyword","keyword_cn_name","rank","search_volume","word_count","product_count","rank_change_of_weekly","cpc","cpc_range","search_conversion_rate","search_conversion_rate_d90","click_conversion_rate_d90","click_of_90d","sales_volume_of_90d","share_click_rate","share_conversion_rate","search_volume_growth_rate_trend","top3_asin","top3_brand","top3_category","season","update"]`, `"yoga mat"`)
}

func TestAmazonGetASINSalesDailyTrendCallsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/amazon.get_asin_sales_daily_trend" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":[{"date":"2026-04-01","sales":12}],"meta":{"capability":"amazon.get_asin_sales_daily_trend"}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"amazon", "get_asin_sales_daily_trend",
		"--asin", "B0CM5JV26D",
		"--begin-date", "2026-04-01",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	assertCompactTableOutputContains(t, stdout.String(), `"columns":["date","sales"]`, `"2026-04-01"`)
}

func TestAmazonGetASINsSalesHistoryCallsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/amazon.get_asins_sales_history" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":[{"asin":"B0CM5JV26D","month":"2026-04","sales":12}],"meta":{"capability":"amazon.get_asins_sales_history"}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"amazon", "get_asins_sales_history",
		"--asins", "B0CM5JV26D,B0D1234567",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	assertCompactTableOutputContains(t, stdout.String(), `"columns":["asin","month","sales"]`, `"2026-04"`)
}

func TestAmazonGetVariantSales30dCallsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/amazon.get_variant_sales_30d" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":[{"asin":"B0CM5JV26D","bought_in_past_month":12}],"meta":{"capability":"amazon.get_variant_sales_30d"}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"amazon", "get_variant_sales_30d",
		"--asin", "B0CM5JV26D",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	assertCompactTableOutputContains(t, stdout.String(), `"columns":["asin","bought_in_past_month","update_time"]`, `"B0CM5JV26D"`)
}

func TestTikTokSearchVideosCallsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/tiktok.search_videos" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":[],"meta":{"capability":"tiktok.search_videos"}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"tiktok", "search_videos",
		"--keyword", "desk setup",
		"--sort-by", "likes",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	assertCompactTableOutputContains(t, stdout.String(), `"columns":["video_id","description","create_time","author","author_handle","play_count","like_count","comment_count","share_count"]`, `"rows":[]`)
}

func TestMetaAdsSearchCallsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/meta_ads.search_ads" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":[],"meta":{"capability":"meta_ads.search_ads","next_page_token":"tok"}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"meta_ads", "search_ads",
		"--q", "fitness app",
		"--media-type", "video",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	if got := stdout.String(); !strings.Contains(got, `"columns":["ad_id","page_name","start_date","end_date","collation_count","is_active"]`) || !strings.Contains(got, `"rows":[]`) {
		t.Fatalf("expected compact table output, got %s", got)
	}
}

func TestMetaAdsGetAdDetailCallsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/meta_ads.get_ad_detail" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":{"ad_id":"477570185419072"},"meta":{"capability":"meta_ads.get_ad_detail"}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"meta_ads", "get_ad_detail",
		"--ad-id", "477570185419072",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	assertDataOnlyOutputContains(t, stdout.String(), `"ad_id":"477570185419072"`)
}

func TestDouyinSearchVideosCallsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/douyin.search_videos" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":[],"meta":{"capability":"douyin.search_videos","cursor":20,"has_more":true}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"douyin", "search_videos",
		"--keyword", "desk setup",
		"--sort-type", "likes",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	if got := stdout.String(); !strings.Contains(got, `"columns":["aweme_id","description","create_time","author","like_count","comment_count","share_count"]`) || !strings.Contains(got, `"rows":[]`) {
		t.Fatalf("expected compact table output, got %s", got)
	}
}

func TestDouyinGetVideoDetailCallsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/douyin.get_video_detail" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":{"aweme_id":"7489123456789012345"},"meta":{"capability":"douyin.get_video_detail"}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"douyin", "get_video_detail",
		"--aweme-id", "7489123456789012345",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	assertDataOnlyOutputContains(t, stdout.String(), `"aweme_id":"7489123456789012345"`)
}

func TestDouyinGetCommentRepliesCallsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/douyin.get_comment_replies" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":[],"meta":{"capability":"douyin.get_comment_replies","cursor":20,"has_more":true}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"douyin", "get_comment_replies",
		"--aweme-id", "7489123456789012345",
		"--comment-id", "cmt_123",
		"--count", "10",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	if got := stdout.String(); !strings.Contains(got, `"columns":["comment_id","text","create_time","author","like_count"]`) || !strings.Contains(got, `"rows":[]`) {
		t.Fatalf("expected compact table output, got %s", got)
	}
}

func TestGoogleAdsSearchAdvertisersCallsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/google_ads.search_advertisers" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":{"advertisers":[],"domains":[]},"meta":{"capability":"google_ads.search_advertisers"}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"google_ads", "search_advertisers",
		"--query", "tesla",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	assertCompactTableOutputContains(t, stdout.String(), `"advertisers":{"columns":["advertiser_id","advertiser_name","ads_count","region","is_verified"]`, `"domains":[]`)
}

func TestGoogleAdsListAdCreativesCallsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/google_ads.list_ad_creatives" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var payload struct {
			Params map[string]any `json:"params"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if got := payload.Params["advertiser_id"]; got != "AR17828074650563772417" {
			t.Fatalf("expected advertiser_id, got %#v", got)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":[],"meta":{"capability":"google_ads.list_ad_creatives"}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"google_ads", "list_ad_creatives",
		"--advertiser-id", "AR17828074650563772417",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	if got := stdout.String(); !strings.Contains(got, `"columns":["advertiser_name","creative_id","format","first_shown_datetime","last_shown_datetime","total_days_shown"]`) || !strings.Contains(got, `"rows":[]`) {
		t.Fatalf("expected compact table output, got %s", got)
	}
}

func TestGoogleAdsListAdCreativesRequiresAdvertiserIDOrDomain(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		t.Fatalf("capability should not be called, got %s", r.URL.Path)
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root := NewRoot(&stdout, &stderr)

	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"google_ads", "list_ad_creatives",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 2 {
		t.Fatalf("expected exit code 2, got %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "requires --advertiser-id or --domain") {
		t.Fatalf("expected any-of-required error, got %s", stdout.String())
	}
}

func TestGoogleAdsGetAdDetailsCallsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/google_ads.get_ad_details" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":{"ad_information":{"format":"text","last_shown_date":"2026-04-20","regions":[]},"variations":[{"title":"Nike","description":"Shop now"}]},"meta":{"capability":"google_ads.get_ad_details"}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"google_ads", "get_ad_details",
		"--advertiser-id", "AR17828074650563772417",
		"--creative-id", "CR08581524376019533825",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	assertDataOnlyOutputContains(t, stdout.String(), `"format":"text"`)
}

func TestRedditSearchCallsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/reddit.search" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":[],"meta":{"capability":"reddit.search","cursor":"t3_next"}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"reddit", "search",
		"--query", "tesla",
		"--search-type", "post",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	if got := stdout.String(); !strings.Contains(got, `"columns":["post_id","title","subreddit","score","num_comments","created_at"]`) || !strings.Contains(got, `"rows":[]`) {
		t.Fatalf("expected compact table output, got %s", got)
	}
}

func TestRedditGetPostDetailCallsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/reddit.get_post_detail" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":{"post_id":"t3_abc123"},"meta":{"capability":"reddit.get_post_detail"}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"reddit", "get_post_detail",
		"--post-id", "t3_abc123",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	assertDataOnlyOutputContains(t, stdout.String(), `"post_id":"t3_abc123"`)
}

func TestXiaohongshuSearchNotesCallsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/xiaohongshu.search_notes" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":[],"meta":{"capability":"xiaohongshu.search_notes","current_page":1,"has_more":true}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"xiaohongshu", "search_notes",
		"--keyword", "护肤",
		"--note-type", "video",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	if got := stdout.String(); !strings.Contains(got, `"columns":["note_id","title","desc","liked_count","collected_count","author"]`) || !strings.Contains(got, `"rows":[]`) {
		t.Fatalf("expected compact table output, got %s", got)
	}
}

func TestXiaohongshuGetNoteDetailCallsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/xiaohongshu.get_note_detail" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var req struct {
			Params map[string]any `json:"params"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		if req.Params["xsec_token"] != "token-from-search" {
			t.Fatalf("expected xsec_token param, got %#v", req.Params)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":{"note_id":"67c1f4f1000000001a01b6d3"},"meta":{"capability":"xiaohongshu.get_note_detail"}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"xiaohongshu", "get_note_detail",
		"--note-id", "67c1f4f1000000001a01b6d3",
		"--xsec-token", "token-from-search",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	assertDataOnlyOutputContains(t, stdout.String(), `"note_id":"67c1f4f1000000001a01b6d3"`)
}
