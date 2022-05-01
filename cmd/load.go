package cmd

import (
	"errors"
	"fmt"

	"github.com/cpendery/figgy/figgy/config"
	"github.com/spf13/cobra"
)

var loadCmd = &cobra.Command{
	Use:   "load [FILE_NAMES]",
	Short: "Figgy the given files",
	Long: "Load the given files into figgied representations in your " +
		"figgy config",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one files to load")
		}
		return nil
	},
	RunE: loadFilesIntoFiggy,
}

func init() {
	rootCmd.AddCommand(loadCmd)
}

func loadFilesIntoFiggy(_ *cobra.Command, args []string) error {
	figgyConfig, err := config.ReadFiggyConfig(config.FiggyConfigName)
	if err != nil {
		return err
	}
	for _, filePath := range args {
		fig, err := config.ReadConfig(filePath)
		if err != nil {
			return fmt.Errorf("Failed to load/parse file: %s", filePath)
		}
		figgyConfig[filePath] = fig
	}
	return config.WriteFiggyConfig(figgyConfig)
}
