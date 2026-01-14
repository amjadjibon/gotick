package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/amjadjibon/gotick/pkg/yfinance"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║           YFinance Go Package - Complete Demo                ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Create ticker
	ticker, err := yfinance.NewTicker("AAPL")
	if err != nil {
		log.Fatalf("Failed to create ticker: %v", err)
	}

	// ===== BASIC DATA =====
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📊 BASIC DATA")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Test 1: Quote
	fmt.Println("\n1. Quote")
	quote, err := ticker.Quote(ctx)
	if err != nil {
		log.Printf("   ❌ Failed: %v", err)
	} else {
		fmt.Printf("   ✅ %s: $%.2f (%.2f%%)\n", quote.Symbol, quote.RegularMarketPrice, quote.RegularMarketChangePercent)
	}

	// Test 2: History
	fmt.Println("\n2. Historical Data (5 days)")
	history, err := ticker.History(ctx, yfinance.HistoryParams{Period: yfinance.Period5d, Interval: yfinance.Interval1d})
	if err != nil {
		log.Printf("   ❌ Failed: %v", err)
	} else {
		fmt.Printf("   ✅ Got %d bars\n", len(history.Bars))
		for i, bar := range history.Bars {
			if i < 3 {
				fmt.Printf("      %s: O=%.2f H=%.2f L=%.2f C=%.2f\n",
					bar.Timestamp.Format("2006-01-02"), bar.Open, bar.High, bar.Low, bar.Close)
			}
		}
	}

	// ===== ANALYST DATA =====
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📈 ANALYST DATA")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Test 3: Recommendations
	fmt.Println("\n3. Analyst Recommendations")
	recs, err := ticker.Recommendations(ctx)
	if err != nil {
		log.Printf("   ❌ Failed: %v", err)
	} else {
		for _, r := range recs {
			fmt.Printf("   ✅ %s: 🟢Buy=%d 🟡Hold=%d 🔴Sell=%d\n",
				r.Period, r.Buy+r.StrongBuy, r.Hold, r.Sell+r.StrongSell)
		}
	}

	// Test 4: Price Targets
	fmt.Println("\n4. Analyst Price Targets")
	pt, err := ticker.AnalystPriceTargets(ctx)
	if err != nil {
		log.Printf("   ❌ Failed: %v", err)
	} else {
		fmt.Printf("   ✅ Current: $%.2f\n", pt.Current)
		fmt.Printf("      Target: Low=$%.2f | Mean=$%.2f | High=$%.2f\n", pt.Low, pt.Mean, pt.High)
		fmt.Printf("      Analysts: %d\n", pt.NumAnalysts)
	}

	// Test 5: Earnings Estimates
	fmt.Println("\n5. Earnings Estimates")
	ee, err := ticker.EarningsEstimates(ctx)
	if err != nil {
		log.Printf("   ❌ Failed: %v", err)
	} else {
		for _, e := range ee {
			fmt.Printf("   ✅ %s: Avg=$%.2f (Low=$%.2f, High=$%.2f)\n", e.Period, e.Avg, e.Low, e.High)
		}
	}

	// ===== HOLDERS DATA =====
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🏛️  HOLDERS DATA")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Test 6: Major Holders
	fmt.Println("\n6. Major Holders Breakdown")
	mh, err := ticker.MajorHolders(ctx)
	if err != nil {
		log.Printf("   ❌ Failed: %v", err)
	} else {
		fmt.Printf("   ✅ Insiders: %.2f%%\n", mh.InsidersPercentHeld*100)
		fmt.Printf("      Institutions: %.2f%%\n", mh.InstitutionsPercentHeld*100)
		fmt.Printf("      Institution Count: %d\n", mh.InstitutionsCount)
	}

	// Test 7: Institutional Holders
	fmt.Println("\n7. Top Institutional Holders")
	ih, err := ticker.InstitutionalHolders(ctx)
	if err != nil {
		log.Printf("   ❌ Failed: %v", err)
	} else {
		count := 5
		if len(ih) < count {
			count = len(ih)
		}
		for i := 0; i < count; i++ {
			fmt.Printf("   ✅ %s: %.2f%%\n", ih[i].Holder, ih[i].PctHeld*100)
		}
	}

	// Test 8: Insider Transactions
	fmt.Println("\n8. Recent Insider Transactions")
	it, err := ticker.InsiderTransactions(ctx)
	if err != nil {
		log.Printf("   ❌ Failed: %v", err)
	} else {
		count := 3
		if len(it) < count {
			count = len(it)
		}
		for i := 0; i < count; i++ {
			fmt.Printf("   ✅ %s (%s)\n      %s\n", it[i].Insider, it[i].Relation, it[i].Transaction)
		}
	}

	// ===== FINANCIALS =====
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("💰 FINANCIALS")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Test 9: Income Statement (Quarterly)
	fmt.Println("\n9. Income Statement (Quarterly)")
	income, err := ticker.IncomeStatement(ctx, true)
	if err != nil {
		log.Printf("   ❌ Failed: %v", err)
	} else {
		fmt.Printf("   ✅ Got %d quarterly periods\n", len(income.Quarterly))
	}

	// Test 10: Balance Sheet (Annual)
	fmt.Println("\n10. Balance Sheet (Annual)")
	balance, err := ticker.BalanceSheet(ctx, false)
	if err != nil {
		log.Printf("   ❌ Failed: %v", err)
	} else {
		fmt.Printf("   ✅ Got %d annual periods\n", len(balance.Annual))
	}

	// ===== OPTIONS & GREEKS =====
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📉 OPTIONS & GREEKS")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Test 11: Options Chain
	fmt.Println("\n11. Options Chain")
	options, err := ticker.Options(ctx, "")
	if err != nil {
		log.Printf("   ❌ Failed: %v", err)
	} else {
		fmt.Printf("   ✅ Underlying: $%.2f\n", options.UnderlyingPrice)
		fmt.Printf("      Expiration Dates: %d available\n", len(options.ExpirationDates))
		fmt.Printf("      Calls: %d, Puts: %d\n", len(options.Calls), len(options.Puts))
	}

	// Test 12: Greeks Calculation
	fmt.Println("\n12. Greeks Calculation (Sample)")
	greeks := yfinance.CalculateGreeks(150, 150, 0.05, 0.25, 0.25, true)
	if greeks != nil {
		fmt.Printf("   ✅ Delta: %.4f\n", greeks.Delta)
		fmt.Printf("      Gamma: %.4f\n", greeks.Gamma)
		fmt.Printf("      Theta: %.4f (daily)\n", greeks.Theta)
		fmt.Printf("      Vega:  %.4f\n", greeks.Vega)
		fmt.Printf("      Rho:   %.4f\n", greeks.Rho)
	}

	// ===== BATCH OPERATIONS =====
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🔄 BATCH OPERATIONS")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Test 13: Tickers Batch
	fmt.Println("\n13. Batch Quotes (AAPL, GOOGL, MSFT, AMZN)")
	tickers, err := yfinance.NewTickers("AAPL", "GOOGL", "MSFT", "AMZN")
	if err != nil {
		log.Printf("   ❌ Failed: %v", err)
	} else {
		quotes, err := tickers.Quotes(ctx)
		if err != nil {
			log.Printf("   ❌ Failed: %v", err)
		} else {
			for sym, q := range quotes {
				fmt.Printf("   ✅ %s: $%.2f (%.2f%%)\n", sym, q.RegularMarketPrice, q.RegularMarketChangePercent)
			}
		}
	}

	// Test 14: Download
	fmt.Println("\n14. Batch Download (History)")
	result, err := yfinance.Download(ctx, yfinance.DownloadParams{
		Symbols:  []string{"AAPL", "GOOGL"},
		Period:   yfinance.Period5d,
		Interval: yfinance.Interval1d,
		Threads:  2,
	})
	if err != nil {
		log.Printf("   ❌ Failed: %v", err)
	} else {
		for sym, data := range result.Data {
			fmt.Printf("   ✅ %s: %d bars\n", sym, len(data.Bars))
		}
	}

	// ===== SEARCH & MARKET =====
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🔍 SEARCH & MARKET")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Test 15: Search
	fmt.Println("\n15. Search: \"Tesla\"")
	results, err := yfinance.Search(ctx, "Tesla", yfinance.WithQuotesCount(3))
	if err != nil {
		log.Printf("   ❌ Failed: %v", err)
	} else {
		for _, r := range results.Quotes {
			fmt.Printf("   ✅ %s: %s\n", r.Symbol, r.ShortName)
		}
	}

	// Test 16: Actions
	fmt.Println("\n16. Corporate Actions (Dividends + Splits)")
	actions, err := ticker.Actions(ctx, yfinance.HistoryParams{Period: yfinance.Period1y})
	if err != nil {
		log.Printf("   ❌ Failed: %v", err)
	} else {
		count := 5
		if len(actions) < count {
			count = len(actions)
		}
		fmt.Printf("   ✅ Found %d actions (showing first %d)\n", len(actions), count)
		for i := 0; i < count; i++ {
			a := actions[i]
			if a.Type == "dividend" {
				fmt.Printf("      💵 %s: Dividend $%.4f\n", a.Date.Format("2006-01-02"), a.Amount)
			} else {
				fmt.Printf("      📊 %s: Split %s\n", a.Date.Format("2006-01-02"), a.Ratio)
			}
		}
	}

	// ===== FUND DATA (ETF) =====
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📦 FUND/ETF DATA (SPY)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Test 17: Fund Holdings
	spy, _ := yfinance.NewTicker("SPY")
	fmt.Println("\n17. Top Fund Holdings")
	holdings, err := spy.FundHoldings(ctx)
	if err != nil {
		log.Printf("   ❌ Failed: %v", err)
	} else {
		count := 5
		if len(holdings) < count {
			count = len(holdings)
		}
		for i := 0; i < count; i++ {
			fmt.Printf("   ✅ %s (%s): %.2f%%\n", holdings[i].Name, holdings[i].Symbol, holdings[i].Percent*100)
		}
	}

	// Test 18: Sector Weightings
	fmt.Println("\n18. Sector Weightings")
	weightings, err := spy.FundSectorWeightings(ctx)
	if err != nil {
		log.Printf("   ❌ Failed: %v", err)
	} else {
		for _, w := range weightings {
			fmt.Printf("   ✅ %s: %.2f%%\n", w.Sector, w.Percent*100)
		}
	}

	fmt.Println("\n╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                     Demo Complete! ✅                        ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
}
