package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/vecherochek/evo-ai-agents-skills-cli/internal/auth"
)

type Client struct {
	baseURL    string
	timeout    time.Duration
	httpClient *http.Client
	auth       auth.IAMAuthServiceInterface
}

func NewClient(baseURL string, timeoutSec int, authService auth.IAMAuthServiceInterface) (*Client, error) {
	base := strings.TrimSpace(baseURL)
	if base == "" {
		return nil, fmt.Errorf("PUBLIC_BFF_URL is required")
	}

	timeout := time.Duration(timeoutSec) * time.Second
	if timeout <= 0 {
		timeout = 60 * time.Second
	}

	return &Client{
		baseURL:    strings.TrimRight(base, "/"),
		timeout:    timeout,
		httpClient: &http.Client{Timeout: timeout},
		auth:       authService,
	}, nil
}

func (c *Client) Get(path string, authHeader string) ([]byte, int, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("create GET request: %w", err)
	}
	finalAuthHeader := strings.TrimSpace(authHeader)
	if finalAuthHeader == "" && c.auth != nil {
		token, tokenErr := c.auth.GetToken(context.Background())
		if tokenErr != nil {
			return nil, 0, fmt.Errorf("get IAM token: %w", tokenErr)
		}
		finalAuthHeader = "Bearer " + token
	}
	if finalAuthHeader != "" {
		req.Header.Set("Authorization", finalAuthHeader)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("perform GET request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read GET response body: %w", err)
	}

	return body, resp.StatusCode, nil
}
