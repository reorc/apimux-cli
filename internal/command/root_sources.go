package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newStaticSourceCommand(use, short, long, example, supports string, subcommands ...*cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:     use,
		Short:   short,
		Long:    long,
		Example: example,
		RunE: func(cmd *cobra.Command, args []string) error {
			return &cliError{
				exitCode: 2,
				code:     "cli_invalid_command",
				message:  fmt.Sprintf("%s supports: %s", use, supports),
			}
		},
	}
	cmd.AddCommand(subcommands...)
	return cmd
}

func (r *Root) newAmazonCommand(runCtx *runContext) *cobra.Command {
	return newStaticSourceCommand(
		"amazon",
		"Amazon product, keyword, and sales data",
		"Amazon product, keyword, and sales data.\n\nUse this command for Amazon product lookup, keyword research, sales history, reviews, and category trend endpoints.",
		"  apimux amazon search_products --keyword 'desk lamp'\n  apimux amazon get_product --asin B0EXAMPLE",
		"expand_keywords, get_keyword_overview, get_keyword_trends, list_asin_keywords, query_aba_keywords, search_category, get_asin_sales_daily_trend, get_asins_sales_history, get_variant_sales_30d, get_product, search_products, get_product_reviews, get_category_best_sellers, get_category_trend",
		newSchemaBoundCapabilityCommand(runCtx, "amazon.expand_keywords", "expand_keywords", "Expand one Amazon keyword into related queries", "amazon expand_keywords"),
		newSchemaBoundCapabilityCommand(runCtx, "amazon.get_keyword_overview", "get_keyword_overview", "Fetch one Amazon keyword overview", "amazon get_keyword_overview"),
		newSchemaBoundCapabilityCommand(runCtx, "amazon.get_keyword_trends", "get_keyword_trends", "Fetch Amazon keyword search trends", "amazon get_keyword_trends"),
		newSchemaBoundCapabilityCommand(runCtx, "amazon.list_asin_keywords", "list_asin_keywords", "List keywords associated with one Amazon ASIN", "amazon list_asin_keywords"),
		newSchemaBoundCapabilityCommand(runCtx, "amazon.query_aba_keywords", "query_aba_keywords", "Query Amazon Brand Analytics keywords", "amazon query_aba_keywords"),
		newSchemaBoundCapabilityCommand(runCtx, "amazon.search_category", "search_category", "Search Amazon categories by name", "amazon search_category"),
		newSchemaBoundCapabilityCommand(runCtx, "amazon.get_asin_sales_daily_trend", "get_asin_sales_daily_trend", "Fetch daily sales trend for one Amazon ASIN", "amazon get_asin_sales_daily_trend"),
		newSchemaBoundCapabilityCommand(runCtx, "amazon.get_asins_sales_history", "get_asins_sales_history", "Fetch monthly sales history for multiple Amazon ASINs", "amazon get_asins_sales_history"),
		newSchemaBoundCapabilityCommand(runCtx, "amazon.get_variant_sales_30d", "get_variant_sales_30d", "Fetch 30-day sales for Amazon variants", "amazon get_variant_sales_30d"),
		newSchemaBoundCapabilityCommand(runCtx, "amazon.get_product", "get_product", "Fetch one Amazon product", "amazon get_product"),
		newSchemaBoundCapabilityCommand(runCtx, "amazon.search_products", "search_products", "Search Amazon products", "amazon search_products"),
		newSchemaBoundCapabilityCommand(runCtx, "amazon.get_product_reviews", "get_product_reviews", "Fetch canonical Amazon product reviews", "amazon get_product_reviews"),
		newSchemaBoundCapabilityCommand(runCtx, "amazon.get_category_best_sellers", "get_category_best_sellers", "Fetch Amazon category best sellers", "amazon get_category_best_sellers"),
		newSchemaBoundCapabilityCommand(runCtx, "amazon.get_category_trend", "get_category_trend", "Fetch Amazon category trend metrics", "amazon get_category_trend"),
	)
}

