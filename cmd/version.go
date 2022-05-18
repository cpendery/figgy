package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of figgy",
	Long:  `All software has versions. This is figgy's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("figgy v0.0.0-beta.1")
	},
}
