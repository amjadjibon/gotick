package tui

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/barchart"
	"github.com/mum4k/termdash/widgets/donut"
	"github.com/mum4k/termdash/widgets/linechart"
	"github.com/mum4k/termdash/widgets/text"
	"github.com/mum4k/termdash/widgets/textinput"
)

type Options struct {
	Symbol   string
	Interval string
	Range    string
}

// Valid ranges and intervals for Yahoo Finance
var (
	validRanges    = []string{"1d", "5d", "1mo", "3mo", "6mo", "1y", "2y", "5y", "ytd", "max"}
	validIntervals = []string{"1m", "5m", "15m", "30m", "1h", "1d", "1wk", "1mo"}

	// Valid intervals for each range (Yahoo Finance limitations)
	rangeToValidIntervals = map[string][]string{
		"1d":  {"1m", "5m", "15m", "30m", "1h"},
		"5d":  {"1m", "5m", "15m", "30m", "1h", "1d"},
		"1mo": {"5m", "15m", "30m", "1h", "1d"},
		"3mo": {"1h", "1d", "1wk"},
		"6mo": {"1h", "1d", "1wk"},
		"1y":  {"1d", "1wk", "1mo"},
		"2y":  {"1d", "1wk", "1mo"},
		"5y":  {"1wk", "1mo"},
		"ytd": {"1d", "1wk", "1mo"},
		"max": {"1wk", "1mo"},
	}
)

// getValidIntervalsForRange returns the valid intervals for a given range
func getValidIntervalsForRange(r string) []string {
	if intervals, ok := rangeToValidIntervals[r]; ok {
		return intervals
	}
	return []string{"1d"} // fallback
}

// isValidIntervalForRange checks if an interval is valid for a given range
func isValidIntervalForRange(interval, rangeVal string) bool {
	validIntervals := getValidIntervalsForRange(rangeVal)
	for _, v := range validIntervals {
		if v == interval {
			return true
		}
	}
	return false
}

// findIndexInSlice finds the index of a string in a slice, returns -1 if not found
func findIndexInSlice(slice []string, val string) int {
	for i, v := range slice {
		if v == val {
			return i
		}
	}
	return -1
}

type App struct {
	ctx             context.Context
	cancel          context.CancelFunc
	currentSymbol   string
	currentInterval string
	currentRange    string
	rangeIdx        int
	intervalIdx     int

	input        *textinput.TextInput
	lc           *linechart.LineChart
	quoteText    *text.Text
	marketText   *text.Text
	newsText     *text.Text
	recBar       *barchart.BarChart
	rangeDonut   *donut.Donut
	settingsText *text.Text
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

	// Find initial indices for range and interval
	for i, r := range validRanges {
		if r == app.currentRange {
			app.rangeIdx = i
			break
		}
	}
	for i, inv := range validIntervals {
		if inv == app.currentInterval {
			app.intervalIdx = i
			break
		}
	}

	t, err := tcell.New()
	if err != nil {
		log.Fatal(err)
	}
	defer t.Close()

	app.ctx, app.cancel = context.WithCancel(context.Background())
	defer app.cancel()

	app.initWidgets()
	c := createLayout(t, app)

	go app.updateDashboard()
	app.updateSettings()

	ticker := time.NewTicker(30 * time.Second)
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

	keyHandler := func(k *terminalapi.Keyboard) {
		switch k.Key {
		case 'q', keyboard.KeyEsc:
			app.cancel()
		case 'r':
			// Next range
			app.rangeIdx = (app.rangeIdx + 1) % len(validRanges)
			app.currentRange = validRanges[app.rangeIdx]
			// Auto-adjust interval if not valid for new range
			if !isValidIntervalForRange(app.currentInterval, app.currentRange) {
				validInts := getValidIntervalsForRange(app.currentRange)
				app.currentInterval = validInts[0]
				app.intervalIdx = findIndexInSlice(validIntervals, app.currentInterval)
			}
			app.updateSettings()
			go app.updateDashboard()
		case 'R':
			// Previous range
			app.rangeIdx = (app.rangeIdx - 1 + len(validRanges)) % len(validRanges)
			app.currentRange = validRanges[app.rangeIdx]
			// Auto-adjust interval if not valid for new range
			if !isValidIntervalForRange(app.currentInterval, app.currentRange) {
				validInts := getValidIntervalsForRange(app.currentRange)
				app.currentInterval = validInts[0]
				app.intervalIdx = findIndexInSlice(validIntervals, app.currentInterval)
			}
			app.updateSettings()
			go app.updateDashboard()
		case 'i':
			// Next valid interval for current range
			validInts := getValidIntervalsForRange(app.currentRange)
			currentIdx := findIndexInSlice(validInts, app.currentInterval)
			if currentIdx == -1 {
				currentIdx = 0
			}
			nextIdx := (currentIdx + 1) % len(validInts)
			app.currentInterval = validInts[nextIdx]
			app.intervalIdx = findIndexInSlice(validIntervals, app.currentInterval)
			app.updateSettings()
			go app.updateDashboard()
		case 'I':
			// Previous valid interval for current range
			validInts := getValidIntervalsForRange(app.currentRange)
			currentIdx := findIndexInSlice(validInts, app.currentInterval)
			if currentIdx == -1 {
				currentIdx = 0
			}
			prevIdx := (currentIdx - 1 + len(validInts)) % len(validInts)
			app.currentInterval = validInts[prevIdx]
			app.intervalIdx = findIndexInSlice(validIntervals, app.currentInterval)
			app.updateSettings()
			go app.updateDashboard()
		}
	}

	if err := termdash.Run(app.ctx, t, c, termdash.KeyboardSubscriber(keyHandler)); err != nil {
		// Note: This will log the error but ticker.Stop() won't run if we use Fatal.
		// In practice, this is acceptable since the app is terminating anyway.
		panic(err)
	}
}

func (app *App) updateSettings() {
	_ = app.settingsText.Write(
		fmt.Sprintf("Range: %s | Interval: %s | [r/R] Range [i/I] Interval [q] Quit",
			app.currentRange, app.currentInterval),
		text.WriteReplace(),
	)
}

func (app *App) initWidgets() {
	app.input = createSearchInput(app)
	app.lc = createPriceChart()
	app.quoteText = createQuoteText()
	app.marketText = createMarketText()
	app.newsText = createNewsText()
	app.recBar = createRecommendationsBar()
	app.rangeDonut = createRangeDonut()
	app.settingsText = createSettingsText()
}
