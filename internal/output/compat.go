package output

import (
	"encoding/json"
	"fmt"
	"strings"
)

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
