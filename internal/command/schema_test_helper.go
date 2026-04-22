package command

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/reorc/apimux-cli/internal/schema"
)

func maybeServeSchema(w http.ResponseWriter, r *http.Request) bool {
	const prefix = "/v1/schema/"
	if !strings.HasPrefix(r.URL.Path, prefix) {
		return false
	}
	capability := strings.TrimPrefix(r.URL.Path, prefix)
	spec, ok := testSchemas()[capability]
	if !ok {
		http.NotFound(w, r)
		return true
	}
	body, _ := json.Marshal(map[string]any{
		"ok":   true,
		"data": spec,
	})
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(body)
	return true
}

func testSchemas() map[string]schema.CapabilitySchema {
	return map[string]schema.CapabilitySchema{
		"google_trends.get_interest_over_time": {
			Name: "google_trends.get_interest_over_time",
			Parameters: []schema.CapabilityParam{
				{Name: "q", Type: "string", Required: true},
				{Name: "time", Type: "string", Default: "today 12-m"},
				{Name: "geo", Type: "string"},
				{Name: "cat", Type: "string"},
				{Name: "gprop", Type: "string", Enum: []string{"images", "news", "froogle", "youtube"}},
				{Name: "tz", Type: "integer"},
			},
		},
		"trendcloud.search_filter_values": {
			Name: "trendcloud.search_filter_values",
			Parameters: []schema.CapabilityParam{
				{Name: "kind", Type: "string", Required: true, Enum: []string{"category", "brand", "series", "sku", "attribute"}},
				{Name: "query", Type: "string", Required: true},
				{Name: "platforms", Type: "array", ItemsType: "string", Encoding: "csv", Enum: []string{"douyin", "jd", "tmall"}},
				{Name: "categories", Type: "array", ItemsType: "string", Encoding: "csv"},
				{Name: "limit", Type: "integer", Default: 10},
			},
		},
		"trendcloud.get_market_trend": {
			Name: "trendcloud.get_market_trend",
			Parameters: []schema.CapabilityParam{
				{Name: "start_month", Type: "string"},
				{Name: "end_month", Type: "string"},
				{Name: "metrics", Type: "array", ItemsType: "string", Encoding: "csv", Enum: []string{"sales", "volume"}},
				{Name: "filters", Type: "object", Encoding: "json", FlagName: "filters-json"},
			},
		},
		"trendcloud.get_top_rankings": {
			Name: "trendcloud.get_top_rankings",
			Parameters: []schema.CapabilityParam{
				{Name: "entity", Type: "string", Required: true, Enum: []string{"brand", "category", "series", "sku", "attribute"}},
				{Name: "metric", Type: "string"},
				{Name: "start_month", Type: "string"},
				{Name: "end_month", Type: "string"},
				{Name: "top_n", Type: "integer", Default: 20, FlagName: "top-n"},
				{Name: "category_level", Type: "string", Enum: []string{"category1", "category2", "category3"}, FlagName: "category-level"},
				{Name: "filters", Type: "object", Encoding: "json", FlagName: "filters-json"},
			},
		},
		"amazon.get_product":                amazonSchema("amazon.get_product", []schema.CapabilityParam{{Name: "asin", Type: "string", Required: true}, {Name: "market", Type: "string"}}),
		"amazon.expand_keywords":            amazonSchema("amazon.expand_keywords", []schema.CapabilityParam{{Name: "keyword", Type: "string", Required: true}, {Name: "market", Type: "string"}}),
		"amazon.get_keyword_overview":       amazonSchema("amazon.get_keyword_overview", []schema.CapabilityParam{{Name: "keyword", Type: "string", Required: true}, {Name: "market", Type: "string"}}),
		"amazon.get_keyword_trends":         amazonSchema("amazon.get_keyword_trends", []schema.CapabilityParam{{Name: "keywords", Type: "array", Required: true, ItemsType: "string", Encoding: "csv"}, {Name: "market", Type: "string"}, {Name: "granularity", Type: "string", Enum: []string{"week", "month"}}}),
		"amazon.list_asin_keywords":         amazonSchema("amazon.list_asin_keywords", []schema.CapabilityParam{{Name: "asin", Type: "string", Required: true}, {Name: "keyword", Type: "string"}, {Name: "market", Type: "string"}}),
		"amazon.query_aba_keywords":         amazonSchema("amazon.query_aba_keywords", []schema.CapabilityParam{{Name: "keyword", Type: "string"}, {Name: "node_ids", Type: "string", FlagName: "node-ids"}, {Name: "page", Type: "integer", Default: 1}, {Name: "page_index", Type: "integer", Default: 1, FlagName: "page-index"}, {Name: "page_size", Type: "integer", Default: 40, FlagName: "page-size"}, {Name: "market", Type: "string"}}),
		"amazon.search_category":            amazonSchema("amazon.search_category", []schema.CapabilityParam{{Name: "name", Type: "string", Required: true}, {Name: "market", Type: "string"}, {Name: "limit", Type: "integer", Default: 20}}),
		"amazon.get_asin_sales_daily_trend": amazonSchema("amazon.get_asin_sales_daily_trend", []schema.CapabilityParam{{Name: "asin", Type: "string", Required: true}, {Name: "begin_date", Type: "string", FlagName: "begin-date"}, {Name: "market", Type: "string"}}),
		"amazon.get_asins_sales_history":    amazonSchema("amazon.get_asins_sales_history", []schema.CapabilityParam{{Name: "asins", Type: "array", Required: true, ItemsType: "string", Encoding: "csv"}, {Name: "market", Type: "string"}}),
		"amazon.get_variant_sales_30d":      amazonSchema("amazon.get_variant_sales_30d", []schema.CapabilityParam{{Name: "asin", Type: "string", Required: true}, {Name: "market", Type: "string"}}),
		"amazon.get_product_reviews":        amazonSchema("amazon.get_product_reviews", []schema.CapabilityParam{{Name: "asin", Type: "string", Required: true}, {Name: "market", Type: "string"}, {Name: "start_date", Type: "string", FlagName: "start-date"}, {Name: "star", Type: "string", Enum: []string{"positive", "negative", "1", "2", "3", "4", "5"}}, {Name: "only_purchase", Type: "boolean", FlagName: "only-purchase"}, {Name: "page_index", Type: "integer", FlagName: "page-index"}}),
		"amazon.get_category_best_sellers":  amazonSchema("amazon.get_category_best_sellers", []schema.CapabilityParam{{Name: "node_id", Type: "string", Required: true, FlagName: "node-id"}, {Name: "market", Type: "string"}, {Name: "query_start", Type: "string", FlagName: "query-start"}, {Name: "query_date", Type: "string", FlagName: "query-date"}, {Name: "query_days", Type: "integer", FlagName: "query-days"}}),
		"amazon.get_category_trend":         amazonSchema("amazon.get_category_trend", []schema.CapabilityParam{{Name: "node_id", Type: "string", Required: true, FlagName: "node-id"}, {Name: "market", Type: "string"}, {Name: "trend_types", Type: "array", Required: true, ItemsType: "string", Encoding: "csv", FlagName: "trend-types"}}),
		"tiktok.search_videos":              {Name: "tiktok.search_videos", Parameters: []schema.CapabilityParam{{Name: "keyword", Type: "string", Required: true}, {Name: "region", Type: "string"}, {Name: "sort_by", Type: "string", Enum: []string{"relevance", "likes", "date"}, FlagName: "sort-by"}, {Name: "publish_time", Type: "string", Enum: []string{"all", "1d", "1w", "1m", "3m", "6m"}, FlagName: "publish-time"}, {Name: "cursor", Type: "integer"}, {Name: "count", Type: "integer"}}},
		"meta_ads.search_ads":               {Name: "meta_ads.search_ads", Parameters: []schema.CapabilityParam{{Name: "q", Type: "string", Required: true}, {Name: "country", Type: "string"}, {Name: "ad_type", Type: "string", Enum: []string{"all", "political_and_issue_ads", "housing_ads", "employment_ads", "credit_ads"}, FlagName: "ad-type"}, {Name: "active_status", Type: "string", Enum: []string{"active", "inactive", "all"}, FlagName: "active-status"}, {Name: "media_type", Type: "string", Enum: []string{"all", "video", "image", "meme", "image_and_meme", "none"}, FlagName: "media-type"}, {Name: "platforms", Type: "string"}, {Name: "start_date", Type: "string", FlagName: "start-date"}, {Name: "end_date", Type: "string", FlagName: "end-date"}, {Name: "next_page_token", Type: "string", FlagName: "next-page-token"}}},
		"meta_ads.get_ad_detail":            {Name: "meta_ads.get_ad_detail", Parameters: []schema.CapabilityParam{{Name: "ad_id", Type: "string", Required: true, FlagName: "ad-id"}}},
		"douyin.search_videos":              {Name: "douyin.search_videos", Parameters: []schema.CapabilityParam{{Name: "keyword", Type: "string", Required: true}, {Name: "sort_type", Type: "string", Enum: []string{"comprehensive", "likes", "latest"}, FlagName: "sort-type"}, {Name: "publish_time", Type: "string", Enum: []string{"all", "1d", "1w", "6m"}, FlagName: "publish-time"}, {Name: "filter_duration", Type: "string", Enum: []string{"all", "under_1m", "1m_5m", "over_5m"}, FlagName: "filter-duration"}, {Name: "content_type", Type: "string", Enum: []string{"all", "video", "image", "article"}, FlagName: "content-type"}, {Name: "cursor", Type: "integer"}}},
		"douyin.get_video_detail":           {Name: "douyin.get_video_detail", Parameters: []schema.CapabilityParam{{Name: "aweme_id", Type: "string", Required: true, FlagName: "aweme-id"}}},
		"douyin.get_video_comments":         {Name: "douyin.get_video_comments", Parameters: []schema.CapabilityParam{{Name: "aweme_id", Type: "string", Required: true, FlagName: "aweme-id"}, {Name: "cursor", Type: "integer"}, {Name: "count", Type: "integer"}}},
		"douyin.get_comment_replies":        {Name: "douyin.get_comment_replies", Parameters: []schema.CapabilityParam{{Name: "aweme_id", Type: "string", Required: true, FlagName: "aweme-id"}, {Name: "comment_id", Type: "string", Required: true, FlagName: "comment-id"}, {Name: "cursor", Type: "integer"}, {Name: "count", Type: "integer", Default: 20}}},
		"reddit.search":                     {Name: "reddit.search", Parameters: []schema.CapabilityParam{{Name: "query", Type: "string", Required: true}, {Name: "search_type", Type: "string", Enum: []string{"post", "community", "comment", "media", "people"}, FlagName: "search-type"}, {Name: "sort", Type: "string", Enum: []string{"relevance", "hot", "top", "new", "comments"}}, {Name: "time_range", Type: "string", Enum: []string{"all", "year", "month", "week", "day", "hour"}, FlagName: "time-range"}, {Name: "after", Type: "string"}}},
		"reddit.get_subreddit_feed":         {Name: "reddit.get_subreddit_feed", Parameters: []schema.CapabilityParam{{Name: "subreddit_name", Type: "string", Required: true, FlagName: "subreddit-name"}, {Name: "sort", Type: "string", Enum: []string{"best", "hot", "new", "top", "controversial", "rising"}}, {Name: "after", Type: "string"}}},
		"reddit.get_post_detail":            {Name: "reddit.get_post_detail", Parameters: []schema.CapabilityParam{{Name: "post_id", Type: "string", Required: true, FlagName: "post-id"}}},
		"reddit.get_post_comments":          {Name: "reddit.get_post_comments", Parameters: []schema.CapabilityParam{{Name: "post_id", Type: "string", Required: true, FlagName: "post-id"}, {Name: "sort_type", Type: "string", Enum: []string{"confidence", "new", "top", "hot", "controversial", "old", "random"}, FlagName: "sort-type"}, {Name: "after", Type: "string"}}},
		"xiaohongshu.search_notes":          {Name: "xiaohongshu.search_notes", Parameters: []schema.CapabilityParam{{Name: "keyword", Type: "string", Required: true}, {Name: "page", Type: "integer"}, {Name: "note_type", Type: "string", Enum: []string{"all", "video", "normal", "live"}, FlagName: "note-type"}, {Name: "time_filter", Type: "string", Enum: []string{"all", "1d", "1w", "6m"}, FlagName: "time-filter"}, {Name: "sort_strategy", Type: "string", Enum: []string{"default", "latest", "likes"}, FlagName: "sort-strategy"}}},
		"xiaohongshu.get_note_detail":       {Name: "xiaohongshu.get_note_detail", Parameters: []schema.CapabilityParam{{Name: "note_id", Type: "string", Required: true, FlagName: "note-id"}, {Name: "xsec_token", Type: "string", FlagName: "xsec-token"}}},
		"xiaohongshu.get_note_comments":     {Name: "xiaohongshu.get_note_comments", Parameters: []schema.CapabilityParam{{Name: "note_id", Type: "string", Required: true, FlagName: "note-id"}, {Name: "cursor", Type: "string"}, {Name: "sort_strategy", Type: "string", Enum: []string{"default", "latest", "likes"}, FlagName: "sort-strategy"}}},
		"google_ads.search_advertisers":     {Name: "google_ads.search_advertisers", Parameters: []schema.CapabilityParam{{Name: "query", Type: "string", Required: true}, {Name: "region", Type: "string"}, {Name: "num_advertisers", Type: "integer", Default: 10, FlagName: "num-advertisers"}, {Name: "num_domains", Type: "integer", Default: 10, FlagName: "num-domains"}}},
		"google_ads.list_ad_creatives": {
			Name:       "google_ads.list_ad_creatives",
			Parameters: []schema.CapabilityParam{{Name: "advertiser_id", Type: "string", FlagName: "advertiser-id"}, {Name: "domain", Type: "string"}, {Name: "region", Type: "string"}, {Name: "platform", Type: "string", Enum: []string{"google_play", "google_maps", "google_search", "youtube", "google_shopping"}}, {Name: "ad_format", Type: "string", Enum: []string{"text", "image", "video"}, FlagName: "ad-format"}, {Name: "time_period", Type: "string", Enum: []string{"last_7_days", "last_30_days", "last_90_days", "last_year"}, FlagName: "time-period"}, {Name: "page_token", Type: "string", FlagName: "page-token"}},
			Rules:      []schema.CapabilityRule{{Kind: "any_of_required", Params: []string{"advertiser_id", "domain"}, Message: "google_ads list_ad_creatives requires --advertiser-id or --domain"}},
		},
		"google_ads.get_ad_details": {Name: "google_ads.get_ad_details", Parameters: []schema.CapabilityParam{{Name: "advertiser_id", Type: "string", Required: true, FlagName: "advertiser-id"}, {Name: "creative_id", Type: "string", Required: true, FlagName: "creative-id"}}},
	}
}

func amazonSchema(name string, params []schema.CapabilityParam) schema.CapabilitySchema {
	return schema.CapabilitySchema{Name: name, Parameters: params}
}
