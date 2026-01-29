package yfinance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Ticker represents a financial instrument and provides methods to fetch its data
type Ticker struct {
	Symbol string
	client *Client
}

// TickerOption is a function that configures Ticker options
type TickerOption func(*Ticker)

// WithClient sets a custom client for the ticker
func WithClient(client *Client) TickerOption {
	return func(t *Ticker) {
		t.client = client
	}
}

// NewTicker creates a new Ticker instance for the given symbol
func NewTicker(symbol string, opts ...TickerOption) (*Ticker, error) {
	if symbol == "" {
		return nil, ErrInvalidSymbol
	}

	ticker := &Ticker{
		Symbol: strings.ToUpper(symbol),
	}

	for _, opt := range opts {
		opt(ticker)
	}

	// Use default client if none provided
	if ticker.client == nil {
		client, err := getDefaultClient()
		if err != nil {
			return nil, err
		}
		ticker.client = client
	}

	return ticker, nil
}

// Quote fetches real-time quote data for the ticker
func (t *Ticker) Quote(ctx context.Context) (*Quote, error) {
	params := url.Values{}
	params.Set("symbols", t.Symbol)

	data, err := t.client.Get(ctx, QuoteURL, params)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
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
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse quote response: %w", err))
	}

	if response.QuoteResponse.Error != nil {
		return nil, NewSymbolError(t.Symbol, &APIError{
			Code:        response.QuoteResponse.Error.Code,
			Description: response.QuoteResponse.Error.Description,
		})
	}

	if len(response.QuoteResponse.Result) == 0 {
		return nil, NewSymbolError(t.Symbol, ErrNotFound)
	}

	return &response.QuoteResponse.Result[0], nil
}

// History fetches historical OHLCV data for the ticker
func (t *Ticker) History(ctx context.Context, params HistoryParams) (*ChartData, error) {
	endpoint := fmt.Sprintf("%s/%s", ChartURL, t.Symbol)

	queryParams := url.Values{}

	// Set period or date range
	//nolint:gocritic // ifElseChain: if-else chain is clearer here
	if !params.Start.IsZero() && !params.End.IsZero() {
		queryParams.Set("period1", strconv.FormatInt(params.Start.Unix(), 10))
		queryParams.Set("period2", strconv.FormatInt(params.End.Unix(), 10))
	} else if params.Period != "" {
		queryParams.Set("range", string(params.Period))
	} else {
		queryParams.Set("range", string(Period1mo)) // Default to 1 month
	}

	// Set interval
	if params.Interval != "" {
		queryParams.Set("interval", string(params.Interval))
	} else {
		queryParams.Set("interval", string(Interval1d)) // Default to daily
	}

	// Set events
	if params.Events != "" {
		queryParams.Set("events", params.Events)
	}

	// Pre/Post market
	if params.PrePost {
		queryParams.Set("includePrePost", "true")
	}

	data, err := t.client.Get(ctx, endpoint, queryParams)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		Chart struct {
			Result []struct {
				Meta       ChartMeta `json:"meta"`
				Timestamp  []int64   `json:"timestamp"`
				Indicators struct {
					Quote []struct {
						Open   []float64 `json:"open"`
						High   []float64 `json:"high"`
						Low    []float64 `json:"low"`
						Close  []float64 `json:"close"`
						Volume []int64   `json:"volume"`
					} `json:"quote"`
					AdjClose []struct {
						AdjClose []float64 `json:"adjclose"`
					} `json:"adjclose"`
				} `json:"indicators"`
			} `json:"result"`
			Error *struct {
				Code        string `json:"code"`
				Description string `json:"description"`
			} `json:"error"`
		} `json:"chart"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse chart response: %w", err))
	}

	if response.Chart.Error != nil {
		return nil, NewSymbolError(t.Symbol, &APIError{
			Code:        response.Chart.Error.Code,
			Description: response.Chart.Error.Description,
		})
	}

	if len(response.Chart.Result) == 0 {
		return nil, NewSymbolError(t.Symbol, ErrNoData)
	}

	result := response.Chart.Result[0]
	chartData := &ChartData{
		Symbol:   t.Symbol,
		Currency: result.Meta.Currency,
		Interval: params.Interval,
		Meta:     &result.Meta,
		Bars:     make([]Bar, len(result.Timestamp)),
	}

	if len(result.Indicators.Quote) > 0 {
		quote := result.Indicators.Quote[0]
		var adjCloses []float64
		if len(result.Indicators.AdjClose) > 0 {
			adjCloses = result.Indicators.AdjClose[0].AdjClose
		}

		for i, ts := range result.Timestamp {
			bar := Bar{
				Timestamp: time.Unix(ts, 0),
			}
			if i < len(quote.Open) {
				bar.Open = quote.Open[i]
			}
			if i < len(quote.High) {
				bar.High = quote.High[i]
			}
			if i < len(quote.Low) {
				bar.Low = quote.Low[i]
			}
			if i < len(quote.Close) {
				bar.Close = quote.Close[i]
			}
			if i < len(quote.Volume) {
				bar.Volume = quote.Volume[i]
			}
			if adjCloses != nil && i < len(adjCloses) {
				bar.AdjClose = adjCloses[i]
			} else {
				bar.AdjClose = bar.Close
			}
			chartData.Bars[i] = bar
		}
	}

	return chartData, nil
}

// Info fetches comprehensive information about the ticker using quoteSummary
func (t *Ticker) Info(ctx context.Context, modules ...string) (*QuoteSummary, error) {
	if len(modules) == 0 {
		modules = DefaultModules()
	}

	endpoint := fmt.Sprintf("%s/%s", QuoteSummaryURL, t.Symbol)
	params := url.Values{}
	params.Set("modules", strings.Join(modules, ","))

	data, err := t.client.Get(ctx, endpoint, params)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		QuoteSummary struct {
			Result []map[string]json.RawMessage `json:"result"`
			Error  *struct {
				Code        string `json:"code"`
				Description string `json:"description"`
			} `json:"error"`
		} `json:"quoteSummary"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse quote summary response: %w", err))
	}

	if response.QuoteSummary.Error != nil {
		return nil, NewSymbolError(t.Symbol, &APIError{
			Code:        response.QuoteSummary.Error.Code,
			Description: response.QuoteSummary.Error.Description,
		})
	}

	if len(response.QuoteSummary.Result) == 0 {
		return nil, NewSymbolError(t.Symbol, ErrNotFound)
	}

	result := response.QuoteSummary.Result[0]
	summary := &QuoteSummary{Symbol: t.Symbol}

	// Parse each module
	if raw, ok := result["assetProfile"]; ok {
		summary.AssetProfile = &AssetProfile{}
		_ = json.Unmarshal(raw, summary.AssetProfile)
	}
	if raw, ok := result["summaryProfile"]; ok {
		summary.SummaryProfile = &SummaryProfile{}
		_ = json.Unmarshal(raw, summary.SummaryProfile)
	}
	if raw, ok := result["summaryDetail"]; ok {
		summary.SummaryDetail = &SummaryDetail{}
		_ = json.Unmarshal(raw, summary.SummaryDetail)
	}
	if raw, ok := result["price"]; ok {
		summary.Price = &PriceInfo{}
		_ = json.Unmarshal(raw, summary.Price)
	}
	if raw, ok := result["defaultKeyStatistics"]; ok {
		summary.KeyStatistics = &KeyStatistics{}
		_ = json.Unmarshal(raw, summary.KeyStatistics)
	}
	if raw, ok := result["financialData"]; ok {
		summary.FinancialData = &FinancialData{}
		_ = json.Unmarshal(raw, summary.FinancialData)
	}
	if raw, ok := result["calendarEvents"]; ok {
		summary.CalendarEvents = &CalendarEvents{}
		_ = json.Unmarshal(raw, summary.CalendarEvents)
	}

	return summary, nil
}

