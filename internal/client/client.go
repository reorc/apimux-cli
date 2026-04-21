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
