package tui

import (
	"log"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/barchart"
	"github.com/mum4k/termdash/widgets/donut"
	"github.com/mum4k/termdash/widgets/linechart"
	"github.com/mum4k/termdash/widgets/text"
	"github.com/mum4k/termdash/widgets/textinput"
)

func createSearchInput(app *App) *textinput.TextInput {
	input, err := textinput.New(
		textinput.Label("Symbol: ", cell.FgColor(cell.ColorNumber(33))),
		textinput.MaxWidthCells(30),
		textinput.PlaceHolder("Enter symbol (e.g. AAPL)"),
		textinput.OnSubmit(func(text string) error {
			if text != "" {
				app.currentSymbol = text
				go app.updateDashboard()
			}
			return nil
		}),
		textinput.ClearOnSubmit(),
	)
	if err != nil {
		log.Fatal(err)
	}
	return input
}

func createPriceChart() *linechart.LineChart {
	lc, err := linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorRed)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorGreen)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorGreen)),
	)
	if err != nil {
		log.Fatal(err)
	}
	return lc
}

func createQuoteText() *text.Text {
	t, err := text.New()
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func createMarketText() *text.Text {
	t, err := text.New()
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func createNewsText() *text.Text {
	t, err := text.New()
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func createRecommendationsBar() *barchart.BarChart {
	bc, err := barchart.New(
		barchart.BarColors([]cell.Color{
			cell.ColorGreen,
			cell.ColorNumber(118),
			cell.ColorYellow,
			cell.ColorRed,
			cell.ColorNumber(88),
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
	return bc
}

func createRangeDonut() *donut.Donut {
	d, err := donut.New(
		donut.CellOpts(cell.FgColor(cell.ColorCyan)),
		donut.Label("Range %", cell.FgColor(cell.ColorWhite)),
	)
	if err != nil {
		log.Fatal(err)
	}
	return d
}

func createSettingsText() *text.Text {
	t, err := text.New()
	if err != nil {
		log.Fatal(err)
	}
	return t
}
