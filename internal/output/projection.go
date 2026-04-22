package output

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Format string

const (
	FormatAuto    Format = "auto"
	FormatData    Format = "data"
	FormatCompact Format = "compact"
)

type UnsupportedProjectionError struct {
	Capability string
	Format     Format
}

func (e *UnsupportedProjectionError) Error() string {
	return fmt.Sprintf("format %q is not supported for capability %q", e.Format, e.Capability)
}

type InvalidProjectionError struct {
	Capability string
	Format     Format
	Cause      error
}

func (e *InvalidProjectionError) Error() string {
	return fmt.Sprintf("failed to apply %q projection for capability %q: %v", e.Format, e.Capability, e.Cause)
}

func (e *InvalidProjectionError) Unwrap() error {
	return e.Cause
}

type projectionRule struct {
	Compact projectionSpec
}

type projectionSpec struct {
	Scalars []fieldRule
	Lists   []listRule
	Tables  []tableRule
}

type fieldRule struct {
	From      string
	To        string
	Transform transformKind
}

type listRule struct {
	From   string
	To     string
	Fields []fieldRule
	Limit  int
}

type tableRule struct {
	From    string
	To      string
	Columns []fieldRule
	Limit   int
}

type transformKind string

const (
	transformNone           transformKind = ""
	transformCount          transformKind = "count"
	transformJoinDimensions transformKind = "join_dimensions"
)

