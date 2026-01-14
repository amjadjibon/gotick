package yfinance

import (
	"context"
	"encoding/json"
	"time"
)

// Action represents a corporate action (dividend or split)
type Action struct {
	Date        time.Time `json:"date"`
	Type        string    `json:"type"`                  // "dividend" or "split"
	Amount      float64   `json:"amount,omitempty"`      // For dividends
	Ratio       string    `json:"ratio,omitempty"`       // For splits (e.g., "4:1")
	Numerator   float64   `json:"numerator,omitempty"`   // For splits
	Denominator float64   `json:"denominator,omitempty"` // For splits
}

// Actions fetches all corporate actions (dividends and splits) for the ticker
func (t *Ticker) Actions(ctx context.Context, params HistoryParams) ([]Action, error) {
	if params.Period == "" {
		params.Period = PeriodMax
	}

	var actions []Action

	// Fetch dividends
	dividends, err := t.Dividends(ctx, params)
	if err == nil {
		for _, d := range dividends {
			actions = append(actions, Action{
				Date:   d.Date,
				Type:   "dividend",
				Amount: d.Amount,
			})
		}
	}

	// Fetch splits
	splits, err := t.Splits(ctx, params)
	if err == nil {
		for _, s := range splits {
			actions = append(actions, Action{
				Date:        s.Date,
				Type:        "split",
				Ratio:       s.Ratio,
				Numerator:   s.Numerator,
				Denominator: s.Denominator,
			})
		}
	}

	// Sort by date (most recent first)
	for i := 0; i < len(actions)-1; i++ {
		for j := i + 1; j < len(actions); j++ {
			if actions[j].Date.After(actions[i].Date) {
				actions[i], actions[j] = actions[j], actions[i]
			}
		}
	}

	return actions, nil
}

// DividendHistory is an alias for Dividends for API compatibility
func (t *Ticker) DividendHistory(ctx context.Context) ([]Dividend, error) {
	return t.Dividends(ctx, HistoryParams{Period: PeriodMax})
}

// SplitHistory is an alias for Splits for API compatibility
func (t *Ticker) SplitHistory(ctx context.Context) ([]Split, error) {
	return t.Splits(ctx, HistoryParams{Period: PeriodMax})
}

// CapitalGains fetches capital gains distributions (for mutual funds)
func (t *Ticker) CapitalGains(ctx context.Context, params HistoryParams) ([]CapitalGain, error) {
	if params.Period == "" {
		params.Period = PeriodMax
	}

	endpoint := ChartURL + "/" + t.Symbol
	queryParams := map[string][]string{
		"range":    {string(params.Period)},
		"interval": {string(Interval1d)},
		"events":   {"capitalGain"},
	}

	data, err := t.client.Get(ctx, endpoint, queryParams)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		Chart struct {
			Result []struct {
				Events struct {
					CapitalGains map[string]struct {
						Amount float64 `json:"amount"`
						Date   int64   `json:"date"`
					} `json:"capitalGains"`
				} `json:"events"`
			} `json:"result"`
		} `json:"chart"`
	}

	if err := parseJSON(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var gains []CapitalGain
	if len(response.Chart.Result) > 0 && response.Chart.Result[0].Events.CapitalGains != nil {
		for _, cg := range response.Chart.Result[0].Events.CapitalGains {
			gains = append(gains, CapitalGain{
				Date:   time.Unix(cg.Date, 0),
				Amount: cg.Amount,
			})
		}
	}

	return gains, nil
}

// CapitalGain represents a capital gains distribution
type CapitalGain struct {
	Date   time.Time `json:"date"`
	Amount float64   `json:"amount"`
}

// parseJSON is a helper to unmarshal JSON data
func parseJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
