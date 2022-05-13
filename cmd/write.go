package cmd

import (
	"sync"

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
	var wg sync.WaitGroup
	errs := make(chan error, len(figgyConfig))
	for _, fig := range figgyConfig {
		wg.Add(1)
		go func(f interface{}) {
			defer wg.Done()
			err := config.WriteConfig(f.(map[string]interface{}))
			if err != nil {
				errs <- err
			}
		}(fig)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		return err
	}

	return nil
}