var projectionRules = map[string]projectionRule{
	"amazon.expand_keywords": {
		Compact: projectionSpec{
			Tables: []tableRule{
				{
					From:  "$root",
					To:    "$root",
					Limit: 20,
					Columns: []fieldRule{
						{From: "keyword", To: "keyword"},
						{From: "est_searches_num", To: "est_searches_num"},
						{From: "searches_rank", To: "searches_rank"},
						{From: "match_types", To: "match_types"},
					},
				},
			},
		},
	},
	"amazon.get_asins_sales_history": {
		Compact: projectionSpec{
			Tables: []tableRule{
				{
					From:  "$root",
					To:    "items",
					Limit: 50,
					Columns: []fieldRule{
						{From: "asin", To: "asin"},
						{From: "month", To: "month"},
						{From: "sales", To: "sales"},
					},
				},
			},
		},
	},
	"amazon.get_keyword_overview": {
		Compact: projectionSpec{
			Scalars: []fieldRule{
				{From: "keyword", To: "keyword"},
				{From: "est_searches_num", To: "est_searches_num"},
				{From: "searches_rank", To: "searches_rank"},
				{From: "brand_ad_asin_num", To: "brand_ad_asin_num"},
				{From: "sp_ad_asin_num", To: "sp_ad_asin_num"},
				{From: "ppc_ad_asin_num", To: "ppc_ad_asin_num"},
				{From: "sale_num", To: "sale_num"},
				{From: "nf_asin_num", To: "nf_asin_num"},
				{From: "ac_asin_num", To: "ac_asin_num"},
				{From: "video_ad_asin_num", To: "video_ad_asin_num"},
				{From: "global_keyword_num", To: "global_keyword_num"},
				{From: "update_time", To: "update_time"},
			},
		},
	},
	"amazon.get_keyword_trends": {
		Compact: projectionSpec{
			Lists: []listRule{
				{
					From:  "$root",
					To:    "items",
					Limit: 20,
					Fields: []fieldRule{
						{From: "keyword", To: "keyword"},
						{From: "est_searches_num_history", To: "est_searches_num_history_map"},
						{From: "searches_rank_history", To: "searches_rank_history_map"},
					},
				},
			},
		},
	},
	"amazon.get_product": {
		Compact: projectionSpec{
			Scalars: []fieldRule{
				{From: "asin", To: "asin"},
				{From: "title", To: "title"},
				{From: "brand_store.text", To: "brand"},
				{From: "price.display", To: "price"},
				{From: "buybox.availability", To: "availability"},
				{From: "rating", To: "rating"},
				{From: "review_count", To: "reviews"},
				{From: "product_url", To: "link"},
				{From: "main_image", To: "main_image"},
				{From: "feature_bullets", To: "feature_bullets"},
				{From: "images", To: "images"},
				{From: "variants", To: "variants"},
				{From: "buybox.price.display", To: "buybox.price"},
				{From: "buybox.original_price.display", To: "buybox.original_price"},
			},
		},
	},
	"amazon.get_product_reviews": {
		Compact: projectionSpec{
			Tables: []tableRule{
				{
					From:  "$root",
					To:    "data",
					Limit: 20,
					Columns: []fieldRule{
						{From: "reviewer_name", To: "consumer_name"},
						{From: "title", To: "title"},
						{From: "star", To: "star"},
						{From: "date", To: "reviews_date"},
						{From: "is_verified_purchase", To: "is_vp"},
						{From: "helpful_votes", To: "helpful"},
						{From: "content", To: "content"},
						{From: "reviewed_country", To: "reviewed_country"},
					},
				},
			},
		},
	},
	"amazon.search_category": {
		Compact: projectionSpec{
			Tables: []tableRule{
				{
					From:  "$root",
					To:    "items",
					Limit: 20,
					Columns: []fieldRule{
						{From: "node_id", To: "node_id"},
						{From: "name", To: "name"},
						{From: "cn_name", To: "cn_name"},
						{From: "path", To: "path"},
					},
				},
			},
		},
	},
	"amazon.search_products": {
		Compact: projectionSpec{
			Tables: []tableRule{
				{
					From:  "$root",
					To:    "items",
					Limit: 10,
					Columns: []fieldRule{
						{From: "asin", To: "asin"},
						{From: "title", To: "title"},
						{From: "product_url", To: "link"},
						{From: "main_image", To: "thumbnail"},
						{From: "rating", To: "rating"},
						{From: "review_count", To: "reviews"},
						{From: "price.display", To: "price"},
					},
				},
			},
		},
	},
	"amazon.list_asin_keywords": {
		Compact: projectionSpec{
			Tables: []tableRule{
				{
					From:  "$root",
					To:    "items",
					Limit: 100,
					Columns: []fieldRule{
						{From: "keyword", To: "keyword"},
						{From: "kw_characters", To: "kw_characters"},
						{From: "conversion_characters", To: "conversion_characters"},
						{From: "exposure_type", To: "exposure_type"},
						{From: "last_rank", To: "last_rank_str"},
						{From: "ad_last_rank", To: "ad_last_rank_str"},
						{From: "est_searches_num", To: "est_searches_num"},
						{From: "searches_rank", To: "searches_rank"},
						{From: "ratio_score", To: "ratio_score"},
					},
				},
			},
		},
	},
	"amazon.query_aba_keywords": {
		Compact: projectionSpec{
			Tables: []tableRule{
				{
					From:  "$root",
					To:    "items",
					Limit: 100,
					Columns: []fieldRule{
						{From: "keyword", To: "keyword"},
						{From: "keyword_cn_name", To: "keyword_cn_name"},
						{From: "rank", To: "rank"},
						{From: "search_volume", To: "search_volume"},
						{From: "word_count", To: "word_count"},
						{From: "product_count", To: "product_count"},
						{From: "rank_change_of_weekly", To: "rank_change_of_weekly"},
						{From: "cpc", To: "cpc"},
						{From: "cpc_range", To: "cpc_range"},
						{From: "search_conversion_rate", To: "search_conversion_rate"},
						{From: "search_conversion_rate_d90", To: "search_conversion_rate_d90"},
						{From: "click_conversion_rate_d90", To: "click_conversion_rate_d90"},
						{From: "click_of_90d", To: "click_of_90d"},
						{From: "sales_volume_of_90d", To: "sales_volume_of_90d"},
						{From: "share_click_rate", To: "share_click_rate"},
						{From: "share_conversion_rate", To: "share_conversion_rate"},
						{From: "search_volume_growth_rate_trend", To: "search_volume_growth_rate_trend"},
						{From: "top3_asin", To: "top3_asin"},
						{From: "top3_brand", To: "top3_brand"},
						{From: "top3_category", To: "top3_category"},
						{From: "season", To: "season"},
						{From: "update", To: "update"},
					},
				},
			},
		},
	},
	"amazon.get_asin_sales_daily_trend": {
		Compact: projectionSpec{
			Tables: []tableRule{
				{
					From:  "$root",
					To:    "items",
					Limit: 180,
					Columns: []fieldRule{
						{From: "date", To: "date"},
						{From: "sales", To: "sales"},
					},
				},
			},
		},
	},
	"amazon.get_variant_sales_30d": {
		Compact: projectionSpec{
			Tables: []tableRule{
				{
					From:  "$root",
					To:    "items",
					Limit: 100,
					Columns: []fieldRule{
						{From: "asin", To: "asin"},
						{From: "bought_in_past_month", To: "bought_in_past_month"},
						{From: "update_time", To: "update_time"},
					},
				},
			},
		},
	},
	"amazon.get_category_best_sellers": {
		Compact: projectionSpec{
			Tables: []tableRule{
				{
					From:  "$root",
					To:    "products",
					Limit: 100,
					Columns: []fieldRule{
						{From: "listing_sales_volume_of_daily", To: "listing_sales_volume_of_daily"},
						{From: "listing_sales_volume_of_month", To: "listing_sales_volume_of_month"},
						{From: "listing_sales_of_daily", To: "listing_sales_of_daily"},
						{From: "listing_sales_of_month", To: "listing_sales_of_month"},
						{From: "asin", To: "asin"},
						{From: "title", To: "title"},
						{From: "brand", To: "brand"},
						{From: "photo", To: "photo"},
						{From: "price", To: "price"},
						{From: "list_price", To: "list_price"},
						{From: "sales_price", To: "sales_price"},
						{From: "coupon", To: "coupon"},
						{From: "seller_count", To: "seller_count"},
						{From: "is_fba", To: "is_fba"},
						{From: "profit", To: "profit"},
						{From: "profit_rate", To: "profit_rate"},
						{From: "online_days", To: "online_days"},
						{From: "rating_count", To: "ratings_count"},
						{From: "rating", To: "ratings"},
						{From: "rank", To: "rank"},
						{From: "category", To: "category"},
					},
				},
			},
		},
	},
	"amazon.get_category_trend": {
		Compact: projectionSpec{
			Tables: []tableRule{
				{
					From:  "$root",
					To:    "items",
					Limit: 120,
					Columns: []fieldRule{
						{From: "month", To: "month"},
						{From: "sales_volume", To: "sales_volume"},
						{From: "brand_count", To: "brand_count"},
						{From: "seller_count", To: "seller_count"},
						{From: "avg_price", To: "avg_price"},
						{From: "avg_rating_count", To: "avg_rating_count"},
						{From: "avg_star", To: "avg_star"},
						{From: "new_product_ratio_1m", To: "new_product_ratio_1m"},
						{From: "new_product_ratio_3m", To: "new_product_ratio_3m"},
						{From: "amazon_self_ratio", To: "amazon_self_ratio"},
						{From: "avg_profit", To: "avg_profit"},
						{From: "top100_share", To: "top100_share"},
						{From: "top3_listing_monopoly", To: "top3_listing_monopoly"},
						{From: "top10_brand_monopoly", To: "top10_brand_monopoly"},
					},
				},
			},
		},
	},
	"douyin.get_comment_replies": {
		Compact: projectionSpec{
			Tables: []tableRule{
				{
					From:  "$root",
					To:    "$root",
					Limit: 20,
					Columns: []fieldRule{
						{From: "comment_id", To: "comment_id"},
						{From: "text", To: "text"},
						{From: "create_time", To: "create_time"},
						{From: "author.nickname", To: "author"},
						{From: "like_count", To: "like_count"},
					},
				},
			},
		},
	},
	"douyin.get_video_comments": {
		Compact: projectionSpec{
			Tables: []tableRule{
				{
					From:  "$root",
					To:    "$root",
					Limit: 20,
					Columns: []fieldRule{
						{From: "comment_id", To: "comment_id"},
						{From: "text", To: "text"},
						{From: "create_time", To: "create_time"},
						{From: "author.nickname", To: "author"},
						{From: "like_count", To: "like_count"},
						{From: "reply_count", To: "reply_count"},
					},
				},
			},
		},
	},
	"douyin.get_video_detail": {
		Compact: projectionSpec{
			Scalars: []fieldRule{
				{From: "aweme_id", To: "aweme_id"},
				{From: "description", To: "description"},
				{From: "create_time", To: "create_time"},
				{From: "share_url", To: "share_url"},
				{From: "author.nickname", To: "author.nickname"},
				{From: "author.region", To: "author.region"},
				{From: "statistics.like_count", To: "statistics.like_count"},
				{From: "statistics.comment_count", To: "statistics.comment_count"},
				{From: "statistics.share_count", To: "statistics.share_count"},
				{From: "video.duration", To: "video.duration"},
				{From: "video.ratio", To: "video.ratio"},
			},
		},
	},
	"douyin.search_videos": {
		Compact: projectionSpec{
			Tables: []tableRule{
				{
					From:  "$root",
					To:    "$root",
					Limit: 10,
					Columns: []fieldRule{
						{From: "aweme_id", To: "aweme_id"},
						{From: "description", To: "description"},
						{From: "create_time", To: "create_time"},
						{From: "author.nickname", To: "author"},
						{From: "statistics.like_count", To: "like_count"},
						{From: "statistics.comment_count", To: "comment_count"},
						{From: "statistics.share_count", To: "share_count"},
					},
				},
			},
		},
	},
	"google_ads.get_ad_details": {
		Compact: projectionSpec{
			Scalars: []fieldRule{
				{From: "ad_information.format", To: "format"},
				{From: "ad_information.first_shown_date", To: "first_shown_date"},
				{From: "ad_information.last_shown_date", To: "last_shown_date"},
				{From: "ad_information.last_shown_datetime", To: "last_shown_datetime"},
				{From: "ad_information.regions", To: "region_count", Transform: transformCount},
			},
			Tables: []tableRule{
				{
					From:  "variations",
					To:    "variations",
					Limit: 5,
					Columns: []fieldRule{
						{From: "title", To: "title"},
						{From: "link", To: "link"},
						{From: "description", To: "description"},
						{From: "displayed_link", To: "displayed_link"},
						{From: "long_headline", To: "long_headline"},
						{From: "call_to_action", To: "call_to_action"},
						{From: "thumbnail", To: "thumbnail"},
						{From: "image", To: "image"},
						{From: "video_link", To: "video_link"},
						{From: "video_id", To: "video_id"},
						{From: "duration", To: "duration"},
						{From: "channel", To: "channel"},
						{From: "is_skippable", To: "is_skippable"},
					},
				},
			},
		},
	},
	"google_ads.list_ad_creatives": {
		Compact: projectionSpec{
			Tables: []tableRule{
				{
					From:  "$root",
					To:    "$root",
					Limit: 10,
					Columns: []fieldRule{
						{From: "position", To: "position"},
						{From: "creative_id", To: "id"},
						{From: "target_domain", To: "target_domain"},
						{From: "advertiser_id", To: "advertiser.id"},
						{From: "advertiser_name", To: "advertiser_name"},
						{From: "advertiser_name", To: "advertiser.name"},
						{From: "format", To: "format"},
						{From: "first_shown_datetime", To: "first_shown_datetime"},
						{From: "last_shown_datetime", To: "last_shown_datetime"},
						{From: "total_days_shown", To: "total_days_shown"},
						{From: "details_link", To: "details_link"},
					},
				},
			},
		},
	},
	"google_ads.search_advertisers": {
		Compact: projectionSpec{
			Scalars: []fieldRule{
				{From: "domains", To: "domain_count", Transform: transformCount},
			},
			Lists: []listRule{
				{
					From:  "domains",
					To:    "domains",
					Limit: 5,
					Fields: []fieldRule{
						{From: "domain", To: "domain"},
					},
				},
			},
			Tables: []tableRule{
				{
					From:  "advertisers",
					To:    "advertisers",
					Limit: 10,
					Columns: []fieldRule{
						{From: "advertiser_id", To: "advertiser_id"},
						{From: "advertiser_name", To: "advertiser_name"},
						{From: "ads_count", To: "ads_count"},
						{From: "region", To: "region"},
						{From: "is_verified", To: "is_verified"},
					},
				},
			},
		},
	},
	"google_trends.get_interest_over_time": {
		Compact: projectionSpec{
			Scalars: []fieldRule{
				{From: "search_parameters.q", To: "query"},
				{From: "search_parameters.geo", To: "geo"},
				{From: "search_parameters.time", To: "time"},
				{From: "search_parameters.data_type", To: "data_type"},
				{From: "search_parameters.gprop", To: "gprop"},
			},
			Lists: []listRule{
				{
					From:  "averages",
					To:    "averages",
					Limit: 10,
					Fields: []fieldRule{
						{From: "query", To: "query"},
						{From: "value", To: "value"},
					},
				},
				{
					From:  "regions",
					To:    "regions",
					Limit: 20,
					Fields: []fieldRule{
						{From: "geo", To: "geo"},
						{From: "name", To: "name"},
						{From: "values", To: "values"},
					},
				},
			},
			Tables: []tableRule{
				{
					From: "timeline_data",
					To:   "timeline",
					Columns: []fieldRule{
						{From: "date", To: "date"},
						{From: "timestamp", To: "timestamp"},
						{From: "values", To: "values"},
					},
				},
			},
		},
	},
	"meta_ads.get_ad_detail": {
		Compact: projectionSpec{
			Scalars: []fieldRule{
				{From: "ad_id", To: "ad_id"},
				{From: "eu_transparency.age_audience.min", To: "age_audience.min"},
				{From: "eu_transparency.age_audience.max", To: "age_audience.max"},
				{From: "eu_transparency.gender_audience", To: "gender_audience"},
				{From: "eu_transparency.location_audience", To: "location_count", Transform: transformCount},
				{From: "verified_voice", To: "verified_voice"},
			},
		},
	},
	"meta_ads.search_ads": {
		Compact: projectionSpec{
			Tables: []tableRule{
				{
					From:  "$root",
					To:    "$root",
					Limit: 10,
					Columns: []fieldRule{
						{From: "ad_id", To: "ad_id"},
						{From: "page_name", To: "page_name"},
						{From: "start_date", To: "start_date"},
						{From: "end_date", To: "end_date"},
						{From: "is_active", To: "is_active"},
						{From: "publisher_platforms", To: "publisher_platform"},
						{From: "snapshot", To: "snapshot"},
						{From: "collation_count", To: "collation_count"},
					},
				},
			},
		},
	},
	"reddit.get_post_comments": {
		Compact: projectionSpec{
			Tables: []tableRule{
				{
					From:  "$root",
					To:    "items",
					Limit: 20,
					Columns: []fieldRule{
						{From: "comment_id", To: "id"},
						{From: "author", To: "author"},
						{From: "text", To: "body"},
						{From: "score", To: "score"},
						{From: "created_time", To: "created_at"},
						{From: "parent_id", To: "parent_id"},
						{From: "depth", To: "depth"},
					},
				},
			},
		},
	},
	"reddit.get_post_detail": {
		Compact: projectionSpec{
			Scalars: []fieldRule{
				{From: "post_id", To: "post_id"},
				{From: "title", To: "title"},
				{From: "subreddit", To: "subreddit"},
				{From: "score", To: "score"},
				{From: "num_comments", To: "num_comments"},
				{From: "created_time", To: "created_at"},
				{From: "selftext", To: "selftext"},
			},
		},
	},
	"reddit.get_subreddit_feed": {
		Compact: projectionSpec{
			Tables: []tableRule{
				{
					From:  "$root",
					To:    "$root",
					Limit: 20,
					Columns: []fieldRule{
						{From: "post_id", To: "post_id"},
						{From: "title", To: "title"},
						{From: "subreddit", To: "subreddit"},
						{From: "score", To: "score"},
						{From: "num_comments", To: "num_comments"},
						{From: "created_time", To: "created_at"},
					},
				},
			},
		},
	},
	"reddit.search": {
		Compact: projectionSpec{
			Tables: []tableRule{
				{
					From:  "$root",
					To:    "$root",
					Limit: 20,
					Columns: []fieldRule{
						{From: "post_id", To: "post_id"},
						{From: "title", To: "title"},
						{From: "subreddit", To: "subreddit"},
						{From: "score", To: "score"},
						{From: "num_comments", To: "num_comments"},
						{From: "created_time", To: "created_at"},
					},
				},
			},
		},
	},
	"tiktok.list_comments": {
		Compact: projectionSpec{
			Tables: []tableRule{
				{
					From:  "$root",
					To:    "$root",
					Limit: 20,
					Columns: []fieldRule{
						{From: "comment_id", To: "comment_id"},
						{From: "text", To: "text"},
						{From: "create_time", To: "create_time"},
						{From: "like_count", To: "like_count"},
						{From: "reply_count", To: "reply_count"},
						{From: "author", To: "author"},
					},
				},
			},
		},
	},
	"tiktok.search_videos": {
		Compact: projectionSpec{
			Tables: []tableRule{
				{
					From:  "$root",
					To:    "$root",
					Limit: 10,
					Columns: []fieldRule{
						{From: "video_id", To: "video_id"},
						{From: "video_url", To: "video_url"},
						{From: "description", To: "description"},
						{From: "create_time", To: "create_time"},
						{From: "like_count", To: "like_count"},
						{From: "comment_count", To: "comment_count"},
						{From: "share_count", To: "share_count"},
						{From: "play_count", To: "play_count"},
						{From: "cover_image", To: "cover_image"},
						{From: "duration", To: "duration"},
						{From: "region", To: "region"},
						{From: "is_ad", To: "is_ad"},
						{From: "author", To: "author"},
					},
				},
			},
		},
	},
	"tiktok.shop_product_info": {
		Compact: projectionSpec{
			Scalars: []fieldRule{
				{From: "product_id", To: "product_id"},
				{From: "product_name", To: "product_name"},
				{From: "status", To: "status"},
				{From: "seller_id", To: "seller_id"},
				{From: "seller_name", To: "seller_name"},
				{From: "sold_count", To: "sold_count"},
				{From: "rating", To: "rating"},
				{From: "original_price", To: "original_price"},
				{From: "real_price", To: "real_price"},
				{From: "discount", To: "discount"},
				{From: "images", To: "images"},
				{From: "is_platform_product", To: "is_platform_product"},
				{From: "review_count", To: "review_count"},
			},
		},
	},
	"tiktok.shop_products": {
		Compact: projectionSpec{
			Tables: []tableRule{
				{
					From:  "$root",
					To:    "$root",
					Limit: 10,
					Columns: []fieldRule{
						{From: "product_id", To: "product_id"},
						{From: "product_name", To: "product_name"},
						{From: "product_cover", To: "product_cover"},
						{From: "product_sold_count", To: "product_sold_count"},
						{From: "format_available_price", To: "format_available_price"},
						{From: "format_origin_price", To: "format_origin_price"},
						{From: "discount", To: "discount"},
					},
				},
			},
		},
	},
	"xiaohongshu.get_note_comments": {
		Compact: projectionSpec{
			Tables: []tableRule{
				{
					From:  "$root",
					To:    "$root",
					Limit: 20,
					Columns: []fieldRule{
						{From: "comment_id", To: "comment_id"},
						{From: "content", To: "content"},
						{From: "like_count", To: "like_count"},
						{From: "create_time", To: "create_time"},
						{From: "author.nickname", To: "author"},
					},
				},
			},
		},
	},
	"xiaohongshu.get_note_detail": {
		Compact: projectionSpec{
			Scalars: []fieldRule{
				{From: "note_id", To: "note_id"},
				{From: "title", To: "title"},
				{From: "desc", To: "desc"},
				{From: "liked_count", To: "liked_count"},
				{From: "collected_count", To: "collected_count"},
				{From: "comment_count", To: "comment_count"},
				{From: "author.nickname", To: "author.nickname"},
				{From: "tags", To: "tags"},
			},
		},
	},
	"xiaohongshu.search_notes": {
		Compact: projectionSpec{
			Tables: []tableRule{
				{
					From:  "$root",
					To:    "$root",
					Limit: 10,
					Columns: []fieldRule{
						{From: "note_id", To: "note_id"},
						{From: "title", To: "title"},
						{From: "desc", To: "desc"},
						{From: "liked_count", To: "liked_count"},
						{From: "collected_count", To: "collected_count"},
						{From: "author.nickname", To: "author"},
					},
				},
			},
		},
	},
}

