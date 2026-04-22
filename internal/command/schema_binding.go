package command

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/reorc/apimux-cli/internal/schema"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type schemaEnvelope struct {
	OK    bool                    `json:"ok"`
	Data  schema.CapabilitySchema `json:"data"`
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type flagBinding struct {
	capability  string
	param       schema.CapabilityParam
	flagName    string
	hasDefault  bool
	stringValue *string
	intValue    *int
	floatValue  *float64
	boolValue   *bool
}

func newSchemaBoundCapabilityCommand(runCtx *runContext, capability, use, short, commandPath string) *cobra.Command {
	return &cobra.Command{
		Use:                use,
		Short:              short,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			spec, err := fetchCapabilitySchema(cmd.Context(), runCtx, capability)
			if err != nil {
				return err
			}
			if len(spec.Parameters) == 0 {
				return &cliError{
					exitCode: 1,
					code:     "cli_schema_missing_parameters",
					message:  fmt.Sprintf("%s does not have configurable parameters", capability),
				}
			}
			params, err := parseSchemaBoundParams(stripPersistentArgs(args), spec, commandPath)
			if err != nil {
				return err
			}
			return callCapability(cmd.Context(), runCtx, capability, params)
		},
	}
}

func fetchCapabilitySchema(ctx context.Context, runCtx *runContext, capability string) (schema.CapabilitySchema, error) {
	if runCtx.verbose {
		outputDiagnostic(runCtx.stderr, "[apimux] GET /v1/schema/%s", capability)
	}
	resp, err := runCtx.client.GetSchema(ctx, capability)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return schema.CapabilitySchema{}, &cliError{exitCode: 3, code: "cli_timeout", message: "schema request timed out"}
		}
		return schema.CapabilitySchema{}, &cliError{exitCode: 1, code: "cli_transport_error", message: err.Error()}
	}
	if runCtx.verbose {
		outputDiagnostic(runCtx.stderr, "[apimux] HTTP %d", resp.StatusCode)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return schema.CapabilitySchema{}, &cliError{
			exitCode: exitCodeForHTTPStatus(resp.StatusCode),
			code:     "cli_schema_fetch_failed",
			message:  fmt.Sprintf("schema lookup failed for %s", capability),
		}
	}

	var env schemaEnvelope
	if err := json.Unmarshal(resp.Body, &env); err != nil {
		return schema.CapabilitySchema{}, &cliError{
			exitCode: 1,
			code:     "cli_schema_decode_failed",
			message:  "schema response must be valid JSON",
		}
	}
	if !env.OK {
		message := strings.TrimSpace(env.Error.Message)
		if message == "" {
			message = fmt.Sprintf("schema lookup failed for %s", capability)
		}
		return schema.CapabilitySchema{}, &cliError{
			exitCode: 1,
			code:     "cli_schema_fetch_failed",
			message:  message,
		}
	}
	return env.Data, nil
}

func parseSchemaBoundParams(args []string, spec schema.CapabilitySchema, commandPath string) (map[string]any, error) {
	flagSet := pflag.NewFlagSet(commandPath, pflag.ContinueOnError)
	flagSet.SetOutput(io.Discard)

	bindings := make([]flagBinding, 0, len(spec.Parameters))
	for _, param := range spec.Parameters {
		binding, err := registerSchemaFlag(flagSet, spec.Name, param)
		if err != nil {
			return nil, err
		}
		bindings = append(bindings, binding)
	}

	if err := flagSet.Parse(args); err != nil {
		return nil, &cliError{exitCode: 2, code: "cli_invalid_flags", message: err.Error()}
	}

	missing := make([]string, 0)
	out := make(map[string]any, len(bindings))
	provided := make(map[string]bool, len(bindings))
	paramLookup := make(map[string]schema.CapabilityParam, len(bindings))
	for _, binding := range bindings {
		paramLookup[binding.param.Name] = binding.param
		value, isProvided, err := binding.resolve(flagSet)
		if err != nil {
			return nil, err
		}
		if !isProvided {
			if binding.param.Required && !binding.hasDefault {
				missing = append(missing, "--"+binding.flagName)
			}
			continue
		}
		out[binding.param.Name] = value
		provided[binding.param.Name] = true
	}

	if len(missing) > 0 {
		return nil, &cliError{
			exitCode: 2,
			code:     "cli_invalid_params",
			message:  fmt.Sprintf("%s requires %s", commandPath, joinFlagNames(missing)),
		}
	}

	if err := validateSchemaRules(spec.Rules, provided, paramLookup, commandPath); err != nil {
		return nil, err
	}

	return out, nil
}

