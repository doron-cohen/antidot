package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/doron-cohen/antidot/internal/dotfile"
	"github.com/doron-cohen/antidot/internal/rules"
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

		err := rules.LoadRulesConfig("rules.yaml")
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

		// TODO: ease eye sore
		for _, dotfile := range dotfiles {
			rule := rules.MatchRule(&dotfile)
			if rule == nil {
				continue
			}

			rule.Pprint()
			// TODO: move to Rule.Apply()
			if !rule.Ignore {
				for _, action := range rule.Actions {
					action.Apply()
				}
			}
		}
	},
}
