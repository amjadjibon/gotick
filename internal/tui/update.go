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
		app.quoteText.Write(fmt.Sprintf("Error creating ticker: %v", err), text.WriteReplace())
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
		app.quoteText.Write(fmt.Sprintf("Error fetching quote: %v", err), text.WriteReplace())
		app.rangeDonut.Percent(0)
		return
	}

	color := cell.ColorGreen
	if quote.RegularMarketChangePercent < 0 {
		color = cell.ColorRed
	}

	app.quoteText.Write(fmt.Sprintf("%s (%s)\n", quote.Symbol, quote.ShortName), text.WriteReplace())
	app.quoteText.Write(fmt.Sprintf("Price:  $%.2f\n", quote.RegularMarketPrice))
	app.quoteText.Write(fmt.Sprintf("Change: $%.2f (%.2f%%)\n", quote.RegularMarketChange, quote.RegularMarketChangePercent), text.WriteCellOpts(cell.FgColor(color)))
	app.quoteText.Write(fmt.Sprintf("Volume: %d\n", quote.RegularMarketVolume))
	app.quoteText.Write(fmt.Sprintf("Cap:    $%.2f B\n", float64(quote.MarketCap)/1e9))
	app.quoteText.Write(fmt.Sprintf("PE:     %.2f\n", quote.TrailingPE))
	app.quoteText.Write(fmt.Sprintf("52w L/H: %.2f - %.2f\n", quote.FiftyTwoWeekLow, quote.FiftyTwoWeekHigh))

	if quote.FiftyTwoWeekHigh > quote.FiftyTwoWeekLow {
		percent := int(((quote.RegularMarketPrice - quote.FiftyTwoWeekLow) / (quote.FiftyTwoWeekHigh - quote.FiftyTwoWeekLow)) * 100)
		if percent < 0 {
			percent = 0
		}
		if percent > 100 {
			percent = 100
		}
		app.rangeDonut.Percent(percent)
	} else {
		app.rangeDonut.Percent(0)
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

	app.lc.Series("min_bound", minLine, linechart.SeriesCellOpts(cell.FgColor(cell.ColorBlack)))
	app.lc.Series("max_bound", maxLine, linechart.SeriesCellOpts(cell.FgColor(cell.ColorBlack)))

	app.lc.Series("Price", prices,
		linechart.SeriesCellOpts(cell.FgColor(cell.ColorYellow)),
		linechart.SeriesXLabels(map[int]string{
			0:                     history.Bars[0].Timestamp.Format("01/02"),
			len(history.Bars) / 2: history.Bars[len(history.Bars)/2].Timestamp.Format("01/02"),
			len(history.Bars) - 1: history.Bars[len(history.Bars)-1].Timestamp.Format("01/02"),
		}),
	)
}

func (app *App) updateMarketSummary() {
	indices, err := yfinance.GetMajorIndices(app.ctx)
	if err != nil {
		app.marketText.Write(fmt.Sprintf("Error: %v", err), text.WriteReplace())
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

		app.marketText.Write(fmt.Sprintf("%-18s %8.2f ", name, idx.RegularMarketPrice))
		app.marketText.Write(fmt.Sprintf("%+6.2f%%\n", idx.RegularMarketChangePercent), text.WriteCellOpts(cell.FgColor(color)))
	}
}

func (app *App) updateNews(t *yfinance.Ticker) {
	news, err := t.News(app.ctx, 5)
	if err != nil {
		app.newsText.Write(fmt.Sprintf("Error: %v", err), text.WriteReplace())
		return
	}

	app.newsText.Reset()
	for _, item := range news {
		app.newsText.Write(fmt.Sprintf("• %s\n", item.Title))
		pubTime := time.Unix(item.PublishTime, 0)
		app.newsText.Write(fmt.Sprintf("  %s - %s\n\n", item.Publisher, pubTime.Format("15:04 01/02")))
	}
}

func (app *App) updateRecommendations(t *yfinance.Ticker) {
	recs, err := t.Recommendations(app.ctx)
	if err != nil || len(recs) == 0 {
		app.recBar.Values([]int{0, 0, 0, 0, 0}, 10)
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

	app.recBar.Values(vals, maxVal)
}
