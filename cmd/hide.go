package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cpendery/figgy/figgy/hiders"
)

var hideCmd = &cobra.Command{
	Use:   "hide [editor]",
	Short: "Hide figgied configs",
	Long: "Hide figgied configs using local settings for the given editor" +
		" at the current depth and all nested levels",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf(
				"requires exactly one arg (available=[%s])",
				strings.Join(cmd.ValidArgs, ","),
			)
		}
		for _, validArg := range cmd.ValidArgs {
			if args[0] == validArg {
				return nil
			}
		}
		return fmt.Errorf(
			"invalid arg specified: %s (available=[%s])",
			args[0],
			strings.Join(cmd.ValidArgs, ","),
		)
	},
	ValidArgs: []string{"vsc"},
	RunE:      hideConfigs,
}

func init() {
	rootCmd.AddCommand(hideCmd)
}

func hideConfigs(_ *cobra.Command, args []string) error {
	switch args[0] {
	case "vsc":
		return hiders.NewVSCodeHider().Hide()
	default:
		return fmt.Errorf("unsupported editor: %s", args[0])
	}
}
