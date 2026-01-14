package yfinance

import (
	"time"
)

// Interval represents the time interval for chart data
type Interval string

// Available intervals for chart data
const (
	Interval1m  Interval = "1m"
	Interval2m  Interval = "2m"
	Interval5m  Interval = "5m"
	Interval15m Interval = "15m"
	Interval30m Interval = "30m"
	Interval60m Interval = "60m"
	Interval90m Interval = "90m"
	Interval1h  Interval = "1h"
	Interval1d  Interval = "1d"
	Interval5d  Interval = "5d"
	Interval1wk Interval = "1wk"
	Interval1mo Interval = "1mo"
	Interval3mo Interval = "3mo"
)

// Period represents the time period for historical data
type Period string

// Available periods for historical data
const (
	Period1d  Period = "1d"
	Period5d  Period = "5d"
	Period1mo Period = "1mo"
	Period3mo Period = "3mo"
	Period6mo Period = "6mo"
	Period1y  Period = "1y"
	Period2y  Period = "2y"
	Period5y  Period = "5y"
	Period10y Period = "10y"
	PeriodYTD Period = "ytd"
	PeriodMax Period = "max"
)

// Quote represents real-time quote data for a security
type Quote struct {
	Symbol                     string  `json:"symbol"`
	ShortName                  string  `json:"shortName"`
	LongName                   string  `json:"longName"`
	Exchange                   string  `json:"exchange"`
	FullExchangeName           string  `json:"fullExchangeName"`
	QuoteType                  string  `json:"quoteType"`
	Currency                   string  `json:"currency"`
	MarketState                string  `json:"marketState"`
	RegularMarketPrice         float64 `json:"regularMarketPrice"`
	RegularMarketChange        float64 `json:"regularMarketChange"`
	RegularMarketChangePercent float64 `json:"regularMarketChangePercent"`
	RegularMarketOpen          float64 `json:"regularMarketOpen"`
	RegularMarketDayHigh       float64 `json:"regularMarketDayHigh"`
	RegularMarketDayLow        float64 `json:"regularMarketDayLow"`
	RegularMarketVolume        int64   `json:"regularMarketVolume"`
	RegularMarketPreviousClose float64 `json:"regularMarketPreviousClose"`
	RegularMarketTime          int64   `json:"regularMarketTime"`
	PreMarketPrice             float64 `json:"preMarketPrice,omitempty"`
	PreMarketChange            float64 `json:"preMarketChange,omitempty"`
	PreMarketChangePercent     float64 `json:"preMarketChangePercent,omitempty"`
	PreMarketTime              int64   `json:"preMarketTime,omitempty"`
	PostMarketPrice            float64 `json:"postMarketPrice,omitempty"`
	PostMarketChange           float64 `json:"postMarketChange,omitempty"`
	PostMarketChangePercent    float64 `json:"postMarketChangePercent,omitempty"`
	PostMarketTime             int64   `json:"postMarketTime,omitempty"`
	Bid                        float64 `json:"bid"`
	BidSize                    int64   `json:"bidSize"`
	Ask                        float64 `json:"ask"`
	AskSize                    int64   `json:"askSize"`
	FiftyTwoWeekHigh           float64 `json:"fiftyTwoWeekHigh"`
	FiftyTwoWeekLow            float64 `json:"fiftyTwoWeekLow"`
	FiftyTwoWeekHighChange     float64 `json:"fiftyTwoWeekHighChange"`
	FiftyTwoWeekLowChange      float64 `json:"fiftyTwoWeekLowChange"`
	FiftyDayAverage            float64 `json:"fiftyDayAverage"`
	TwoHundredDayAverage       float64 `json:"twoHundredDayAverage"`
	MarketCap                  int64   `json:"marketCap"`
	TrailingPE                 float64 `json:"trailingPE"`
	ForwardPE                  float64 `json:"forwardPE"`
	DividendYield              float64 `json:"dividendYield"`
	DividendRate               float64 `json:"dividendRate"`
	EpsTrailingTwelveMonths    float64 `json:"epsTrailingTwelveMonths"`
	EpsForward                 float64 `json:"epsForward"`
	SharesOutstanding          int64   `json:"sharesOutstanding"`
	AverageDailyVolume3Month   int64   `json:"averageDailyVolume3Month"`
	AverageDailyVolume10Day    int64   `json:"averageDailyVolume10Day"`
}

