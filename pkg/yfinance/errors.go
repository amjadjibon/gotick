package yfinance

import (
	"errors"
	"fmt"
)

// Sentinel errors for common error conditions
var (
	// ErrInvalidSymbol is returned when an invalid ticker symbol is provided
	ErrInvalidSymbol = errors.New("yfinance: invalid symbol")

	// ErrNotFound is returned when requested data is not found
	ErrNotFound = errors.New("yfinance: data not found")

	// ErrRateLimited is returned when rate limit is exceeded
	ErrRateLimited = errors.New("yfinance: rate limit exceeded")

	// ErrAuthentication is returned when authentication fails
	ErrAuthentication = errors.New("yfinance: authentication failed")

	// ErrNetwork is returned for network-related errors
	ErrNetwork = errors.New("yfinance: network error")

	// ErrInvalidResponse is returned when API returns invalid response
	ErrInvalidResponse = errors.New("yfinance: invalid response")

	// ErrNoData is returned when no data is available for the request
	ErrNoData = errors.New("yfinance: no data available")

	// ErrInvalidInterval is returned when an invalid interval is specified
	ErrInvalidInterval = errors.New("yfinance: invalid interval")

	// ErrInvalidPeriod is returned when an invalid period is specified
	ErrInvalidPeriod = errors.New("yfinance: invalid period")

	// ErrWebSocketClosed is returned when WebSocket connection is closed
	ErrWebSocketClosed = errors.New("yfinance: websocket connection closed")
)

// APIError represents an error returned by the Yahoo Finance API
type APIError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	StatusCode  int    `json:"-"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("yfinance: API error [%s]: %s", e.Code, e.Description)
	}
	return fmt.Sprintf("yfinance: API error (status %d): %s", e.StatusCode, e.Description)
}

// RequestError wraps an error with request context
type RequestError struct {
	Endpoint string
	Method   string
	Err      error
}

// Error implements the error interface
func (e *RequestError) Error() string {
	return fmt.Sprintf("yfinance: %s %s: %v", e.Method, e.Endpoint, e.Err)
}

// Unwrap returns the underlying error
func (e *RequestError) Unwrap() error {
	return e.Err
}

// SymbolError represents an error for a specific symbol
type SymbolError struct {
	Symbol string
	Err    error
}

// Error implements the error interface
func (e *SymbolError) Error() string {
	return fmt.Sprintf("yfinance: symbol %s: %v", e.Symbol, e.Err)
}

// Unwrap returns the underlying error
func (e *SymbolError) Unwrap() error {
	return e.Err
}

// NewSymbolError creates a new SymbolError
func NewSymbolError(symbol string, err error) *SymbolError {
	return &SymbolError{Symbol: symbol, Err: err}
}

// IsNotFound checks if the error is a not found error
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsRateLimited checks if the error is a rate limit error
func IsRateLimited(err error) bool {
	return errors.Is(err, ErrRateLimited)
}

// IsAuthError checks if the error is an authentication error
func IsAuthError(err error) bool {
	return errors.Is(err, ErrAuthentication)
}

// IsNetworkError checks if the error is a network error
func IsNetworkError(err error) bool {
	return errors.Is(err, ErrNetwork)
}
