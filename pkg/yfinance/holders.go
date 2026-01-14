package yfinance

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// MajorHolders represents the major holders breakdown
type MajorHolders struct {
	InsidersPercentHeld          float64 `json:"insidersPercentHeld"`
	InstitutionsPercentHeld      float64 `json:"institutionsPercentHeld"`
	InstitutionsFloatPercentHeld float64 `json:"institutionsFloatPercentHeld"`
	InstitutionsCount            int     `json:"institutionsCount"`
}

// Holder represents an institutional or fund holder
type Holder struct {
	Holder       string    `json:"holder"`
	Shares       int64     `json:"shares"`
	DateReported time.Time `json:"dateReported"`
	PercentOut   float64   `json:"pctOut"`
	Value        int64     `json:"value"`
	PctHeld      float64   `json:"pctHeld"`
}

// InsiderTransaction represents an insider transaction
type InsiderTransaction struct {
	Insider     string    `json:"filerName"`
	Relation    string    `json:"filerRelation"`
	URL         string    `json:"filerUrl"`
	Transaction string    `json:"transactionText"`
	Shares      int64     `json:"shares"`
	Value       int64     `json:"value"`
	StartDate   time.Time `json:"startDate"`
	Ownership   string    `json:"ownership"`
}

// InsiderHolder represents an insider holder
type InsiderHolder struct {
	Name                   string    `json:"name"`
	Relation               string    `json:"relation"`
	URL                    string    `json:"url"`
	TransactionDescription string    `json:"transactionDescription"`
	LatestTransDate        time.Time `json:"latestTransDate"`
	PositionDirect         int64     `json:"positionDirect"`
	PositionDirectDate     time.Time `json:"positionDirectDate"`
	PositionIndirect       int64     `json:"positionIndirect"`
	PositionIndirectDate   time.Time `json:"positionIndirectDate"`
}

// InsiderPurchases represents insider purchase activity summary
type InsiderPurchases struct {
	Purchases          int64   `json:"purchases"`
	Sales              int64   `json:"sales"`
	NetSharesPurchased int64   `json:"netSharesPurchased"`
	TotalInsiderShares int64   `json:"totalInsiderShares"`
	PercentNetShares   float64 `json:"pctNetShares"`
	PercentBuyShares   float64 `json:"pctBuyShares"`
	PercentSellShares  float64 `json:"pctSellShares"`
}

// MajorHolders fetches major holders breakdown
func (t *Ticker) MajorHolders(ctx context.Context) (*MajorHolders, error) {
	endpoint := fmt.Sprintf("%s/%s", QuoteSummaryURL, t.Symbol)
	params := buildModulesParams(ModuleMajorHoldersBreakdown)

	data, err := t.client.Get(ctx, endpoint, params)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		QuoteSummary struct {
			Result []struct {
				MajorHoldersBreakdown struct {
					InsidersPercentHeld          RawValue `json:"insidersPercentHeld"`
					InstitutionsPercentHeld      RawValue `json:"institutionsPercentHeld"`
					InstitutionsFloatPercentHeld RawValue `json:"institutionsFloatPercentHeld"`
					InstitutionsCount            RawValue `json:"institutionsCount"`
				} `json:"majorHoldersBreakdown"`
			} `json:"result"`
		} `json:"quoteSummary"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse major holders: %w", err))
	}

	if len(response.QuoteSummary.Result) == 0 {
		return nil, NewSymbolError(t.Symbol, ErrNoData)
	}

	mh := response.QuoteSummary.Result[0].MajorHoldersBreakdown
	return &MajorHolders{
		InsidersPercentHeld:          mh.InsidersPercentHeld.Raw,
		InstitutionsPercentHeld:      mh.InstitutionsPercentHeld.Raw,
		InstitutionsFloatPercentHeld: mh.InstitutionsFloatPercentHeld.Raw,
		InstitutionsCount:            int(mh.InstitutionsCount.Raw),
	}, nil
}

// InstitutionalHolders fetches institutional holders
func (t *Ticker) InstitutionalHolders(ctx context.Context) ([]Holder, error) {
	endpoint := fmt.Sprintf("%s/%s", QuoteSummaryURL, t.Symbol)
	params := buildModulesParams(ModuleInstitutionOwnership)

	data, err := t.client.Get(ctx, endpoint, params)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		QuoteSummary struct {
			Result []struct {
				InstitutionOwnership struct {
					OwnershipList []struct {
						Organization string   `json:"organization"`
						PctHeld      RawValue `json:"pctHeld"`
						Position     RawValue `json:"position"`
						Value        RawValue `json:"value"`
						ReportDate   RawValue `json:"reportDate"`
					} `json:"ownershipList"`
				} `json:"institutionOwnership"`
			} `json:"result"`
		} `json:"quoteSummary"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse institutional holders: %w", err))
	}

	if len(response.QuoteSummary.Result) == 0 {
		return nil, NewSymbolError(t.Symbol, ErrNoData)
	}

	var holders []Holder
	for _, h := range response.QuoteSummary.Result[0].InstitutionOwnership.OwnershipList {
		holders = append(holders, Holder{
			Holder:       h.Organization,
			Shares:       int64(h.Position.Raw),
			Value:        int64(h.Value.Raw),
			PctHeld:      h.PctHeld.Raw,
			DateReported: time.Unix(int64(h.ReportDate.Raw), 0),
		})
	}

	return holders, nil
}

