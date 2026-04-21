package command

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/reorc/apimux-cli/internal/buildinfo"
	"github.com/reorc/apimux-cli/internal/client"
	"github.com/reorc/apimux-cli/internal/config"
	"github.com/reorc/apimux-cli/internal/output"
	"github.com/reorc/apimux-cli/internal/update"
	"github.com/spf13/cobra"
)

type Root struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

type runContext struct {
	stdout   io.Writer
	stderr   io.Writer
	verbose  bool
	debug    bool
	output   output.BodyOutput
	cfg      config.Config
	client   *client.Client
	exitCode int
}

type cliError struct {
	exitCode int
	code     string
	message  string
}

func (e *cliError) Error() string {
	return e.message
}

func NewRoot(stdout, stderr io.Writer) *Root {
	return NewRootWithIO(os.Stdin, stdout, stderr)
}

func NewRootWithIO(stdin io.Reader, stdout, stderr io.Writer) *Root {
	return &Root{stdin: stdin, stdout: stdout, stderr: stderr}
}

func (r *Root) Execute(ctx context.Context, args []string) (int, error) {
	cfg, err := config.Load()
	if err != nil {
		return 1, err
	}

	runCtx := &runContext{
		stdout: r.stdout,
		stderr: r.stderr,
		cfg:    cfg,
	}
	if err := applyPersistentArgs(runCtx, args); err != nil {
		var cliErr *cliError
		if errors.As(err, &cliErr) {
			_ = output.Renderer{Stdout: r.stdout, Stderr: r.stderr}.WriteLocalError(cliErr.message, cliErr.code)
			return cliErr.exitCode, nil
		}
		return 1, err
	}

	cmd := r.newCommand(runCtx)
	cmd.SetArgs(args)
	cmd.SetOut(r.stdout)
	cmd.SetErr(io.Discard)

	if err := cmd.ExecuteContext(ctx); err != nil {
		renderer := output.Renderer{Stdout: r.stdout, Stderr: r.stderr}

		var cliErr *cliError
		if errors.As(err, &cliErr) {
			_ = renderer.WriteLocalError(cliErr.message, cliErr.code)
			return cliErr.exitCode, nil
		}

		switch {
		case strings.Contains(err.Error(), "unknown command"):
			_ = renderer.WriteLocalError("unknown command", "cli_unknown_command")
			return 2, nil
		case strings.Contains(err.Error(), "unknown flag"), strings.Contains(err.Error(), "required flag"):
			_ = renderer.WriteLocalError(err.Error(), "cli_invalid_flags")
			return 2, nil
		default:
			return 1, err
		}
	}

	return runCtx.exitCode, nil
}

