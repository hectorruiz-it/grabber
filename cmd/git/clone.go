/*
Copyright © 2024 HÉCTOR <EMAIL ADDRESS>
*/
package git_actions

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	common "github.com/hectorruiz-it/grabber/cmd"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

var (
	AUTH_METHODS = [...]string{"basic", "ssh", "token"}
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clones a git repository using your credentials configuration",
	Long:  `Clones a git repository using your credentials configuration`,
	// Args:  cobra.ExactArgs(1),
	Example: `  grabber clone git@github.com:lockedinspace/letme.git --profile github
  grabber clone https://github.com/lockedinspace/letme.git --profile github-with-token`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}

		profile, err := cmd.Flags().GetString("profile")
		common.CheckAndReturnError(err)

		clone(args[0], profile)
	},
}

func init() {
	common.RootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().StringP("profile", "p", "", "grabber profile to use")
}

func clone(repository string, profile string) {
	check, profileMap := common.ReadCheckExistsProfiles(profile, false)
	if !check {
		err := errors.New("grabber: profile `" + profile + "` does not exist")
		common.CheckAndReturnError(err)
	}

	sshRegex := regexp.MustCompile(`^git@`)
	httpsRegex := regexp.MustCompile(`^https://`)

	id := strings.Split(repository, "/")
	directoryDotgit := id[len(id)-1]
	directory := strings.Split(directoryDotgit, ".")
	service := "grabber"

	switch {
	case httpsRegex.MatchString(repository):
		password, err := keyring.Get(service, profile+"-profile")
		common.CheckAndReturnError(err)

		_, err = git.PlainClone("./"+directory[0], false, &git.CloneOptions{
			Auth: &http.BasicAuth{
				Username: "git",
				Password: password,
			},
			URL:      repository,
			Progress: os.Stdout,
		})
		common.CheckAndReturnError(err)
		addRepositoryToConfig(profile, repository, directory[0])

	case sshRegex.MatchString(repository):
		if profileMap[AUTH_METHODS[1]] {
			sshProfiles := common.ReadSshProfilesFile()
			privateKey, err := sshProfiles.Section(profile).GetKey("private_key")
			common.CheckAndReturnError(err)

			password, err := keyring.Get(service, profile+"-profile")
			common.CheckAndReturnError(err)

			publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKey.Value(), password)
			common.CheckAndReturnError(err)

			_, err = git.PlainClone("./"+directory[0], false, &git.CloneOptions{
				Auth:     publicKeys,
				Progress: os.Stdout,
				URL:      repository,
			})
			common.CheckAndReturnError(err)
			addRepositoryToConfig(profile, repository, directory[0])
		} else {
			err := errors.New("grabber: profile `" + profile + "` is not an ssh profile.")
			common.CheckAndReturnError(err)
		}
	}
}

func addRepositoryToConfig(profile string, repository string, directory string) {
	config := common.ReadGrabberConfig()
	homeDir := common.GetHomeDirectory()
	currentDir, err := os.Getwd()
	common.CheckAndReturnError(err)

	for i := range config.Profiles {
		if config.Profiles[i].Profile == profile {
			config.Profiles[i].Repositories = append(config.Profiles[i].Repositories, common.Repository{
				Path: path.Join(currentDir, directory),
				Name: repository,
			})
		}
	}

	data, err := json.Marshal(config)
	common.CheckAndReturnError(err)
	err = os.WriteFile(homeDir+common.MAPPINGS_FILE_PATH, data, 0700)
	common.CheckAndReturnError(err)
}
