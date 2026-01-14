package yfinance

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// FinancialStatement represents a financial statement (income, balance, cashflow)
type FinancialStatement struct {
	Symbol    string                     `json:"symbol"`
	Annual    []FinancialStatementPeriod `json:"annual"`
	Quarterly []FinancialStatementPeriod `json:"quarterly"`
}

// FinancialStatementPeriod represents a single period's financial data
type FinancialStatementPeriod struct {
	Date     time.Time          `json:"date"`
	EndDate  string             `json:"endDate"`
	Currency string             `json:"currency"`
	Data     map[string]float64 `json:"data"`
}

// IncomeStatement fetches income statement data
func (t *Ticker) IncomeStatement(ctx context.Context, quarterly bool) (*FinancialStatement, error) {
	var module string
	if quarterly {
		module = ModuleIncomeStatementHistoryQuarterly
	} else {
		module = ModuleIncomeStatementHistory
	}
	return t.fetchFinancialStatement(ctx, module, "incomeStatementHistory", "incomeStatementHistoryQuarterly", quarterly)
}

// BalanceSheet fetches balance sheet data
func (t *Ticker) BalanceSheet(ctx context.Context, quarterly bool) (*FinancialStatement, error) {
	var module string
	if quarterly {
		module = ModuleBalanceSheetHistoryQuarterly
	} else {
		module = ModuleBalanceSheetHistory
	}
	return t.fetchFinancialStatement(ctx, module, "balanceSheetHistory", "balanceSheetHistoryQuarterly", quarterly)
}

// CashFlow fetches cash flow statement data
func (t *Ticker) CashFlow(ctx context.Context, quarterly bool) (*FinancialStatement, error) {
	var module string
	if quarterly {
		module = ModuleCashFlowStatementHistoryQuarterly
	} else {
		module = ModuleCashFlowStatementHistory
	}
	return t.fetchFinancialStatement(ctx, module, "cashFlowStatementHistory", "cashFlowStatementHistoryQuarterly", quarterly)
}

// AllFinancialStatements fetches all financial statements at once
func (t *Ticker) AllFinancialStatements(ctx context.Context, quarterly bool) (*AllFinancials, error) {
	modules := FinancialModules()

	endpoint := fmt.Sprintf("%s/%s", QuoteSummaryURL, t.Symbol)
	params := buildModulesParams(modules...)

	data, err := t.client.Get(ctx, endpoint, params)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		QuoteSummary struct {
			Result []json.RawMessage `json:"result"`
		} `json:"quoteSummary"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse financials: %w", err))
	}

	if len(response.QuoteSummary.Result) == 0 {
		return nil, NewSymbolError(t.Symbol, ErrNoData)
	}

	all := &AllFinancials{Symbol: t.Symbol}

	// Parse each statement type
	income, _ := t.IncomeStatement(ctx, quarterly)
	if income != nil {
		all.IncomeStatement = income
	}

	balance, _ := t.BalanceSheet(ctx, quarterly)
	if balance != nil {
		all.BalanceSheet = balance
	}

	cashflow, _ := t.CashFlow(ctx, quarterly)
	if cashflow != nil {
		all.CashFlow = cashflow
	}

	return all, nil
}

// AllFinancials contains all three financial statements
type AllFinancials struct {
	Symbol          string              `json:"symbol"`
	IncomeStatement *FinancialStatement `json:"incomeStatement"`
	BalanceSheet    *FinancialStatement `json:"balanceSheet"`
	CashFlow        *FinancialStatement `json:"cashFlow"`
}

// fetchFinancialStatement is a helper to fetch and parse financial statements
func (t *Ticker) fetchFinancialStatement(ctx context.Context, module, annualKey, quarterlyKey string, quarterly bool) (*FinancialStatement, error) {
	endpoint := fmt.Sprintf("%s/%s", QuoteSummaryURL, t.Symbol)
	params := buildModulesParams(module)

	data, err := t.client.Get(ctx, endpoint, params)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		QuoteSummary struct {
			Result []map[string]json.RawMessage `json:"result"`
		} `json:"quoteSummary"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse financial statement: %w", err))
	}

	if len(response.QuoteSummary.Result) == 0 {
		return nil, NewSymbolError(t.Symbol, ErrNoData)
	}

	result := response.QuoteSummary.Result[0]

	// Get the appropriate key based on quarterly flag
	key := annualKey
	if quarterly {
		key = quarterlyKey
	}

	rawData, ok := result[key]
	if !ok {
		return nil, NewSymbolError(t.Symbol, ErrNoData)
	}

	// Try different key names based on statement type
	var parsed map[string]json.RawMessage
	if err := json.Unmarshal(rawData, &parsed); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse financial data: %w", err))
	}

	// Find the statements array (different key names for different statement types)
	var statementsRaw json.RawMessage
	for _, possibleKey := range []string{"incomeStatementHistory", "balanceSheetStatements", "cashflowStatements"} {
		if raw, ok := parsed[possibleKey]; ok {
			statementsRaw = raw
			break
		}
	}

	if statementsRaw == nil {
		return nil, NewSymbolError(t.Symbol, ErrNoData)
	}

	var statements []map[string]interface{}
	if err := json.Unmarshal(statementsRaw, &statements); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse statements: %w", err))
	}

	fs := &FinancialStatement{Symbol: t.Symbol}
	periods := make([]FinancialStatementPeriod, 0, len(statements))

	for _, stmt := range statements {
		period := FinancialStatementPeriod{
			Data: make(map[string]float64),
		}

		// Parse end date
		if endDate, ok := stmt["endDate"].(map[string]interface{}); ok {
			if raw, ok := endDate["raw"].(float64); ok {
				period.Date = time.Unix(int64(raw), 0)
			}
			if fmt, ok := endDate["fmt"].(string); ok {
				period.EndDate = fmt
			}
		}

		// Parse all numeric fields
		for key, value := range stmt {
			if key == "endDate" || key == "maxAge" {
				continue
			}
			if valMap, ok := value.(map[string]interface{}); ok {
				if raw, ok := valMap["raw"].(float64); ok {
					period.Data[key] = raw
				}
			}
		}

		periods = append(periods, period)
	}

	if quarterly {
		fs.Quarterly = periods
	} else {
		fs.Annual = periods
	}

	return fs, nil
}

// GetFinancialMetric extracts a specific metric from financial statements
func GetFinancialMetric(fs *FinancialStatement, metric string, quarterly bool) []float64 {
	var periods []FinancialStatementPeriod
	if quarterly {
		periods = fs.Quarterly
	} else {
		periods = fs.Annual
	}

	values := make([]float64, 0, len(periods))
	for _, period := range periods {
		if val, ok := period.Data[metric]; ok {
			values = append(values, val)
		}
	}
	return values
}
