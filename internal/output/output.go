package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type BodyOutput string

const (
	BodyOutputCompact    BodyOutput = "compact"
	BodyOutputPretty     BodyOutput = "pretty"
	BodyOutputData       BodyOutput = "data"
	BodyOutputDataPretty BodyOutput = "data-pretty"
)

type Renderer struct {
	Stdout io.Writer
	Stderr io.Writer
}

type envelope struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error json.RawMessage `json:"error"`
	Meta  map[string]any  `json:"meta"`
}

func ParseBodyOutput(value string) (BodyOutput, bool) {
	switch strings.TrimSpace(value) {
	case "", string(BodyOutputCompact):
		return BodyOutputCompact, true
	case string(BodyOutputPretty):
		return BodyOutputPretty, true
	case string(BodyOutputData):
		return BodyOutputData, true
	case string(BodyOutputDataPretty):
		return BodyOutputDataPretty, true
	default:
		return "", false
	}
}

func (o BodyOutput) pretty() bool {
	return o == BodyOutputPretty || o == BodyOutputDataPretty
}

func (o BodyOutput) projectionFormat() Format {
	switch o {
	case BodyOutputData, BodyOutputDataPretty:
		return FormatData
	default:
		return FormatAuto
	}
}

func (r Renderer) WriteCapabilityResponse(body []byte, output BodyOutput, debug bool) error {
	var env envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return err
	}

	if debug {
		debugPayload := map[string]any{
			"ok": env.OK,
		}
		if len(env.Data) > 0 && string(env.Data) != "null" {
			debugPayload["data"] = json.RawMessage(env.Data)
		}
		if len(env.Error) > 0 && string(env.Error) != "null" {
			debugPayload["error"] = json.RawMessage(env.Error)
		}
		if env.Meta != nil {
			delete(env.Meta, "source")
			debugPayload["meta"] = env.Meta
		}
		sanitized, err := json.Marshal(debugPayload)
		if err != nil {
			return err
		}
		return r.writeJSON(sanitized, false)
	}

	if len(env.Error) > 0 && string(env.Error) != "null" {
		payload, err := json.Marshal(map[string]json.RawMessage{"error": env.Error})
		if err != nil {
			return err
		}
		return r.writeJSON(payload, output.pretty())
	}
	if len(env.Data) == 0 {
		return r.writeJSON([]byte("null"), output.pretty())
	}

	capability, _ := env.Meta["capability"].(string)
	if custom, ok, err := projectCapabilityWithMeta(capability, env.Data, env.Meta, output.projectionFormat()); err != nil {
		return err
	} else if ok {
		metadata := extractCriticalMetadata(env.Meta)
		if len(metadata) > 0 {
			var projectedData any
			if err := json.Unmarshal(custom, &projectedData); err != nil {
				return err
			}
			wrapped, err := json.Marshal(map[string]any{
				"data": projectedData,
				"meta": metadata,
			})
			if err != nil {
				return err
			}
			return r.writeJSON(wrapped, output.pretty())
		}
		return r.writeJSON(custom, output.pretty())
	}
	projected, err := projectCapability(capability, env.Data, output.projectionFormat())
	if err != nil {
		return err
	}

	// Extract critical metadata for compact mode
	metadata := extractCriticalMetadata(env.Meta)
	if len(metadata) > 0 {
		// Wrap projected data with metadata
		var projectedData any
		if err := json.Unmarshal(projected, &projectedData); err != nil {
			return err
		}
		wrapper := map[string]any{
			"data": projectedData,
			"meta": metadata,
		}
		wrapped, err := json.Marshal(wrapper)
		if err != nil {
			return err
		}
		return r.writeJSON(wrapped, output.pretty())
	}

	return r.writeJSON(projected, output.pretty())
}

func (r Renderer) writeJSON(body []byte, pretty bool) error {
	var out bytes.Buffer
	if pretty {
		if err := json.Indent(&out, body, "", "  "); err != nil {
			return err
		}
		out.WriteByte('\n')
	} else {
		out.Write(body)
		if len(body) == 0 || body[len(body)-1] != '\n' {
			out.WriteByte('\n')
		}
	}
	_, err := r.Stdout.Write(out.Bytes())
	return err
}

func (r Renderer) Diagnostic(format string, args ...any) {
	if r.Stderr == nil {
		return
	}
	fmt.Fprintf(r.Stderr, format+"\n", args...)
}

func (r Renderer) WriteLocalError(message string, code string) error {
	body, err := json.Marshal(map[string]any{
		"error": map[string]any{
			"type":    "internal",
			"code":    code,
			"message": message,
		},
	})
	if err != nil {
		return err
	}
	return r.writeJSON(body, false)
}

// extractCriticalMetadata extracts pagination and partial-failure metadata
// that agents need to see even in compact mode
func extractCriticalMetadata(meta map[string]any) map[string]any {
	if meta == nil {
		return nil
	}

	critical := make(map[string]any)
	capability, _ := meta["capability"].(string)

	// Pagination metadata
	if cursor, ok := meta["cursor"]; ok {
		critical["cursor"] = cursor
	}
	if hasMore, ok := meta["has_more"]; ok && !suppressAmazonCompactMeta(capability) {
		critical["has_more"] = hasMore
	}
	if currentPage, ok := meta["current_page"]; ok && !suppressAmazonCompactMeta(capability) {
		critical["current_page"] = currentPage
	}
	if nextPage, ok := meta["next_page"]; ok {
		critical["next_page"] = nextPage
	}
	if total, ok := meta["total"]; ok && !suppressAmazonCompactMeta(capability) {
		critical["total"] = total
	}

	// Partial-failure metadata
	if partial, ok := meta["partial"]; ok && partial == true {
		critical["partial"] = partial
		if subrequestCount, ok := meta["subrequest_count"]; ok {
			critical["subrequest_count"] = subrequestCount
		}
		if subrequests, ok := meta["subrequests"]; ok {
			critical["subrequests"] = subrequests
		}
	}

	return critical
}

func suppressAmazonCompactMeta(capability string) bool {
	switch capability {
	case "amazon.search_category",
		"amazon.list_asin_keywords",
		"amazon.query_aba_keywords",
		"amazon.get_asin_sales_daily_trend",
		"amazon.get_asins_sales_history",
		"amazon.get_variant_sales_30d",
		"amazon.get_product_reviews",
		"amazon.get_category_trend":
		return true
	default:
		return false
	}
}