func ParseFormat(value string) (Format, bool) {
	switch strings.TrimSpace(value) {
	case "", string(FormatAuto):
		return FormatAuto, true
	case string(FormatData):
		return FormatData, true
	case string(FormatCompact):
		return FormatCompact, true
	default:
		return "", false
	}
}

func projectCapability(capability string, data json.RawMessage, format Format) ([]byte, error) {
	if format == FormatData {
		return data, nil
	}

	if format == FormatAuto {
		if _, ok := projectionRules[capability]; !ok {
			return data, nil
		}
		format = FormatCompact
	}

	rule, ok := projectionRules[capability]
	if !ok {
		return nil, &UnsupportedProjectionError{Capability: capability, Format: format}
	}

	var payload any
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, &InvalidProjectionError{Capability: capability, Format: format, Cause: err}
	}

	projected, err := applyProjectionSpec(payload, rule.Compact)
	if err != nil {
		return nil, &InvalidProjectionError{Capability: capability, Format: format, Cause: err}
	}
	body, err := json.Marshal(projected)
	if err != nil {
		return nil, &InvalidProjectionError{Capability: capability, Format: format, Cause: err}
	}
	return body, nil
}

func applyProjectionSpec(payload any, spec projectionSpec) (any, error) {
	if len(spec.Scalars) == 0 && len(spec.Lists) == 1 && len(spec.Tables) == 0 && spec.Lists[0].To == "$root" {
		return projectList(getRoot(payload), spec.Lists[0])
	}
	if len(spec.Scalars) == 0 && len(spec.Tables) == 1 && len(spec.Lists) == 0 && spec.Tables[0].To == "$root" {
		return projectTable(getRoot(payload), spec.Tables[0])
	}

	out := map[string]any{}
	for _, scalar := range spec.Scalars {
		value, ok, err := resolveFieldValue(payload, scalar)
		if err != nil {
			return nil, err
		}
		if ok {
			setPathValue(out, scalar.To, value)
		}
	}
	for _, list := range spec.Lists {
		value, ok, err := resolvePath(payload, list.From)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		projected, err := projectList(value, list)
		if err != nil {
			return nil, err
		}
		setPathValue(out, list.To, projected)
	}
	for _, table := range spec.Tables {
		value, ok, err := resolvePath(payload, table.From)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		projected, err := projectTable(value, table)
		if err != nil {
			return nil, err
		}
		setPathValue(out, table.To, projected)
	}
	return out, nil
}

