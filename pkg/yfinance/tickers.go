package yfinance

import (
	"context"
	"sync"
)

// Tickers provides batch operations for multiple tickers
type Tickers struct {
	symbols []string
	tickers map[string]*Ticker
	client  *Client
	mu      sync.RWMutex
}

// NewTickers creates a new Tickers instance for batch operations
func NewTickers(symbols ...string) (*Tickers, error) {
	client, err := getDefaultClient()
	if err != nil {
		return nil, err
	}

	t := &Tickers{
		symbols: symbols,
		tickers: make(map[string]*Ticker),
		client:  client,
	}

	// Pre-create ticker instances
	for _, symbol := range symbols {
		ticker, err := NewTicker(symbol, WithClient(client))
		if err != nil {
			continue
		}
		t.tickers[symbol] = ticker
	}

	return t, nil
}

// Symbols returns the list of symbols
func (t *Tickers) Symbols() []string {
	return t.symbols
}

// Ticker returns a specific ticker by symbol
func (t *Tickers) Ticker(symbol string) (*Ticker, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	ticker, ok := t.tickers[symbol]
	return ticker, ok
}

// Quotes fetches quotes for all tickers
func (t *Tickers) Quotes(ctx context.Context) (map[string]*Quote, error) {
	quotes, err := QuoteMultiple(ctx, t.symbols)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*Quote)
	for i := range quotes {
		result[quotes[i].Symbol] = &quotes[i]
	}
	return result, nil
}

// History fetches historical data for all tickers
func (t *Tickers) History(ctx context.Context, params HistoryParams) (map[string]*ChartData, error) {
	result := make(map[string]*ChartData)
	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, len(t.symbols))

	for _, symbol := range t.symbols {
		wg.Add(1)
		go func(sym string) {
			defer wg.Done()
			ticker, ok := t.tickers[sym]
			if !ok {
				return
			}
			history, err := ticker.History(ctx, params)
			if err != nil {
				errChan <- err
				return
			}
			mu.Lock()
			result[sym] = history
			mu.Unlock()
		}(symbol)
	}

	wg.Wait()
	close(errChan)

	// Return first error if any
	for err := range errChan {
		return result, err
	}

	return result, nil
}

// Info fetches company info for all tickers
func (t *Tickers) Info(ctx context.Context, modules ...string) (map[string]*QuoteSummary, error) {
	result := make(map[string]*QuoteSummary)
	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, len(t.symbols))

	for _, symbol := range t.symbols {
		wg.Add(1)
		go func(sym string) {
			defer wg.Done()
			ticker, ok := t.tickers[sym]
			if !ok {
				return
			}
			info, err := ticker.Info(ctx, modules...)
			if err != nil {
				errChan <- err
				return
			}
			mu.Lock()
			result[sym] = info
			mu.Unlock()
		}(symbol)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		return result, err
	}

	return result, nil
}

// Recommendations fetches analyst recommendations for all tickers
func (t *Tickers) Recommendations(ctx context.Context) (map[string][]RecommendationTrend, error) {
	result := make(map[string][]RecommendationTrend)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, symbol := range t.symbols {
		wg.Add(1)
		go func(sym string) {
			defer wg.Done()
			ticker, ok := t.tickers[sym]
			if !ok {
				return
			}
			recs, err := ticker.Recommendations(ctx)
			if err != nil {
				return
			}
			mu.Lock()
			result[sym] = recs
			mu.Unlock()
		}(symbol)
	}

	wg.Wait()
	return result, nil
}

// MajorHolders fetches major holders for all tickers
func (t *Tickers) MajorHolders(ctx context.Context) (map[string]*MajorHolders, error) {
	result := make(map[string]*MajorHolders)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, symbol := range t.symbols {
		wg.Add(1)
		go func(sym string) {
			defer wg.Done()
			ticker, ok := t.tickers[sym]
			if !ok {
				return
			}
			holders, err := ticker.MajorHolders(ctx)
			if err != nil {
				return
			}
			mu.Lock()
			result[sym] = holders
			mu.Unlock()
		}(symbol)
	}

	wg.Wait()
	return result, nil
}
