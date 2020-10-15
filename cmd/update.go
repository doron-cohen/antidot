package cmd

import (
	"log"

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
		log.Printf("Updating rules...")
		err := utils.Download(rulesSource, rulesFilePath)
		if err != nil {
			log.Fatalf("Failed to update rules: %v", err)
		}

		log.Printf("Rules updated")
	},
}
