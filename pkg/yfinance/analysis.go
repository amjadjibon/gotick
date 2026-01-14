package yfinance

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// RecommendationTrend represents analyst recommendation trends
type RecommendationTrend struct {
	Period     string `json:"period"`
	StrongBuy  int    `json:"strongBuy"`
	Buy        int    `json:"buy"`
	Hold       int    `json:"hold"`
	Sell       int    `json:"sell"`
	StrongSell int    `json:"strongSell"`
}

// PriceTarget represents analyst price targets
type PriceTarget struct {
	Current     float64 `json:"currentPrice"`
	Low         float64 `json:"targetLowPrice"`
	High        float64 `json:"targetHighPrice"`
	Mean        float64 `json:"targetMeanPrice"`
	Median      float64 `json:"targetMedianPrice"`
	NumAnalysts int     `json:"numberOfAnalystOpinions"`
}

// EarningsEstimate represents earnings estimates for a period
type EarningsEstimate struct {
	Period     string  `json:"period"`
	EndDate    string  `json:"endDate"`
	Avg        float64 `json:"avg"`
	Low        float64 `json:"low"`
	High       float64 `json:"high"`
	YearAgoEps float64 `json:"yearAgoEps"`
	NumOfEst   int     `json:"numberOfAnalysts"`
	Growth     float64 `json:"growth"`
}

// RevenueEstimate represents revenue estimates for a period
type RevenueEstimate struct {
	Period         string  `json:"period"`
	EndDate        string  `json:"endDate"`
	Avg            int64   `json:"avg"`
	Low            int64   `json:"low"`
	High           int64   `json:"high"`
	YearAgoRevenue int64   `json:"yearAgoRevenue"`
	NumOfEst       int     `json:"numberOfAnalysts"`
	Growth         float64 `json:"growth"`
}

// EPSTrend represents EPS trend data
type EPSTrend struct {
	Period        string  `json:"period"`
	EndDate       string  `json:"endDate"`
	Current       float64 `json:"current"`
	SevenDaysAgo  float64 `json:"7daysAgo"`
	ThirtyDaysAgo float64 `json:"30daysAgo"`
	SixtyDaysAgo  float64 `json:"60daysAgo"`
	NinetyDaysAgo float64 `json:"90daysAgo"`
}

// EPSRevision represents EPS revision data
type EPSRevision struct {
	Period     string `json:"period"`
	EndDate    string `json:"endDate"`
	UpLast7    int    `json:"upLast7days"`
	UpLast30   int    `json:"upLast30days"`
	DownLast7  int    `json:"downLast7days"`
	DownLast30 int    `json:"downLast30days"`
}

// EarningsHistoryItem represents a historical earnings record
type EarningsHistoryItem struct {
	Quarter         string  `json:"quarter"`
	EpsActual       float64 `json:"epsActual"`
	EpsEstimate     float64 `json:"epsEstimate"`
	EpsDifference   float64 `json:"epsDifference"`
	SurprisePercent float64 `json:"surprisePercent"`
}

// GrowthEstimate represents growth estimates
type GrowthEstimate struct {
	Period   string  `json:"period"`
	Growth   float64 `json:"growth"`
	Industry float64 `json:"industryGrowth,omitempty"`
	Sector   float64 `json:"sectorGrowth,omitempty"`
	SP500    float64 `json:"sp500Growth,omitempty"`
}