func (r *Root) newCommand(runCtx *runContext) *cobra.Command {
	baseURL := runCtx.cfg.BaseURL
	outputMode := string(runCtx.output)
	if outputMode == "" {
		outputMode = string(output.BodyOutputCompact)
	}
	verbose := runCtx.verbose
	debug := runCtx.debug

	rootCmd := &cobra.Command{
		Use:           "apimux",
		Short:         "APIMux data API client",
		Long:          "APIMux data API client.\n\nUse source commands for common tasks grouped by data source, use capability to call any endpoint by name, and use schema to inspect available endpoints and parameters before calling them.",
		Example:       "  apimux schema list\n  apimux schema show tiktok.search_videos\n  apimux capability call tiktok.search_videos --params-json '{\"keyword\":\"laptop\"}'\n  apimux tiktok search_videos --keyword laptop --count 10",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cfg := runCtx.cfg
			cfg.BaseURL = strings.TrimSpace(baseURL)
			if cfg.BaseURL == "" {
				cfg.BaseURL = "http://127.0.0.1:8081"
			}

			bodyOutput, ok := output.ParseBodyOutput(outputMode)
			if !ok {
				return &cliError{
					exitCode: 2,
					code:     "cli_invalid_output",
					message:  "output must be one of: compact, pretty, data, data-pretty",
				}
			}

			runCtx.verbose = runCtx.verbose || verbose
			runCtx.debug = runCtx.debug || debug
			runCtx.output = bodyOutput
			runCtx.cfg = cfg
			runCtx.client = client.New(client.Config{
				BaseURL: cfg.BaseURL,
				APIKey:  cfg.APIKey,
			})
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", runCtx.cfg.BaseURL, "APIMux service base URL")
	rootCmd.PersistentFlags().StringVar(&outputMode, "output", outputMode, "Agent-facing body output: compact, pretty, data, or data-pretty")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", verbose, "Write request diagnostics to stderr")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", debug, "Emit full response envelope for debugging")
	rootCmd.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return &cliError{exitCode: 2, code: "cli_invalid_flags", message: err.Error()}
	})

	rootCmd.AddCommand(
		r.newVersionCommand(),
		r.newConfigCommand(),
		r.newUpgradeCommand(),
		r.newSchemaCommand(runCtx),
		r.newCapabilityCommand(runCtx),
		r.newAmazonCommand(runCtx),
		r.newDouyinCommand(runCtx),
		r.newGoogleAdsCommand(runCtx),
		r.newGoogleTrendsCommand(runCtx),
		r.newMetaAdsCommand(runCtx),
		r.newRedditCommand(runCtx),
		r.newTiktokCommand(runCtx),
		r.newTrendCloudCommand(runCtx),
		r.newXiaohongshuCommand(runCtx),
	)

	return rootCmd
}

func (r *Root) newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show CLI version",
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := json.MarshalIndent(buildinfo.Current(), "", "  ")
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(r.stdout, string(body))
			return err
		},
	}
}

func (r *Root) newConfigCommand() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Show or persist local CLI config",
		RunE: func(cmd *cobra.Command, args []string) error {
			loaded, err := config.LoadDetailed()
			if err != nil {
				return err
			}
			body, err := json.MarshalIndent(map[string]any{
				"path":     loaded.Path,
				"base_url": loaded.Config.BaseURL,
				"api_key":  redact(loaded.Config.APIKey),
				"sources":  loaded.Sources,
			}, "", "  ")
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(r.stdout, string(body))
			return err
		},
	}

	var baseURL string
	var apiKey string
	setCmd := &cobra.Command{
		Use:   "set",
		Short: "Persist CLI config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(baseURL) == "" && strings.TrimSpace(apiKey) == "" {
				return &cliError{
					exitCode: 2,
					code:     "cli_invalid_config",
					message:  "config set requires --base-url and/or --api-key",
				}
			}
			if err := config.Save(config.Config{BaseURL: baseURL, APIKey: apiKey}); err != nil {
				return err
			}
			_, err := fmt.Fprintln(r.stdout, `{"ok":true}`)
			return err
		},
	}
	setCmd.Flags().StringVar(&baseURL, "base-url", "", "Persist the APIMux service base URL")
	setCmd.Flags().StringVar(&apiKey, "api-key", "", "Persist the APIMux API key")

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Set up CLI credentials and server URL",
		RunE: func(cmd *cobra.Command, args []string) error {
			loaded, err := config.LoadDetailed()
			if err != nil {
				return err
			}

			reader := bufio.NewReader(r.stdin)
			baseURL, err := promptValue(reader, r.stdout, "APIMux service base URL", loaded.Config.BaseURL, false)
			if err != nil {
				return err
			}
			apiKey, err := promptValue(reader, r.stdout, "APIMux API key", "", true)
			if err != nil {
				return err
			}
			if strings.TrimSpace(apiKey) == "" {
				return &cliError{
					exitCode: 2,
					code:     "cli_invalid_config",
					message:  "interactive config requires a non-empty API key",
				}
			}
			if err := config.Save(config.Config{BaseURL: baseURL, APIKey: apiKey}); err != nil {
				return err
			}

			updated, err := config.LoadDetailed()
			if err != nil {
				return err
			}
			body, err := json.MarshalIndent(map[string]any{
				"ok":       true,
				"path":     updated.Path,
				"base_url": updated.Config.BaseURL,
				"api_key":  redact(updated.Config.APIKey),
				"sources":  updated.Sources,
			}, "", "  ")
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(r.stdout, string(body))
			return err
		},
	}

	configCmd.AddCommand(setCmd, initCmd)

	return configCmd
}

