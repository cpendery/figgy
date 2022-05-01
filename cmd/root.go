package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "figgy",
	Short: "A tool for unifying your configs",
	Long: `A command line utility for unifying all your configuration files.
	Complete documentation is available at https://github.com/cpendery/figgy`,
	Args: cobra.MinimumNArgs(1),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
