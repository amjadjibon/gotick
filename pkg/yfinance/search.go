package yfinance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// SearchOption is a function that configures search options
type SearchOption func(*searchConfig)

type searchConfig struct {
	QuotesCount int
	NewsCount   int
	Region      string
	Lang        string
}

// WithQuotesCount sets the number of quotes to return
func WithQuotesCount(count int) SearchOption {
	return func(c *searchConfig) {
		c.QuotesCount = count
	}
}

// WithNewsCount sets the number of news to return
func WithNewsCount(count int) SearchOption {
	return func(c *searchConfig) {
		c.NewsCount = count
	}
}

// WithRegion sets the region for search
func WithRegion(region string) SearchOption {
	return func(c *searchConfig) {
		c.Region = region
	}
}

// WithLang sets the language for search
func WithLang(lang string) SearchOption {
	return func(c *searchConfig) {
		c.Lang = lang
	}
}

// Search searches for symbols and companies matching the query
func Search(ctx context.Context, query string, opts ...SearchOption) (*SearchResult, error) {
	if query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	client, err := getDefaultClient()
	if err != nil {
		return nil, err
	}

	return SearchWithClient(ctx, client, query, opts...)
}

// SearchWithClient searches using a specific client
func SearchWithClient(ctx context.Context, client *Client, query string, opts ...SearchOption) (*SearchResult, error) {
	config := &searchConfig{
		QuotesCount: 10,
		NewsCount:   0,
		Region:      "US",
		Lang:        "en",
	}

	for _, opt := range opts {
		opt(config)
	}

	params := url.Values{}
	params.Set("q", query)
	params.Set("quotesCount", strconv.Itoa(config.QuotesCount))
	params.Set("newsCount", strconv.Itoa(config.NewsCount))
	params.Set("region", config.Region)
	params.Set("lang", config.Lang)

	data, err := client.Get(ctx, SearchURL, params)
	if err != nil {
		return nil, err
	}

	var response struct {
		Quotes []SearchQuote `json:"quotes"`
		News   []NewsItem    `json:"news"`
		Count  int           `json:"count"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
	}

	return &SearchResult{
		Query:  query,
		Quotes: response.Quotes,
		News:   response.News,
		Count:  len(response.Quotes),
	}, nil
}

// Lookup performs a symbol lookup
func Lookup(ctx context.Context, query string, lookupType string) (*LookupResult, error) {
	if query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	client, err := getDefaultClient()
	if err != nil {
		return nil, err
	}

	return LookupWithClient(ctx, client, query, lookupType)
}

// LookupWithClient performs lookup using a specific client
func LookupWithClient(ctx context.Context, client *Client, query string, lookupType string) (*LookupResult, error) {
	params := url.Values{}
	params.Set("query", query)
	if lookupType != "" {
		params.Set("type", lookupType)
	}

	data, err := client.Get(ctx, LookupURL, params)
	if err != nil {
		return nil, err
	}

	var response struct {
		Finance struct {
			Result []struct {
				Documents []LookupItem `json:"documents"`
			} `json:"result"`
		} `json:"finance"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse lookup response: %w", err)
	}

	result := &LookupResult{
		Query: query,
		Items: []LookupItem{},
	}

	if len(response.Finance.Result) > 0 {
		result.Items = response.Finance.Result[0].Documents
		result.Count = len(result.Items)
	}

	return result, nil
}

// QuoteMultiple fetches quotes for multiple symbols at once
func QuoteMultiple(ctx context.Context, symbols []string) ([]Quote, error) {
	if len(symbols) == 0 {
		return nil, fmt.Errorf("symbols cannot be empty")
	}

	client, err := getDefaultClient()
	if err != nil {
		return nil, err
	}

	return QuoteMultipleWithClient(ctx, client, symbols)
}

// QuoteMultipleWithClient fetches multiple quotes using a specific client
func QuoteMultipleWithClient(ctx context.Context, client *Client, symbols []string) ([]Quote, error) {
	params := url.Values{}
	params.Set("symbols", joinSymbols(symbols))

	data, err := client.Get(ctx, QuoteURL, params)
	if err != nil {
		return nil, err
	}

	var response struct {
		QuoteResponse struct {
			Result []Quote `json:"result"`
			Error  *struct {
				Code        string `json:"code"`
				Description string `json:"description"`
			} `json:"error"`
		} `json:"quoteResponse"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse quote response: %w", err)
	}

	if response.QuoteResponse.Error != nil {
		return nil, &APIError{
			Code:        response.QuoteResponse.Error.Code,
			Description: response.QuoteResponse.Error.Description,
		}
	}

	return response.QuoteResponse.Result, nil
}

func joinSymbols(symbols []string) string {
	result := ""
	for i, s := range symbols {
		if i > 0 {
			result += ","
		}
		result += s
	}
	return result
}
