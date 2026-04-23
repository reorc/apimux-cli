package output

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
						{From: "author.user_id", To: "author_id"},
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
						{From: "author.user_id", To: "author_id"},
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
			PassThrough: true,
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
						{From: "share_url", To: "share_url"},
						{From: "author.user_id", To: "author_id"},
						{From: "author.nickname", To: "author"},
						{From: "statistics.like_count", To: "like_count"},
						{From: "statistics.comment_count", To: "comment_count"},
						{From: "statistics.share_count", To: "share_count"},
						{From: "statistics.play_count", To: "play_count"},
						{From: "video.duration", To: "duration"},
						{From: "video.ratio", To: "ratio"},
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
						{From: "permalink", To: "permalink"},
						{From: "parent_id", To: "parent_id"},
						{From: "depth", To: "depth"},
					},
				},
			},
		},
	},
	"reddit.get_post_detail": {
		Compact: projectionSpec{
			PassThrough: true,
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
						{From: "author", To: "author"},
						{From: "score", To: "score"},
						{From: "upvote_ratio", To: "upvote_ratio"},
						{From: "num_comments", To: "num_comments"},
						{From: "created_time", To: "created_at"},
						{From: "url", To: "url"},
						{From: "permalink", To: "permalink"},
						{From: "thumbnail", To: "thumbnail"},
						{From: "is_video", To: "is_video"},
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
						{From: "author", To: "author"},
						{From: "score", To: "score"},
						{From: "upvote_ratio", To: "upvote_ratio"},
						{From: "num_comments", To: "num_comments"},
						{From: "created_time", To: "created_at"},
						{From: "url", To: "url"},
						{From: "permalink", To: "permalink"},
						{From: "thumbnail", To: "thumbnail"},
						{From: "is_video", To: "is_video"},
					},
				},
			},
		},
	},
	"trendcloud.get_market_trend": {
		Compact: projectionSpec{
			Lists: []listRule{
				{
					From: "$root",
					To:   "$root",
				},
			},
		},
	},
	"trendcloud.get_top_rankings": {
		Compact: projectionSpec{
			Lists: []listRule{
				{
					From: "$root",
					To:   "$root",
				},
			},
		},
	},
	"trendcloud.search_filter_values": {
		Compact: projectionSpec{
			Lists: []listRule{
				{
					From: "$root",
					To:   "$root",
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
						{From: "reply_count", To: "reply_count"},
						{From: "create_time", To: "create_time"},
						{From: "user_id", To: "user_id"},
						{From: "nickname", To: "author"},
					},
				},
			},
		},
	},
	"xiaohongshu.get_note_detail": {
		Compact: projectionSpec{
			PassThrough: true,
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
						{From: "xsec_token", To: "xsec_token"},
						{From: "title", To: "title"},
						{From: "description", To: "description"},
						{From: "type", To: "type"},
						{From: "like_count", To: "like_count"},
						{From: "collect_count", To: "collect_count"},
						{From: "comment_count", To: "comment_count"},
						{From: "author.user_id", To: "author_id"},
						{From: "author.nickname", To: "author"},
					},
				},
			},
		},
	},
}
