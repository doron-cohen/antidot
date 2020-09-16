package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/doron-cohen/antidot/internal/action"
	"github.com/doron-cohen/antidot/internal/dirs"
	"github.com/doron-cohen/antidot/internal/dotfile"
)

func init() {
	rootCmd.AddCommand(cleanCmd)
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up dotfiles from your $HOME",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Cleaning up!")
		go action.LoadRulesConfig("rules.yaml")
		userHomeDir, err := dirs.GetHomeDir()
		if err != nil {
			log.Fatalln(err)
		}

		dotfiles, err := dotfile.Detect(userHomeDir)
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Found %d dotfiles in %s\n", len(dotfiles), userHomeDir)
		// TODO: block here until LoadRulesConfig succeeds
		for _, dotfile := range dotfiles {
			action.MatchActions(&dotfile)
		}
	},
}
