package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

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

		rulesConfig, err := rules.LoadRulesConfig(rulesFilePath)
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

		sh, err := shell.Get(shellOverride)
		tui.FatalIfError("Failed to detect shell", err)

		kvStore, err := shell.LoadKeyValueStore("")
		tui.FatalIfError("Failed to load key value store", err)
		actx := rules.ActionContext{KeyValueStore: kvStore}

		appliedRule := false
		defer func() {
			if appliedRule {
				if err := shell.DumpAliases(sh, kvStore); err != nil {
					tui.Warn("Failed to dump aliases")
				}

				if err = shell.DumpExports(sh, kvStore); err != nil {
					tui.Warn("Failed to dump exports")
				}
			}
		}()

		for _, rule := range rulesConfig.Rules {
			fullPath := filepath.Join(userHomeDir, rule.Dotfile.Name)
			match, err := rule.Dotfile.MatchPath(fullPath)
			if err != nil {
				tui.Warn("Failed to find dotfile for rule %s", rule.Name)
			}

			if match {
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
		}

		if appliedRule {
			tui.Print(
				"Cleanup finished - run %s to take effect",
				tui.ApplyStyle(tui.Blue, "eval \"$(antidot init)\""),
			)
		} else {
			tui.Print("No dotfiles detected in home directory. You're all clean!")
		}
	},
}
