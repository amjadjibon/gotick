package tui

import (
	"fmt"
	"time"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/linechart"
	"github.com/mum4k/termdash/widgets/text"

	"github.com/amjadjibon/gotick/pkg/yfinance"
)

func (app *App) updateDashboard() {
	t, err := yfinance.NewTicker(app.currentSymbol)
	if err != nil {
		_ = app.quoteText.Write(fmt.Sprintf("Error creating ticker: %v", err), text.WriteReplace())
		return
	}

	app.updateQuote(t)
	app.updateChart(t)
	app.updateMarketSummary()
	app.updateNews(t)
	app.updateRecommendations(t)
}

func (app *App) updateQuote(t *yfinance.Ticker) {
	quote, err := t.Quote(app.ctx)
	if err != nil {
		_ = app.quoteText.Write(fmt.Sprintf("Error fetching quote: %v", err), text.WriteReplace())
		_ = app.rangeDonut.Percent(0)
		return
	}

	color := cell.ColorGreen
	if quote.RegularMarketChangePercent < 0 {
		color = cell.ColorRed
	}

	_ = app.quoteText.Write(fmt.Sprintf("%s (%s)\n", quote.Symbol, quote.ShortName), text.WriteReplace())
	_ = app.quoteText.Write(fmt.Sprintf("Price:  $%.2f\n", quote.RegularMarketPrice))
	_ = app.quoteText.Write(fmt.Sprintf("Change: $%.2f (%.2f%%)\n", quote.RegularMarketChange, quote.RegularMarketChangePercent), text.WriteCellOpts(cell.FgColor(color)))
	_ = app.quoteText.Write(fmt.Sprintf("Volume: %d\n", quote.RegularMarketVolume))
	_ = app.quoteText.Write(fmt.Sprintf("Cap:    $%.2f B\n", float64(quote.MarketCap)/1e9))
	_ = app.quoteText.Write(fmt.Sprintf("PE:     %.2f\n", quote.TrailingPE))
	_ = app.quoteText.Write(fmt.Sprintf("52w L/H: %.2f - %.2f\n", quote.FiftyTwoWeekLow, quote.FiftyTwoWeekHigh))

	if quote.FiftyTwoWeekHigh > quote.FiftyTwoWeekLow {
		percent := int(((quote.RegularMarketPrice - quote.FiftyTwoWeekLow) / (quote.FiftyTwoWeekHigh - quote.FiftyTwoWeekLow)) * 100)
		if percent < 0 {
			percent = 0
		}
		if percent > 100 {
			percent = 100
		}
		_ = app.rangeDonut.Percent(percent)
	} else {
		_ = app.rangeDonut.Percent(0)
	}
}

