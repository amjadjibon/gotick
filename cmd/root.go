package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func init() {
	// TODO: this is a placeholder for now for adding subcommands and flags
}

var rootCmd = &cobra.Command{
	Use: "gotick",
	// TODO: change short and long descriptions to something more meaningful
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmd.Help(); err != nil {
			return err
		}
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
