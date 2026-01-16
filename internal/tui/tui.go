package tui

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/align"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/barchart"
	"github.com/mum4k/termdash/widgets/donut"
	"github.com/mum4k/termdash/widgets/linechart"
	"github.com/mum4k/termdash/widgets/text"
	"github.com/mum4k/termdash/widgets/textinput"

	"github.com/amjadjibon/gotick/pkg/yfinance"
)

// Options holds configuration options for the TUI
type Options struct {
	Symbol   string
	Interval string
	Range    string
}

// App holds the dashboard application state
type App struct {
	ctx             context.Context
	cancel          context.CancelFunc
	currentSymbol   string
	currentInterval string
	currentRange    string

	// Widgets
	input      *textinput.TextInput
	lc         *linechart.LineChart
	quoteText  *text.Text
	marketText *text.Text
	newsText   *text.Text
	recBar     *barchart.BarChart
	rangeDonut *donut.Donut
}

func Run(opts Options) {
	app := &App{
		currentSymbol:   opts.Symbol,
		currentInterval: opts.Interval,
		currentRange:    opts.Range,
	}

	if app.currentSymbol == "" {
		app.currentSymbol = "AAPL"
	}

	t, err := tcell.New()
	if err != nil {
		log.Fatal(err)
	}
	defer t.Close()

	app.ctx, app.cancel = context.WithCancel(context.Background())
	defer app.cancel()

	// --- Widgets ---

	// Input for symbol search
	app.input, err = textinput.New(
		textinput.Label("Symbol: ", cell.FgColor(cell.ColorNumber(33))),
		textinput.MaxWidthCells(30),
		textinput.PlaceHolder("Enter symbol (e.g. AAPL)"),
		textinput.OnSubmit(func(text string) error {
			if text != "" {
				app.currentSymbol = strings.ToUpper(text)
				go app.updateDashboard()
			}
			return nil
		}),
		textinput.ClearOnSubmit(),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Price Chart
	app.lc, err = linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorRed)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorGreen)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorGreen)),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Quote Details
	app.quoteText, err = text.New()
	if err != nil {
		log.Fatal(err)
	}

	// Market Summary
	app.marketText, err = text.New()
	if err != nil {
		log.Fatal(err)
	}

	// News Feed
	app.newsText, err = text.New()
	if err != nil {
		log.Fatal(err)
	}

	// Analyst Recommendations (Bar Chart)
	app.recBar, err = barchart.New(
		barchart.BarColors([]cell.Color{
			cell.ColorGreen,       // Strong Buy
			cell.ColorNumber(118), // Buy (Light Green)
			cell.ColorYellow,      // Hold
			cell.ColorRed,         // Sell
			cell.ColorNumber(88),  // Strong Sell (Dark Red)
		}),
		barchart.ValueColors([]cell.Color{
			cell.ColorGreen,
			cell.ColorNumber(118),
			cell.ColorYellow,
			cell.ColorRed,
			cell.ColorNumber(88),
			cell.ColorNumber(88),
		}),
		barchart.ShowValues(),
		barchart.Labels([]string{"S.Buy", "Buy", "Hold", "Sell", "S.Sell"}),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 52-Week Range (Donut)
	app.rangeDonut, err = donut.New(
		donut.CellOpts(cell.FgColor(cell.ColorCyan)),
		donut.Label("Range %", cell.FgColor(cell.ColorWhite)),
	)
	if err != nil {
		log.Fatal(err)
	}

	// --- Layout ---

	c, err := container.New(
		t,
		container.Border(linestyle.Light),
		container.BorderTitle(" YFinance Go Terminal "),
		container.SplitHorizontal(
			container.Top(
				container.SplitVertical(
					container.Left(
						container.SplitHorizontal(
							container.Top(
								container.PlaceWidget(app.input),
								container.AlignHorizontal(align.HorizontalLeft),
								container.Border(linestyle.Light),
								container.BorderTitle(" Search "),
							),
							container.Bottom(
								container.PlaceWidget(app.lc),
								container.Border(linestyle.Light),
								container.BorderTitle(" Price History (1 Year) "),
							),
							container.SplitFixed(3),
						),
					),
					container.Right(
						container.SplitHorizontal(
							container.Top(
								container.PlaceWidget(app.marketText),
								container.Border(linestyle.Light),
								container.BorderTitle(" Market Summary "),
							),
							container.Bottom(
								container.PlaceWidget(app.newsText),
								container.Border(linestyle.Light),
								container.BorderTitle(" News Feed "),
							),
							container.SplitPercent(40),
						),
					),
					container.SplitPercent(65),
				),
			),
			container.Bottom(
				container.SplitVertical(
					container.Left(
						container.PlaceWidget(app.quoteText),
						container.Border(linestyle.Light),
						container.BorderTitle(" Quote Info "),
					),
					container.Right(
						container.SplitVertical(
							container.Left(
								container.PlaceWidget(app.rangeDonut),
								container.Border(linestyle.Light),
								container.BorderTitle(" 52-Week Range "),
							),
							container.Right(
								container.PlaceWidget(app.recBar),
								container.Border(linestyle.Light),
								container.BorderTitle(" Analyst Recommendations "),
							),
							container.SplitPercent(40),
						),
					),
					container.SplitPercent(30),
				),
			),
			container.SplitPercent(70),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	// --- Data Refresh ---

	// Initial load
	go app.updateDashboard()

	// Periodic update
	ticker := time.NewTicker(30 * time.Second) // Refresh every 30s
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				app.updateDashboard()
			case <-app.ctx.Done():
				return
			}
		}
	}()

	// --- Run ---

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == keyboard.KeyEsc {
			app.cancel()
		}
	}

	if err := termdash.Run(app.ctx, t, c, termdash.KeyboardSubscriber(quitter)); err != nil {
		log.Fatal(err)
	}
}