func registerSchemaFlag(flagSet *pflag.FlagSet, capability string, param schema.CapabilityParam) (flagBinding, error) {
	binding := flagBinding{
		capability: capability,
		param:      param,
		flagName:   schemaFlagName(param),
		hasDefault: param.Default != nil,
	}
	usage := param.Description
	if usage == "" {
		usage = param.Name
	}

	switch param.Type {
	case "string":
		def, _ := param.Default.(string)
		binding.stringValue = flagSet.String(binding.flagName, def, usage)
	case "integer":
		def, err := toInt(param.Default)
		if err != nil {
			return flagBinding{}, &cliError{exitCode: 1, code: "cli_schema_invalid_default", message: fmt.Sprintf("invalid default for %s", param.Name)}
		}
		binding.intValue = flagSet.Int(binding.flagName, def, usage)
	case "number":
		def, err := toFloat(param.Default)
		if err != nil {
			return flagBinding{}, &cliError{exitCode: 1, code: "cli_schema_invalid_default", message: fmt.Sprintf("invalid default for %s", param.Name)}
		}
		binding.floatValue = flagSet.Float64(binding.flagName, def, usage)
	case "boolean":
		def, err := toBool(param.Default)
		if err != nil {
			return flagBinding{}, &cliError{exitCode: 1, code: "cli_schema_invalid_default", message: fmt.Sprintf("invalid default for %s", param.Name)}
		}
		binding.boolValue = flagSet.Bool(binding.flagName, def, usage)
	case "array", "object":
		def, _ := param.Default.(string)
		binding.stringValue = flagSet.String(binding.flagName, def, usage)
	default:
		return flagBinding{}, &cliError{
			exitCode: 1,
			code:     "cli_schema_unsupported_type",
			message:  fmt.Sprintf("unsupported schema param type %q for %s", param.Type, param.Name),
		}
	}

	return binding, nil
}

func (b flagBinding) resolve(flagSet *pflag.FlagSet) (any, bool, error) {
	flag := flagSet.Lookup(b.flagName)
	if flag == nil {
		return nil, false, &cliError{exitCode: 1, code: "cli_schema_missing_flag", message: fmt.Sprintf("flag %s was not registered", b.flagName)}
	}
	changed := flag.Changed
	if !changed && !b.hasDefault {
		return nil, false, nil
	}

	switch b.param.Type {
	case "string":
		value := strings.TrimSpace(derefString(b.stringValue))
		if value == "" && b.param.Required && !b.hasDefault {
			return nil, false, nil
		}
		value = normalizeCompatStringValue(b.capability, b.param.Name, value)
		if err := validateEnum(value, b.param); err != nil {
			return nil, false, err
		}
		return value, true, nil
	case "integer":
		value := derefInt(b.intValue)
		if err := validateEnum(strconv.Itoa(value), b.param); err != nil {
			return nil, false, err
		}
		if !changed && !b.hasDefault {
			return nil, false, nil
		}
		return value, true, nil
	case "number":
		value := derefFloat(b.floatValue)
		if !changed && !b.hasDefault {
			return nil, false, nil
		}
		return value, true, nil
	case "boolean":
		value := derefBool(b.boolValue)
		if !changed && !b.hasDefault {
			return nil, false, nil
		}
		return value, true, nil
	case "array":
		raw := strings.TrimSpace(derefString(b.stringValue))
		if raw == "" {
			if b.param.Required && !b.hasDefault {
				return nil, false, nil
			}
			return nil, false, nil
		}
		switch b.param.Encoding {
		case "", "csv":
			values := splitCSV(raw)
			for i, value := range values {
				values[i] = normalizeCompatStringValue(b.capability, b.param.Name, value)
				if err := validateEnum(values[i], b.param); err != nil {
					return nil, false, err
				}
			}
			return values, true, nil
		default:
			return nil, false, &cliError{exitCode: 1, code: "cli_schema_unsupported_encoding", message: fmt.Sprintf("unsupported encoding %q for %s", b.param.Encoding, b.param.Name)}
		}
	case "object":
		raw := strings.TrimSpace(derefString(b.stringValue))
		if raw == "" {
			return nil, false, nil
		}
		switch b.param.Encoding {
		case "", "json":
			payload, err := parseObjectFlag(raw, b.flagName)
			if err != nil {
				return nil, false, err
			}
			return payload, true, nil
		default:
			return nil, false, &cliError{exitCode: 1, code: "cli_schema_unsupported_encoding", message: fmt.Sprintf("unsupported encoding %q for %s", b.param.Encoding, b.param.Name)}
		}
	default:
		return nil, false, &cliError{exitCode: 1, code: "cli_schema_unsupported_type", message: fmt.Sprintf("unsupported schema param type %q for %s", b.param.Type, b.param.Name)}
	}
}

func schemaFlagName(param schema.CapabilityParam) string {
	if strings.TrimSpace(param.FlagName) != "" {
		return strings.TrimSpace(param.FlagName)
	}
	return strings.ReplaceAll(param.Name, "_", "-")
}

func validateEnum(value string, param schema.CapabilityParam) error {
	if len(param.Enum) == 0 || value == "" {
		return nil
	}
	for _, candidate := range param.Enum {
		if value == candidate {
			return nil
		}
	}
	return &cliError{
		exitCode: 2,
		code:     "cli_invalid_params",
		message:  fmt.Sprintf("--%s must be one of: %s", schemaFlagName(param), strings.Join(param.Enum, ", ")),
	}
}

