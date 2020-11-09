package cmd

import (
	"github.com/doron-cohen/antidot/internal/tui"
	"github.com/doron-cohen/antidot/internal/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var rulesSource = "https://raw.githubusercontent.com/doron-cohen/antidot/master/rules.yaml"

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update rules file",
	Run: func(cmd *cobra.Command, args []string) {
		tui.Debug("Updating rules...")
		err := utils.Download(rulesSource, rulesFilePath)
		tui.FatalIfError("Failed to update rules", err)

		tui.Print("Rules updated successfully")
	},
}
