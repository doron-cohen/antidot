package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/doron-cohen/antidot/internal/dotfile"
	"github.com/doron-cohen/antidot/internal/rules"
	"github.com/doron-cohen/antidot/internal/shell"
	"github.com/doron-cohen/antidot/internal/tui"
	"github.com/doron-cohen/antidot/internal/utils"
)

func init() {
	cleanCmd.Flags().StringVarP(
		&shellOverride, "shell", "s", "", "Which shell syntax to apply rules in",
	)
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

		utils.ApplyDefaultXdgEnv()

		dotfiles, err := dotfile.Detect(userHomeDir)
		tui.FatalIfError("Failed to detect dotfiles in home dir", err)
		if len(dotfiles) == 0 {
			tui.Print("No dotfiles detected in home directory. You're all clean!")
			return
		}

		tui.Debug("Found %d dotfiles in %s", len(dotfiles), userHomeDir)

		kvStore, err := shell.LoadKeyValueStore("")
		tui.FatalIfError("Failed to load key value store", err)
		actx := rules.ActionContext{KeyValueStore: kvStore}

		appliedRule := false
		for _, dotfile := range dotfiles {
			rule := rules.MatchRule(&dotfile)
			if rule == nil {
				continue
			}

			rule.Pprint()
			if rule.Ignore {
				continue
			}

			confirmed := tui.Confirm(fmt.Sprintf("Apply rule %s?", rule.Name))
			if confirmed {
				rule.Apply(&actx)
				appliedRule = true
			}

			tui.Print("") // one line space
		}

		if appliedRule {
			tui.Print(
				"Cleanup finished - run %s to take effect",
				tui.ApplyStyle(tui.Blue, "eval \"$(antidot init)\""),
			)
		}
	},
}
