package cmd

import (
	"github.com/spf13/cobra"

	sh "github.com/doron-cohen/antidot/internal/shell"
	"github.com/doron-cohen/antidot/internal/tui"
)

func init() {
	initCmd.Flags().StringVarP(
		&shellOverride, "shell", "s", "", "Which shell to render the init script to",
	)
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize antidot for aliases and env vars to work",
	Run: func(cmd *cobra.Command, args []string) {
		shell, err := sh.Get(shellOverride)
		tui.FatalIfError("", err)

		kv, err := sh.LoadKeyValueStore("")
		tui.FatalIfError("", err)

		tui.Print(shell.RenderInit(kv))
	},
}
