# yfinance - Yahoo Finance API Client for Go

A comprehensive Go package for accessing Yahoo Finance APIs, providing real-time quotes, historical data, financials, options, market data, analyst data, and more.

## Installation

```bash
go get github.com/amjadjibon/gotick/pkg/yfinance
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/amjadjibon/gotick/pkg/yfinance"
)

func main() {
    ctx := context.Background()

    ticker, err := yfinance.NewTicker("AAPL")
    if err != nil {
        log.Fatal(err)
    }

    quote, err := ticker.Quote(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("AAPL: $%.2f (%.2f%%)\n", 
        quote.RegularMarketPrice, 
        quote.RegularMarketChangePercent)
}
```

## Features

### Ticker API

```go
ticker, _ := yfinance.NewTicker("AAPL")

// Real-time quote
quote, _ := ticker.Quote(ctx)

// Historical data
history, _ := ticker.History(ctx, yfinance.HistoryParams{
    Period:   yfinance.Period1mo,
    Interval: yfinance.Interval1d,
})

// Company info
info, _ := ticker.Info(ctx)

// Options chain
options, _ := ticker.Options(ctx, "")

// Dividends and splits
dividends, _ := ticker.Dividends(ctx, yfinance.HistoryParams{})
splits, _ := ticker.Splits(ctx, yfinance.HistoryParams{})

// News
news, _ := ticker.News(ctx, 10)
```

### Analysis API

```go
// Analyst recommendations
recs, _ := ticker.Recommendations(ctx)
// Returns: Period, StrongBuy, Buy, Hold, Sell, StrongSell

// Analyst price targets
targets, _ := ticker.AnalystPriceTargets(ctx)
// Returns: Current, Low, Mean, Median, High, NumAnalysts

// Earnings estimates
earnings, _ := ticker.EarningsEstimates(ctx)
// Returns: Period, Avg, Low, High, YearAgoEps, Growth

// Revenue estimates
revenue, _ := ticker.RevenueEstimates(ctx)

// EPS trends and revisions
epsTrend, _ := ticker.EPSTrends(ctx)
epsRevisions, _ := ticker.EPSRevisions(ctx)

// Earnings history
history, _ := ticker.EarningsHistoryData(ctx)

// Growth estimates
growth, _ := ticker.GrowthEstimates(ctx)
```

### Holders API

```go
// Major holders breakdown
major, _ := ticker.MajorHolders(ctx)
// Returns: InsidersPercentHeld, InstitutionsPercentHeld, InstitutionsCount

// Institutional holders (Vanguard, BlackRock, etc.)
institutional, _ := ticker.InstitutionalHolders(ctx)
// Returns: Holder, Shares, Value, PctHeld, DateReported

// Mutual fund holders
funds, _ := ticker.MutualFundHolders(ctx)

// Insider transactions
transactions, _ := ticker.InsiderTransactions(ctx)
// Returns: Insider, Relation, Transaction, Shares, Value, StartDate

// Insider roster
roster, _ := ticker.InsiderRosterHolders(ctx)

// Insider purchases summary
purchases, _ := ticker.InsiderPurchasesData(ctx)
```

### Search

```go
results, _ := yfinance.Search(ctx, "Apple", yfinance.WithQuotesCount(10))
lookup, _ := yfinance.Lookup(ctx, "AAPL", "equity")
```

### Multiple Quotes

```go
quotes, _ := yfinance.QuoteMultiple(ctx, []string{"AAPL", "GOOGL", "MSFT"})
```

### Market Data

```go
summary, _ := yfinance.GetMarketSummary(ctx)
indices, _ := yfinance.GetMajorIndices(ctx)
futures, _ := yfinance.GetMajorFutures(ctx)
crypto, _ := yfinance.GetMajorCrypto(ctx)
trending, _ := yfinance.GetTrending(ctx, "US", 10)
```

### Screener

```go
active, _ := yfinance.ScreenMostActive(ctx, 10)
gainers, _ := yfinance.ScreenGainers(ctx, 10)
losers, _ := yfinance.ScreenLosers(ctx, 10)
tech, _ := yfinance.ScreenBySector(ctx, yfinance.SectorTechnology, 10)
dividend, _ := yfinance.ScreenHighDividend(ctx, 0.03, 10)
```

### Calendar Events

```go
earnings, _ := yfinance.GetEarningsCalendar(ctx, yfinance.CalendarParams{})
ipos, _ := yfinance.GetIPOCalendar(ctx, yfinance.CalendarParams{})
splits, _ := yfinance.GetSplitsCalendar(ctx, yfinance.CalendarParams{})
```

### Sector & Industry

```go
sectors, _ := yfinance.GetSectors(ctx)
industries, _ := yfinance.GetIndustries(ctx)
```

### News

```go
news, _ := yfinance.GetNews(ctx, []string{"AAPL", "GOOGL"}, 10)
latest, _ := yfinance.GetLatestNews(ctx, 20)
```

### WebSocket Streaming

```go
stream := yfinance.NewStream([]string{"AAPL", "GOOGL"})
err := stream.Connect(ctx)
if err != nil {
    log.Fatal(err)
}
defer stream.Close()

for msg := range stream.Messages() {
    fmt.Printf("%s: $%.2f\n", msg.ID, msg.Price)
}
```

## Available Intervals

| Interval | Constant |
|----------|----------|
| 1 minute | `Interval1m` |
| 5 minutes | `Interval5m` |
| 15 minutes | `Interval15m` |
| 1 hour | `Interval1h` |
| 1 day | `Interval1d` |
| 1 week | `Interval1wk` |
| 1 month | `Interval1mo` |

## Available Periods

| Period | Constant |
|--------|----------|
| 1 day | `Period1d` |
| 5 days | `Period5d` |
| 1 month | `Period1mo` |
| 3 months | `Period3mo` |
| 6 months | `Period6mo` |
| 1 year | `Period1y` |
| 5 years | `Period5y` |
| Max | `PeriodMax` |

## Error Handling

```go
quote, err := ticker.Quote(ctx)
if err != nil {
    if yfinance.IsNotFound(err) {
        fmt.Println("Symbol not found")
    } else if yfinance.IsRateLimited(err) {
        fmt.Println("Rate limited, try again later")
    } else {
        fmt.Printf("Error: %v\n", err)
    }
}
```

## Custom Client

```go
client, _ := yfinance.NewClient(
    yfinance.WithTimeout(60 * time.Second),
    yfinance.WithUserAgent("MyApp/1.0"),
)

ticker, _ := yfinance.NewTicker("AAPL", yfinance.WithClient(client))
```

## License

MIT License
