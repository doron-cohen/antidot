package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/doron-cohen/antidot/internal/dotfile"
	"github.com/doron-cohen/antidot/internal/rules"
	"github.com/doron-cohen/antidot/internal/tui"
	"github.com/doron-cohen/antidot/internal/utils"
)

func init() {
	rootCmd.AddCommand(cleanCmd)
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up dotfiles from your $HOME",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Cleaning up!")

		_, err := rules.LoadRulesConfig("rules.yaml")
		if err != nil {
			log.Fatalln("Failed to read rules file: ", err)
		}

		userHomeDir, err := utils.GetHomeDir()
		if err != nil {
			log.Fatalln("Unable to detect user home dir: ", err)
		}

		dotfiles, err := dotfile.Detect(userHomeDir)
		if err != nil {
			log.Fatalln("Failed to detect dotfiles in home dir: ", err)
		}

		log.Printf("Found %d dotfiles in %s\n", len(dotfiles), userHomeDir)

		foundRules := make([]*rules.Rule, 0)
		for _, dotfile := range dotfiles {
			rule := rules.MatchRule(&dotfile)
			if rule == nil {
				continue
			}

			rule.Pprint()
			foundRules = append(foundRules, rule)
		}

		confirmed, err := tui.Confirm("Apply rules?")
		if err != nil {
			log.Fatalf("Failed to read input from stdin: %v", err)
		}

		if !confirmed {
			log.Println("User cancelled. No action was preformed")
			return
		}

		log.Println("Applying rules")
		for _, rule := range foundRules {
			rule.Apply()
		}
	},
}