func projectList(value any, rule listRule) ([]any, error) {
	items, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("expected list at %q", rule.From)
	}
	limit := boundedLimit(len(items), rule.Limit)
	out := make([]any, 0, limit)
	for _, item := range items[:limit] {
		if len(rule.Fields) == 0 {
			out = append(out, item)
			continue
		}
		row := map[string]any{}
		for _, field := range rule.Fields {
			resolved, ok, err := resolveFieldValue(item, field)
			if err != nil {
				return nil, err
			}
			if ok {
				setPathValue(row, field.To, resolved)
			}
		}
		out = append(out, row)
	}
	return out, nil
}

func projectTable(value any, rule tableRule) (map[string]any, error) {
	items, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("expected list at %q", rule.From)
	}
	limit := boundedLimit(len(items), rule.Limit)
	columns := make([]string, 0, len(rule.Columns))
	rows := make([][]any, 0, limit)
	for _, column := range rule.Columns {
		columns = append(columns, column.To)
	}
	for _, item := range items[:limit] {
		row := make([]any, 0, len(rule.Columns))
		for _, column := range rule.Columns {
			resolved, ok, err := resolveFieldValue(item, column)
			if err != nil {
				return nil, err
			}
			if ok {
				row = append(row, resolved)
			} else {
				row = append(row, nil)
			}
		}
		rows = append(rows, row)
	}
	return map[string]any{
		"columns": columns,
		"rows":    rows,
	}, nil
}

