package command

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAmazonSearchCategoryCallsService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if maybeServeSchema(w, r) {
			return
		}
		if r.URL.Path != "/v1/capabilities/amazon.search_category" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":[{"node_id":"123","name":"Cell Phones"}],"meta":{"capability":"amazon.search_category","total":1}}`))
	}))
	defer server.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	root := NewRoot(&stdout, &stderr)
	exitCode, err := root.Execute(context.Background(), []string{
		"--base-url", server.URL,
		"amazon", "search_category",
		"--name", "cell phone",
		"--market", "US",
		"--limit", "5",
	})
	if err != nil {
		t.Fatalf("execute root: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"columns":["node_id","name","cn_name","path"]`) || !strings.Contains(stdout.String(), `"123"`) {
		t.Fatalf("expected compact table output, got %s", stdout.String())
	}
	if strings.Contains(stdout.String(), `"meta"`) || strings.Contains(stdout.String(), `"ok"`) {
		t.Fatalf("expected data-only output without envelope fields, got %s", stdout.String())
	}
}
