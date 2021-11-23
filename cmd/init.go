package cmd

import (
	"github.com/spf13/cobra"

	sh "github.com/doron-cohen/antidot/internal/shell"
	"github.com/doron-cohen/antidot/internal/tui"
)

func init() {
	initCmd.Flags().StringVarP(
		&shellOverride, "shell", "s", "", "What shell to print an init script for - One of: bash zsh fish",
	)
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Print shell code to initialize aliases and environment variables based on your current shell, unless -s is passed",
	Run: func(cmd *cobra.Command, args []string) {
		shell, err := sh.Get(shellOverride)
		tui.FatalIfError("", err)
		script, err := sh.GetShellScript(shell)
		tui.FatalIfError("", err)
		tui.Print(script)
	},
}
