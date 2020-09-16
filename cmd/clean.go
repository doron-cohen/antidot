package cmd

import (
	"fmt"

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
		var userHomeDir = dirs.GetHomeDir()
		var dotfiles = dotfile.Detect(userHomeDir)
		fmt.Println(dotfiles)
	},
}
