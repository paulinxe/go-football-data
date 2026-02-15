package football_data

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Option is a function that can be used to modify a target (Functional Options pattern).
type Option[T any] func(*T)
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func New(apiKey string, options ...Option[Client]) *Client {
	client := &Client{
		baseURL:    "https://api.football-data.org/v4",
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}

	for _, option := range options {
		option(client)
	}

	return client
}

// WithBaseURL allows you to change the baseURL to a custom one. It can be useful when testing.
func WithBaseURL(baseURL string) Option[Client] {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func (c *Client) get(ctx context.Context, path string, queryParams *url.Values) ([]byte, error) {
	fullURL := c.baseURL + path
	if queryParams != nil {
		fullURL += "?" + queryParams.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Auth-Token", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, fmt.Errorf("request canceled: %w", err)
		}

		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("context timeout: %w", err)
		}

		return nil, fmt.Errorf("failed to get matches: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, &HTTPError{
			StatusCode: resp.StatusCode,
			Body:       body,
		}
	}

	return body, nil
}