func resolveFieldValue(payload any, rule fieldRule) (any, bool, error) {
	value, ok, err := resolvePath(payload, rule.From)
	if err != nil || !ok {
		return nil, ok, err
	}
	switch rule.Transform {
	case transformNone:
		return value, true, nil
	case transformCount:
		items, ok := value.([]any)
		if !ok {
			return nil, false, nil
		}
		return len(items), true, nil
	case transformJoinDimensions:
		items, ok := value.([]any)
		if !ok {
			return nil, false, nil
		}
		parts := make([]string, 0, len(items))
		for _, item := range items {
			name, _, _ := resolvePath(item, "name")
			val, _, _ := resolvePath(item, "value")
			nameStr := strings.TrimSpace(toString(name))
			valStr := strings.TrimSpace(toString(val))
			switch {
			case nameStr != "" && valStr != "":
				parts = append(parts, nameStr+"="+valStr)
			case valStr != "":
				parts = append(parts, valStr)
			case nameStr != "":
				parts = append(parts, nameStr)
			}
		}
		return strings.Join(parts, "; "), true, nil
	default:
		return nil, false, errors.New("unknown transform")
	}
}

func resolvePath(payload any, path string) (any, bool, error) {
	if path == "" || path == "$root" {
		return payload, true, nil
	}
	current := payload
	for _, segment := range strings.Split(path, ".") {
		switch typed := current.(type) {
		case map[string]any:
			next, ok := typed[segment]
			if !ok {
				return nil, false, nil
			}
			current = next
		case []any:
			index, err := parseIndex(segment, len(typed))
			if err != nil {
				return nil, false, err
			}
			current = typed[index]
		default:
			return nil, false, nil
		}
	}
	return current, true, nil
}

