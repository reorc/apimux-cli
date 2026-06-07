package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Config struct {
	BaseURL string
	APIKey  string
}

type Response struct {
	StatusCode int
	Body       []byte
}

type CLILoginStartResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

type CLILoginPollResponse struct {
	APIKey   string `json:"api_key"`
	APIKeyID string `json:"api_key_id"`
	Status   string `json:"status"`
	Interval int    `json:"interval"`
	Error    string `json:"error"`
}

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func New(cfg Config) *Client {
	return &Client{
		baseURL: strings.TrimRight(cfg.BaseURL, "/"),
		apiKey:  cfg.APIKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func NewWithHTTPClient(cfg Config, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &Client{
		baseURL:    strings.TrimRight(cfg.BaseURL, "/"),
		apiKey:     cfg.APIKey,
		httpClient: httpClient,
	}
}

func (c *Client) ExecuteCapability(ctx context.Context, capability string, params map[string]any) (Response, error) {
	body, err := json.Marshal(map[string]any{"params": params})
	if err != nil {
		return Response{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/capabilities/"+capability, bytes.NewReader(body))
	if err != nil {
		return Response{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	applyCallerContextHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return Response{}, err
		}
		return Response{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
	if err != nil {
		return Response{}, err
	}

	return Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
	}, nil
}

// callerCtxEnvPrefix matches the apimux-side header prefix one-to-one.
// Any environment variable starting with this is forwarded as an
// X-Apimux-Caller-Ctx-* header to the capability request. The naming is
// deliberately generic — apimux does not interpret the keys, the
// configured webhook receiver decides what schema it expects.
const callerCtxEnvPrefix = "APIMUX_CALLER_CTX_"

// applyCallerContextHeaders scans the process environment for every
// variable prefixed with APIMUX_CALLER_CTX_ and copies its value into the
// matching X-Apimux-Caller-Ctx-* HTTP header on req. Empty / whitespace-only
// values are dropped — apimux-service treats absent and empty header
// values identically and persists SQL NULL on api_call_logs.caller_context
// in either case, so emitting blanks would only pollute the wire.
//
// Header naming: strip the env prefix, replace '_' with '-', and let Go's
// http.Header canonicalizer normalize the casing. apimux extracts the key
// back out by lowercasing and swapping '-' for '_', so:
//
//	APIMUX_CALLER_CTX_USER_ID   -> X-Apimux-Caller-Ctx-USER-ID
//	                            -> canonical X-Apimux-Caller-Ctx-User-Id
//	                            -> jsonb key "user_id"
func applyCallerContextHeaders(req *http.Request) {
	for _, kv := range os.Environ() {
		eq := strings.IndexByte(kv, '=')
		if eq <= 0 {
			continue
		}
		name, value := kv[:eq], kv[eq+1:]
		if !strings.HasPrefix(name, callerCtxEnvPrefix) {
			continue
		}
		v := strings.TrimSpace(value)
		if v == "" {
			continue
		}
		suffix := name[len(callerCtxEnvPrefix):]
		if suffix == "" {
			continue
		}
		req.Header.Set("X-Apimux-Caller-Ctx-"+strings.ReplaceAll(suffix, "_", "-"), v)
	}
}

func (c *Client) ListSchemas(ctx context.Context) (Response, error) {
	return c.doRequest(ctx, http.MethodGet, c.baseURL+"/v1/schema", nil)
}

func (c *Client) GetSchema(ctx context.Context, capability string) (Response, error) {
	return c.doRequest(ctx, http.MethodGet, c.baseURL+"/v1/schema/"+capability, nil)
}

func (c *Client) StartCLILogin(ctx context.Context, deviceName string) (CLILoginStartResponse, error) {
	body, err := json.Marshal(map[string]any{
		"deviceName": deviceName,
	})
	if err != nil {
		return CLILoginStartResponse{}, err
	}

	resp, err := c.doRequest(ctx, http.MethodPost, c.baseURL+"/api/cli-auth/start", body)
	if err != nil {
		return CLILoginStartResponse{}, err
	}
	if resp.StatusCode >= 400 {
		return CLILoginStartResponse{}, fmt.Errorf("start login failed with status %d", resp.StatusCode)
	}

	var payload CLILoginStartResponse
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return CLILoginStartResponse{}, err
	}
	return payload, nil
}

func (c *Client) PollCLILogin(ctx context.Context, deviceCode string) (CLILoginPollResponse, int, error) {
	body, err := json.Marshal(map[string]any{
		"device_code": deviceCode,
	})
	if err != nil {
		return CLILoginPollResponse{}, 0, err
	}

	resp, err := c.doRequest(ctx, http.MethodPost, c.baseURL+"/api/cli-auth/poll", body)
	if err != nil {
		return CLILoginPollResponse{}, 0, err
	}

	var payload CLILoginPollResponse
	if len(resp.Body) > 0 {
		if err := json.Unmarshal(resp.Body, &payload); err != nil {
			return CLILoginPollResponse{}, resp.StatusCode, err
		}
	}
	return payload, resp.StatusCode, nil
}

func (c *Client) doRequest(ctx context.Context, method, targetURL string, body []byte) (Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, targetURL, bytes.NewReader(body))
	if err != nil {
		return Response{}, err
	}
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return Response{}, err
		}
		return Response{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
	if err != nil {
		return Response{}, err
	}

	return Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
	}, nil
}
