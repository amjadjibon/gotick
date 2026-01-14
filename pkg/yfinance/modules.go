package yfinance

// QuoteSummaryModules defines all available modules for the quoteSummary endpoint.
// Each module provides different types of information about a security.
var QuoteSummaryModules = []string{
	// Profile modules
	"summaryProfile",
	"summaryDetail",
	"assetProfile",
	"fundProfile",
	"price",
	"quoteType",
	"esgScores",

	// Financial statement modules
	"incomeStatementHistory",
	"incomeStatementHistoryQuarterly",
	"balanceSheetHistory",
	"balanceSheetHistoryQuarterly",
	"cashFlowStatementHistory",
	"cashFlowStatementHistoryQuarterly",

	// Key statistics and financial data
	"defaultKeyStatistics",
	"financialData",

	// Calendar and events
	"calendarEvents",
	"secFilings",

	// Analyst data
	"upgradeDowngradeHistory",
	"recommendationTrend",

	// Ownership modules
	"institutionOwnership",
	"fundOwnership",
	"majorDirectHolders",
	"majorHoldersBreakdown",
	"insiderTransactions",
	"insiderHolders",
	"netSharePurchaseActivity",

	// Earnings modules
	"earnings",
	"earningsHistory",
	"earningsTrend",

	// Trend modules
	"industryTrend",
	"indexTrend",
	"sectorTrend",

	// Other
	"futuresChain",
}

// Module category constants for easier grouping
const (
	// Profile modules
	ModuleSummaryProfile = "summaryProfile"
	ModuleSummaryDetail  = "summaryDetail"
	ModuleAssetProfile   = "assetProfile"
	ModuleFundProfile    = "fundProfile"
	ModulePrice          = "price"
	ModuleQuoteType      = "quoteType"
	ModuleESGScores      = "esgScores"

	// Financial statement modules
	ModuleIncomeStatementHistory            = "incomeStatementHistory"
	ModuleIncomeStatementHistoryQuarterly   = "incomeStatementHistoryQuarterly"
	ModuleBalanceSheetHistory               = "balanceSheetHistory"
	ModuleBalanceSheetHistoryQuarterly      = "balanceSheetHistoryQuarterly"
	ModuleCashFlowStatementHistory          = "cashFlowStatementHistory"
	ModuleCashFlowStatementHistoryQuarterly = "cashFlowStatementHistoryQuarterly"

	// Key statistics and financial data
	ModuleDefaultKeyStatistics = "defaultKeyStatistics"
	ModuleFinancialData        = "financialData"

	// Calendar and events
	ModuleCalendarEvents = "calendarEvents"
	ModuleSecFilings     = "secFilings"

	// Analyst data
	ModuleUpgradeDowngradeHistory = "upgradeDowngradeHistory"
	ModuleRecommendationTrend     = "recommendationTrend"

	// Ownership modules
	ModuleInstitutionOwnership     = "institutionOwnership"
	ModuleFundOwnership            = "fundOwnership"
	ModuleMajorDirectHolders       = "majorDirectHolders"
	ModuleMajorHoldersBreakdown    = "majorHoldersBreakdown"
	ModuleInsiderTransactions      = "insiderTransactions"
	ModuleInsiderHolders           = "insiderHolders"
	ModuleNetSharePurchaseActivity = "netSharePurchaseActivity"

	// Earnings modules
	ModuleEarnings        = "earnings"
	ModuleEarningsHistory = "earningsHistory"
	ModuleEarningsTrend   = "earningsTrend"

	// Trend modules
	ModuleIndustryTrend = "industryTrend"
	ModuleIndexTrend    = "indexTrend"
	ModuleSectorTrend   = "sectorTrend"

	// Other
	ModuleFuturesChain = "futuresChain"
)

// DefaultModules returns a commonly used set of modules for basic info.
func DefaultModules() []string {
	return []string{
		ModuleSummaryProfile,
		ModuleSummaryDetail,
		ModulePrice,
		ModuleQuoteType,
		ModuleDefaultKeyStatistics,
		ModuleFinancialData,
	}
}

// FinancialModules returns all financial statement related modules.
func FinancialModules() []string {
	return []string{
		ModuleIncomeStatementHistory,
		ModuleIncomeStatementHistoryQuarterly,
		ModuleBalanceSheetHistory,
		ModuleBalanceSheetHistoryQuarterly,
		ModuleCashFlowStatementHistory,
		ModuleCashFlowStatementHistoryQuarterly,
	}
}

// OwnershipModules returns all ownership related modules.
func OwnershipModules() []string {
	return []string{
		ModuleInstitutionOwnership,
		ModuleFundOwnership,
		ModuleMajorDirectHolders,
		ModuleMajorHoldersBreakdown,
		ModuleInsiderTransactions,
		ModuleInsiderHolders,
		ModuleNetSharePurchaseActivity,
	}
}

// EarningsModules returns all earnings related modules.
func EarningsModules() []string {
	return []string{
		ModuleEarnings,
		ModuleEarningsHistory,
		ModuleEarningsTrend,
	}
}
