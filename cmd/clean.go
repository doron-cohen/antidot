package cmd

import (
	"fmt"
	"os"

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
		tui.Debug("Cleaning up!")

		_, err := rules.LoadRulesConfig(rulesFilePath)
		if err != nil {
			if _, rulesMissing := err.(*rules.MissingRulesFile); rulesMissing {
				tui.Print("Couldn't find rules file. Please run `antidot update`.")
				os.Exit(2)
			}
			tui.FatalIfError("Failed to read rules file", err)
		}

		userHomeDir, err := utils.GetHomeDir()
		tui.FatalIfError("Unable to detect user home dir", err)

		dotfiles, err := dotfile.Detect(userHomeDir)
		tui.FatalIfError("Failed to detect dotfiles in home dir", err)

		tui.Debug("Found %d dotfiles in %s", len(dotfiles), userHomeDir)

		for _, dotfile := range dotfiles {
			rule := rules.MatchRule(&dotfile)
			if rule == nil {
				continue
			}

			rule.Pprint()
			confirmed := tui.Confirm(fmt.Sprintf("Apply rule %s?", rule.Name))
			if confirmed {
				rule.Apply()
			}

			tui.Print("") // one line space
		}
	},
}