func (r *Root) newUpgradeCommand() *cobra.Command {
	var check bool
	upgradeCmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Inspect or run CLI updates",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !check {
				return &cliError{
					exitCode: 2,
					code:     "cli_upgrade_not_implemented",
					message:  "upgrade currently supports --check only",
				}
			}

			manifestURL := releaseManifestURL()
			result, err := update.Check(cmd.Context(), &http.Client{}, buildinfo.Current().Version, manifestURL)
			if err != nil {
				return &cliError{
					exitCode: 1,
					code:     "cli_upgrade_check_failed",
					message:  err.Error(),
				}
			}

			body, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(r.stdout, string(body))
			return err
		},
	}
	upgradeCmd.Flags().BoolVar(&check, "check", false, "Check whether a newer CLI version is available")
	return upgradeCmd
}

func (r *Root) newSchemaCommand(runCtx *runContext) *cobra.Command {
	schemaCmd := &cobra.Command{
		Use:     "schema",
		Short:   "Browse available API endpoints",
		Long:    "Browse available API endpoints and inspect their parameters.\n\nUse 'schema list' to see available endpoints and 'schema show <endpoint>' to inspect the parameters and options for one endpoint before calling it.",
		Example: "  apimux schema list\n  apimux schema show reddit.search\n  apimux schema show google_ads.list_ad_creatives",
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := runCtx.client.ListSchemas(cmd.Context())
			return writeServiceResponse(runCtx, resp, err)
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all available endpoints",
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := runCtx.client.ListSchemas(cmd.Context())
			return writeServiceResponse(runCtx, resp, err)
		},
	}

	showCmd := &cobra.Command{
		Use:   "show <capability>",
		Short: "Show endpoint parameters and options",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := runCtx.client.GetSchema(cmd.Context(), args[0])
			return writeServiceResponse(runCtx, resp, err)
		},
	}

	schemaCmd.AddCommand(listCmd, showCmd)
	return schemaCmd
}

func (r *Root) newCapabilityCommand(runCtx *runContext) *cobra.Command {
	capabilityCmd := &cobra.Command{
		Use:     "capability",
		Short:   "Call any endpoint by name",
		Long:    "Call any endpoint by name with a JSON parameter object.\n\nUse this command when you already know the endpoint name. For source-grouped shortcuts, use commands such as amazon, reddit, or tiktok. Use schema to discover endpoint names and parameters first.",
		Example: "  apimux capability call reddit.search --params-json '{\"query\":\"openai\",\"sort\":\"hot\"}'\n  apimux capability call amazon.search_products --params-json '{\"keyword\":\"desk lamp\"}'",
		RunE: func(cmd *cobra.Command, args []string) error {
			return &cliError{
				exitCode: 2,
				code:     "cli_invalid_command",
				message:  "capability supports: call <capability> --params-json '{...}'",
			}
		},
	}

	var paramsJSON string
	callCmd := &cobra.Command{
		Use:   "call <capability>",
		Short: "Call an endpoint with JSON parameters",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var params map[string]any
			if err := json.Unmarshal([]byte(paramsJSON), &params); err != nil {
				return &cliError{
					exitCode: 2,
					code:     "cli_invalid_params",
					message:  "params-json must be a valid JSON object",
				}
			}
			return callCapability(cmd.Context(), runCtx, args[0], params)
		},
	}
	callCmd.Flags().StringVar(&paramsJSON, "params-json", "{}", "Capability params as JSON object")
	capabilityCmd.AddCommand(callCmd)

	return capabilityCmd
}