// updateDashboard fetches data and updates all widgets
func (app *App) updateDashboard() {
	// Create ticker
	t, err := yfinance.NewTicker(app.currentSymbol)
	if err != nil {
		app.quoteText.Write(fmt.Sprintf("Error creating ticker: %v", err), text.WriteReplace())
		return
	}

	// WaitGroup to fetch data concurrently? Ideally yes, but for now sequential to avoid race conditions on text widgets

	// 1. Update Quote & Donut
	quote, err := t.Quote(app.ctx)
	if err != nil {
		app.quoteText.Write(fmt.Sprintf("Error fetching quote: %v", err), text.WriteReplace())
		app.rangeDonut.Percent(0)
	} else {
		// Text
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

		// Donut Percentage
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

	// 2. Update Chart (History)
	historyParams := yfinance.HistoryParams{
		Period:   yfinance.Period(app.currentRange),
		Interval: yfinance.Interval(app.currentInterval),
	}
	history, err := t.History(app.ctx, historyParams)
	if err != nil {
		// Log error to quote text just so user sees it
		// qt.Write(fmt.Sprintf("\nHistory Error: %v", err))
	} else if len(history.Bars) > 0 {
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

		// Center the graph by creating a symmetric margin around the mid-point of data
		// Calculate data range
		rangeVal := maxP - minP
		if rangeVal == 0 {
			rangeVal = maxP * 0.1
		}

		// Use a margin that is 50% of the range on both top and bottom
		// This forces the actual data to occupy the middle ~50% of the chart
		padding := rangeVal * 1.0

		upperBound := maxP + padding
		lowerBound := minP - padding
		if lowerBound < 0 {
			lowerBound = 0
		}

		// Create full-length arrays for bounds to avoid drawing diagonal artifact lines
		// We use two series: one flat line at min, one flat line at max
		minLine := make([]float64, len(prices))
		maxLine := make([]float64, len(prices))
		for i := range prices {
			minLine[i] = lowerBound
			maxLine[i] = upperBound
		}

		app.lc.Series("min_bound", minLine, linechart.SeriesCellOpts(cell.FgColor(cell.ColorBlack)))
		app.lc.Series("max_bound", maxLine, linechart.SeriesCellOpts(cell.FgColor(cell.ColorBlack)))

		if err := app.lc.Series("Price", prices,
			linechart.SeriesCellOpts(cell.FgColor(cell.ColorYellow)),
			linechart.SeriesXLabels(map[int]string{
				0:                     history.Bars[0].Timestamp.Format("01/02"),
				len(history.Bars) / 2: history.Bars[len(history.Bars)/2].Timestamp.Format("01/02"),
				len(history.Bars) - 1: history.Bars[len(history.Bars)-1].Timestamp.Format("01/02"),
			}),
		); err != nil {
			// Ignore series error
		}
	}

	// 3. Update Market Summary
	indices, err := yfinance.GetMajorIndices(app.ctx)
	if err != nil {
		app.marketText.Write(fmt.Sprintf("Error: %v", err), text.WriteReplace())
	} else {
		app.marketText.Reset()
		for _, idx := range indices {
			// Some indices might fail individually, skip them
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

	// 4. Update News
	news, err := t.News(app.ctx, 5)
	if err != nil {
		app.newsText.Write(fmt.Sprintf("Error: %v", err), text.WriteReplace())
	} else {
		app.newsText.Reset()
		for _, item := range news {
			app.newsText.Write(fmt.Sprintf("• %s\n", item.Title))
			pubTime := time.Unix(item.PublishTime, 0)
			app.newsText.Write(fmt.Sprintf("  %s - %s\n\n", item.Publisher, pubTime.Format("15:04 01/02")))
		}
	}

	// 5. Update Recommendations
	recs, err := t.Recommendations(app.ctx)
	if err != nil || len(recs) == 0 {
		// Clear or show empty
		app.recBar.Values([]int{0, 0, 0, 0, 0}, 10)
	} else {
		latest := recs[0] // Trends are sorted by period, first is usually current Month
		maxVal := 0
		vals := []int{latest.StrongBuy, latest.Buy, latest.Hold, latest.Sell, latest.StrongSell}
		for _, v := range vals {
			if v > maxVal {
				maxVal = v
			}
		}
		if maxVal == 0 {
			maxVal = 10 // Prevent scale error
		} else {
			maxVal += 2 // Add headroom
		}
		app.recBar.Values(vals, maxVal)
	}
}
