package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "antidot",
	Short: "Clean your $HOME from those pesky dotfiles",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute(appVersion string) {
	rootCmd.Version = appVersion
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
