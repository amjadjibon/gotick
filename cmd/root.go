package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/amjadjibon/gotick/internal/tui"
)

var (
	symbol    string
	interval  string
	timeRange string
)

func init() {
	rootCmd.Flags().StringVarP(&symbol, "symbol", "s", "AAPL", "Stock symbol to display")
	rootCmd.Flags().StringVarP(&interval, "interval", "i", "1d", "Chart interval (e.g. 1d, 1h, 5m)")
	rootCmd.Flags().StringVarP(&timeRange, "range", "r", "1y", "Chart time range (e.g. 1y, 5d, 1mo)")
}

var rootCmd = &cobra.Command{
	Use:   "gotick",
	Short: "Real-time terminal stock ticker",
	Long: `A terminal-based stock ticker and dashboard using Yahoo Finance data.
Displays real-time price, history chart, market summary, news, and analyst recommendations.`,
	Run: func(cmd *cobra.Command, args []string) {
		tui.Run(tui.Options{
			Symbol:   symbol,
			Interval: interval,
			Range:    timeRange,
		})
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