func (r *Root) newGoogleTrendsCommand(runCtx *runContext) *cobra.Command {
	googleCmd := &cobra.Command{
		Use:     "google_trends",
		Short:   "Google Trends interest and search data",
		Long:    "Google Trends interest and search data.\n\nUse this command to fetch Google Trends time-series data for one or more queries.",
		Example: "  apimux google_trends get_interest_over_time --query openai --time_range today_12-m",
		RunE: func(cmd *cobra.Command, args []string) error {
			return &cliError{
				exitCode: 2,
				code:     "cli_invalid_command",
				message:  "google_trends supports: get_interest_over_time",
			}
		},
	}
	googleCmd.AddCommand(newSchemaBoundCapabilityCommand(
		runCtx,
		"google_trends.get_interest_over_time",
		"get_interest_over_time",
		"Fetch Google Trends interest over time",
		"google_trends get_interest_over_time",
	))
	return googleCmd
}

func (r *Root) newTrendCloudCommand(runCtx *runContext) *cobra.Command {
	trendCmd := &cobra.Command{
		Use:     "trendcloud",
		Short:   "TrendCloud market trends and rankings",
		Long:    "TrendCloud market trends and rankings.\n\nUse this command to search TrendCloud filters, fetch market trends, and inspect ranking data.",
		Example: "  apimux trendcloud get_market_trend --category 'electronics'",
		RunE: func(cmd *cobra.Command, args []string) error {
			return &cliError{
				exitCode: 2,
				code:     "cli_invalid_command",
				message:  "trendcloud supports: search_filter_values, get_market_trend, get_top_rankings",
			}
		},
	}
	trendCmd.AddCommand(
		newSchemaBoundCapabilityCommand(
			runCtx,
			"trendcloud.search_filter_values",
			"search_filter_values",
			"Search TrendCloud filter candidates",
			"trendcloud search_filter_values",
		),
		newSchemaBoundCapabilityCommand(
			runCtx,
			"trendcloud.get_market_trend",
			"get_market_trend",
			"Fetch TrendCloud market trend series",
			"trendcloud get_market_trend",
		),
		newSchemaBoundCapabilityCommand(
			runCtx,
			"trendcloud.get_top_rankings",
			"get_top_rankings",
			"Fetch TrendCloud rankings by entity and metric",
			"trendcloud get_top_rankings",
		),
	)
	return trendCmd
}

func (r *Root) newDouyinCommand(runCtx *runContext) *cobra.Command {
	return newStaticSourceCommand(
		"douyin",
		"Douyin video search and comments",
		"Douyin video search and comments.\n\nUse this command to search Douyin videos and inspect video details, comments, and replies.",
		"  apimux douyin search_videos --keyword 美食\n  apimux douyin get_video_comments --aweme-id 1234567890",
		"search_videos, get_video_detail, get_video_comments, get_comment_replies",
		newSchemaBoundCapabilityCommand(runCtx, "douyin.search_videos", "search_videos", "Search Douyin videos", "douyin search_videos"),
		newSchemaBoundCapabilityCommand(runCtx, "douyin.get_video_detail", "get_video_detail", "Fetch one Douyin video detail", "douyin get_video_detail"),
		newSchemaBoundCapabilityCommand(runCtx, "douyin.get_video_comments", "get_video_comments", "List Douyin video comments", "douyin get_video_comments"),
		newSchemaBoundCapabilityCommand(runCtx, "douyin.get_comment_replies", "get_comment_replies", "List Douyin comment replies", "douyin get_comment_replies"),
	)
}

func (r *Root) newGoogleAdsCommand(runCtx *runContext) *cobra.Command {
	return newStaticSourceCommand(
		"google_ads",
		"Google Ads advertiser and creative data",
		"Google Ads advertiser and creative data.\n\nUse this command to search advertisers, list ad creatives, and fetch creative details from Google Ads transparency endpoints.",
		"  apimux google_ads search_advertisers --query openai\n  apimux google_ads list_ad_creatives --domain openai.com",
		"search_advertisers, list_ad_creatives, get_ad_details",
		newSchemaBoundCapabilityCommand(runCtx, "google_ads.search_advertisers", "search_advertisers", "Search Google Ads advertisers", "google_ads search_advertisers"),
		newSchemaBoundCapabilityCommand(runCtx, "google_ads.list_ad_creatives", "list_ad_creatives", "List Google Ads ad creatives", "google_ads list_ad_creatives"),
		newSchemaBoundCapabilityCommand(runCtx, "google_ads.get_ad_details", "get_ad_details", "Fetch one Google Ads creative detail", "google_ads get_ad_details"),
	)
}

