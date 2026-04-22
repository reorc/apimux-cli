package output

import "encoding/json"

type compatProjector func(data json.RawMessage, meta map[string]any) ([]byte, error)

type compactMetaPolicy struct {
	suppressHasMore     bool
	suppressCurrentPage bool
	suppressTotal       bool
}

var compatProjectors = map[string]compatProjector{
	"amazon.get_asin_sales_daily_trend": projectAmazonASINSalesDailyTrend,
	"amazon.get_asins_sales_history":    projectAmazonASINSalesHistory,
	"amazon.get_variant_sales_30d":      projectAmazonVariantSales30d,
	"amazon.search_category":            projectAmazonSearchCategory,
	"amazon.get_category_trend":         projectAmazonCategoryTrend,
	"google_ads.search_advertisers":     projectGoogleAdsSearchAdvertisers,
	"tiktok.shop_products":              projectTikTokShopProducts,
}

var compactMetaPolicies = map[string]compactMetaPolicy{
	"amazon.search_category":         {suppressHasMore: true, suppressCurrentPage: true, suppressTotal: true},
	"amazon.list_asin_keywords":      {suppressHasMore: true, suppressCurrentPage: true, suppressTotal: true},
	"amazon.query_aba_keywords":      {suppressHasMore: true, suppressCurrentPage: true, suppressTotal: true},
	"amazon.get_asin_sales_daily_trend": {suppressHasMore: true, suppressCurrentPage: true, suppressTotal: true},
	"amazon.get_asins_sales_history": {suppressHasMore: true, suppressCurrentPage: true, suppressTotal: true},
	"amazon.get_variant_sales_30d":   {suppressHasMore: true, suppressCurrentPage: true, suppressTotal: true},
	"amazon.get_product_reviews":     {suppressHasMore: true, suppressCurrentPage: true, suppressTotal: true},
	"amazon.get_category_trend":      {suppressHasMore: true, suppressCurrentPage: true, suppressTotal: true},
}

func renderCapabilityOutput(capability string, data json.RawMessage, meta map[string]any, format Format) ([]byte, error) {
	projected, err := projectCapabilityWithCompat(capability, data, meta, format)
	if err != nil {
		return nil, err
	}

	metadata := extractCriticalMetadata(capability, meta)
	if len(metadata) == 0 {
		return projected, nil
	}

	var projectedData any
	if err := json.Unmarshal(projected, &projectedData); err != nil {
		return nil, err
	}

	return json.Marshal(map[string]any{
		"data": projectedData,
		"meta": metadata,
	})
}

func projectCapabilityWithCompat(capability string, data json.RawMessage, meta map[string]any, format Format) ([]byte, error) {
	if format == FormatData {
		return data, nil
	}
	if projector, ok := compatProjectors[capability]; ok {
		return projector(data, meta)
	}
	return projectCapability(capability, data, format)
}

func projectCapabilityWithMeta(capability string, data json.RawMessage, meta map[string]any, format Format) ([]byte, bool, error) {
	if format == FormatData {
		return nil, false, nil
	}
	projector, ok := compatProjectors[capability]
	if !ok {
		return nil, false, nil
	}
	body, err := projector(data, meta)
	return body, true, err
}

func extractCriticalMetadata(capability string, meta map[string]any) map[string]any {
	if meta == nil {
		return nil
	}

	policy := compactMetaPolicies[capability]
	critical := make(map[string]any)

	if cursor, ok := meta["cursor"]; ok {
		critical["cursor"] = cursor
	}
	if hasMore, ok := meta["has_more"]; ok && !policy.suppressHasMore {
		critical["has_more"] = hasMore
	}
	if currentPage, ok := meta["current_page"]; ok && !policy.suppressCurrentPage {
		critical["current_page"] = currentPage
	}
	if nextPage, ok := meta["next_page"]; ok {
		critical["next_page"] = nextPage
	}
	if total, ok := meta["total"]; ok && !policy.suppressTotal {
		critical["total"] = total
	}
	if partial, ok := meta["partial"]; ok && partial == true {
		critical["partial"] = partial
		if subrequestCount, ok := meta["subrequest_count"]; ok {
			critical["subrequest_count"] = subrequestCount
		}
		if subrequests, ok := meta["subrequests"]; ok {
			critical["subrequests"] = subrequests
		}
	}

	if len(critical) == 0 {
		return nil
	}
	return critical
}
