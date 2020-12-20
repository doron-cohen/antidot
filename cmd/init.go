package cmd

import (
	"github.com/spf13/cobra"

	"github.com/doron-cohen/antidot/internal/tui"
	"github.com/doron-cohen/antidot/internal/utils"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize antidot for aliases and env vars to work",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: detect shell and generate appropriate script
		envFilePath, err := utils.GetEnvFile()
		tui.FatalIfError("Failed to get env file path", err)

		aliasFilePath, err := utils.GetAliasFile()
		tui.FatalIfError("Failed to get alias file path", err)

		tui.Print(`%s

if [ -f "%s" ]; then source "%s"; fi
if [ -f "%s" ]; then source "%s"; fi`,
			utils.XdgVarsExport(),
			envFilePath,
			envFilePath,
			aliasFilePath,
			aliasFilePath,
		)
	},
}