func parseIndex(segment string, length int) (int, error) {
	var index int
	if _, err := fmt.Sscanf(segment, "%d", &index); err != nil {
		return 0, fmt.Errorf("expected numeric index, got %q", segment)
	}
	if index < 0 || index >= length {
		return 0, fmt.Errorf("index %d out of range", index)
	}
	return index, nil
}

func setPathValue(target map[string]any, path string, value any) {
	if path == "" || path == "$root" {
		return
	}
	parts := strings.Split(path, ".")
	current := target
	for _, segment := range parts[:len(parts)-1] {
		next, ok := current[segment].(map[string]any)
		if !ok {
			next = map[string]any{}
			current[segment] = next
		}
		current = next
	}
	current[parts[len(parts)-1]] = value
}

func boundedLimit(length int, limit int) int {
	if limit <= 0 || limit > length {
		return length
	}
	return limit
}

func getRoot(value any) any {
	return value
}

func projectAmazonASINSalesDailyTrend(data json.RawMessage, meta map[string]any) ([]byte, error) {
	projected, err := projectCapability("amazon.get_asin_sales_daily_trend", data, FormatCompact)
	if err != nil {
		return nil, err
	}
	var items any
	if err := json.Unmarshal(projected, &items); err != nil {
		return nil, err
	}
	out := map[string]any{
		"items": items.(map[string]any)["items"],
	}
	if asin, ok := meta["asin"]; ok {
		out["asin"] = asin
	}
	body, err := json.Marshal(out)
	return body, err
}

