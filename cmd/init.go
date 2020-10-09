package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

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
		if err != nil {
			log.Fatalf("Failed to get env file path: %v", err)
		}

		aliasFilePath, err := utils.GetAliasFile()
		if err != nil {
			log.Fatalf("Failed to get alias file path: %v", err)
		}

		fmt.Printf(`if [[ "$ANTIDOT_INIT" != "1" ]]; then
%s
  source %s
  source %s

  export ANTIDOT_INIT=1
fi`,
			utils.IndentLines(utils.XdgVarsExport()),
			envFilePath,
			aliasFilePath,
		)
	},
}
