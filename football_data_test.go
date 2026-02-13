package football_data

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_new_client_uses_default_baseURL(t *testing.T) {
	apiKey := "test-key"
	client := New(apiKey)

	if client.apiKey != apiKey {
		t.Errorf("expected apiKey %s, got %s", apiKey, client.apiKey)
	}

	expectedBaseURL := "https://api.football-data.org/v4"
	if client.baseURL != expectedBaseURL {
		t.Errorf("expected default baseURL %s, got %s", expectedBaseURL, client.baseURL)
	}
}

func Test_with_base_url_allows_to_change_the_base_url(t *testing.T) {
	customURL := "https://api.test.com"
	client := New("key", WithBaseURL(customURL))

	if client.baseURL != customURL {
		t.Errorf("expected baseURL %s, got %s", customURL, client.baseURL)
	}
}

func Test_get_works_as_expected(t *testing.T) {
	tests := []struct {
		name           string
		status         int
		responseBody   string
		ctx            context.Context
		wantError      bool
		expectedError  string
		checkHTTPError bool
	}{
		{
			name:         "Success 200 OK",
			status:       http.StatusOK,
			responseBody: `{"status": "ok"}`,
			ctx:          context.Background(),
			wantError:    false,
		},
		{
			name:           "HTTP Error 403 Forbidden",
			status:         http.StatusForbidden,
			responseBody:   `{"message": "Your API token is invalid"}`,
			ctx:            context.Background(),
			wantError:      true,
			checkHTTPError: true,
		},
		{
			name:         "Context Canceled",
			status:       http.StatusOK,
			responseBody: `{}`,
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			wantError:     true,
			expectedError: "request canceled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("X-Auth-Token") != "test-api-key" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// Initialize client with the test server's URL
			client := New("test-api-key", WithBaseURL(server.URL))

			// Execute the get method
			body, err := client.get(tt.ctx, "/test-path")

			// Validation logic
			if (err != nil) != tt.wantError {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantError)
				return
			}

			if tt.checkHTTPError {
				var httpErr *HTTPError
				if !errors.As(err, &httpErr) {
					t.Errorf("expected error to be of type *HTTPError, got %T", err)
				} else if httpErr.StatusCode != tt.status {
					t.Errorf("expected status code %d, got %d", tt.status, httpErr.StatusCode)
				}
			}

			if !tt.wantError && string(body) != tt.responseBody {
				t.Errorf("expected body %s, got %s", tt.responseBody, string(body))
			}
		})
	}
}
