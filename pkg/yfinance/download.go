package yfinance

import (
	"context"
	"sync"
	"time"
)

// DownloadParams defines parameters for batch downloading historical data
type DownloadParams struct {
	Symbols  []string  // Symbols to download
	Period   Period    // Time period (e.g., "1mo", "1y")
	Interval Interval  // Data interval (e.g., "1d", "1h")
	Start    time.Time // Start date
	End      time.Time // End date
	PrePost  bool      // Include pre/post market data
	Actions  bool      // Include dividends and splits
	Progress bool      // Show progress (not implemented in Go)
	Threads  int       // Number of concurrent downloads
}

// DownloadResult contains downloaded data for multiple symbols
type DownloadResult struct {
	Data   map[string]*ChartData
	Errors map[string]error
}

// Download fetches historical data for multiple symbols concurrently
func Download(ctx context.Context, params DownloadParams) (*DownloadResult, error) {
	if len(params.Symbols) == 0 {
		return nil, ErrInvalidSymbol
	}

	// Default threads
	if params.Threads <= 0 {
		params.Threads = 5
	}

	result := &DownloadResult{
		Data:   make(map[string]*ChartData),
		Errors: make(map[string]error),
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, params.Threads) // Semaphore for concurrency limit

	client, err := getDefaultClient()
	if err != nil {
		return nil, err
	}

	for _, symbol := range params.Symbols {
		wg.Add(1)
		go func(sym string) {
			defer wg.Done()
			sem <- struct{}{}        // Acquire
			defer func() { <-sem }() // Release

			ticker, err := NewTicker(sym, WithClient(client))
			if err != nil {
				mu.Lock()
				result.Errors[sym] = err
				mu.Unlock()
				return
			}

			histParams := HistoryParams{
				Period:   params.Period,
				Interval: params.Interval,
				Start:    params.Start,
				End:      params.End,
				PrePost:  params.PrePost,
			}

			if params.Actions {
				histParams.Events = "div,split"
			}

			data, err := ticker.History(ctx, histParams)
			if err != nil {
				mu.Lock()
				result.Errors[sym] = err
				mu.Unlock()
				return
			}

			mu.Lock()
			result.Data[sym] = data
			mu.Unlock()
		}(symbol)
	}

	wg.Wait()
	return result, nil
}

// DownloadQuotes fetches quotes for multiple symbols
func DownloadQuotes(ctx context.Context, symbols []string) (map[string]*Quote, error) {
	quotes, err := QuoteMultiple(ctx, symbols)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*Quote)
	for i := range quotes {
		result[quotes[i].Symbol] = &quotes[i]
	}
	return result, nil
}

// DownloadInfo fetches company info for multiple symbols
func DownloadInfo(ctx context.Context, symbols []string, modules ...string) (map[string]*QuoteSummary, error) {
	tickers, err := NewTickers(symbols...)
	if err != nil {
		return nil, err
	}
	return tickers.Info(ctx, modules...)
}
