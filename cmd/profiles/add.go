/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/

package profiles

import (
	"fmt"
	"os"
	"path/filepath"

	common "github.com/hectorruiz-it/grabber/cmd"
	"gopkg.in/ini.v1"

	"github.com/spf13/cobra"
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
}

func createProfile(newProfile string, basic bool, ssh bool, token bool) {
	common.ReadCheckExistsProfiles(newProfile, true)
	homeDir := common.GetHomeDirectory()

	var appliedProfiles *ini.File
	switch {
	case basic:
		basicProfiles := common.ReadBasicProfilesFile()
		appliedProfiles = basicProfile(basicProfiles, newProfile)
		appliedProfiles.SaveTo(homeDir + common.PROFILES_BASIC_FILE_PATH)
	case ssh:
		sshProfiles := common.ReadSshProfilesFile()
		appliedProfiles = sshProfile(sshProfiles, newProfile)
		appliedProfiles.SaveTo(homeDir + common.PROFILES_SSH_FILE_PATH)
	case token:
		tokenProfiles := common.ReadTokenProfilesFile()
		appliedProfiles = tokenProfile(tokenProfiles, newProfile)
		appliedProfiles.SaveTo(homeDir + common.PROFILES_TOKEN_FILE_PATH)
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

func basicProfile(profiles *ini.File, newProfile string) *ini.File {
	section, err := profiles.NewSection(newProfile)
	common.CheckAndReturnError(err)

loop:
	for {
		var username string
		fmt.Print("> Enter git username: ")
		fmt.Scanln(&username)

		var password string
		fmt.Print("> Enter user password: ")
		fmt.Scanln(&password)

		fmt.Println("grabber: Your profile configuration is the following: ")
		fmt.Println("> Profile ID:", newProfile)
		fmt.Println("> Username:", username)
		fmt.Println("> Password:", password)

		var apply string
		fmt.Print("grabber: do you want to apply this configuration [Y/n]: ")
		fmt.Scanln(&apply)

		switch {
		case apply == "n" || apply == "N":
			fmt.Println("grabber: configuration not applied.")
			continue loop
		default:
			var err error
			_, err = section.NewKey("username", username)
			common.CheckAndReturnError(err)
			_, err = section.NewKey("password", password)
			common.CheckAndReturnError(err)
			break loop
		}
	}
	return profiles
}

func sshProfile(profiles *ini.File, newProfile string) *ini.File {
	section, err := profiles.NewSection(newProfile)
	common.CheckAndReturnError(err)

loop:
	for {
		var privateSshKeyPath string
		fmt.Print("> Enter SSH private key absolute path: ")
		fmt.Scanln(&privateSshKeyPath)

		if err := checkPathAbsoluteExists(privateSshKeyPath, "ssh key"); err != nil {
			common.CheckAndReturnError(err)
		}

		if err := checkPathAbsoluteExists(privateSshKeyPath+".pub", "public ssh key"); err != nil {
			common.CheckAndReturnError(err)
		}

		fmt.Println("grabber: Your profile configuration is the following: ")
		fmt.Println("> Profile ID:", newProfile)
		fmt.Println("> Private Key Path:", privateSshKeyPath)
		fmt.Println("> Public Key Path:", privateSshKeyPath+".pub")

		var apply string
		fmt.Print("grabber: do you want to apply this configuration [Y/n]: ")
		fmt.Scanln(&apply)

		switch {
		case apply == "n" || apply == "N":
			fmt.Println("grabber: configuration not applied.")
			continue loop
		default:
			var err error
			_, err = section.NewKey("private_key", privateSshKeyPath)
			common.CheckAndReturnError(err)
			_, err = section.NewKey("public_key", privateSshKeyPath+".pub")
			common.CheckAndReturnError(err)
			break loop
		}
	}
	return profiles
}

func tokenProfile(profiles *ini.File, newProfile string) *ini.File {
	section, err := profiles.NewSection(newProfile)
	common.CheckAndReturnError(err)

loop:
	for {
		var username string
		fmt.Print("> Enter git username: ")
		fmt.Scanln(&username)

		var token string
		fmt.Print("> Enter user token: ")
		fmt.Scanln(&token)

		fmt.Println("grabber: Your profile configuration is the following: ")
		fmt.Println("> Profile ID:", newProfile)
		fmt.Println("> Username:", username)
		fmt.Println("> Token:", token)

		var apply string
		fmt.Print("grabber: do you want to apply this configuration [Y/n]: ")
		fmt.Scanln(&apply)

		switch {
		case apply == "n" || apply == "N":
			fmt.Println("grabber: configuration not applied.")
			continue loop
		default:
			var err error
			_, err = section.NewKey("username", username)
			common.CheckAndReturnError(err)
			_, err = section.NewKey("token", token)
			common.CheckAndReturnError(err)
			break loop
		}
	}
	return profiles
}
