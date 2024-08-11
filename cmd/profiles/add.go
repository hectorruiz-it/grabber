/*
Copyright © 2024 HÉCTOR <EMAIL ADDRESS>
*/

package profiles

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	common "github.com/hectorruiz-it/grabber/cmd"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
	"gopkg.in/ini.v1"
)

var addCommand = &cobra.Command{
	Use:   "add-profile profile",
	Short: "Creates a new grabber profile",
	Long: `Creates a new grabber profile. A profile is an ID that has mapped a pair of
private an public keys, token or user auth`,
	Args:    cobra.ExactArgs(1),
	Example: `  grabber add-profile github --ssh`,
	Run: func(cmd *cobra.Command, args []string) {
		basic, err := cmd.Flags().GetBool("basic")
		common.CheckAndReturnError(err)

		ssh, err := cmd.Flags().GetBool("ssh")
		common.CheckAndReturnError(err)

		token, err := cmd.Flags().GetBool("token")
		common.CheckAndReturnError(err)
		createProfile(args[0], basic, ssh, token)

	},
}

func init() {
	common.RootCmd.AddCommand(addCommand)
	addCommand.Flags().Bool("ssh", false, "add ssh key authentication")
	addCommand.Flags().Bool("token", false, "add token based authentication")
	addCommand.Flags().Bool("basic", false, "add user password authentication")
	addCommand.MarkFlagsMutuallyExclusive("ssh", "token", "basic")
	addCommand.MarkFlagsOneRequired("ssh", "token", "basic")
}

func createProfile(newProfile string, basic bool, ssh bool, token bool) {
	common.ReadCheckExistsProfiles(newProfile, true)
	homeDir := common.GetHomeDirectory()

	var appliedProfiles *ini.File
	model := profileTui(basic, ssh, token)

	switch {
	case basic:
		basicProfiles := common.ReadBasicProfilesFile()
		appliedProfiles = addConfigUsernameProfile(basicProfiles, newProfile, model)
		appliedProfiles.SaveTo(homeDir + common.PROFILES_BASIC_FILE_PATH)
		storeOnKeychain(newProfile, model.inputs[1].Value())
		addProfileToConfig(newProfile, model.profileType)
	case ssh:
		sshProfiles := common.ReadSshProfilesFile()
		appliedProfiles = addConfigSshProfile(sshProfiles, newProfile, model)
		appliedProfiles.SaveTo(homeDir + common.PROFILES_SSH_FILE_PATH)
		storeOnKeychain(newProfile, model.inputs[2].Value())
		addProfileToConfig(newProfile, model.profileType)
	case token:
		tokenProfiles := common.ReadTokenProfilesFile()
		appliedProfiles = addConfigUsernameProfile(tokenProfiles, newProfile, model)
		appliedProfiles.SaveTo(homeDir + common.PROFILES_TOKEN_FILE_PATH)
		storeOnKeychain(newProfile, model.inputs[1].Value())
		addProfileToConfig(newProfile, model.profileType)
	}

	fmt.Println("grabber: " + newProfile + " profile configured.")
}

func checkProfileExists(profiles *ini.File, newProfile string) error {
	switch profiles.HasSection(newProfile) {
	case true:
		err := fmt.Errorf("grabber: profile already exists.")
		return err
	}
	return nil
}

func checkPathAbsoluteExists(keyPath string, id string) error {
	if !filepath.IsAbs(keyPath) {
		return fmt.Errorf("grabber: path is not absolute.")
	}

	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return fmt.Errorf("grabber: " + id + " does not exist on the given path.")
	}
	return nil
}

func addConfigSshProfile(profiles *ini.File, newProfile string, model model) *ini.File {
	section, err := profiles.NewSection(newProfile)
	common.CheckAndReturnError(err)

	_, err = section.NewKey("private_key", model.inputs[0].Value())
	common.CheckAndReturnError(err)
	_, err = section.NewKey("public_key", model.inputs[1].Value())
	common.CheckAndReturnError(err)

	return profiles
}

func addConfigUsernameProfile(profiles *ini.File, newProfile string, model model) *ini.File {
	section, err := profiles.NewSection(newProfile)
	common.CheckAndReturnError(err)

	_, err = section.NewKey("username", model.inputs[0].Value())
	common.CheckAndReturnError(err)

	return profiles
}

func storeOnKeychain(profile string, password string) {
	service := "grabber"

	err := keyring.Set(service, profile+"-profile", password)
	common.CheckAndReturnError(err)
}

func addProfileToConfig(profile string, profileType string) {
	config := common.ReadGrabberConfig()
	homeDir := common.GetHomeDirectory()

	config.Profiles = append(config.Profiles, common.Profile{
		Profile:      profile,
		Type:         profileType,
		Repositories: []common.Repository{},
	})

	data, err := json.Marshal(config)
	common.CheckAndReturnError(err)

	err = os.WriteFile(homeDir+common.MAPPINGS_FILE_PATH, data, 0700)
	common.CheckAndReturnError(err)
}