func (app *App) updateChart(t *yfinance.Ticker) {
	historyParams := yfinance.HistoryParams{
		Period:   yfinance.Period(app.currentRange),
		Interval: yfinance.Interval(app.currentInterval),
	}

	history, err := t.History(app.ctx, historyParams)
	if err != nil || len(history.Bars) == 0 {
		return
	}

	var prices []float64
	minP := history.Bars[0].Close
	maxP := history.Bars[0].Close

	for _, bar := range history.Bars {
		val := bar.Close
		prices = append(prices, val)
		if val < minP {
			minP = val
		}
		if val > maxP {
			maxP = val
		}
	}

	rangeVal := maxP - minP
	if rangeVal == 0 {
		rangeVal = maxP * 0.1
	}

	padding := rangeVal * 1.0
	upperBound := maxP + padding
	lowerBound := minP - padding
	if lowerBound < 0 {
		lowerBound = 0
	}

	minLine := make([]float64, len(prices))
	maxLine := make([]float64, len(prices))
	for i := range prices {
		minLine[i] = lowerBound
		maxLine[i] = upperBound
	}

	_ = app.lc.Series("min_bound", minLine, linechart.SeriesCellOpts(cell.FgColor(cell.ColorBlack)))
	_ = app.lc.Series("max_bound", maxLine, linechart.SeriesCellOpts(cell.FgColor(cell.ColorBlack)))

	// Create X-axis labels based on time range
	xLabels := make(map[int]string)
	numBars := len(history.Bars)

	// Track which months/days we've already labeled to avoid duplicates
	lastMonth := -1
	lastDay := -1

	for i, bar := range history.Bars {
		switch app.currentRange {
		case "1y", "2y", "5y", "10y", "ytd", "max":
			// For yearly ranges, show month numbers (1, 2, 3, ...)
			month := int(bar.Timestamp.Month())
			if month != lastMonth {
				xLabels[i] = fmt.Sprintf("%d", month)
				lastMonth = month
			}
		case "1mo", "3mo", "6mo":
			// For monthly ranges, show day numbers (1, 5, 10, ...)
			day := bar.Timestamp.Day()
			if day != lastDay && (day == 1 || day%5 == 0) {
				xLabels[i] = fmt.Sprintf("%d", day)
				lastDay = day
			}
		case "1d", "5d":
			// For daily ranges, show hours
			hour := bar.Timestamp.Hour()
			if i == 0 || i == numBars-1 || hour%2 == 0 {
				xLabels[i] = bar.Timestamp.Format("15:04")
			}
		default:
			// Fallback: show evenly distributed dates
			if i == 0 || i == numBars/2 || i == numBars-1 {
				xLabels[i] = bar.Timestamp.Format("01/02")
			}
		}
	}

	_ = app.lc.Series("Price", prices,
		linechart.SeriesCellOpts(cell.FgColor(cell.ColorYellow)),
		linechart.SeriesXLabels(xLabels),
	)
}

func (app *App) updateMarketSummary() {
	indices, err := yfinance.GetMajorIndices(app.ctx)
	if err != nil {
		_ = app.marketText.Write(fmt.Sprintf("Error: %v", err), text.WriteReplace())
		return
	}

	app.marketText.Reset()
	for _, idx := range indices {
		if idx.Symbol == "" || idx.RegularMarketPrice == 0 {
			continue
		}

		color := cell.ColorGreen
		if idx.RegularMarketChange < 0 {
			color = cell.ColorRed
		}

		name := idx.ShortName
		if len(name) > 15 {
			name = name[:15] + "..."
		}

		_ = app.marketText.Write(fmt.Sprintf("%-18s %8.2f ", name, idx.RegularMarketPrice))
		_ = app.marketText.Write(fmt.Sprintf("%+6.2f%%\n", idx.RegularMarketChangePercent), text.WriteCellOpts(cell.FgColor(color)))
	}
}

func (app *App) updateNews(t *yfinance.Ticker) {
	news, err := t.News(app.ctx, 5)
	if err != nil {
		_ = app.newsText.Write(fmt.Sprintf("Error: %v", err), text.WriteReplace())
		return
	}

	app.newsText.Reset()
	for _, item := range news {
		_ = app.newsText.Write(fmt.Sprintf("• %s\n", item.Title))
		pubTime := time.Unix(item.PublishTime, 0)
		_ = app.newsText.Write(fmt.Sprintf("  %s - %s\n\n", item.Publisher, pubTime.Format("15:04 01/02")))
	}
}

func (app *App) updateRecommendations(t *yfinance.Ticker) {
	recs, err := t.Recommendations(app.ctx)
	if err != nil || len(recs) == 0 {
		_ = app.recBar.Values([]int{0, 0, 0, 0, 0}, 10)
		return
	}

	latest := recs[0]
	maxVal := 0
	vals := []int{latest.StrongBuy, latest.Buy, latest.Hold, latest.Sell, latest.StrongSell}
	for _, v := range vals {
		if v > maxVal {
			maxVal = v
		}
	}

	if maxVal == 0 {
		maxVal = 10
	} else {
		maxVal += 2
	}

	_ = app.recBar.Values(vals, maxVal)
}