func normalizeCompatStringValue(capability, paramName, value string) string {
	if value == "" {
		return value
	}
	if capability == "reddit.search" && paramName == "sort" {
		return strings.ToLower(value)
	}

	var aliases map[string]string
	switch capability {
	case "douyin.search_videos":
		switch paramName {
		case "sort_type":
			aliases = map[string]string{"0": "comprehensive", "1": "likes", "2": "latest"}
		case "publish_time":
			aliases = map[string]string{"0": "all", "1": "1d", "7": "1w", "180": "6m"}
		}
	case "tiktok.search_videos":
		switch paramName {
		case "sort_by":
			aliases = map[string]string{"0": "relevance", "1": "likes", "2": "date"}
		case "publish_time":
			aliases = map[string]string{"0": "all", "1": "1d", "7": "1w", "30": "1m", "90": "3m", "180": "6m"}
		}
	case "xiaohongshu.search_notes":
		switch paramName {
		case "sort_strategy":
			aliases = map[string]string{
				"general":               "default",
				"time_descending":       "latest",
				"popularity_descending": "likes",
			}
		case "note_type":
			aliases = map[string]string{"0": "all", "1": "normal", "2": "video"}
		}
	}
	if aliases == nil {
		return value
	}
	if mapped, ok := aliases[value]; ok {
		return mapped
	}
	return value
}

func validateSchemaRules(rules []schema.CapabilityRule, provided map[string]bool, params map[string]schema.CapabilityParam, commandPath string) error {
	for _, rule := range rules {
		switch rule.Kind {
		case "any_of_required":
			satisfied := false
			flags := make([]string, 0, len(rule.Params))
			for _, name := range rule.Params {
				param, ok := params[name]
				if !ok {
					return &cliError{
						exitCode: 1,
						code:     "cli_schema_invalid_rule",
						message:  fmt.Sprintf("schema rule %q references unknown param %q", rule.Kind, name),
					}
				}
				flags = append(flags, "--"+schemaFlagName(param))
				if provided[name] {
					satisfied = true
				}
			}
			if satisfied {
				continue
			}
			message := strings.TrimSpace(rule.Message)
			if message == "" {
				message = fmt.Sprintf("%s requires %s", commandPath, joinAlternativeFlagNames(flags))
			}
			return &cliError{
				exitCode: 2,
				code:     "cli_invalid_params",
				message:  message,
			}
		default:
			return &cliError{
				exitCode: 1,
				code:     "cli_schema_unsupported_rule",
				message:  fmt.Sprintf("unsupported schema rule %q", rule.Kind),
			}
		}
	}
	return nil
}

func joinFlagNames(flags []string) string {
	switch len(flags) {
	case 0:
		return ""
	case 1:
		return flags[0]
	case 2:
		return flags[0] + " and " + flags[1]
	default:
		return strings.Join(flags[:len(flags)-1], ", ") + ", and " + flags[len(flags)-1]
	}
}

func joinAlternativeFlagNames(flags []string) string {
	switch len(flags) {
	case 0:
		return ""
	case 1:
		return flags[0]
	case 2:
		return flags[0] + " or " + flags[1]
	default:
		return strings.Join(flags[:len(flags)-1], ", ") + ", or " + flags[len(flags)-1]
	}
}

func stripPersistentArgs(args []string) []string {
	out := make([]string, 0, len(args))
	for idx := 0; idx < len(args); idx++ {
		arg := strings.TrimSpace(args[idx])
		switch {
		case arg == "--base-url", arg == "--output":
			if idx+1 < len(args) {
				idx++
			}
		case strings.HasPrefix(arg, "--base-url="), strings.HasPrefix(arg, "--output="):
			continue
		case arg == "--verbose", arg == "--debug":
			continue
		default:
			out = append(out, args[idx])
		}
	}
	return out
}

func toInt(value any) (int, error) {
	if value == nil {
		return 0, nil
	}
	switch typed := value.(type) {
	case int:
		return typed, nil
	case int32:
		return int(typed), nil
	case int64:
		return int(typed), nil
	case float64:
		return int(typed), nil
	default:
		return 0, fmt.Errorf("unsupported int default %T", value)
	}
}

func toFloat(value any) (float64, error) {
	if value == nil {
		return 0, nil
	}
	switch typed := value.(type) {
	case float64:
		return typed, nil
	case float32:
		return float64(typed), nil
	case int:
		return float64(typed), nil
	case int64:
		return float64(typed), nil
	default:
		return 0, fmt.Errorf("unsupported float default %T", value)
	}
}

func toBool(value any) (bool, error) {
	if value == nil {
		return false, nil
	}
	typed, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf("unsupported bool default %T", value)
	}
	return typed, nil
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func derefInt(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}

func derefFloat(value *float64) float64 {
	if value == nil {
		return 0
	}
	return *value
}

func derefBool(value *bool) bool {
	if value == nil {
		return false
	}
	return *value
}

func outputDiagnostic(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format+"\n", args...)
}
