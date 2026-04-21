package command

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/reorc/apimux-cli/internal/buildinfo"
	"github.com/reorc/apimux-cli/internal/client"
	"github.com/reorc/apimux-cli/internal/output"
)

func promptValue(reader *bufio.Reader, stdout io.Writer, label string, fallback string, secret bool) (string, error) {
	if _, err := fmt.Fprintf(stdout, "%s", label); err != nil {
		return "", err
	}
	if strings.TrimSpace(fallback) != "" {
		if _, err := fmt.Fprintf(stdout, " [%s]", fallback); err != nil {
			return "", err
		}
	}
	if _, err := fmt.Fprint(stdout, ": "); err != nil {
		return "", err
	}

	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}

	value := strings.TrimSpace(line)
	if value == "" {
		value = strings.TrimSpace(fallback)
	}
	if secret {
		return value, nil
	}
	return value, nil
}

func releaseManifestURL() string {
	if value := strings.TrimSpace(os.Getenv("APIMUX_CLI_RELEASE_MANIFEST_URL")); value != "" {
		return value
	}
	return buildinfo.Current().ReleaseManifestURL
}

func callCapability(ctx context.Context, runCtx *runContext, capability string, params map[string]any) error {
	if runCtx.verbose {
		output.Renderer{Stdout: runCtx.stdout, Stderr: runCtx.stderr}.Diagnostic("[apimux] POST /v1/capabilities/%s", capability)
	}

	resp, err := runCtx.client.ExecuteCapability(ctx, capability, params)
	return writeServiceResponse(runCtx, resp, err)
}

func writeServiceResponse(runCtx *runContext, resp client.Response, err error) error {
	renderer := output.Renderer{Stdout: runCtx.stdout, Stderr: runCtx.stderr}
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return &cliError{
				exitCode: 3,
				code:     "cli_timeout",
				message:  "request timed out",
			}
		}
		return &cliError{
			exitCode: 1,
			code:     "cli_transport_error",
			message:  err.Error(),
		}
	}

	if runCtx.verbose {
		renderer.Diagnostic("[apimux] HTTP %d", resp.StatusCode)
	}
	if err := renderer.WriteCapabilityResponse(resp.Body, runCtx.output, runCtx.debug); err != nil {
		var invalid *output.InvalidProjectionError
		if errors.As(err, &invalid) {
			return &cliError{
				exitCode: 1,
				code:     "cli_projection_failed",
				message:  err.Error(),
			}
		}
		return err
	}
	runCtx.exitCode = exitCodeForHTTPStatus(resp.StatusCode)
	return nil
}

func parseObjectFlag(value, flagName string) (map[string]any, error) {
	var payload map[string]any
	if err := json.Unmarshal([]byte(value), &payload); err != nil {
		return nil, &cliError{
			exitCode: 2,
			code:     "cli_invalid_params",
			message:  fmt.Sprintf("%s must be a valid JSON object", flagName),
		}
	}
	return payload, nil
}

func exitCodeForHTTPStatus(status int) int {
	switch {
	case status >= 200 && status < 300:
		return 0
	case status == 400:
		return 2
	case status == 401:
		return 3
	case status == 402:
		return 4
	case status == 404:
		return 5
	case status == 504:
		return 6
	default:
		return 1
	}
}

func redact(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if len(value) <= 8 {
		return "********"
	}
	return value[:4] + "..." + value[len(value)-4:]
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func applyPersistentArgs(runCtx *runContext, args []string) error {
	baseURL := runCtx.cfg.BaseURL
	bodyOutput := string(output.BodyOutputCompact)
	verbose := false
	debug := false

	for idx := 0; idx < len(args); idx++ {
		arg := strings.TrimSpace(args[idx])
		switch {
		case arg == "--base-url":
			if idx+1 >= len(args) {
				return &cliError{exitCode: 2, code: "cli_invalid_flags", message: "flag needs an argument: --base-url"}
			}
			idx++
			baseURL = strings.TrimSpace(args[idx])
		case strings.HasPrefix(arg, "--base-url="):
			baseURL = strings.TrimSpace(strings.TrimPrefix(arg, "--base-url="))
		case arg == "--output":
			if idx+1 >= len(args) {
				return &cliError{exitCode: 2, code: "cli_invalid_flags", message: "flag needs an argument: --output"}
			}
			idx++
			bodyOutput = strings.TrimSpace(args[idx])
		case strings.HasPrefix(arg, "--output="):
			bodyOutput = strings.TrimSpace(strings.TrimPrefix(arg, "--output="))
		case arg == "--verbose":
			verbose = true
		case arg == "--debug":
			debug = true
		}
	}

	if strings.TrimSpace(baseURL) == "" {
		baseURL = "http://127.0.0.1:8081"
	}
	parsedOutput, ok := output.ParseBodyOutput(bodyOutput)
	if !ok {
		return &cliError{
			exitCode: 2,
			code:     "cli_invalid_output",
			message:  "output must be one of: compact, pretty, data, data-pretty",
		}
	}

	runCtx.verbose = verbose
	runCtx.debug = debug
	runCtx.output = parsedOutput
	runCtx.cfg.BaseURL = baseURL
	runCtx.client = client.New(client.Config{
		BaseURL: baseURL,
		APIKey:  runCtx.cfg.APIKey,
	})
	return nil
}
