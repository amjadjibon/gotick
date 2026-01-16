package tui

import (
	"context"
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

type App struct {
	ctx             context.Context
	cancel          context.CancelFunc
	currentSymbol   string
	currentInterval string
	currentRange    string

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

	app.initWidgets()
	c := createLayout(t, app)

	go app.updateDashboard()

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

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == keyboard.KeyEsc {
			app.cancel()
		}
	}

	if err := termdash.Run(app.ctx, t, c, termdash.KeyboardSubscriber(quitter)); err != nil {
		log.Fatal(err)
	}
}

func (app *App) initWidgets() {
	app.input = createSearchInput(app)
	app.lc = createPriceChart()
	app.quoteText = createQuoteText()
	app.marketText = createMarketText()
	app.newsText = createNewsText()
	app.recBar = createRecommendationsBar()
	app.rangeDonut = createRangeDonut()
}
