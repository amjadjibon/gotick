package yfinance

import (
	"context"
	"encoding/json"
	"fmt"
)

// FundHolding represents a holding in an ETF or mutual fund
type FundHolding struct {
	Symbol  string  `json:"symbol"`
	Name    string  `json:"holdingName"`
	Percent float64 `json:"holdingPercent"`
	Shares  int64   `json:"shares,omitempty"`
	Value   int64   `json:"value,omitempty"`
}

// FundSectorWeighting represents sector allocation
type FundSectorWeighting struct {
	Sector  string  `json:"sector"`
	Percent float64 `json:"weight"`
}

// FundOverview represents fund overview data
type FundOverview struct {
	Category                  string  `json:"category"`
	FundFamily                string  `json:"fundFamily"`
	LegalType                 string  `json:"legalType"`
	TotalAssets               int64   `json:"totalAssets"`
	YTDReturn                 float64 `json:"ytdReturn"`
	TrailingThreeMonthReturns float64 `json:"trailingThreeMonthReturns"`
	TrailingThreeYearReturns  float64 `json:"trailingThreeYearReturns"`
	TrailingFiveYearReturns   float64 `json:"trailingFiveYearReturns"`
	ExpenseRatio              float64 `json:"annualReportExpenseRatio"`
	Turnover                  float64 `json:"turnover"`
	TotalHoldings             int     `json:"holdings"`
	Top10HoldingsPercent      float64 `json:"top10HoldingsPercent"`
}

// FundData represents all fund-specific data
type FundData struct {
	Overview         *FundOverview         `json:"overview"`
	Holdings         []FundHolding         `json:"holdings"`
	SectorWeightings []FundSectorWeighting `json:"sectorWeightings"`
	BondHoldings     map[string]float64    `json:"bondHoldings,omitempty"`
	EquityHoldings   map[string]float64    `json:"equityHoldings,omitempty"`
}

// FundHoldings fetches holdings for an ETF or mutual fund
func (t *Ticker) FundHoldings(ctx context.Context) ([]FundHolding, error) {
	endpoint := fmt.Sprintf("%s/%s", QuoteSummaryURL, t.Symbol)
	params := buildModulesParams(ModuleTopHoldings)

	data, err := t.client.Get(ctx, endpoint, params)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		QuoteSummary struct {
			Result []struct {
				TopHoldings struct {
					Holdings []struct {
						Symbol  string   `json:"symbol"`
						Name    string   `json:"holdingName"`
						Percent RawValue `json:"holdingPercent"`
					} `json:"holdings"`
				} `json:"topHoldings"`
			} `json:"result"`
		} `json:"quoteSummary"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse fund holdings: %w", err))
	}

	if len(response.QuoteSummary.Result) == 0 {
		return nil, NewSymbolError(t.Symbol, ErrNoData)
	}

	var holdings []FundHolding
	for _, h := range response.QuoteSummary.Result[0].TopHoldings.Holdings {
		holdings = append(holdings, FundHolding{
			Symbol:  h.Symbol,
			Name:    h.Name,
			Percent: h.Percent.Raw,
		})
	}

	return holdings, nil
}

// FundSectorWeightings fetches sector weightings for an ETF or mutual fund
func (t *Ticker) FundSectorWeightings(ctx context.Context) ([]FundSectorWeighting, error) {
	endpoint := fmt.Sprintf("%s/%s", QuoteSummaryURL, t.Symbol)
	params := buildModulesParams(ModuleTopHoldings)

	data, err := t.client.Get(ctx, endpoint, params)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		QuoteSummary struct {
			Result []struct {
				TopHoldings struct {
					SectorWeightings []map[string]RawValue `json:"sectorWeightings"`
				} `json:"topHoldings"`
			} `json:"result"`
		} `json:"quoteSummary"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse sector weightings: %w", err))
	}

	if len(response.QuoteSummary.Result) == 0 {
		return nil, NewSymbolError(t.Symbol, ErrNoData)
	}

	var weightings []FundSectorWeighting
	for _, sw := range response.QuoteSummary.Result[0].TopHoldings.SectorWeightings {
		for sector, weight := range sw {
			weightings = append(weightings, FundSectorWeighting{
				Sector:  sector,
				Percent: weight.Raw,
			})
		}
	}

	return weightings, nil
}

// FundProfile fetches fund profile/overview data
func (t *Ticker) FundProfile(ctx context.Context) (*FundOverview, error) {
	endpoint := fmt.Sprintf("%s/%s", QuoteSummaryURL, t.Symbol)
	params := buildModulesParams(ModuleFundProfile)

	data, err := t.client.Get(ctx, endpoint, params)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		QuoteSummary struct {
			Result []struct {
				FundProfile struct {
					CategoryName string `json:"categoryName"`
					FundFamily   string `json:"family"`
					LegalType    string `json:"legalType"`
				} `json:"fundProfile"`
			} `json:"result"`
		} `json:"quoteSummary"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse fund profile: %w", err))
	}

	if len(response.QuoteSummary.Result) == 0 {
		return nil, NewSymbolError(t.Symbol, ErrNoData)
	}

	fp := response.QuoteSummary.Result[0].FundProfile
	return &FundOverview{
		Category:   fp.CategoryName,
		FundFamily: fp.FundFamily,
		LegalType:  fp.LegalType,
	}, nil
}

// FundPerformance fetches fund performance data
func (t *Ticker) FundPerformance(ctx context.Context) (*FundOverview, error) {
	endpoint := fmt.Sprintf("%s/%s", QuoteSummaryURL, t.Symbol)
	params := buildModulesParams(ModuleFundPerformance)

	data, err := t.client.Get(ctx, endpoint, params)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		QuoteSummary struct {
			Result []struct {
				FundPerformance struct {
					TrailingReturns []struct {
						Period string   `json:"period"`
						Return RawValue `json:"return"`
					} `json:"trailingReturns"`
				} `json:"fundPerformance"`
			} `json:"result"`
		} `json:"quoteSummary"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse fund performance: %w", err))
	}

	if len(response.QuoteSummary.Result) == 0 {
		return nil, NewSymbolError(t.Symbol, ErrNoData)
	}

	overview := &FundOverview{}
	for _, tr := range response.QuoteSummary.Result[0].FundPerformance.TrailingReturns {
		switch tr.Period {
		case "ytd":
			overview.YTDReturn = tr.Return.Raw
		case "3m":
			overview.TrailingThreeMonthReturns = tr.Return.Raw
		case "3y":
			overview.TrailingThreeYearReturns = tr.Return.Raw
		case "5y":
			overview.TrailingFiveYearReturns = tr.Return.Raw
		}
	}

	return overview, nil
}
