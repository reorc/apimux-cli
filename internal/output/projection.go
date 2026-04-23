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
	PassThrough bool
	Scalars     []fieldRule
	Lists       []listRule
	Tables      []tableRule
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
	if spec.PassThrough {
		return payload, nil
	}
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