func (r *Root) newMetaAdsCommand(runCtx *runContext) *cobra.Command {
	return newStaticSourceCommand(
		"meta_ads",
		"Meta Ads Library search and details",
		"Meta Ads Library search and details.\n\nUse this command to search the Meta Ads Library and fetch one ad detail by ID.",
		"  apimux meta_ads search_ads --q openai\n  apimux meta_ads get_ad_detail --ad-id 123456789",
		"search_ads, get_ad_detail",
		newSchemaBoundCapabilityCommand(runCtx, "meta_ads.search_ads", "search_ads", "Search Meta Ads Library ads", "meta_ads search_ads"),
		newSchemaBoundCapabilityCommand(runCtx, "meta_ads.get_ad_detail", "get_ad_detail", "Fetch one Meta ad detail", "meta_ads get_ad_detail"),
	)
}

func (r *Root) newRedditCommand(runCtx *runContext) *cobra.Command {
	return newStaticSourceCommand(
		"reddit",
		"Reddit post search and comments",
		"Reddit post search and comments.\n\nUse this command to search Reddit posts, inspect subreddit feeds, and fetch post details and comments.",
		"  apimux reddit search --query openai\n  apimux reddit get_post_comments --post-id t3_abcdef",
		"search, get_subreddit_feed, get_post_detail, get_post_comments",
		newSchemaBoundCapabilityCommand(runCtx, "reddit.search", "search", "Search Reddit posts", "reddit search"),
		newSchemaBoundCapabilityCommand(runCtx, "reddit.get_subreddit_feed", "get_subreddit_feed", "List one subreddit feed", "reddit get_subreddit_feed"),
		newSchemaBoundCapabilityCommand(runCtx, "reddit.get_post_detail", "get_post_detail", "Fetch one Reddit post detail", "reddit get_post_detail"),
		newSchemaBoundCapabilityCommand(runCtx, "reddit.get_post_comments", "get_post_comments", "List Reddit post comments", "reddit get_post_comments"),
	)
}

func (r *Root) newTiktokCommand(runCtx *runContext) *cobra.Command {
	return newStaticSourceCommand(
		"tiktok",
		"TikTok video, comment, and shop data",
		"TikTok video, comment, and shop data.\n\nUse this command to search TikTok videos, list comments, and query TikTok Shop products and product details.",
		"  apimux tiktok search_videos --keyword laptop\n  apimux tiktok shop_products --seller-id 123456",
		"search_videos, list_comments, shop_products, shop_product_info",
		newSchemaBoundCapabilityCommand(runCtx, "tiktok.search_videos", "search_videos", "Search TikTok videos", "tiktok search_videos"),
		newSchemaBoundCapabilityCommand(runCtx, "tiktok.list_comments", "list_comments", "List TikTok video comments", "tiktok list_comments"),
		newSchemaBoundCapabilityCommand(runCtx, "tiktok.shop_products", "shop_products", "List TikTok Shop seller products", "tiktok shop_products"),
		newSchemaBoundCapabilityCommand(runCtx, "tiktok.shop_product_info", "shop_product_info", "Fetch one TikTok Shop product detail", "tiktok shop_product_info"),
	)
}

func (r *Root) newXiaohongshuCommand(runCtx *runContext) *cobra.Command {
	return newStaticSourceCommand(
		"xiaohongshu",
		"Xiaohongshu note search and comments",
		"Xiaohongshu note search and comments.\n\nUse this command to search Xiaohongshu notes and inspect note details and comment threads.",
		"  apimux xiaohongshu search_notes --keyword 护肤\n  apimux xiaohongshu get_note_detail --note-id 64cdef1234567890abcdef12",
		"search_notes, get_note_detail, get_note_comments",
		newSchemaBoundCapabilityCommand(runCtx, "xiaohongshu.search_notes", "search_notes", "Search Xiaohongshu notes", "xiaohongshu search_notes"),
		newSchemaBoundCapabilityCommand(runCtx, "xiaohongshu.get_note_detail", "get_note_detail", "Fetch one Xiaohongshu note detail", "xiaohongshu get_note_detail"),
		newSchemaBoundCapabilityCommand(runCtx, "xiaohongshu.get_note_comments", "get_note_comments", "List Xiaohongshu note comments", "xiaohongshu get_note_comments"),
	)
}