// Recommendations fetches analyst recommendation trends
func (t *Ticker) Recommendations(ctx context.Context) ([]RecommendationTrend, error) {
	endpoint := fmt.Sprintf("%s/%s", QuoteSummaryURL, t.Symbol)
	params := buildModulesParams(ModuleRecommendationTrend)

	data, err := t.client.Get(ctx, endpoint, params)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		QuoteSummary struct {
			Result []struct {
				RecommendationTrend struct {
					Trend []RecommendationTrend `json:"trend"`
				} `json:"recommendationTrend"`
			} `json:"result"`
		} `json:"quoteSummary"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse recommendations: %w", err))
	}

	if len(response.QuoteSummary.Result) == 0 {
		return nil, NewSymbolError(t.Symbol, ErrNoData)
	}

	return response.QuoteSummary.Result[0].RecommendationTrend.Trend, nil
}

// AnalystPriceTargets fetches analyst price targets
func (t *Ticker) AnalystPriceTargets(ctx context.Context) (*PriceTarget, error) {
	endpoint := fmt.Sprintf("%s/%s", QuoteSummaryURL, t.Symbol)
	params := buildModulesParams(ModuleFinancialData)

	data, err := t.client.Get(ctx, endpoint, params)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		QuoteSummary struct {
			Result []struct {
				FinancialData struct {
					CurrentPrice            RawValue `json:"currentPrice"`
					TargetLowPrice          RawValue `json:"targetLowPrice"`
					TargetHighPrice         RawValue `json:"targetHighPrice"`
					TargetMeanPrice         RawValue `json:"targetMeanPrice"`
					TargetMedianPrice       RawValue `json:"targetMedianPrice"`
					NumberOfAnalystOpinions RawValue `json:"numberOfAnalystOpinions"`
				} `json:"financialData"`
			} `json:"result"`
		} `json:"quoteSummary"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse price targets: %w", err))
	}

	if len(response.QuoteSummary.Result) == 0 {
		return nil, NewSymbolError(t.Symbol, ErrNoData)
	}

	fd := response.QuoteSummary.Result[0].FinancialData
	return &PriceTarget{
		Current:     fd.CurrentPrice.Raw,
		Low:         fd.TargetLowPrice.Raw,
		High:        fd.TargetHighPrice.Raw,
		Mean:        fd.TargetMeanPrice.Raw,
		Median:      fd.TargetMedianPrice.Raw,
		NumAnalysts: int(fd.NumberOfAnalystOpinions.Raw),
	}, nil
}

// RawValue represents a Yahoo Finance value with raw and formatted versions
type RawValue struct {
	Raw float64 `json:"raw"`
	Fmt string  `json:"fmt"`
}

// EarningsEstimates fetches earnings estimates for upcoming periods
func (t *Ticker) EarningsEstimates(ctx context.Context) ([]EarningsEstimate, error) {
	endpoint := fmt.Sprintf("%s/%s", QuoteSummaryURL, t.Symbol)
	params := buildModulesParams(ModuleEarningsTrend)

	data, err := t.client.Get(ctx, endpoint, params)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		QuoteSummary struct {
			Result []struct {
				EarningsTrend struct {
					Trend []struct {
						Period           string `json:"period"`
						EndDate          string `json:"endDate"`
						EarningsEstimate struct {
							Avg        RawValue `json:"avg"`
							Low        RawValue `json:"low"`
							High       RawValue `json:"high"`
							YearAgoEps RawValue `json:"yearAgoEps"`
							NumOfEst   RawValue `json:"numberOfAnalysts"`
							Growth     RawValue `json:"growth"`
						} `json:"earningsEstimate"`
					} `json:"trend"`
				} `json:"earningsTrend"`
			} `json:"result"`
		} `json:"quoteSummary"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse earnings estimates: %w", err))
	}

	if len(response.QuoteSummary.Result) == 0 {
		return nil, NewSymbolError(t.Symbol, ErrNoData)
	}

	var estimates []EarningsEstimate
	for _, trend := range response.QuoteSummary.Result[0].EarningsTrend.Trend {
		estimates = append(estimates, EarningsEstimate{
			Period:     trend.Period,
			EndDate:    trend.EndDate,
			Avg:        trend.EarningsEstimate.Avg.Raw,
			Low:        trend.EarningsEstimate.Low.Raw,
			High:       trend.EarningsEstimate.High.Raw,
			YearAgoEps: trend.EarningsEstimate.YearAgoEps.Raw,
			NumOfEst:   int(trend.EarningsEstimate.NumOfEst.Raw),
			Growth:     trend.EarningsEstimate.Growth.Raw,
		})
	}

	return estimates, nil
}

// RevenueEstimates fetches revenue estimates for upcoming periods
func (t *Ticker) RevenueEstimates(ctx context.Context) ([]RevenueEstimate, error) {
	endpoint := fmt.Sprintf("%s/%s", QuoteSummaryURL, t.Symbol)
	params := buildModulesParams(ModuleEarningsTrend)

	data, err := t.client.Get(ctx, endpoint, params)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		QuoteSummary struct {
			Result []struct {
				EarningsTrend struct {
					Trend []struct {
						Period          string `json:"period"`
						EndDate         string `json:"endDate"`
						RevenueEstimate struct {
							Avg            RawValue `json:"avg"`
							Low            RawValue `json:"low"`
							High           RawValue `json:"high"`
							YearAgoRevenue RawValue `json:"yearAgoRevenue"`
							NumOfEst       RawValue `json:"numberOfAnalysts"`
							Growth         RawValue `json:"growth"`
						} `json:"revenueEstimate"`
					} `json:"trend"`
				} `json:"earningsTrend"`
			} `json:"result"`
		} `json:"quoteSummary"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse revenue estimates: %w", err))
	}

	if len(response.QuoteSummary.Result) == 0 {
		return nil, NewSymbolError(t.Symbol, ErrNoData)
	}

	var estimates []RevenueEstimate
	for _, trend := range response.QuoteSummary.Result[0].EarningsTrend.Trend {
		estimates = append(estimates, RevenueEstimate{
			Period:         trend.Period,
			EndDate:        trend.EndDate,
			Avg:            int64(trend.RevenueEstimate.Avg.Raw),
			Low:            int64(trend.RevenueEstimate.Low.Raw),
			High:           int64(trend.RevenueEstimate.High.Raw),
			YearAgoRevenue: int64(trend.RevenueEstimate.YearAgoRevenue.Raw),
			NumOfEst:       int(trend.RevenueEstimate.NumOfEst.Raw),
			Growth:         trend.RevenueEstimate.Growth.Raw,
		})
	}

	return estimates, nil
}

