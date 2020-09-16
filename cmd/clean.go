package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

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
		fmt.Println("Cleaning up!")
		userHomeDir, err := dirs.GetHomeDir()
		if err != nil {
			log.Fatalln(err)
		}

		dotfiles, err := dotfile.Detect(userHomeDir)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(dotfiles)
	},
}
