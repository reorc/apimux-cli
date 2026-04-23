package command

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/reorc/apimux-cli/internal/config"
	"github.com/spf13/cobra"
)

func (r *Root) newAuthCommand(runCtx *runContext) *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate the CLI",
	}

	var noBrowser bool
	var deviceName string
	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Start browser-assisted CLI login",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if strings.TrimSpace(deviceName) == "" {
				host, _ := os.Hostname()
				deviceName = host
			}

			start, err := runCtx.client.StartCLILogin(ctx, deviceName)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(r.stdout, "Visit: %s\n", start.VerificationURIComplete)
			_, _ = fmt.Fprintf(r.stdout, "Code: %s\n", start.UserCode)
			if noBrowser {
				_, _ = fmt.Fprintln(r.stdout, "Open the URL above in your browser to approve access.")
			} else if err := openBrowser(start.VerificationURIComplete); err != nil {
				_, _ = fmt.Fprintf(r.stderr, "open browser failed: %v\n", err)
			}
			_, _ = fmt.Fprintln(r.stdout, "Waiting for authorization...")

			interval := time.Duration(start.Interval) * time.Second
			if interval <= 0 {
				interval = 5 * time.Second
			}
			deadline := time.Now().Add(time.Duration(start.ExpiresIn) * time.Second)

			for {
				if time.Now().After(deadline) {
					return &cliError{
						exitCode: 1,
						code:     "cli_auth_expired",
						message:  "login session expired before approval",
					}
				}

				pollCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
				result, statusCode, err := runCtx.client.PollCLILogin(pollCtx, start.DeviceCode)
				cancel()
				if err != nil {
					return err
				}

				switch statusCode {
				case 200:
					if strings.TrimSpace(result.APIKey) == "" {
						return &cliError{
							exitCode: 1,
							code:     "cli_auth_missing_key",
							message:  "login completed without an API key",
						}
					}
					if err := config.Save(config.Config{BaseURL: runCtx.cfg.BaseURL, APIKey: result.APIKey}); err != nil {
						return err
					}
					_, _ = fmt.Fprintln(r.stdout, "Authorized. API key saved to local config.")
					return nil
				case 202:
					time.Sleep(interval)
					continue
				case 403:
					return &cliError{
						exitCode: 1,
						code:     "cli_auth_denied",
						message:  "login request was denied",
					}
				case 404:
					return &cliError{
						exitCode: 1,
						code:     "cli_auth_invalid_device_code",
						message:  "login request is invalid",
					}
				case 409:
					return &cliError{
						exitCode: 1,
						code:     "cli_auth_consumed",
						message:  "login request has already been consumed",
					}
				case 410:
					return &cliError{
						exitCode: 1,
						code:     "cli_auth_expired",
						message:  "login session expired before approval",
					}
				default:
					return &cliError{
						exitCode: 1,
						code:     "cli_auth_failed",
						message:  "login polling failed",
					}
				}
			}
		},
	}
	loginCmd.Flags().BoolVar(&noBrowser, "no-browser", false, "Print the URL but do not open a browser")
	loginCmd.Flags().StringVar(&deviceName, "device-name", "", "Label used when creating the CLI API key")

	authCmd.AddCommand(loginCmd)
	return authCmd
}

func openBrowser(target string) error {
	var command string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		command = "open"
		args = []string{target}
	case "linux":
		command = "xdg-open"
		args = []string{target}
	case "windows":
		command = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", target}
	default:
		return fmt.Errorf("unsupported platform for browser launch: %s", runtime.GOOS)
	}

	return exec.Command(command, args...).Start()
}
