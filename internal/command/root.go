package command

import (
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
		r.newAuthCommand(runCtx),
		r.newConfigCommand(),
		r.newCompletionCommand(),
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
		Use:        "init",
		Short:      "Deprecated: use auth login or config set",
		Deprecated: "use `apimux auth login`; for CI/manual keys use `apimux config set --api-key ... --base-url ...`",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintln(r.stdout, "`apimux config init` is deprecated. Use `apimux auth login` for onboarding. For CI/manual keys, use `apimux config set --api-key ... --base-url ...`.")
			return err
		},
	}

	configCmd.AddCommand(setCmd, initCmd)

	return configCmd
}

func (r *Root) newCompletionCommand() *cobra.Command {
	completionCmd := &cobra.Command{
		Use:   "completion",
		Short: "Generate shell completion scripts",
		Long:  "Generate shell completion scripts for apimux.\n\nLoad the generated script from your shell profile to enable command and flag completion.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return &cliError{
				exitCode: 2,
				code:     "cli_invalid_command",
				message:  "completion supports: bash, zsh, fish, powershell",
			}
		},
	}

	completionCmd.AddCommand(
		&cobra.Command{
			Use:   "bash",
			Short: "Generate bash completion script",
			RunE: func(cmd *cobra.Command, args []string) error {
				return cmd.Root().GenBashCompletion(r.stdout)
			},
		},
		&cobra.Command{
			Use:   "zsh",
			Short: "Generate zsh completion script",
			RunE: func(cmd *cobra.Command, args []string) error {
				return cmd.Root().GenZshCompletion(r.stdout)
			},
		},
		&cobra.Command{
			Use:   "fish",
			Short: "Generate fish completion script",
			RunE: func(cmd *cobra.Command, args []string) error {
				return cmd.Root().GenFishCompletion(r.stdout, true)
			},
		},
		&cobra.Command{
			Use:   "powershell",
			Short: "Generate PowerShell completion script",
			RunE: func(cmd *cobra.Command, args []string) error {
				return cmd.Root().GenPowerShellCompletion(r.stdout)
			},
		},
	)

	return completionCmd
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
		Long:    "Browse available API endpoints and inspect their parameters.\n\nUse 'schema capabilities' to quickly discover available endpoint names, 'schema show <endpoint>' to inspect one endpoint before calling it, and 'schema list' when you need the full schema payload.",
		Example: "  apimux schema capabilities\n  apimux schema show reddit.search\n  apimux schema list",
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := runCtx.client.ListSchemas(cmd.Context())
			return writeServiceResponse(runCtx, resp, err)
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List the full schema payload",
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

	capabilitiesCmd := &cobra.Command{
		Use:   "capabilities",
		Short: "List capability names for quick discovery",
		Long:  "List available capability names without the full schema payload.\n\nUse this for fast discovery before drilling into one capability with 'schema show <capability>'.",
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := runCtx.client.ListSchemas(cmd.Context())
			asJSON, _ := cmd.Flags().GetBool("json")
			return writeSchemaNames(runCtx, resp, err, asJSON)
		},
	}
	capabilitiesCmd.Flags().Bool("json", false, "Output as JSON array instead of one name per line")

	schemaCmd.AddCommand(listCmd, showCmd, capabilitiesCmd)
	return schemaCmd
}

func writeSchemaNames(runCtx *runContext, resp client.Response, err error, asJSON bool) error {
	if err != nil {
		return writeServiceResponse(runCtx, resp, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return writeServiceResponse(runCtx, resp, nil)
	}

	names, ok := extractSchemaNames(resp.Body)
	if !ok {
		return writeServiceResponse(runCtx, resp, nil)
	}

	if asJSON {
		body, err := json.Marshal(names)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(runCtx.stdout, string(body))
		return err
	}

	for _, name := range names {
		if _, err := fmt.Fprintln(runCtx.stdout, name); err != nil {
			return err
		}
	}
	return nil
}

func extractSchemaNames(body []byte) ([]string, bool) {
	var listEnvelope struct {
		OK   bool `json:"ok"`
		Data struct {
			Capabilities []struct {
				Name string `json:"name"`
			} `json:"capabilities"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &listEnvelope); err == nil && listEnvelope.OK && len(listEnvelope.Data.Capabilities) > 0 {
		names := make([]string, 0, len(listEnvelope.Data.Capabilities))
		for _, capability := range listEnvelope.Data.Capabilities {
			names = append(names, capability.Name)
		}
		return names, true
	}

	var legacyEnvelope struct {
		OK   bool `json:"ok"`
		Data []struct {
			Name string `json:"name"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &legacyEnvelope); err == nil && legacyEnvelope.OK {
		names := make([]string, 0, len(legacyEnvelope.Data))
		for _, capability := range legacyEnvelope.Data {
			names = append(names, capability.Name)
		}
		return names, true
	}

	return nil, false
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
