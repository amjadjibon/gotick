package yfinance

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetMarketSummary fetches market summary data
func GetMarketSummary(ctx context.Context) (*MarketSummary, error) {
	client, err := getDefaultClient()
	if err != nil {
		return nil, err
	}

	return GetMarketSummaryWithClient(ctx, client)
}

// GetMarketSummaryWithClient fetches market summary using a specific client
func GetMarketSummaryWithClient(ctx context.Context, client *Client) (*MarketSummary, error) {
	data, err := client.Get(ctx, MarketSummaryURL, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		MarketSummaryResponse struct {
			Result []MarketIndex `json:"result"`
			Error  *struct {
				Code        string `json:"code"`
				Description string `json:"description"`
			} `json:"error"`
		} `json:"marketSummaryResponse"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse market summary response: %w", err)
	}

	if response.MarketSummaryResponse.Error != nil {
		return nil, &APIError{
			Code:        response.MarketSummaryResponse.Error.Code,
			Description: response.MarketSummaryResponse.Error.Description,
		}
	}

	summary := &MarketSummary{
		Markets: response.MarketSummaryResponse.Result,
	}

	if len(summary.Markets) > 0 {
		summary.MarketState = summary.Markets[0].MarketState
	}

	return summary, nil
}

// GetMarketTime fetches market time information for an exchange
func GetMarketTime(ctx context.Context, exchange string) (*MarketTime, error) {
	client, err := getDefaultClient()
	if err != nil {
		return nil, err
	}

	return GetMarketTimeWithClient(ctx, client, exchange)
}

// GetMarketTimeWithClient fetches market time using a specific client
func GetMarketTimeWithClient(ctx context.Context, client *Client, exchange string) (*MarketTime, error) {
	data, err := client.Get(ctx, MarketTimeURL, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Finance struct {
			Result []MarketTime `json:"result"`
			Error  *struct {
				Code        string `json:"code"`
				Description string `json:"description"`
			} `json:"error"`
		} `json:"finance"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse market time response: %w", err)
	}

	if response.Finance.Error != nil {
		return nil, &APIError{
			Code:        response.Finance.Error.Code,
			Description: response.Finance.Error.Description,
		}
	}

	// Find matching exchange or return first
	for _, mt := range response.Finance.Result {
		if exchange == "" || mt.Exchange == exchange {
			return &mt, nil
		}
	}

	if len(response.Finance.Result) > 0 {
		return &response.Finance.Result[0], nil
	}

	return nil, ErrNotFound
}

// GetTrending fetches trending tickers
func GetTrending(ctx context.Context, region string, count int) ([]Quote, error) {
	client, err := getDefaultClient()
	if err != nil {
		return nil, err
	}

	return GetTrendingWithClient(ctx, client, region, count)
}

// GetTrendingWithClient fetches trending tickers using a specific client
func GetTrendingWithClient(ctx context.Context, client *Client, region string, count int) ([]Quote, error) {
	if region == "" {
		region = "US"
	}
	if count <= 0 {
		count = 10
	}

	// Use screener to get trending
	criteria := ScreenCriteria{
		Size:      count,
		Region:    region,
		SortField: "dayvolume",
		SortType:  "DESC",
	}

	result, err := ScreenWithClient(ctx, client, criteria)
	if err != nil {
		return nil, err
	}

	return result.Quotes, nil
}

// Major index symbols
var (
	// US Indices
	IndexSP500    = "^GSPC"
	IndexDowJones = "^DJI"
	IndexNasdaq   = "^IXIC"
	IndexRussell  = "^RUT"
	IndexVIX      = "^VIX"

	// International Indices
	IndexFTSE100  = "^FTSE"
	IndexDAX      = "^GDAXI"
	IndexNikkei   = "^N225"
	IndexHangSeng = "^HSI"
	IndexShanghai = "000001.SS"

	// Futures
	FuturesGold     = "GC=F"
	FuturesSilver   = "SI=F"
	FuturesCrudeOil = "CL=F"
	FuturesNatGas   = "NG=F"
	FuturesSP500    = "ES=F"
	FuturesNasdaq   = "NQ=F"

	// Crypto
	CryptoBTC = "BTC-USD"
	CryptoETH = "ETH-USD"
)

// GetMajorIndices fetches quotes for major US indices
func GetMajorIndices(ctx context.Context) ([]Quote, error) {
	symbols := []string{
		IndexSP500,
		IndexDowJones,
		IndexNasdaq,
		IndexRussell,
		IndexVIX,
	}
	return QuoteMultiple(ctx, symbols)
}

// GetMajorFutures fetches quotes for major futures
func GetMajorFutures(ctx context.Context) ([]Quote, error) {
	symbols := []string{
		FuturesGold,
		FuturesSilver,
		FuturesCrudeOil,
		FuturesNatGas,
		FuturesSP500,
		FuturesNasdaq,
	}
	return QuoteMultiple(ctx, symbols)
}

// GetMajorCrypto fetches quotes for major cryptocurrencies
func GetMajorCrypto(ctx context.Context) ([]Quote, error) {
	symbols := []string{
		CryptoBTC,
		CryptoETH,
	}
	return QuoteMultiple(ctx, symbols)
}
