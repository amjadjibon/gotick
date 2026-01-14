package yfinance

// IncomeStatementKeys defines available fields for income statement data.
// These keys are used when querying financial fundamentals from Yahoo Finance.
var IncomeStatementKeys = []string{
	"TotalRevenue",
	"OperatingRevenue",
	"CostOfRevenue",
	"GrossProfit",
	"OperatingExpense",
	"SellingGeneralAndAdministration",
	"ResearchAndDevelopment",
	"OperatingIncome",
	"NetNonOperatingInterestIncomeExpense",
	"InterestIncomeNonOperating",
	"InterestExpenseNonOperating",
	"OtherIncomeExpense",
	"PretaxIncome",
	"TaxProvision",
	"NetIncome",
	"NetIncomeCommonStockholders",
	"DilutedEPS",
	"BasicEPS",
	"DilutedAverageShares",
	"BasicAverageShares",
	"EBITDA",
	"EBIT",
	"NormalizedEBITDA",
	"NormalizedIncome",
	"ReconciledDepreciation",
	"ReconciledCostOfRevenue",
}

// BalanceSheetKeys defines available fields for balance sheet data.
// These keys are used when querying financial fundamentals from Yahoo Finance.
var BalanceSheetKeys = []string{
	"TotalAssets",
	"CurrentAssets",
	"CashAndCashEquivalents",
	"CashCashEquivalentsAndShortTermInvestments",
	"Receivables",
	"AccountsReceivable",
	"Inventory",
	"PrepaidAssets",
	"OtherCurrentAssets",
	"TotalNonCurrentAssets",
	"NetPPE",
	"GrossPPE",
	"AccumulatedDepreciation",
	"GoodwillAndOtherIntangibleAssets",
	"Goodwill",
	"OtherIntangibleAssets",
	"InvestmentsAndAdvances",
	"LongTermEquityInvestment",
	"OtherNonCurrentAssets",
	"TotalLiabilitiesNetMinorityInterest",
	"CurrentLiabilities",
	"AccountsPayable",
	"PayablesAndAccruedExpenses",
	"CurrentDebtAndCapitalLeaseObligation",
	"CurrentDebt",
	"OtherCurrentLiabilities",
	"TotalNonCurrentLiabilitiesNetMinorityInterest",
	"LongTermDebtAndCapitalLeaseObligation",
	"LongTermDebt",
	"OtherNonCurrentLiabilities",
	"TotalEquityGrossMinorityInterest",
	"StockholdersEquity",
	"CommonStockEquity",
	"RetainedEarnings",
	"AdditionalPaidInCapital",
	"TreasuryStock",
	"TotalDebt",
	"NetDebt",
	"WorkingCapital",
	"TangibleBookValue",
	"InvestedCapital",
}

// CashFlowKeys defines available fields for cash flow statement data.
// These keys are used when querying financial fundamentals from Yahoo Finance.
var CashFlowKeys = []string{
	// Operating Activities
	"OperatingCashFlow",
	"NetIncomeFromContinuingOperations",
	"DepreciationAmortizationDepletion",
	"Depreciation",
	"AmortizationOfIntangibles",
	"DeferredIncomeTax",
	"DeferredTax",
	"StockBasedCompensation",
	"ChangeInWorkingCapital",
	"ChangeInReceivables",
	"ChangeInInventory",
	"ChangeInPayablesAndAccruedExpense",
	"ChangeInAccountPayable",
	"ChangeInOtherWorkingCapital",
	"OtherNonCashItems",
	// Investing Activities
	"InvestingCashFlow",
	"CapitalExpenditure",
	"PurchaseOfPPE",
	"PurchaseOfInvestment",
	"SaleOfInvestment",
	"NetInvestmentPurchaseAndSale",
	"PurchaseOfBusiness",
	"NetBusinessPurchaseAndSale",
	"NetIntangiblesPurchaseAndSale",
	"OtherInvestingChanges",
	// Financing Activities
	"FinancingCashFlow",
	"NetIssuancePaymentsOfDebt",
	"NetLongTermDebtIssuance",
	"LongTermDebtIssuance",
	"LongTermDebtPayments",
	"NetShortTermDebtIssuance",
	"NetCommonStockIssuance",
	"CommonStockIssuance",
	"CommonStockPayments",
	"RepurchaseOfCapitalStock",
	"CashDividendsPaid",
	"CommonStockDividendPaid",
	"NetOtherFinancingCharges",
	"FreeCashFlow",
	// Cash Position
	"BeginningCashPosition",
	"EndCashPosition",
	"ChangesInCash",
	"EffectOfExchangeRateChanges",
}

// AllFinancialKeys returns all financial statement keys combined.
func AllFinancialKeys() []string {
	result := make([]string, 0, len(IncomeStatementKeys)+len(BalanceSheetKeys)+len(CashFlowKeys))
	result = append(result, IncomeStatementKeys...)
	result = append(result, BalanceSheetKeys...)
	result = append(result, CashFlowKeys...)
	return result
}