// MutualFundHolders fetches mutual fund holders
func (t *Ticker) MutualFundHolders(ctx context.Context) ([]Holder, error) {
	endpoint := fmt.Sprintf("%s/%s", QuoteSummaryURL, t.Symbol)
	params := buildModulesParams(ModuleFundOwnership)

	data, err := t.client.Get(ctx, endpoint, params)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		QuoteSummary struct {
			Result []struct {
				FundOwnership struct {
					OwnershipList []struct {
						Organization string   `json:"organization"`
						PctHeld      RawValue `json:"pctHeld"`
						Position     RawValue `json:"position"`
						Value        RawValue `json:"value"`
						ReportDate   RawValue `json:"reportDate"`
					} `json:"ownershipList"`
				} `json:"fundOwnership"`
			} `json:"result"`
		} `json:"quoteSummary"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse mutual fund holders: %w", err))
	}

	if len(response.QuoteSummary.Result) == 0 {
		return nil, NewSymbolError(t.Symbol, ErrNoData)
	}

	var holders []Holder
	for _, h := range response.QuoteSummary.Result[0].FundOwnership.OwnershipList {
		holders = append(holders, Holder{
			Holder:       h.Organization,
			Shares:       int64(h.Position.Raw),
			Value:        int64(h.Value.Raw),
			PctHeld:      h.PctHeld.Raw,
			DateReported: time.Unix(int64(h.ReportDate.Raw), 0),
		})
	}

	return holders, nil
}

// InsiderTransactions fetches insider transactions
func (t *Ticker) InsiderTransactions(ctx context.Context) ([]InsiderTransaction, error) {
	endpoint := fmt.Sprintf("%s/%s", QuoteSummaryURL, t.Symbol)
	params := buildModulesParams(ModuleInsiderTransactions)

	data, err := t.client.Get(ctx, endpoint, params)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		QuoteSummary struct {
			Result []struct {
				InsiderTransactions struct {
					Transactions []struct {
						FilerName       string   `json:"filerName"`
						FilerRelation   string   `json:"filerRelation"`
						FilerURL        string   `json:"filerUrl"`
						TransactionText string   `json:"transactionText"`
						Shares          RawValue `json:"shares"`
						Value           RawValue `json:"value"`
						StartDate       RawValue `json:"startDate"`
						Ownership       string   `json:"ownership"`
					} `json:"transactions"`
				} `json:"insiderTransactions"`
			} `json:"result"`
		} `json:"quoteSummary"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse insider transactions: %w", err))
	}

	if len(response.QuoteSummary.Result) == 0 {
		return nil, NewSymbolError(t.Symbol, ErrNoData)
	}

	var transactions []InsiderTransaction
	for _, tx := range response.QuoteSummary.Result[0].InsiderTransactions.Transactions {
		transactions = append(transactions, InsiderTransaction{
			Insider:     tx.FilerName,
			Relation:    tx.FilerRelation,
			URL:         tx.FilerURL,
			Transaction: tx.TransactionText,
			Shares:      int64(tx.Shares.Raw),
			Value:       int64(tx.Value.Raw),
			StartDate:   time.Unix(int64(tx.StartDate.Raw), 0),
			Ownership:   tx.Ownership,
		})
	}

	return transactions, nil
}

