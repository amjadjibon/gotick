// Package yfinance provides a comprehensive Go client for Yahoo Finance APIs.
//
// This package offers access to real-time and historical stock market data,
// including quotes, charts, options, financials, news, and more.
//
// # Quick Start
//
// Create a ticker and fetch data:
//
//	ticker, err := yfinance.NewTicker("AAPL")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Get real-time quote
//	quote, err := ticker.Quote(context.Background())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Price: $%.2f\n", quote.RegularMarketPrice)
//
//	// Get historical data
//	history, err := ticker.History(context.Background(), yfinance.HistoryParams{
//	    Period:   yfinance.Period1mo,
//	    Interval: yfinance.Interval1d,
//	})
//
// # Features
//
//   - Real-time quotes for stocks, ETFs, indices, forex, and crypto
//   - Historical OHLCV data with configurable intervals
//   - Company information and financial data
//   - Options chain data
//   - Stock screening and search
//   - Market summaries and trending tickers
//   - Calendar events (earnings, IPOs, splits)
//   - Real-time WebSocket streaming
//
// # Authentication
//
// The package handles Yahoo Finance authentication automatically using
// cookies and crumb tokens. No API key is required.
//
// # Thread Safety
//
// The Client and Ticker types are safe for concurrent use.
package yfinance
