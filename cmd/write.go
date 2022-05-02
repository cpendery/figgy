package cmd

import (
	"github.com/cpendery/figgy/figgy/config"
	"github.com/spf13/cobra"
)

var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "Write figgied configs",
	Long:  "Write figgied configs out to the original files specified",
	RunE:  writeFiggiedFiles,
}

func init() {
	rootCmd.AddCommand(writeCmd)
}

func writeFiggiedFiles(_ *cobra.Command, _ []string) error {
	figgyConfig, err := config.ReadFiggyConfig(config.FiggyConfigName)
	if err != nil {
		return err
	}
	for _, fig := range figgyConfig {
		err := config.WriteConfig(fig.(map[string]interface{}))
		if err != nil {
			return err
		}
	}
	return nil
}