// InsiderRosterHolders fetches insider roster holders
func (t *Ticker) InsiderRosterHolders(ctx context.Context) ([]InsiderHolder, error) {
	endpoint := fmt.Sprintf("%s/%s", QuoteSummaryURL, t.Symbol)
	params := buildModulesParams(ModuleInsiderHolders)

	data, err := t.client.Get(ctx, endpoint, params)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		QuoteSummary struct {
			Result []struct {
				InsiderHolders struct {
					Holders []struct {
						Name                   string   `json:"name"`
						Relation               string   `json:"relation"`
						URL                    string   `json:"url"`
						TransactionDescription string   `json:"transactionDescription"`
						LatestTransDate        RawValue `json:"latestTransDate"`
						PositionDirect         RawValue `json:"positionDirect"`
						PositionDirectDate     RawValue `json:"positionDirectDate"`
						PositionIndirect       RawValue `json:"positionIndirect"`
						PositionIndirectDate   RawValue `json:"positionIndirectDate"`
					} `json:"holders"`
				} `json:"insiderHolders"`
			} `json:"result"`
		} `json:"quoteSummary"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse insider roster holders: %w", err))
	}

	if len(response.QuoteSummary.Result) == 0 {
		return nil, NewSymbolError(t.Symbol, ErrNoData)
	}

	var holders []InsiderHolder
	for _, h := range response.QuoteSummary.Result[0].InsiderHolders.Holders {
		holders = append(holders, InsiderHolder{
			Name:                   h.Name,
			Relation:               h.Relation,
			URL:                    h.URL,
			TransactionDescription: h.TransactionDescription,
			LatestTransDate:        time.Unix(int64(h.LatestTransDate.Raw), 0),
			PositionDirect:         int64(h.PositionDirect.Raw),
			PositionDirectDate:     time.Unix(int64(h.PositionDirectDate.Raw), 0),
			PositionIndirect:       int64(h.PositionIndirect.Raw),
			PositionIndirectDate:   time.Unix(int64(h.PositionIndirectDate.Raw), 0),
		})
	}

	return holders, nil
}

// InsiderPurchasesData fetches insider purchase activity summary
func (t *Ticker) InsiderPurchasesData(ctx context.Context) (*InsiderPurchases, error) {
	endpoint := fmt.Sprintf("%s/%s", QuoteSummaryURL, t.Symbol)
	params := buildModulesParams(ModuleNetSharePurchaseActivity)

	data, err := t.client.Get(ctx, endpoint, params)
	if err != nil {
		return nil, NewSymbolError(t.Symbol, err)
	}

	var response struct {
		QuoteSummary struct {
			Result []struct {
				NetSharePurchaseActivity struct {
					Purchases                RawValue `json:"buyInfoShares"`
					Sales                    RawValue `json:"sellInfoShares"`
					NetSharesPurchased       RawValue `json:"netPercentInsiderShares"`
					TotalInsiderShares       RawValue `json:"totalInsiderShares"`
					BuyPercentInsiderShares  RawValue `json:"buyPercentInsiderShares"`
					SellPercentInsiderShares RawValue `json:"sellPercentInsiderShares"`
					NetInfoShares            RawValue `json:"netInfoShares"`
				} `json:"netSharePurchaseActivity"`
			} `json:"result"`
		} `json:"quoteSummary"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, NewSymbolError(t.Symbol, fmt.Errorf("failed to parse insider purchases: %w", err))
	}

	if len(response.QuoteSummary.Result) == 0 {
		return nil, NewSymbolError(t.Symbol, ErrNoData)
	}

	ip := response.QuoteSummary.Result[0].NetSharePurchaseActivity
	return &InsiderPurchases{
		Purchases:          int64(ip.Purchases.Raw),
		Sales:              int64(ip.Sales.Raw),
		NetSharesPurchased: int64(ip.NetInfoShares.Raw),
		TotalInsiderShares: int64(ip.TotalInsiderShares.Raw),
		PercentNetShares:   ip.NetSharesPurchased.Raw,
		PercentBuyShares:   ip.BuyPercentInsiderShares.Raw,
		PercentSellShares:  ip.SellPercentInsiderShares.Raw,
	}, nil
}
