package yfinance

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"
)

// ClientOption is a function that configures Client options
type ClientOption func(*Client)

// Client represents a Yahoo Finance API client with authentication
type Client struct {
	httpClient *http.Client
	userAgent  string
	crumb      string
	crumbMu    sync.RWMutex
	timeout    time.Duration
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithUserAgent sets a custom User-Agent header
func WithUserAgent(ua string) ClientOption {
	return func(c *Client) {
		c.userAgent = ua
	}
}

// WithTimeout sets the request timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// NewClient creates a new Yahoo Finance API client
func NewClient(opts ...ClientOption) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	client := &Client{
		httpClient: &http.Client{
			Jar:     jar,
			Timeout: 30 * time.Second,
		},
		userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		timeout:   30 * time.Second,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

// authenticate obtains cookies and crumb token for authenticated requests
func (c *Client) authenticate(ctx context.Context) error {
	// First, get cookies from fc.yahoo.com
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, CookieURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create cookie request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get cookies: %w", err)
	}
	defer resp.Body.Close()

	// Then, get the crumb
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, CrumbURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create crumb request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err = c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get crumb: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get crumb: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read crumb response: %w", err)
	}

	crumb := strings.TrimSpace(string(body))
	if crumb == "" {
		return ErrAuthentication
	}

	c.crumbMu.Lock()
	c.crumb = crumb
	c.crumbMu.Unlock()

	return nil
}

// ensureAuthenticated ensures the client has valid authentication
func (c *Client) ensureAuthenticated(ctx context.Context) error {
	c.crumbMu.RLock()
	crumb := c.crumb
	c.crumbMu.RUnlock()

	if crumb == "" {
		return c.authenticate(ctx)
	}
	return nil
}

// getCrumb returns the current crumb value
func (c *Client) getCrumb() string {
	c.crumbMu.RLock()
	defer c.crumbMu.RUnlock()
	return c.crumb
}

// Get performs a GET request to the specified URL
func (c *Client) Get(ctx context.Context, endpoint string, params url.Values) ([]byte, error) {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Add crumb to params
	if params == nil {
		params = url.Values{}
	}
	crumb := c.getCrumb()
	if crumb != "" {
		params.Set("crumb", crumb)
	}

	reqURL := endpoint
	if len(params) > 0 {
		reqURL = fmt.Sprintf("%s?%s", endpoint, params.Encode())
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, &RequestError{Endpoint: endpoint, Method: "GET", Err: err}
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &RequestError{Endpoint: endpoint, Method: "GET", Err: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &RequestError{Endpoint: endpoint, Method: "GET", Err: err}
	}

	// Handle error responses
	if resp.StatusCode == http.StatusUnauthorized {
		// Try to re-authenticate
		c.crumbMu.Lock()
		c.crumb = ""
		c.crumbMu.Unlock()
		return nil, ErrAuthentication
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrRateLimited
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if resp.StatusCode >= 400 {
		var apiErr APIError
		if json.Unmarshal(body, &apiErr) == nil && apiErr.Description != "" {
			apiErr.StatusCode = resp.StatusCode
			return nil, &apiErr
		}
		return nil, &APIError{StatusCode: resp.StatusCode, Description: string(body)}
	}

	return body, nil
}

// Post performs a POST request to the specified URL
func (c *Client) Post(ctx context.Context, endpoint string, params url.Values, body interface{}) ([]byte, error) {
	if err := c.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}

	// Add crumb to params
	if params == nil {
		params = url.Values{}
	}
	crumb := c.getCrumb()
	if crumb != "" {
		params.Set("crumb", crumb)
	}

	reqURL := endpoint
	if len(params) > 0 {
		reqURL = fmt.Sprintf("%s?%s", endpoint, params.Encode())
	}

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = strings.NewReader(string(jsonBody))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, reqBody)
	if err != nil {
		return nil, &RequestError{Endpoint: endpoint, Method: "POST", Err: err}
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &RequestError{Endpoint: endpoint, Method: "POST", Err: err}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &RequestError{Endpoint: endpoint, Method: "POST", Err: err}
	}

	// Handle error responses
	if resp.StatusCode == http.StatusUnauthorized {
		c.crumbMu.Lock()
		c.crumb = ""
		c.crumbMu.Unlock()
		return nil, ErrAuthentication
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrRateLimited
	}

	if resp.StatusCode >= 400 {
		var apiErr APIError
		if json.Unmarshal(respBody, &apiErr) == nil && apiErr.Description != "" {
			apiErr.StatusCode = resp.StatusCode
			return nil, &apiErr
		}
		return nil, &APIError{StatusCode: resp.StatusCode, Description: string(respBody)}
	}

	return respBody, nil
}

// defaultClient is a package-level default client
var (
	defaultClient     *Client
	defaultClientOnce sync.Once
	defaultClientErr  error
)

// getDefaultClient returns the default client, creating it if necessary
func getDefaultClient() (*Client, error) {
	defaultClientOnce.Do(func() {
		defaultClient, defaultClientErr = NewClient()
	})
	return defaultClient, defaultClientErr
}

// SetDefaultClient sets the package-level default client
func SetDefaultClient(client *Client) {
	defaultClient = client
	defaultClientErr = nil
}
