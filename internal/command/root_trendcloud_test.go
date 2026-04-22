package command

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func trendCloudSchemaResponse(capability string) string {
	switch capability {
	case "trendcloud.search_filter_values":
		return `{"ok":true,"data":{"name":"trendcloud.search_filter_values","parameters":[{"name":"kind","type":"string","required":true,"description":"Filter kind","enum":["category","brand","series","sku","attribute"]},{"name":"query","type":"string","required":true,"description":"Search keyword or path fragment"},{"name":"platforms","type":"array","description":"Platform scope","items_type":"string","encoding":"csv","enum":["douyin","jd","tmall"]},{"name":"categories","type":"array","description":"Category scope hint","items_type":"string","encoding":"csv"},{"name":"limit","type":"integer","description":"Result limit","default":10}]}}`
	case "trendcloud.get_market_trend":
		return `{"ok":true,"data":{"name":"trendcloud.get_market_trend","parameters":[{"name":"start_month","type":"string","description":"Start month in YYYY-MM"},{"name":"end_month","type":"string","description":"End month in YYYY-MM"},{"name":"metrics","type":"array","description":"Requested trend metrics","items_type":"string","encoding":"csv","enum":["sales","volume"]},{"name":"filters","type":"object","description":"Structured filters","encoding":"json","flag_name":"filters-json"}]}}`
	case "trendcloud.get_top_rankings":
		return `{"ok":true,"data":{"name":"trendcloud.get_top_rankings","parameters":[{"name":"entity","type":"string","required":true,"description":"Ranking entity","enum":["brand","category","series","sku","attribute"]},{"name":"metric","type":"string","description":"Entity-specific ranking metric"},{"name":"start_month","type":"string","description":"Start month in YYYY-MM"},{"name":"end_month","type":"string","description":"End month in YYYY-MM"},{"name":"top_n","type":"integer","description":"Result size","default":20},{"name":"category_level","type":"string","description":"Category level","enum":["category1","category2","category3"]},{"name":"filters","type":"object","description":"Structured filters","encoding":"json","flag_name":"filters-json"}]}}`
	default:
		return `{"ok":false,"error":{"code":"unknown_capability","message":"unknown capability"}}`
	}
}

func TestTrendCloudSearchFilterValuesCallsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/schema/trendcloud.search_filter_values":
			_, _ = w.Write([]byte(trendCloudSchemaResponse("trendcloud.search_filter_values")))
			return
		case "/v1/capabilities/trendcloud.search_filter_values":
			_, _ = w.Write([]byte(`{"ok":true,"data":[{"label":"饮料 > 咖啡","path":["饮料","咖啡"]}],"meta":{"capability":"trendcloud.search_filter_values","total":1}}`))
			return
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"trendcloud", "search_filter_values",
		"--kind", "category",
		"--query", "咖啡",
		"--platforms", "douyin,jd",
		"--limit", "5",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stdout=%s stderr=%s", exitCode, stdout.String(), stderr.String())
	}
	// JSON encoder escapes > as >
	if !strings.Contains(stdout.String(), `"label":"饮料`) || !strings.Contains(stdout.String(), `咖啡"`) {
		t.Fatalf("expected output to contain projected data, got %s", stdout.String())
	}
	if strings.Contains(stdout.String(), `"ok"`) {
		t.Fatalf("expected compact output without ok field, got %s", stdout.String())
	}
}

func TestTrendCloudGetMarketTrendCallsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/schema/trendcloud.get_market_trend":
			_, _ = w.Write([]byte(trendCloudSchemaResponse("trendcloud.get_market_trend")))
			return
		case "/v1/capabilities/trendcloud.get_market_trend":
			_, _ = w.Write([]byte(`{"ok":true,"data":[{"period":"2026-01","sales":1234.56}],"meta":{"capability":"trendcloud.get_market_trend","metrics":["sales"]}}`))
			return
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"trendcloud", "get_market_trend",
		"--start-month", "2026-01",
		"--end-month", "2026-03",
		"--metrics", "sales,volume",
		"--filters-json", `{"platforms":["douyin"],"brands":["瑞幸"]}`,
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stdout=%s stderr=%s", exitCode, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), `"sales":1234.56`) {
		t.Fatalf("expected data-only output, got %s", stdout.String())
	}
	if strings.Contains(stdout.String(), `"meta"`) || strings.Contains(stdout.String(), `"ok"`) {
		t.Fatalf("expected data-only output without envelope fields, got %s", stdout.String())
	}
}

func TestTrendCloudGetTopRankingsCallsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/schema/trendcloud.get_top_rankings":
			_, _ = w.Write([]byte(trendCloudSchemaResponse("trendcloud.get_top_rankings")))
			return
		case "/v1/capabilities/trendcloud.get_top_rankings":
			_, _ = w.Write([]byte(`{"ok":true,"data":[{"rank":1,"label":"瑞幸","sales":2000.12}],"meta":{"capability":"trendcloud.get_top_rankings","entity":"brand","metric":"sales"}}`))
			return
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"trendcloud", "get_top_rankings",
		"--entity", "brand",
		"--metric", "sales",
		"--top-n", "10",
		"--filters-json", `{"platforms":["tmall"]}`,
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stdout=%s stderr=%s", exitCode, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), `"label":"瑞幸"`) {
		t.Fatalf("expected data-only output, got %s", stdout.String())
	}
	if strings.Contains(stdout.String(), `"meta"`) || strings.Contains(stdout.String(), `"ok"`) {
		t.Fatalf("expected data-only output without envelope fields, got %s", stdout.String())
	}
}

func TestTrendCloudSearchFilterValuesRequiresKindAndQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/schema/trendcloud.search_filter_values" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(trendCloudSchemaResponse("trendcloud.search_filter_values")))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"trendcloud", "search_filter_values",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 2 {
		t.Fatalf("expected exit code 2, got %d, stdout=%s stderr=%s", exitCode, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), `"trendcloud search_filter_values requires --kind and --query"`) {
		t.Fatalf("expected required-flag error, got %s", stdout.String())
	}
}
