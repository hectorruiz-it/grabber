/*
Copyright © 2024 HÉCTOR <EMAIL ADDRESS>
*/
package git_actions

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	common "github.com/hectorruiz-it/grabber/cmd"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		push()
	},
}

func init() {
	common.RootCmd.AddCommand(pushCmd)
}

func push() {
	r, err := git.PlainOpen(".")
	common.CheckAndReturnError(err)

	remotes, err := r.Remotes()
	common.CheckAndReturnError(err)

	remoteURLs := remotes[0].Config().URLs
	// fmt.Println(remoteURLs)
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

		err = r.Push(&git.PushOptions{
			Auth:     publicKeys,
			Progress: os.Stdout,
		})

		common.CheckAndReturnError(err)

	case "token":
		password, err := keyring.Get("grabber", profile+"-profile")
		common.CheckAndReturnError(err)

		err = r.Push(&git.PushOptions{
			Auth: &http.BasicAuth{
				Username: "git",
				Password: password,
			},
			Progress: os.Stdout,
		})
		common.CheckAndReturnError(err)
	}

}