func projectAmazonASINSalesHistory(data json.RawMessage, meta map[string]any) ([]byte, error) {
	projected, err := projectCapability("amazon.get_asins_sales_history", data, FormatCompact)
	if err != nil {
		return nil, err
	}
	var payload map[string]any
	if err := json.Unmarshal(projected, &payload); err != nil {
		return nil, err
	}
	out := map[string]any{
		"items": payload["items"],
	}
	if queried, ok := meta["queried_asins"]; ok {
		out["queried_asins"] = queried
	}
	if withoutHistory, ok := meta["asins_without_history"]; ok {
		out["asins_without_history"] = withoutHistory
	}
	body, err := json.Marshal(out)
	return body, err
}

func projectAmazonVariantSales30d(data json.RawMessage, meta map[string]any) ([]byte, error) {
	projected, err := projectCapability("amazon.get_variant_sales_30d", data, FormatCompact)
	if err != nil {
		return nil, err
	}
	var payload map[string]any
	if err := json.Unmarshal(projected, &payload); err != nil {
		return nil, err
	}
	out := map[string]any{
		"items": payload["items"],
	}
	if queried, ok := meta["queried_asin"]; ok {
		out["queried_asin"] = queried
	}
	body, err := json.Marshal(out)
	return body, err
}

func projectAmazonSearchCategory(data json.RawMessage, _ map[string]any) ([]byte, error) {
	projected, err := projectCapability("amazon.search_category", data, FormatCompact)
	if err != nil {
		return nil, err
	}
	var payload map[string]any
	if err := json.Unmarshal(projected, &payload); err != nil {
		return nil, err
	}
	body, err := json.Marshal(map[string]any{
		"items": payload["items"],
	})
	return body, err
}

