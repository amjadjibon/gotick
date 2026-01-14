package yfinance

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

// RetryConfig configures retry behavior
type RetryConfig struct {
	MaxRetries     int           // Maximum number of retries
	InitialBackoff time.Duration // Initial backoff duration
	MaxBackoff     time.Duration // Maximum backoff duration
	BackoffFactor  float64       // Backoff multiplier (e.g., 2.0 for exponential)
	Jitter         float64       // Random jitter factor (0-1)
	RetryOnStatus  []int         // HTTP status codes to retry on
}

// DefaultRetryConfig returns sensible defaults for retry
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 1 * time.Second,
		MaxBackoff:     30 * time.Second,
		BackoffFactor:  2.0,
		Jitter:         0.1,
		RetryOnStatus:  []int{429, 500, 502, 503, 504},
	}
}

// ProxyConfig configures proxy settings
type ProxyConfig struct {
	URL      string // Proxy URL (e.g., "http://proxy:8080")
	Username string // Optional username
	Password string // Optional password
}

// WithRetry configures retry behavior for the client
func WithRetry(config RetryConfig) ClientOption {
	return func(c *Client) {
		c.retryConfig = &config
	}
}

// WithProxy configures proxy settings for the client
func WithProxy(config ProxyConfig) ClientOption {
	return func(c *Client) {
		c.proxyConfig = &config
	}
}

// WithProxyURL is a convenience function to set just the proxy URL
func WithProxyURL(proxyURL string) ClientOption {
	return func(c *Client) {
		c.proxyConfig = &ProxyConfig{URL: proxyURL}
	}
}

// doWithRetry executes a request with retry logic
func (c *Client) doWithRetry(ctx context.Context, req *http.Request) (*http.Response, error) {
	config := c.retryConfig
	if config == nil {
		// Use defaults if not configured
		defaultConfig := DefaultRetryConfig()
		config = &defaultConfig
	}

	var lastErr error
	backoff := config.InitialBackoff

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		// Clone request for retry
		reqClone := req.Clone(ctx)

		resp, err := c.httpClient.Do(reqClone)
		if err != nil {
			lastErr = err
			if attempt < config.MaxRetries {
				waitTime := calculateBackoff(backoff, config.MaxBackoff, config.Jitter)
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(waitTime):
				}
				backoff = time.Duration(float64(backoff) * config.BackoffFactor)
				continue
			}
			return nil, fmt.Errorf("request failed after %d retries: %w", config.MaxRetries, lastErr)
		}

		// Check if we should retry based on status code
		if shouldRetry(resp.StatusCode, config.RetryOnStatus) && attempt < config.MaxRetries {
			resp.Body.Close()
			waitTime := calculateBackoff(backoff, config.MaxBackoff, config.Jitter)

			// Check for Retry-After header
			if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
				if duration, err := time.ParseDuration(retryAfter + "s"); err == nil {
					waitTime = duration
				}
			}

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(waitTime):
			}
			backoff = time.Duration(float64(backoff) * config.BackoffFactor)
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", config.MaxRetries, lastErr)
}

// calculateBackoff calculates the backoff duration with jitter
func calculateBackoff(base, max time.Duration, jitter float64) time.Duration {
	backoff := base
	if backoff > max {
		backoff = max
	}

	if jitter > 0 {
		// Add random jitter
		jitterRange := float64(backoff) * jitter
		backoff = time.Duration(float64(backoff) + (rand.Float64()*2-1)*jitterRange)
	}

	return backoff
}

// shouldRetry checks if the status code is in the retry list
func shouldRetry(statusCode int, retryOnStatus []int) bool {
	for _, code := range retryOnStatus {
		if statusCode == code {
			return true
		}
	}
	return false
}

// configureProxy applies proxy configuration to the HTTP client
func (c *Client) configureProxy() {
	if c.proxyConfig == nil || c.proxyConfig.URL == "" {
		return
	}

	proxyURL, err := url.Parse(c.proxyConfig.URL)
	if err != nil {
		return
	}

	// Add authentication if provided
	if c.proxyConfig.Username != "" {
		proxyURL.User = url.UserPassword(c.proxyConfig.Username, c.proxyConfig.Password)
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	c.httpClient.Transport = transport
}

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	tokens         float64
	maxTokens      float64
	refillRate     float64 // tokens per second
	lastRefillTime time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerSecond float64, burst int) *RateLimiter {
	return &RateLimiter{
		tokens:         float64(burst),
		maxTokens:      float64(burst),
		refillRate:     requestsPerSecond,
		lastRefillTime: time.Now(),
	}
}

// Wait blocks until a token is available
func (rl *RateLimiter) Wait(ctx context.Context) error {
	for {
		rl.refill()
		if rl.tokens >= 1 {
			rl.tokens--
			return nil
		}

		waitTime := time.Duration((1 - rl.tokens) / rl.refillRate * float64(time.Second))
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitTime):
		}
	}
}

// refill adds tokens based on elapsed time
func (rl *RateLimiter) refill() {
	now := time.Now()
	elapsed := now.Sub(rl.lastRefillTime).Seconds()
	rl.tokens = math.Min(rl.maxTokens, rl.tokens+elapsed*rl.refillRate)
	rl.lastRefillTime = now
}

// WithRateLimiter configures rate limiting for the client
func WithRateLimiter(requestsPerSecond float64, burst int) ClientOption {
	return func(c *Client) {
		c.rateLimiter = NewRateLimiter(requestsPerSecond, burst)
	}
}
