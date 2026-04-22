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
		"product_url":"https://example.com/product",
		"brand_store":{"text":"Acme"},
		"price":{"display":"$19.99"},
		"buybox":{"availability":"In Stock","price":{"display":"$19.99"},"original_price":{"display":"$24.99"}},
		"rating":4.5,
		"review_count":123,
		"main_image":"https://example.com/main.jpg",
		"images":[{"variant":"MAIN","link":"a"},{"variant":"PT01","link":"b"}],
		"variants":[{"asin":"B001-A","title":"Black","dimensions":[{"name":"Color","value":"Black"}]},{"asin":"B001-B","title":"White"}],
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
	if got["brand"] != "Acme" || got["reviews"] != float64(123) || got["link"] != "https://example.com/product" {
		t.Fatalf("unexpected compact projection: %#v", got)
	}
	bullets, _ := got["feature_bullets"].([]any)
	if len(bullets) != 4 {
		t.Fatalf("expected full bullets, got %#v", got["feature_bullets"])
	}
	images, _ := got["images"].([]any)
	variants, _ := got["variants"].([]any)
	if len(images) != 2 || len(variants) != 2 {
		t.Fatalf("expected nested arrays to stay uncolumnar, got %#v", got)
	}
}

func TestProjectCapabilityCompactSearchProducts(t *testing.T) {
	payload := json.RawMessage(`[
		{"position":1,"asin":"A1","title":"Desk Lamp","product_url":"https://example.com/a1","main_image":"https://example.com/a1.jpg","price":{"display":"$19.99"},"rating":4.5,"review_count":10},
		{"position":2,"asin":"A2","title":"Floor Lamp","product_url":"https://example.com/a2","main_image":"https://example.com/a2.jpg","price":{"display":"$29.99"},"rating":4.7,"review_count":20}
	]`)

	body, err := projectCapability("amazon.search_products", payload, FormatCompact)
	if err != nil {
		t.Fatalf("projectCapability() error = %v", err)
	}

	var got struct {
		Items struct {
			Columns []string `json:"columns"`
			Rows    [][]any  `json:"rows"`
		} `json:"items"`
	}
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("unmarshal compact projection: %v", err)
	}
	if len(got.Items.Columns) != 7 || len(got.Items.Rows) != 2 {
		t.Fatalf("unexpected table projection: %#v", got)
	}
	if got.Items.Columns[0] != "asin" || got.Items.Columns[2] != "link" || got.Items.Columns[6] != "price" {
		t.Fatalf("unexpected columns: %#v", got.Items.Columns)
	}
}