func projectAmazonCategoryTrend(data json.RawMessage, meta map[string]any) ([]byte, error) {
	var rows []map[string]any
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, err
	}

	metrics := extractStringList(meta, "metrics")
	columns := []string{"month"}
	seen := map[string]bool{"month": true}
	for _, metric := range metrics {
		if metric == "" || seen[metric] {
			continue
		}
		seen[metric] = true
		columns = append(columns, metric)
	}
	if len(columns) == 1 {
		for _, fallback := range projectionRules["amazon.get_category_trend"].Compact.Tables[0].Columns {
			if fallback.To == "month" || seen[fallback.To] {
				continue
			}
			seen[fallback.To] = true
			columns = append(columns, fallback.To)
		}
	}

	tableRows := make([][]any, 0, len(rows))
	for _, row := range rows {
		record := make([]any, 0, len(columns))
		for _, column := range columns {
			record = append(record, row[column])
		}
		tableRows = append(tableRows, record)
	}

	body, err := json.Marshal(map[string]any{
		"items": map[string]any{
			"columns": columns,
			"rows":    tableRows,
		},
	})
	return body, err
}

func projectGoogleAdsSearchAdvertisers(data json.RawMessage, _ map[string]any) ([]byte, error) {
	var payload map[string]any
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}

	out := map[string]any{}

	if advertisers, ok := payload["advertisers"].([]any); ok {
		if len(advertisers) == 0 {
			out["advertisers"] = nil
		} else {
			table, err := projectTable(advertisers, tableRule{
				From:  "advertisers",
				To:    "advertisers",
				Limit: 10,
				Columns: []fieldRule{
					{From: "advertiser_name", To: "name"},
					{From: "advertiser_id", To: "id"},
					{From: "region", To: "region"},
					{From: "ads_count", To: "ads_count"},
					{From: "is_verified", To: "is_verified"},
				},
			})
			if err != nil {
				return nil, err
			}
			out["advertisers"] = table
		}
	}

	if domains, ok := payload["domains"].([]any); ok {
		list, err := projectList(domains, listRule{
			From:  "domains",
			To:    "domains",
			Limit: 10,
			Fields: []fieldRule{
				{From: "domain", To: "name"},
			},
		})
		if err != nil {
			return nil, err
		}
		out["domains"] = list
	}

	body, err := json.Marshal(out)
	return body, err
}

func projectTikTokShopProducts(data json.RawMessage, _ map[string]any) ([]byte, error) {
	projected, err := projectCapability("tiktok.shop_products", data, FormatCompact)
	if err != nil {
		return nil, err
	}

	var items any
	if err := json.Unmarshal(projected, &items); err != nil {
		return nil, err
	}

	var rawItems []any
	if err := json.Unmarshal(data, &rawItems); err != nil {
		return nil, err
	}

	body, err := json.Marshal(map[string]any{
		"items":     items,
		"items_len": len(rawItems),
	})
	return body, err
}

func extractStringList(meta map[string]any, key string) []string {
	if meta == nil {
		return nil
	}
	raw, ok := meta[key]
	if !ok {
		return nil
	}
	items, ok := raw.([]any)
	if !ok {
		if typed, ok := raw.([]string); ok {
			return append([]string(nil), typed...)
		}
		return nil
	}
	out := make([]string, 0, len(items))
	for _, item := range items {
		if text, ok := item.(string); ok && strings.TrimSpace(text) != "" {
			out = append(out, text)
		}
	}
	return out
}

func toString(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	default:
		return fmt.Sprint(typed)
	}
}
