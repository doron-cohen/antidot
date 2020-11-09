package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/doron-cohen/antidot/internal/tui"
	"github.com/doron-cohen/antidot/internal/utils"
)

var rulesFilePath string

var rootCmd = &cobra.Command{
	Use:   "antidot",
	Short: "Clean your $HOME from those pesky dotfiles",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&tui.Verbose, "verbose", "v", false, "Output verbosity")
	rootCmd.PersistentFlags().StringVarP(&rulesFilePath, "rules", "r", utils.GetRulesFilePath(), "Rules file path")
}

func Execute(appVersion string) {
	rootCmd.Version = appVersion
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
