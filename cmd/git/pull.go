/*
Copyright © 2024 HÉCTOR <EMAIL ADDRESS>
*/
package git_actions

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	common "github.com/hectorruiz-it/grabber/cmd"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
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
		// fmt.Println("pull called")
		pull()
	},
}

func init() {
	common.RootCmd.AddCommand(pullCmd)
}

func pull() {
	r, err := git.PlainOpen("./")
	common.CheckAndReturnError(err)

	remotes, err := r.Remotes()
	common.CheckAndReturnError(err)
	remoteURLs := remotes[0].Config().URLs

	w, err := r.Worktree()

	common.CheckAndReturnError(err)

	profile, authMethod, err := common.GetProfileByRepository(remoteURLs[0])
	common.CheckAndReturnError(err)

	switch authMethod {
	case "ssh":
		sshProfiles := common.ReadSshProfilesFile()
		section := sshProfiles.Section(profile)
		privateKey := section.Key("private_key")

		password, err := keyring.Get("grabber", profile+"-profile")
		common.CheckAndReturnError(err)

		publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKey.Value(), password)
		common.CheckAndReturnError(err)
		err = w.Pull(&git.PullOptions{
			Auth:     publicKeys,
			Progress: os.Stdout,
		})

		switch err {
		case nil:
			break
		case git.NoErrAlreadyUpToDate:
			fmt.Println("grabber: already up to date")
			os.Exit(0)
		default:
			err = fmt.Errorf("grabber: %w", err)
			common.CheckAndReturnError(err)
		}
	case "token":
		password, err := keyring.Get("grabber", profile+"-profile")
		common.CheckAndReturnError(err)
		err = w.Pull(&git.PullOptions{
			Auth: &http.BasicAuth{
				Username: "git",
				Password: password,
			},
			Progress: os.Stdout,
		})

		switch err {
		case nil:
			break
		case git.NoErrAlreadyUpToDate:
			fmt.Println("grabber: already up to date")
			os.Exit(0)
		default:
			err = fmt.Errorf("grabber: %w", err)
			common.CheckAndReturnError(err)
		}
	}
}
