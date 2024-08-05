/*
Copyright © 2024 HÉCTOR <EMAIL ADDRESS>
*/
package grabber

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pull called")
		pull()
	},
}

func init() {
	RootCmd.AddCommand(pullCmd)
}

func pull() {
	repository, err := git.PlainOpen("./")
	CheckAndReturnError(err)

	w, err := repository.Worktree()
	CheckAndReturnError(err)

	err = w.Pull(&git.PullOptions{})

	switch err {
	case git.NoErrAlreadyUpToDate:
		fmt.Println("grabber: already up to date")
		os.Exit(0)
	default:
		CheckAndReturnError(err)
	}
}
