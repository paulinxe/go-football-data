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

// HTTPError is an error that is returned when the HTTP response status is not 200 OK.
// It can be consumed as follows:
// var httpErr *football_data.HTTPError
//
//	if errors.As(err, &httpErr) {
//	    log.Printf("status: %d", httpErr.StatusCode)
//	}
type HTTPError struct {
	StatusCode int
	Body       []byte
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP error: %d %s", e.StatusCode, http.StatusText(e.StatusCode))
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
		return nil, fmt.Errorf("failed to create request: %w", err) // TODO: use a custom error?
	}

	req.Header.Set("X-Auth-Token", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, fmt.Errorf("request canceled: %w", err) // TODO: use a custom error?
		}

		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("context timeout: %w", err) // TODO: use a custom error?
		}

		return nil, fmt.Errorf("failed to get matches: %w", err) // TODO: use a custom error?
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err) // TODO: use a custom error
	}

	if resp.StatusCode != 200 {
		return nil, &HTTPError{
			StatusCode: resp.StatusCode,
			Body:       body,
		}
	}

	return body, nil
}
