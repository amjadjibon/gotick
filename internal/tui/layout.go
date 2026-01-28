package tui

import (
	"log"

	"github.com/mum4k/termdash/align"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/terminalapi"
)

func createLayout(t terminalapi.Terminal, app *App) *container.Container {
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
								container.SplitVertical(
									container.Left(
										container.PlaceWidget(app.input),
										container.AlignHorizontal(align.HorizontalLeft),
										container.Border(linestyle.Light),
										container.BorderTitle(" Search "),
									),
									container.Right(
										container.PlaceWidget(app.settingsText),
										container.Border(linestyle.Light),
										container.BorderTitle(" Settings "),
									),
									container.SplitPercent(40),
								),
							),
							container.Bottom(
								container.PlaceWidget(app.lc),
								container.Border(linestyle.Light),
								container.BorderTitle(" Price History "),
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
	return c
}
