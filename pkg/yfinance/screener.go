package yfinance

import (
	"context"
	"encoding/json"
	"fmt"
)

// Screen performs stock screening based on criteria
func Screen(ctx context.Context, criteria ScreenCriteria) (*ScreenResult, error) {
	client, err := getDefaultClient()
	if err != nil {
		return nil, err
	}

	return ScreenWithClient(ctx, client, criteria)
}

// ScreenWithClient performs screening using a specific client
func ScreenWithClient(ctx context.Context, client *Client, criteria ScreenCriteria) (*ScreenResult, error) {
	// Set defaults
	if criteria.Size == 0 {
		criteria.Size = 25
	}
	if criteria.Region == "" {
		criteria.Region = "us"
	}

	data, err := client.Post(ctx, ScreenerURL, nil, criteria)
	if err != nil {
		return nil, err
	}

	var response struct {
		Finance struct {
			Result []struct {
				Count  int     `json:"count"`
				Total  int     `json:"total"`
				Quotes []Quote `json:"quotes"`
			} `json:"result"`
			Error *struct {
				Code        string `json:"code"`
				Description string `json:"description"`
			} `json:"error"`
		} `json:"finance"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse screener response: %w", err)
	}

	if response.Finance.Error != nil {
		return nil, &APIError{
			Code:        response.Finance.Error.Code,
			Description: response.Finance.Error.Description,
		}
	}

	result := &ScreenResult{}
	if len(response.Finance.Result) > 0 {
		r := response.Finance.Result[0]
		result.Count = r.Count
		result.Total = r.Total
		result.Quotes = r.Quotes
	}

	return result, nil
}

// Predefined screener queries

// ScreenMostActive returns the most actively traded stocks
func ScreenMostActive(ctx context.Context, size int) (*ScreenResult, error) {
	criteria := ScreenCriteria{
		Size:      size,
		SortField: "dayvolume",
		SortType:  "DESC",
		Query: map[string]interface{}{
			"operator": "and",
			"operands": []map[string]interface{}{
				{
					"operator": "eq",
					"operands": []interface{}{"region", "us"},
				},
			},
		},
	}
	return Screen(ctx, criteria)
}

// ScreenGainers returns top gaining stocks
func ScreenGainers(ctx context.Context, size int) (*ScreenResult, error) {
	criteria := ScreenCriteria{
		Size:      size,
		SortField: "percentchange",
		SortType:  "DESC",
		Query: map[string]interface{}{
			"operator": "and",
			"operands": []map[string]interface{}{
				{
					"operator": "eq",
					"operands": []interface{}{"region", "us"},
				},
				{
					"operator": "gt",
					"operands": []interface{}{"percentchange", 0},
				},
			},
		},
	}
	return Screen(ctx, criteria)
}

// ScreenLosers returns top losing stocks
func ScreenLosers(ctx context.Context, size int) (*ScreenResult, error) {
	criteria := ScreenCriteria{
		Size:      size,
		SortField: "percentchange",
		SortType:  "ASC",
		Query: map[string]interface{}{
			"operator": "and",
			"operands": []map[string]interface{}{
				{
					"operator": "eq",
					"operands": []interface{}{"region", "us"},
				},
				{
					"operator": "lt",
					"operands": []interface{}{"percentchange", 0},
				},
			},
		},
	}
	return Screen(ctx, criteria)
}

// ScreenByMarketCap screens stocks by market cap range
func ScreenByMarketCap(ctx context.Context, minCap, maxCap int64, size int) (*ScreenResult, error) {
	operands := []map[string]interface{}{
		{
			"operator": "eq",
			"operands": []interface{}{"region", "us"},
		},
	}

	if minCap > 0 {
		operands = append(operands, map[string]interface{}{
			"operator": "gte",
			"operands": []interface{}{"intradaymarketcap", minCap},
		})
	}

	if maxCap > 0 {
		operands = append(operands, map[string]interface{}{
			"operator": "lte",
			"operands": []interface{}{"intradaymarketcap", maxCap},
		})
	}

	criteria := ScreenCriteria{
		Size:      size,
		SortField: "intradaymarketcap",
		SortType:  "DESC",
		Query: map[string]interface{}{
			"operator": "and",
			"operands": operands,
		},
	}
	return Screen(ctx, criteria)
}

// ScreenBySector screens stocks by sector
func ScreenBySector(ctx context.Context, sector string, size int) (*ScreenResult, error) {
	criteria := ScreenCriteria{
		Size:      size,
		SortField: "intradaymarketcap",
		SortType:  "DESC",
		Query: map[string]interface{}{
			"operator": "and",
			"operands": []map[string]interface{}{
				{
					"operator": "eq",
					"operands": []interface{}{"region", "us"},
				},
				{
					"operator": "eq",
					"operands": []interface{}{"sector", sector},
				},
			},
		},
	}
	return Screen(ctx, criteria)
}

// ScreenHighDividend screens for high dividend yield stocks
func ScreenHighDividend(ctx context.Context, minYield float64, size int) (*ScreenResult, error) {
	criteria := ScreenCriteria{
		Size:      size,
		SortField: "dividendyield",
		SortType:  "DESC",
		Query: map[string]interface{}{
			"operator": "and",
			"operands": []map[string]interface{}{
				{
					"operator": "eq",
					"operands": []interface{}{"region", "us"},
				},
				{
					"operator": "gte",
					"operands": []interface{}{"dividendyield", minYield},
				},
			},
		},
	}
	return Screen(ctx, criteria)
}
