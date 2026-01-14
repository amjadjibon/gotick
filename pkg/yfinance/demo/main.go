package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/amjadjibon/gotick/pkg/yfinance"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	fmt.Println("=== YFinance Go Package Demo ===")
	fmt.Println()

	// Create ticker
	ticker, err := yfinance.NewTicker("AAPL")
	if err != nil {
		log.Fatalf("Failed to create ticker: %v", err)
	}

	// Test 1: Quote
	fmt.Println("1. Testing Quote...")
	quote, err := ticker.Quote(ctx)
	if err != nil {
		log.Printf("Failed to get quote: %v", err)
	} else {
		fmt.Printf("   %s: $%.2f (%.2f%%)\n", quote.Symbol, quote.RegularMarketPrice, quote.RegularMarketChangePercent)
	}
	fmt.Println()

	// Test 2: History
	fmt.Println("2. Testing History...")
	history, err := ticker.History(ctx, yfinance.HistoryParams{Period: yfinance.Period5d, Interval: yfinance.Interval1d})
	if err != nil {
		log.Printf("Failed to get history: %v", err)
	} else {
		fmt.Printf("   Got %d bars\n", len(history.Bars))
	}
	fmt.Println()

	// Test 3: Recommendations (NEW)
	fmt.Println("3. Testing Recommendations...")
	recs, err := ticker.Recommendations(ctx)
	if err != nil {
		log.Printf("Failed to get recommendations: %v", err)
	} else {
		for _, r := range recs {
			fmt.Printf("   %s: Buy=%d, Hold=%d, Sell=%d\n", r.Period, r.Buy+r.StrongBuy, r.Hold, r.Sell+r.StrongSell)
		}
	}
	fmt.Println()

	// Test 4: Analyst Price Targets (NEW)
	fmt.Println("4. Testing Analyst Price Targets...")
	pt, err := ticker.AnalystPriceTargets(ctx)
	if err != nil {
		log.Printf("Failed to get price targets: %v", err)
	} else {
		fmt.Printf("   Current: $%.2f\n", pt.Current)
		fmt.Printf("   Target Low: $%.2f, Mean: $%.2f, High: $%.2f\n", pt.Low, pt.Mean, pt.High)
		fmt.Printf("   Number of Analysts: %d\n", pt.NumAnalysts)
	}
	fmt.Println()

	// Test 5: Major Holders (NEW)
	fmt.Println("5. Testing Major Holders...")
	mh, err := ticker.MajorHolders(ctx)
	if err != nil {
		log.Printf("Failed to get major holders: %v", err)
	} else {
		fmt.Printf("   Insiders: %.2f%%\n", mh.InsidersPercentHeld*100)
		fmt.Printf("   Institutions: %.2f%%\n", mh.InstitutionsPercentHeld*100)
		fmt.Printf("   Institution Count: %d\n", mh.InstitutionsCount)
	}
	fmt.Println()

	// Test 6: Institutional Holders (NEW)
	fmt.Println("6. Testing Institutional Holders (top 3)...")
	ih, err := ticker.InstitutionalHolders(ctx)
	if err != nil {
		log.Printf("Failed to get institutional holders: %v", err)
	} else {
		count := 3
		if len(ih) < count {
			count = len(ih)
		}
		for i := 0; i < count; i++ {
			fmt.Printf("   %s: %d shares (%.2f%%)\n", ih[i].Holder, ih[i].Shares, ih[i].PctHeld*100)
		}
	}
	fmt.Println()

	// Test 7: Earnings Estimates (NEW)
	fmt.Println("7. Testing Earnings Estimates...")
	ee, err := ticker.EarningsEstimates(ctx)
	if err != nil {
		log.Printf("Failed to get earnings estimates: %v", err)
	} else {
		for _, e := range ee {
			fmt.Printf("   %s: Avg=$%.2f (Low=$%.2f, High=$%.2f)\n", e.Period, e.Avg, e.Low, e.High)
		}
	}
	fmt.Println()

	// Test 8: Insider Transactions (NEW)
	fmt.Println("8. Testing Insider Transactions (top 3)...")
	it, err := ticker.InsiderTransactions(ctx)
	if err != nil {
		log.Printf("Failed to get insider transactions: %v", err)
	} else {
		count := 3
		if len(it) < count {
			count = len(it)
		}
		for i := 0; i < count; i++ {
			fmt.Printf("   %s (%s): %s\n", it[i].Insider, it[i].Relation, it[i].Transaction)
		}
	}
	fmt.Println()

	// Test 9: Multiple Quotes
	fmt.Println("9. Testing Multiple Quotes...")
	quotes, err := yfinance.QuoteMultiple(ctx, []string{"AAPL", "GOOGL", "MSFT", "AMZN"})
	if err != nil {
		log.Printf("Failed to get multiple quotes: %v", err)
	} else {
		for _, q := range quotes {
			fmt.Printf("   %s: $%.2f (%.2f%%)\n", q.Symbol, q.RegularMarketPrice, q.RegularMarketChangePercent)
		}
	}
	fmt.Println()

	// Test 10: Search
	fmt.Println("10. Testing Search...")
	results, err := yfinance.Search(ctx, "Tesla", yfinance.WithQuotesCount(3))
	if err != nil {
		log.Printf("Failed to search: %v", err)
	} else {
		for _, r := range results.Quotes {
			fmt.Printf("   %s: %s\n", r.Symbol, r.ShortName)
		}
	}
	fmt.Println()

	fmt.Println("=== Demo Complete ===")
}