// Bar represents a single OHLCV bar
type Bar struct {
	Timestamp time.Time `json:"timestamp"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	AdjClose  float64   `json:"adjClose"`
	Volume    int64     `json:"volume"`
}

// ChartData represents historical chart data
type ChartData struct {
	Symbol   string     `json:"symbol"`
	Currency string     `json:"currency"`
	Interval Interval   `json:"interval"`
	Bars     []Bar      `json:"bars"`
	Meta     *ChartMeta `json:"meta,omitempty"`
}

// ChartMeta contains metadata about chart data
type ChartMeta struct {
	Currency             string  `json:"currency"`
	ExchangeName         string  `json:"exchangeName"`
	InstrumentType       string  `json:"instrumentType"`
	FirstTradeDate       int64   `json:"firstTradeDate"`
	RegularMarketTime    int64   `json:"regularMarketTime"`
	GMTOffset            int     `json:"gmtoffset"`
	Timezone             string  `json:"timezone"`
	ExchangeTimezoneName string  `json:"exchangeTimezoneName"`
	RegularMarketPrice   float64 `json:"regularMarketPrice"`
	ChartPreviousClose   float64 `json:"chartPreviousClose"`
	PriceHint            int     `json:"priceHint"`
	DataGranularity      string  `json:"dataGranularity"`
	Range                string  `json:"range"`
}

// HistoryParams defines parameters for fetching historical data
type HistoryParams struct {
	Period   Period    `json:"period,omitempty"`
	Interval Interval  `json:"interval,omitempty"`
	Start    time.Time `json:"start,omitempty"`
	End      time.Time `json:"end,omitempty"`
	PrePost  bool      `json:"prepost,omitempty"`
	Events   string    `json:"events,omitempty"` // "div", "split", "div,split"
}

// QuoteSummary represents comprehensive quote information
type QuoteSummary struct {
	Symbol         string          `json:"symbol"`
	AssetProfile   *AssetProfile   `json:"assetProfile,omitempty"`
	SummaryProfile *SummaryProfile `json:"summaryProfile,omitempty"`
	SummaryDetail  *SummaryDetail  `json:"summaryDetail,omitempty"`
	Price          *PriceInfo      `json:"price,omitempty"`
	KeyStatistics  *KeyStatistics  `json:"defaultKeyStatistics,omitempty"`
	FinancialData  *FinancialData  `json:"financialData,omitempty"`
	CalendarEvents *CalendarEvents `json:"calendarEvents,omitempty"`
}

// AssetProfile contains company profile information
type AssetProfile struct {
	Address1            string    `json:"address1"`
	Address2            string    `json:"address2,omitempty"`
	City                string    `json:"city"`
	State               string    `json:"state,omitempty"`
	Zip                 string    `json:"zip"`
	Country             string    `json:"country"`
	Phone               string    `json:"phone"`
	Website             string    `json:"website"`
	Industry            string    `json:"industry"`
	Sector              string    `json:"sector"`
	LongBusinessSummary string    `json:"longBusinessSummary"`
	FullTimeEmployees   int       `json:"fullTimeEmployees"`
	CompanyOfficers     []Officer `json:"companyOfficers"`
}

// Officer represents a company officer
type Officer struct {
	Name             string `json:"name"`
	Title            string `json:"title"`
	Age              int    `json:"age,omitempty"`
	YearBorn         int    `json:"yearBorn,omitempty"`
	TotalPay         int64  `json:"totalPay,omitempty"`
	ExercisedValue   int64  `json:"exercisedValue,omitempty"`
	UnexercisedValue int64  `json:"unexercisedValue,omitempty"`
}

// SummaryProfile contains summary profile information
type SummaryProfile struct {
	Address1            string `json:"address1"`
	City                string `json:"city"`
	State               string `json:"state,omitempty"`
	Zip                 string `json:"zip"`
	Country             string `json:"country"`
	Phone               string `json:"phone"`
	Website             string `json:"website"`
	Industry            string `json:"industry"`
	Sector              string `json:"sector"`
	LongBusinessSummary string `json:"longBusinessSummary"`
	FullTimeEmployees   int    `json:"fullTimeEmployees"`
}

// SummaryDetail contains summary detail information
type SummaryDetail struct {
	PreviousClose               float64 `json:"previousClose"`
	Open                        float64 `json:"open"`
	DayLow                      float64 `json:"dayLow"`
	DayHigh                     float64 `json:"dayHigh"`
	RegularMarketPreviousClose  float64 `json:"regularMarketPreviousClose"`
	RegularMarketOpen           float64 `json:"regularMarketOpen"`
	RegularMarketDayLow         float64 `json:"regularMarketDayLow"`
	RegularMarketDayHigh        float64 `json:"regularMarketDayHigh"`
	DividendRate                float64 `json:"dividendRate"`
	DividendYield               float64 `json:"dividendYield"`
	PayoutRatio                 float64 `json:"payoutRatio"`
	FiveYearAvgDividendYield    float64 `json:"fiveYearAvgDividendYield"`
	Beta                        float64 `json:"beta"`
	TrailingPE                  float64 `json:"trailingPE"`
	ForwardPE                   float64 `json:"forwardPE"`
	Volume                      int64   `json:"volume"`
	RegularMarketVolume         int64   `json:"regularMarketVolume"`
	AverageVolume               int64   `json:"averageVolume"`
	AverageVolume10days         int64   `json:"averageVolume10days"`
	AverageDailyVolume10Day     int64   `json:"averageDailyVolume10Day"`
	Bid                         float64 `json:"bid"`
	Ask                         float64 `json:"ask"`
	BidSize                     int64   `json:"bidSize"`
	AskSize                     int64   `json:"askSize"`
	MarketCap                   int64   `json:"marketCap"`
	FiftyTwoWeekLow             float64 `json:"fiftyTwoWeekLow"`
	FiftyTwoWeekHigh            float64 `json:"fiftyTwoWeekHigh"`
	FiftyDayAverage             float64 `json:"fiftyDayAverage"`
	TwoHundredDayAverage        float64 `json:"twoHundredDayAverage"`
	TrailingAnnualDividendRate  float64 `json:"trailingAnnualDividendRate"`
	TrailingAnnualDividendYield float64 `json:"trailingAnnualDividendYield"`
	Currency                    string  `json:"currency"`
}

// PriceInfo contains price information
type PriceInfo struct {
	Symbol                     string  `json:"symbol"`
	ShortName                  string  `json:"shortName"`
	LongName                   string  `json:"longName"`
	Exchange                   string  `json:"exchange"`
	ExchangeName               string  `json:"exchangeName"`
	QuoteType                  string  `json:"quoteType"`
	Currency                   string  `json:"currency"`
	CurrencySymbol             string  `json:"currencySymbol"`
	MarketState                string  `json:"marketState"`
	RegularMarketPrice         float64 `json:"regularMarketPrice"`
	RegularMarketChange        float64 `json:"regularMarketChange"`
	RegularMarketChangePercent float64 `json:"regularMarketChangePercent"`
	RegularMarketOpen          float64 `json:"regularMarketOpen"`
	RegularMarketDayHigh       float64 `json:"regularMarketDayHigh"`
	RegularMarketDayLow        float64 `json:"regularMarketDayLow"`
	RegularMarketVolume        int64   `json:"regularMarketVolume"`
	RegularMarketPreviousClose float64 `json:"regularMarketPreviousClose"`
	RegularMarketTime          int64   `json:"regularMarketTime"`
	PreMarketPrice             float64 `json:"preMarketPrice,omitempty"`
	PreMarketChange            float64 `json:"preMarketChange,omitempty"`
	PreMarketChangePercent     float64 `json:"preMarketChangePercent,omitempty"`
	PostMarketPrice            float64 `json:"postMarketPrice,omitempty"`
	PostMarketChange           float64 `json:"postMarketChange,omitempty"`
	PostMarketChangePercent    float64 `json:"postMarketChangePercent,omitempty"`
	MarketCap                  int64   `json:"marketCap"`
}

// KeyStatistics contains key statistics
type KeyStatistics struct {
	EnterpriseValue         int64   `json:"enterpriseValue"`
	ForwardPE               float64 `json:"forwardPE"`
	ProfitMargins           float64 `json:"profitMargins"`
	FloatShares             int64   `json:"floatShares"`
	SharesOutstanding       int64   `json:"sharesOutstanding"`
	SharesShort             int64   `json:"sharesShort"`
	SharesShortPriorMonth   int64   `json:"sharesShortPriorMonth"`
	ShortRatio              float64 `json:"shortRatio"`
	ShortPercentOfFloat     float64 `json:"shortPercentOfFloat"`
	PercentInsiders         float64 `json:"heldPercentInsiders"`
	PercentInstitutions     float64 `json:"heldPercentInstitutions"`
	Beta                    float64 `json:"beta"`
	BookValue               float64 `json:"bookValue"`
	PriceToBook             float64 `json:"priceToBook"`
	EarningsQuarterlyGrowth float64 `json:"earningsQuarterlyGrowth"`
	NetIncomeToCommon       int64   `json:"netIncomeToCommon"`
	TrailingEps             float64 `json:"trailingEps"`
	ForwardEps              float64 `json:"forwardEps"`
	PegRatio                float64 `json:"pegRatio"`
	LastSplitFactor         string  `json:"lastSplitFactor"`
	LastSplitDate           int64   `json:"lastSplitDate"`
	EnterpriseToRevenue     float64 `json:"enterpriseToRevenue"`
	EnterpriseToEbitda      float64 `json:"enterpriseToEbitda"`
	FiftyTwoWeekChange      float64 `json:"52WeekChange"`
	SandP52WeekChange       float64 `json:"SandP52WeekChange"`
}

// FinancialData contains financial data
type FinancialData struct {
	CurrentPrice            float64 `json:"currentPrice"`
	TargetHighPrice         float64 `json:"targetHighPrice"`
	TargetLowPrice          float64 `json:"targetLowPrice"`
	TargetMeanPrice         float64 `json:"targetMeanPrice"`
	TargetMedianPrice       float64 `json:"targetMedianPrice"`
	RecommendationMean      float64 `json:"recommendationMean"`
	RecommendationKey       string  `json:"recommendationKey"`
	NumberOfAnalystOpinions int     `json:"numberOfAnalystOpinions"`
	TotalCash               int64   `json:"totalCash"`
	TotalCashPerShare       float64 `json:"totalCashPerShare"`
	Ebitda                  int64   `json:"ebitda"`
	TotalDebt               int64   `json:"totalDebt"`
	QuickRatio              float64 `json:"quickRatio"`
	CurrentRatio            float64 `json:"currentRatio"`
	TotalRevenue            int64   `json:"totalRevenue"`
	DebtToEquity            float64 `json:"debtToEquity"`
	RevenuePerShare         float64 `json:"revenuePerShare"`
	ReturnOnAssets          float64 `json:"returnOnAssets"`
	ReturnOnEquity          float64 `json:"returnOnEquity"`
	GrossProfits            int64   `json:"grossProfits"`
	FreeCashflow            int64   `json:"freeCashflow"`
	OperatingCashflow       int64   `json:"operatingCashflow"`
	EarningsGrowth          float64 `json:"earningsGrowth"`
	RevenueGrowth           float64 `json:"revenueGrowth"`
	GrossMargins            float64 `json:"grossMargins"`
	EbitdaMargins           float64 `json:"ebitdaMargins"`
	OperatingMargins        float64 `json:"operatingMargins"`
	ProfitMargins           float64 `json:"profitMargins"`
	FinancialCurrency       string  `json:"financialCurrency"`
}

// CalendarEvents contains calendar events
type CalendarEvents struct {
	Earnings *EarningsInfo `json:"earnings,omitempty"`
	Dividend *DividendInfo `json:"dividend,omitempty"`
}

// EarningsInfo contains earnings information
type EarningsInfo struct {
	EarningsDate    []int64 `json:"earningsDate"`
	EarningsAverage float64 `json:"earningsAverage"`
	EarningsLow     float64 `json:"earningsLow"`
	EarningsHigh    float64 `json:"earningsHigh"`
	RevenueAverage  int64   `json:"revenueAverage"`
	RevenueLow      int64   `json:"revenueLow"`
	RevenueHigh     int64   `json:"revenueHigh"`
}

// DividendInfo contains dividend information
type DividendInfo struct {
	ExDividendDate int64 `json:"exDividendDate"`
	DividendDate   int64 `json:"dividendDate"`
}

// OptionChain represents options data for a security
type OptionChain struct {
	Symbol          string    `json:"symbol"`
	UnderlyingPrice float64   `json:"underlyingPrice"`
	ExpirationDates []int64   `json:"expirationDates"`
	Strikes         []float64 `json:"strikes"`
	Calls           []Option  `json:"calls"`
	Puts            []Option  `json:"puts"`
}

// Option represents a single option contract
type Option struct {
	ContractSymbol    string  `json:"contractSymbol"`
	Strike            float64 `json:"strike"`
	Currency          string  `json:"currency"`
	LastPrice         float64 `json:"lastPrice"`
	Change            float64 `json:"change"`
	PercentChange     float64 `json:"percentChange"`
	Volume            int64   `json:"volume"`
	OpenInterest      int64   `json:"openInterest"`
	Bid               float64 `json:"bid"`
	Ask               float64 `json:"ask"`
	ContractSize      string  `json:"contractSize"`
	Expiration        int64   `json:"expiration"`
	LastTradeDate     int64   `json:"lastTradeDate"`
	ImpliedVolatility float64 `json:"impliedVolatility"`
	InTheMoney        bool    `json:"inTheMoney"`
}

// Financial represents financial statement data
type Financial struct {
	Symbol    string                      `json:"symbol"`
	Timestamp []int64                     `json:"timestamp"`
	Data      map[string][]FinancialValue `json:"data"`
}

// FinancialValue represents a single financial value
type FinancialValue struct {
	Raw           float64 `json:"raw"`
	Fmt           string  `json:"fmt"`
	ReportedValue float64 `json:"reportedValue,omitempty"`
	AsOfDate      string  `json:"asOfDate,omitempty"`
	PeriodType    string  `json:"periodType,omitempty"`
}

// SearchResult represents search results
type SearchResult struct {
	Query  string        `json:"query"`
	Quotes []SearchQuote `json:"quotes"`
	News   []NewsItem    `json:"news,omitempty"`
	Count  int           `json:"count"`
}

// SearchQuote represents a single search result quote
type SearchQuote struct {
	Symbol         string  `json:"symbol"`
	ShortName      string  `json:"shortname"`
	LongName       string  `json:"longname"`
	Exchange       string  `json:"exchange"`
	ExchDisp       string  `json:"exchDisp"`
	TypeDisp       string  `json:"typeDisp"`
	QuoteType      string  `json:"quoteType"`
	Industry       string  `json:"industry,omitempty"`
	Sector         string  `json:"sector,omitempty"`
	Score          float64 `json:"score"`
	IsYahooFinance bool    `json:"isYahooFinance"`
}

// LookupResult represents lookup results
type LookupResult struct {
	Query string       `json:"query"`
	Items []LookupItem `json:"items"`
	Count int          `json:"count"`
}

// LookupItem represents a single lookup result
type LookupItem struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Exchange string `json:"exchange"`
	Type     string `json:"type"`
}

// NewsItem represents a news article
type NewsItem struct {
	UUID        string      `json:"uuid"`
	Title       string      `json:"title"`
	Publisher   string      `json:"publisher"`
	Link        string      `json:"link"`
	Thumbnail   interface{} `json:"thumbnail,omitempty"`
	PublishTime int64       `json:"providerPublishTime"`
	Type        string      `json:"type"`
	Symbols     []string    `json:"relatedTickers,omitempty"`
}

// MarketSummary represents market summary data
type MarketSummary struct {
	MarketState string        `json:"marketState"`
	Region      string        `json:"region"`
	Markets     []MarketIndex `json:"markets"`
}

// MarketIndex represents a market index
type MarketIndex struct {
	Symbol                     string  `json:"symbol"`
	ShortName                  string  `json:"shortName"`
	FullExchangeName           string  `json:"fullExchangeName"`
	RegularMarketPrice         float64 `json:"regularMarketPrice"`
	RegularMarketChange        float64 `json:"regularMarketChange"`
	RegularMarketChangePercent float64 `json:"regularMarketChangePercent"`
	RegularMarketTime          int64   `json:"regularMarketTime"`
	MarketState                string  `json:"marketState"`
}

// MarketTime represents market time information
type MarketTime struct {
	Exchange    string `json:"exchange"`
	Timezone    string `json:"timezone"`
	GMTOffset   int    `json:"gmtoffset"`
	MarketState string `json:"marketState"`
	CurrentTime int64  `json:"currentTime"`
	OpenTime    int64  `json:"openTime"`
	CloseTime   int64  `json:"closeTime"`
}

// Sector represents sector information
type Sector struct {
	Name      string `json:"name"`
	Key       string `json:"key"`
	Companies int    `json:"companies"`
}

// Industry represents industry information
type Industry struct {
	Name      string `json:"name"`
	Key       string `json:"key"`
	Sector    string `json:"sector"`
	Companies int    `json:"companies"`
}

// ScreenCriteria defines criteria for stock screening
type ScreenCriteria struct {
	Region    string                 `json:"region,omitempty"`
	Offset    int                    `json:"offset,omitempty"`
	Size      int                    `json:"size,omitempty"`
	SortField string                 `json:"sortField,omitempty"`
	SortType  string                 `json:"sortType,omitempty"`
	Query     map[string]interface{} `json:"query,omitempty"`
}

// ScreenResult represents screener results
type ScreenResult struct {
	Count  int     `json:"count"`
	Total  int     `json:"total"`
	Quotes []Quote `json:"quotes"`
}

// EarningsEvent represents an earnings calendar event
type EarningsEvent struct {
	Symbol           string  `json:"symbol"`
	CompanyShortName string  `json:"companyShortName"`
	EarningsDate     int64   `json:"earningsDate"`
	EpsEstimate      float64 `json:"epsEstimate,omitempty"`
	EpsActual        float64 `json:"epsActual,omitempty"`
	EpsSurprise      float64 `json:"epsSurprise,omitempty"`
	StartDateTime    int64   `json:"startDateTime,omitempty"`
}

// IPOEvent represents an IPO calendar event
type IPOEvent struct {
	Symbol      string  `json:"symbol"`
	CompanyName string  `json:"companyName"`
	Exchange    string  `json:"exchange"`
	PricingDate int64   `json:"pricingDate"`
	PriceFrom   float64 `json:"priceFrom,omitempty"`
	PriceTo     float64 `json:"priceTo,omitempty"`
	Currency    string  `json:"currency"`
	Actions     string  `json:"actions,omitempty"`
}

// EconomicEvent represents an economic calendar event
type EconomicEvent struct {
	EventName  string  `json:"eventName"`
	EventTime  int64   `json:"eventTime"`
	Country    string  `json:"country"`
	Actual     float64 `json:"actual,omitempty"`
	Estimate   float64 `json:"estimate,omitempty"`
	Previous   float64 `json:"previous,omitempty"`
	Importance string  `json:"importance"`
}

// SplitEvent represents a stock split calendar event
type SplitEvent struct {
	Symbol           string `json:"symbol"`
	CompanyShortName string `json:"companyShortName"`
	SplitDate        int64  `json:"splitDate"`
	SplitRatio       string `json:"splitRatio"`
}

// CalendarParams defines parameters for calendar queries
type CalendarParams struct {
	Start  time.Time `json:"start,omitempty"`
	End    time.Time `json:"end,omitempty"`
	Region string    `json:"region,omitempty"`
	Size   int       `json:"size,omitempty"`
}

// StreamMessage represents a real-time WebSocket message
type StreamMessage struct {
	ID            string  `json:"id"`
	Price         float64 `json:"price"`
	Time          int64   `json:"time"`
	Currency      string  `json:"currency"`
	Exchange      string  `json:"exchange"`
	QuoteType     int     `json:"quoteType"`
	MarketHours   int     `json:"marketHours"`
	ChangePercent float64 `json:"changePercent"`
	Change        float64 `json:"change"`
	DayVolume     int64   `json:"dayVolume"`
	DayHigh       float64 `json:"dayHigh"`
	DayLow        float64 `json:"dayLow"`
	PreviousClose float64 `json:"previousClose"`
	Bid           float64 `json:"bid"`
	BidSize       int64   `json:"bidSize"`
	Ask           float64 `json:"ask"`
	AskSize       int64   `json:"askSize"`
	OpenPrice     float64 `json:"openPrice"`
	ShortName     string  `json:"shortName"`
}
