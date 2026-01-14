// Package yfinance provides a Go client for Yahoo Finance APIs.
package yfinance

// Base URLs for Yahoo Finance API endpoints
const (
	// BaseURL is the primary Yahoo Finance API endpoint
	BaseURL = "https://query2.finance.yahoo.com"
	// Query1URL is an alternative Yahoo Finance API endpoint
	Query1URL = "https://query1.finance.yahoo.com"
	// RootURL is the main Yahoo Finance website
	RootURL = "https://finance.yahoo.com"
)

// Authentication endpoints
const (
	// CookieURL is used to obtain authentication cookies
	CookieURL = "https://fc.yahoo.com"
	// CrumbURL is used to get the CSRF crumb token
	CrumbURL = BaseURL + "/v1/test/getcrumb"
)

// CSRF Consent endpoints (fallback authentication)
const (
	// ConsentURL is the Yahoo GUCE consent page
	ConsentURL = "https://guce.yahoo.com/consent"
	// CollectConsentURL is used to collect user consent
	CollectConsentURL = "https://consent.yahoo.com/v2/collectConsent"
	// CopyConsentURL is used to copy consent across domains
	CopyConsentURL = "https://guce.yahoo.com/copyConsent"
	// CrumbCSRFURL is an alternative crumb endpoint for CSRF
	CrumbCSRFURL = BaseURL + "/v1/test/getcrumb"
)

// Data API endpoints
const (
	// ChartURL provides historical chart/OHLCV data
	ChartURL = BaseURL + "/v8/finance/chart"
	// QuoteSummaryURL provides comprehensive quote information
	QuoteSummaryURL = BaseURL + "/v10/finance/quoteSummary"
	// QuoteURL provides real-time quote data
	QuoteURL = Query1URL + "/v7/finance/quote"
	// OptionsURL provides options chain data
	OptionsURL = BaseURL + "/v7/finance/options"
	// FundamentalsURL provides fundamental financial timeseries data
	FundamentalsURL = BaseURL + "/ws/fundamentals-timeseries/v1/finance/timeseries"
)

// Search and Discovery endpoints
const (
	// SearchURL provides symbol and company search
	SearchURL = BaseURL + "/v1/finance/search"
	// LookupURL provides symbol lookup functionality
	LookupURL = Query1URL + "/v1/finance/lookup"
	// ScreenerURL provides stock screening functionality
	ScreenerURL = Query1URL + "/v1/finance/screener"
)

// Market Data endpoints
const (
	// MarketSummaryURL provides market overview data
	MarketSummaryURL = Query1URL + "/v6/finance/quote/marketSummary"
	// MarketTimeURL provides market trading hours information
	MarketTimeURL = Query1URL + "/v6/finance/markettime"
)

// Sector and Industry endpoints
const (
	// SectorURL provides sector information
	SectorURL = Query1URL + "/v1/finance/sectors"
	// IndustryURL provides industry information
	IndustryURL = Query1URL + "/v1/finance/industries"
)

// Calendar endpoints
const (
	// CalendarURL provides calendar events (earnings, IPO, splits, etc.)
	CalendarURL = Query1URL + "/v1/finance/visualization"
)

// News endpoints
const (
	// NewsURL provides financial news
	NewsURL = RootURL + "/xhr/ncp"
)

// Real-time data endpoints
const (
	// WebSocketURL provides real-time streaming quotes
	WebSocketURL = "wss://streamer.finance.yahoo.com/?version=2"
)