// Options fetches options chain data for the ticker
func (t *Ticker) Options(ctx context.Context, expiration string) (*OptionChain, error) {
	endpoint := fmt.Sprintf("%s/%s", OptionsURL, t.Symbol)
	params := url.Values{}
	if expiration != "" {
		params.Set("date", expiration)
	}

	data, err := t.client.Get(ctx, endpoint, params)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		OptionChain struct {
			Result []struct {
				UnderlyingSymbol string    `json:"underlyingSymbol"`
				ExpirationDates  []int64   `json:"expirationDates"`
				Strikes          []float64 `json:"strikes"`
				Quote            Quote     `json:"quote"`
				Options          []struct {
					ExpirationDate int64    `json:"expirationDate"`
					Calls          []Option `json:"calls"`
					Puts           []Option `json:"puts"`
				} `json:"options"`
			} `json:"result"`
			Error *struct {
				Code        string `json:"code"`
				Description string `json:"description"`
			} `json:"error"`
		} `json:"optionChain"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse options response: %w", err))
	}

	if response.OptionChain.Error != nil {
		return nil, NewSymbolError(t.Symbol, &APIError{
			Code:        response.OptionChain.Error.Code,
			Description: response.OptionChain.Error.Description,
		})
	}

	if len(response.OptionChain.Result) == 0 {
		return nil, NewSymbolError(t.Symbol, ErrNoData)
	}

	result := response.OptionChain.Result[0]
	chain := &OptionChain{
		Symbol:          t.Symbol,
		UnderlyingPrice: result.Quote.RegularMarketPrice,
		ExpirationDates: result.ExpirationDates,
		Strikes:         result.Strikes,
	}

	if len(result.Options) > 0 {
		chain.Calls = result.Options[0].Calls
		chain.Puts = result.Options[0].Puts
	}

	return chain, nil
}

// Financials fetches financial statement data for the ticker
func (t *Ticker) Financials(ctx context.Context, keys []string, period string) (*Financial, error) {
	if len(keys) == 0 {
		keys = AllFinancialKeys()
	}

	// Build type parameter
	types := make([]string, len(keys))
	for i, key := range keys {
		if period == "quarterly" {
			types[i] = "quarterly" + key
		} else {
			types[i] = "annual" + key
		}
	}

	endpoint := fmt.Sprintf("%s/%s", FundamentalsURL, t.Symbol)
	params := url.Values{}
	params.Set("type", strings.Join(types, ","))
	params.Set("merge", "false")
	params.Set("padTimeSeries", "true")

	data, err := t.client.Get(ctx, endpoint, params)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		Timeseries struct {
			Result []struct {
				Meta      map[string]interface{}   `json:"meta"`
				Timestamp []int64                  `json:"timestamp"`
				Data      map[string][]interface{} `json:"-"`
			} `json:"result"`
			Error *struct {
				Code        string `json:"code"`
				Description string `json:"description"`
			} `json:"error"`
		} `json:"timeseries"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse financials response: %w", err))
	}

	if response.Timeseries.Error != nil {
		return nil, NewSymbolError(t.Symbol, &APIError{
			Code:        response.Timeseries.Error.Code,
			Description: response.Timeseries.Error.Description,
		})
	}

	financial := &Financial{
		Symbol: t.Symbol,
		Data:   make(map[string][]FinancialValue),
	}

	if len(response.Timeseries.Result) > 0 {
		result := response.Timeseries.Result[0]
		financial.Timestamp = result.Timestamp
	}

	return financial, nil
}

