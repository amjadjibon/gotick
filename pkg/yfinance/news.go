package yfinance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// GetNews fetches financial news for given symbols
func GetNews(ctx context.Context, symbols []string, count int) ([]NewsItem, error) {
	client, err := getDefaultClient()
	if err != nil {
		return nil, err
	}
	return GetNewsWithClient(ctx, client, symbols, count)
}

// GetNewsWithClient fetches news using a specific client
func GetNewsWithClient(ctx context.Context, client *Client, symbols []string, count int) ([]NewsItem, error) {
	if count <= 0 {
		count = 10
	}

	params := url.Values{}
	if len(symbols) > 0 {
		params.Set("q", joinSymbols(symbols))
	}
	params.Set("newsCount", strconv.Itoa(count))
	params.Set("quotesCount", "0")

	data, err := client.Get(ctx, SearchURL, params)
	if err != nil {
		return nil, err
	}

	var response struct {
		News []NewsItem `json:"news"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse news response: %w", err)
	}

	return response.News, nil
}

// GetLatestNews fetches the latest financial news
func GetLatestNews(ctx context.Context, count int) ([]NewsItem, error) {
	return GetNews(ctx, nil, count)
}

// GetSymbolNews fetches news for a specific symbol
func GetSymbolNews(ctx context.Context, symbol string, count int) ([]NewsItem, error) {
	return GetNews(ctx, []string{symbol}, count)
}