func TestProjectCapabilityCompactAmazonSpecialShapes(t *testing.T) {
	tests := []struct {
		name       string
		capability string
		payload    json.RawMessage
		want       string
	}{
		{
			name:       "asin sales history wraps items",
			capability: "amazon.get_asins_sales_history",
			payload:    json.RawMessage(`[{"asin":"A1","month":"2026-01","sales":12}]`),
			want:       `"items":{"columns":["asin","month","sales"]`,
		},
		{
			name:       "category best sellers uses products table",
			capability: "amazon.get_category_best_sellers",
			payload:    json.RawMessage(`[{"asin":"A1","title":"Desk Lamp","brand":"Acme"}]`),
			want:       `"products":{"columns":["listing_sales_volume_of_daily"`,
		},
		{
			name:       "product reviews uses data table",
			capability: "amazon.get_product_reviews",
			payload:    json.RawMessage(`[{"reviewer_name":"Alice","star":4,"title":"Great","date":"2026-04-01","is_verified_purchase":true,"helpful_votes":2,"content":"Solid","reviewed_country":"US"}]`),
			want:       `"data":{"columns":["consumer_name","title","star","reviews_date","is_vp","helpful","content","reviewed_country"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := projectCapability(tt.capability, tt.payload, FormatCompact)
			if err != nil {
				t.Fatalf("projectCapability() error = %v", err)
			}
			if !strings.Contains(string(body), tt.want) {
				t.Fatalf("unexpected compact projection: %s", string(body))
			}
		})
	}
}

func TestProjectCapabilityWithMetaAmazonGetCategoryTrendFiltersRequestedMetrics(t *testing.T) {
	payload := json.RawMessage(`[
		{"month":"2026-01","sales_volume":1200,"brand_count":15,"seller_count":8},
		{"month":"2026-02","sales_volume":1350,"brand_count":18,"seller_count":9}
	]`)

	body, ok, err := projectCapabilityWithMeta("amazon.get_category_trend", payload, map[string]any{
		"metrics": []any{"sales_volume", "brand_count"},
	}, FormatCompact)
	if err != nil {
		t.Fatalf("projectCapabilityWithMeta() error = %v", err)
	}
	if !ok {
		t.Fatal("expected custom meta projection")
	}

	var got struct {
		Items struct {
			Columns []string `json:"columns"`
			Rows    [][]any  `json:"rows"`
		} `json:"items"`
	}
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("unmarshal compact projection: %v", err)
	}
	if len(got.Items.Columns) != 3 {
		t.Fatalf("expected month + requested metrics only, got %#v", got.Items.Columns)
	}
	if got.Items.Columns[0] != "month" || got.Items.Columns[1] != "sales_volume" || got.Items.Columns[2] != "brand_count" {
		t.Fatalf("unexpected columns: %#v", got.Items.Columns)
	}
	if len(got.Items.Rows) != 2 || len(got.Items.Rows[0]) != 3 {
		t.Fatalf("unexpected rows: %#v", got.Items.Rows)
	}
}

func TestProjectCapabilityCompactGoogleTrends(t *testing.T) {
	payload := json.RawMessage(`{
		"search_parameters":{"q":"AI","geo":"US","time":"today 12-m","data_type":"TIMESERIES","gprop":"youtube"},
		"averages":[{"query":"AI","value":55}],
		"regions":[{"geo":"US-CA","name":"California","values":[{"query":"AI","value":60}]}],
		"timeline_data":[
			{"date":"Jan","timestamp":"1","values":[{"query":"AI","value":50}]},
			{"date":"Feb","timestamp":"2","values":[{"query":"AI","value":60}]}
		]
	}`)

	body, err := projectCapability("google_trends.get_interest_over_time", payload, FormatCompact)
	if err != nil {
		t.Fatalf("projectCapability() error = %v", err)
	}

	text := string(body)
	if !strings.Contains(text, `"timeline":{"columns":["date","timestamp","values"]`) {
		t.Fatalf("expected compact timeline projection, got %s", text)
	}
	if !strings.Contains(text, `"query":"AI"`) {
		t.Fatalf("expected scalar context preserved, got %s", text)
	}
	if !strings.Contains(text, `"data_type":"TIMESERIES"`) || !strings.Contains(text, `"gprop":"youtube"`) {
		t.Fatalf("expected compatibility scalar fields preserved, got %s", text)
	}
	if !strings.Contains(text, `"averages":[{"query":"AI","value":55}]`) {
		t.Fatalf("expected averages preserved, got %s", text)
	}
	if !strings.Contains(text, `"regions":[{"geo":"US-CA","name":"California","values":[{"query":"AI","value":60}]}]`) {
		t.Fatalf("expected regions preserved, got %s", text)
	}
}

func TestProjectCapabilityCompactTikTokSearchVideosKeepsKamayColumns(t *testing.T) {
	payload := json.RawMessage(`[
		{
			"video_id":"v1",
			"video_url":"https://example.com/v1.mp4",
			"description":"coffee maker review",
			"create_time":"2026-04-22T00:00:00Z",
			"like_count":12,
			"comment_count":3,
			"share_count":2,
			"play_count":100,
			"cover_image":"https://example.com/c.jpg",
			"duration":30,
			"region":"US",
			"is_ad":false,
			"author":{"user_id":"u1","unique_id":"user1","nickname":"User 1"}
		}
	]`)

	body, err := projectCapability("tiktok.search_videos", payload, FormatCompact)
	if err != nil {
		t.Fatalf("projectCapability() error = %v", err)
	}
	text := string(body)
	if !strings.Contains(text, `"columns":["video_id","video_url","description","create_time","like_count","comment_count","share_count","play_count","cover_image","duration","region","is_ad","author"]`) {
		t.Fatalf("expected Kamay-compatible TikTok video columns, got %s", text)
	}
}

func TestProjectCapabilityCompactTikTokCommentsKeepsAuthorObject(t *testing.T) {
	payload := json.RawMessage(`[
		{
			"comment_id":"c1",
			"text":"nice",
			"create_time":"2026-04-22T00:00:00Z",
			"like_count":2,
			"reply_count":1,
			"author":{"unique_id":"user1","nickname":"User 1","avatar_url":"https://example.com/a.jpg"}
		}
	]`)

	body, err := projectCapability("tiktok.list_comments", payload, FormatCompact)
	if err != nil {
		t.Fatalf("projectCapability() error = %v", err)
	}
	text := string(body)
	if !strings.Contains(text, `"columns":["comment_id","text","create_time","like_count","reply_count","author"]`) {
		t.Fatalf("expected Kamay-compatible TikTok comment columns, got %s", text)
	}
	if !strings.Contains(text, `"nickname":"User 1"`) {
		t.Fatalf("expected nested author object preserved, got %s", text)
	}
}

func TestProjectCapabilityCompactMetaAdsKeepsSnapshotAndPlatforms(t *testing.T) {
	payload := json.RawMessage(`[
		{
			"ad_id":"a1",
			"page_name":"Coffee Brand",
			"start_date":"2026-01-01T00:00:00Z",
			"end_date":"2026-01-31T00:00:00Z",
			"is_active":true,
			"publisher_platforms":["facebook","instagram"],
			"snapshot":{"body":"hello","title":"world"},
			"collation_count":2
		}
	]`)

	body, err := projectCapability("meta_ads.search_ads", payload, FormatCompact)
	if err != nil {
		t.Fatalf("projectCapability() error = %v", err)
	}
	text := string(body)
	if !strings.Contains(text, `"columns":["ad_id","page_name","start_date","end_date","is_active","publisher_platform","snapshot","collation_count"]`) {
		t.Fatalf("expected Kamay-compatible meta ads columns, got %s", text)
	}
	if !strings.Contains(text, `"facebook"`) || !strings.Contains(text, `"body":"hello"`) {
		t.Fatalf("expected publisher platforms and snapshot preserved, got %s", text)
	}
}

func TestProjectCapabilityCompactGoogleAdsCreativesKeepsKamayColumns(t *testing.T) {
	payload := json.RawMessage(`[
		{
			"position":1,
			"creative_id":"CR1",
			"target_domain":"example.com",
			"advertiser_id":"AR1",
			"advertiser_name":"Example Inc.",
			"first_shown_datetime":"2026-01-01T00:00:00Z",
			"last_shown_datetime":"2026-01-31T00:00:00Z",
			"total_days_shown":30,
			"format":"video",
			"details_link":"https://example.com/detail"
		}
	]`)

	body, err := projectCapability("google_ads.list_ad_creatives", payload, FormatCompact)
	if err != nil {
		t.Fatalf("projectCapability() error = %v", err)
	}
	text := string(body)
	if !strings.Contains(text, `"columns":["position","id","target_domain","advertiser.id","advertiser_name","advertiser.name","format","first_shown_datetime","last_shown_datetime","total_days_shown","details_link"]`) {
		t.Fatalf("expected expanded Google Ads creative columns, got %s", text)
	}
}

func TestProjectCapabilityCompactGoogleAdsDetailsKeepsVariationFields(t *testing.T) {
	payload := json.RawMessage(`{
		"ad_information":{"format":"video","first_shown_date":"2026-01-01","last_shown_date":"2026-01-31","last_shown_datetime":"2026-01-31T00:00:00Z","regions":[{"code":"US"}]},
		"variations":[
			{
				"title":"Title",
				"link":"https://example.com",
				"description":"Desc",
				"displayed_link":"example.com",
				"long_headline":"Long",
				"call_to_action":"Shop now",
				"thumbnail":"https://example.com/t.jpg",
				"image":"https://example.com/i.jpg",
				"video_link":"https://example.com/v.mp4",
				"video_id":"vid1",
				"duration":"0:30",
				"channel":"Channel",
				"is_skippable":true
			}
		]
	}`)

	body, err := projectCapability("google_ads.get_ad_details", payload, FormatCompact)
	if err != nil {
		t.Fatalf("projectCapability() error = %v", err)
	}
	text := string(body)
	if !strings.Contains(text, `"columns":["title","link","description","displayed_link","long_headline","call_to_action","thumbnail","image","video_link","video_id","duration","channel","is_skippable"]`) {
		t.Fatalf("expected expanded Google Ads detail columns, got %s", text)
	}
}

func TestProjectCapabilityCompactRedditGetPostCommentsCompatShape(t *testing.T) {
	payload := json.RawMessage(`[
		{"comment_id":"c1","author":"alice","text":"first","score":10,"created_time":"2026-04-22T03:00:00Z","parent_id":"t3_post","depth":0}
	]`)

	body, err := projectCapability("reddit.get_post_comments", payload, FormatCompact)
	if err != nil {
		t.Fatalf("projectCapability() error = %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("unmarshal compact projection: %v", err)
	}
	items, ok := got["items"].(map[string]any)
	if !ok {
		t.Fatalf("expected compat items table, got %#v", got)
	}
	columns, _ := items["columns"].([]any)
	rows, _ := items["rows"].([]any)
	if len(columns) != 7 || len(rows) != 1 {
		t.Fatalf("unexpected compat table shape: %#v", items)
	}
	if columns[0] != "id" || columns[2] != "body" || columns[4] != "created_at" {
		t.Fatalf("unexpected compat columns: %#v", columns)
	}
}

func TestProjectCapabilityCompactRedditGetPostDetailUsesCreatedAt(t *testing.T) {
	payload := json.RawMessage(`{
		"post_id":"t3_abc123",
		"title":"Title",
		"subreddit":"golang",
		"score":10,
		"num_comments":5,
		"created_time":"2026-04-22T03:00:00Z",
		"selftext":"Body"
	}`)

	body, err := projectCapability("reddit.get_post_detail", payload, FormatCompact)
	if err != nil {
		t.Fatalf("projectCapability() error = %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("unmarshal compact projection: %v", err)
	}
	if got["created_at"] != "2026-04-22T03:00:00Z" {
		t.Fatalf("expected created_at in compact output, got %#v", got)
	}
	if _, exists := got["created_time"]; exists {
		t.Fatalf("did not expect created_time in compact output, got %#v", got)
	}
}

func TestProjectCapabilityCompactTiktokShopProductsCanonicalColumns(t *testing.T) {
	payload := json.RawMessage(`[
		{"product_id":"p1","product_name":"Desk Lamp","product_cover":"https://cdn.example/p1.jpg","product_sold_count":12,"format_available_price":"$9.99","format_origin_price":"$12.99","discount":"20% off"}
	]`)

	body, err := projectCapability("tiktok.shop_products", payload, FormatCompact)
	if err != nil {
		t.Fatalf("projectCapability() error = %v", err)
	}

	var got struct {
		Columns []string `json:"columns"`
		Rows    [][]any  `json:"rows"`
	}
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("unmarshal compact projection: %v", err)
	}
	if len(got.Columns) != 7 || len(got.Rows) != 1 {
		t.Fatalf("unexpected shop_products projection: %#v", got)
	}
	if got.Columns[1] != "product_name" || got.Columns[4] != "format_available_price" {
		t.Fatalf("unexpected shop_products columns: %#v", got.Columns)
	}
}

func TestProjectCapabilityCompactTiktokShopProductInfoIncludesReviewCount(t *testing.T) {
	payload := json.RawMessage(`{
		"product_id":"p1",
		"product_name":"Desk Lamp",
		"status":1,
		"seller_id":"seller-1",
		"seller_name":"Seller",
		"sold_count":12,
		"rating":4.9,
		"original_price":"$12.99",
		"real_price":"$9.99",
		"discount":"20% off",
		"images":["https://cdn.example/p1.jpg"],
		"is_platform_product":true,
		"review_count":87
	}`)

	body, err := projectCapability("tiktok.shop_product_info", payload, FormatCompact)
	if err != nil {
		t.Fatalf("projectCapability() error = %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("unmarshal compact projection: %v", err)
	}
	if got["seller_id"] != "seller-1" || got["review_count"] != float64(87) {
		t.Fatalf("expected seller_id and review_count in compact output, got %#v", got)
	}
	if got["is_platform_product"] != true {
		t.Fatalf("expected is_platform_product in compact output, got %#v", got)
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
		"amazon.list_asin_keywords",
		"amazon.query_aba_keywords",
		"amazon.get_asin_sales_daily_trend",
		"amazon.get_variant_sales_30d",
		"amazon.get_category_best_sellers",
		"amazon.get_category_trend",
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