// News fetches news articles related to the ticker
func (t *Ticker) News(ctx context.Context, count int) ([]NewsItem, error) {
	if count <= 0 {
		count = 10
	}

	// Use search endpoint to get news
	params := url.Values{}
	params.Set("q", t.Symbol)
	params.Set("newsCount", strconv.Itoa(count))
	params.Set("quotesCount", "0")

	data, err := t.client.Get(ctx, SearchURL, params)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		News []NewsItem `json:"news"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse news response: %w", err))
	}

	return response.News, nil
}

// Dividends fetches historical dividend data
func (t *Ticker) Dividends(ctx context.Context, params HistoryParams) ([]Dividend, error) {
	if params.Period == "" {
		params.Period = PeriodMax
	}
	params.Events = "div"

	endpoint := fmt.Sprintf("%s/%s", ChartURL, t.Symbol)
	queryParams := url.Values{}
	queryParams.Set("range", string(params.Period))
	queryParams.Set("interval", string(Interval1d))
	queryParams.Set("events", "div")

	data, err := t.client.Get(ctx, endpoint, queryParams)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		Chart struct {
			Result []struct {
				Events struct {
					Dividends map[string]struct {
						Amount float64 `json:"amount"`
						Date   int64   `json:"date"`
					} `json:"dividends"`
				} `json:"events"`
			} `json:"result"`
		} `json:"chart"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse dividends response: %w", err))
	}

	var dividends []Dividend
	if len(response.Chart.Result) > 0 && response.Chart.Result[0].Events.Dividends != nil {
		for _, div := range response.Chart.Result[0].Events.Dividends {
			dividends = append(dividends, Dividend{
				Date:   time.Unix(div.Date, 0),
				Amount: div.Amount,
			})
		}
	}

	return dividends, nil
}

// Splits fetches historical stock split data
func (t *Ticker) Splits(ctx context.Context, params HistoryParams) ([]Split, error) {
	if params.Period == "" {
		params.Period = PeriodMax
	}

	endpoint := fmt.Sprintf("%s/%s", ChartURL, t.Symbol)
	queryParams := url.Values{}
	queryParams.Set("range", string(params.Period))
	queryParams.Set("interval", string(Interval1d))
	queryParams.Set("events", "split")

	data, err := t.client.Get(ctx, endpoint, queryParams)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		Chart struct {
			Result []struct {
				Events struct {
					Splits map[string]struct {
						Date        int64   `json:"date"`
						Numerator   float64 `json:"numerator"`
						Denominator float64 `json:"denominator"`
						SplitRatio  string  `json:"splitRatio"`
					} `json:"splits"`
				} `json:"events"`
			} `json:"result"`
		} `json:"chart"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse splits response: %w", err))
	}

	var splits []Split
	if len(response.Chart.Result) > 0 && response.Chart.Result[0].Events.Splits != nil {
		for _, s := range response.Chart.Result[0].Events.Splits {
			splits = append(splits, Split{
				Date:        time.Unix(s.Date, 0),
				Numerator:   s.Numerator,
				Denominator: s.Denominator,
				Ratio:       s.SplitRatio,
			})
		}
	}

	return splits, nil
}

// Dividend represents a dividend payment
type Dividend struct {
	Date   time.Time `json:"date"`
	Amount float64   `json:"amount"`
}

// Split represents a stock split
type Split struct {
	Date        time.Time `json:"date"`
	Numerator   float64   `json:"numerator"`
	Denominator float64   `json:"denominator"`
	Ratio       string    `json:"ratio"`
}