// EPSTrends fetches EPS trend data
func (t *Ticker) EPSTrends(ctx context.Context) ([]EPSTrend, error) {
	endpoint := fmt.Sprintf("%s/%s", QuoteSummaryURL, t.Symbol)
	params := buildModulesParams(ModuleEarningsTrend)

	data, err := t.client.Get(ctx, endpoint, params)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		QuoteSummary struct {
			Result []struct {
				EarningsTrend struct {
					Trend []struct {
						Period   string `json:"period"`
						EndDate  string `json:"endDate"`
						EpsTrend struct {
							Current    RawValue `json:"current"`
							SevenDays  RawValue `json:"7daysAgo"`
							ThirtyDays RawValue `json:"30daysAgo"`
							SixtyDays  RawValue `json:"60daysAgo"`
							NinetyDays RawValue `json:"90daysAgo"`
						} `json:"epsTrend"`
					} `json:"trend"`
				} `json:"earningsTrend"`
			} `json:"result"`
		} `json:"quoteSummary"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse EPS trends: %w", err))
	}

	if len(response.QuoteSummary.Result) == 0 {
		return nil, NewSymbolError(t.Symbol, ErrNoData)
	}

	var trends []EPSTrend
	for _, t := range response.QuoteSummary.Result[0].EarningsTrend.Trend {
		trends = append(trends, EPSTrend{
			Period:        t.Period,
			EndDate:       t.EndDate,
			Current:       t.EpsTrend.Current.Raw,
			SevenDaysAgo:  t.EpsTrend.SevenDays.Raw,
			ThirtyDaysAgo: t.EpsTrend.ThirtyDays.Raw,
			SixtyDaysAgo:  t.EpsTrend.SixtyDays.Raw,
			NinetyDaysAgo: t.EpsTrend.NinetyDays.Raw,
		})
	}

	return trends, nil
}

// EPSRevisions fetches EPS revision data
func (t *Ticker) EPSRevisions(ctx context.Context) ([]EPSRevision, error) {
	endpoint := fmt.Sprintf("%s/%s", QuoteSummaryURL, t.Symbol)
	params := buildModulesParams(ModuleEarningsTrend)

	data, err := t.client.Get(ctx, endpoint, params)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		QuoteSummary struct {
			Result []struct {
				EarningsTrend struct {
					Trend []struct {
						Period       string `json:"period"`
						EndDate      string `json:"endDate"`
						EpsRevisions struct {
							UpLast7    RawValue `json:"upLast7days"`
							UpLast30   RawValue `json:"upLast30days"`
							DownLast7  RawValue `json:"downLast7days"`
							DownLast30 RawValue `json:"downLast30days"`
						} `json:"epsRevisions"`
					} `json:"trend"`
				} `json:"earningsTrend"`
			} `json:"result"`
		} `json:"quoteSummary"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse EPS revisions: %w", err))
	}

	if len(response.QuoteSummary.Result) == 0 {
		return nil, NewSymbolError(t.Symbol, ErrNoData)
	}

	var revisions []EPSRevision
	for _, t := range response.QuoteSummary.Result[0].EarningsTrend.Trend {
		revisions = append(revisions, EPSRevision{
			Period:     t.Period,
			EndDate:    t.EndDate,
			UpLast7:    int(t.EpsRevisions.UpLast7.Raw),
			UpLast30:   int(t.EpsRevisions.UpLast30.Raw),
			DownLast7:  int(t.EpsRevisions.DownLast7.Raw),
			DownLast30: int(t.EpsRevisions.DownLast30.Raw),
		})
	}

	return revisions, nil
}

