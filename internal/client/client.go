package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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