// EarningsHistoryData fetches historical earnings data
func (t *Ticker) EarningsHistoryData(ctx context.Context) ([]EarningsHistoryItem, error) {
	endpoint := fmt.Sprintf("%s/%s", QuoteSummaryURL, t.Symbol)
	params := buildModulesParams(ModuleEarningsHistory)

	data, err := t.client.Get(ctx, endpoint, params)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		QuoteSummary struct {
			Result []struct {
				EarningsHistory struct {
					History []struct {
						Quarter         RawValue `json:"fiscalQuarter"`
						EpsActual       RawValue `json:"epsActual"`
						EpsEstimate     RawValue `json:"epsEstimate"`
						EpsDifference   RawValue `json:"epsDifference"`
						SurprisePercent RawValue `json:"surprisePercent"`
					} `json:"history"`
				} `json:"earningsHistory"`
			} `json:"result"`
		} `json:"quoteSummary"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse earnings history: %w", err))
	}

	if len(response.QuoteSummary.Result) == 0 {
		return nil, NewSymbolError(t.Symbol, ErrNoData)
	}

	var history []EarningsHistoryItem
	for _, h := range response.QuoteSummary.Result[0].EarningsHistory.History {
		history = append(history, EarningsHistoryItem{
			Quarter:         h.Quarter.Fmt,
			EpsActual:       h.EpsActual.Raw,
			EpsEstimate:     h.EpsEstimate.Raw,
			EpsDifference:   h.EpsDifference.Raw,
			SurprisePercent: h.SurprisePercent.Raw,
		})
	}

	return history, nil
}

// GrowthEstimates fetches growth estimates
func (t *Ticker) GrowthEstimates(ctx context.Context) ([]GrowthEstimate, error) {
	endpoint := fmt.Sprintf("%s/%s", QuoteSummaryURL, t.Symbol)
	params := buildModulesParams(ModuleEarningsTrend)

	data, err := t.client.Get(ctx, endpoint, params)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		QuoteSummary struct {
			Result []struct {
				EarningsTrend struct {
					Trend []struct {
						Period string   `json:"period"`
						Growth RawValue `json:"growth"`
					} `json:"trend"`
				} `json:"earningsTrend"`
			} `json:"result"`
		} `json:"quoteSummary"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse growth estimates: %w", err))
	}

	if len(response.QuoteSummary.Result) == 0 {
		return nil, NewSymbolError(t.Symbol, ErrNoData)
	}

	var estimates []GrowthEstimate
	for _, t := range response.QuoteSummary.Result[0].EarningsTrend.Trend {
		estimates = append(estimates, GrowthEstimate{
			Period: t.Period,
			Growth: t.Growth.Raw,
		})
	}

	return estimates, nil
}

// buildModulesParams creates query params for quoteSummary modules
func buildModulesParams(modules ...string) map[string][]string {
	return map[string][]string{
		"modules": {strings.Join(modules, ",")},
	}
}
