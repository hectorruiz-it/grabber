/*
Copyright © 2024 NAME HERE list<EMAIL ADDRESS>
*/
package profiles

import (
	"fmt"
	"sort"

	"github.com/fatih/color"
	common "github.com/hectorruiz-it/grabber/cmd"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

var list = &cobra.Command{
	Use:   "list",
	Short: "List grabber configured authentication methods and repositories.",
	Long:  `List grabber configured authentication methods and repositories.`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var listProfiles = &cobra.Command{
	Use:   "profiles",
	Short: "List grabber configured profiles.",
	Long:  `List grabber configured profiles.`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		getProfiles()
	},
}

var listRepositories = &cobra.Command{
	Use:   "repositories",
	Short: "List grabber repositories using a given profile.",
	Long:  `List grabber repositories using a given profile.`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		profile, err := cmd.Flags().GetString("profile")
		common.CheckAndReturnError(err)
		err = fmt.Errorf("grabber: profile is not configured.")
		if !common.ReadCheckExistsProfiles(profile, false) {
			common.CheckAndReturnError(err)
		}

		var repositories []string
		config := common.ReadGrabberConfig()
		for _, id := range config.Profiles {
			// fmt.Println(id.Profile)
			if id.Profile == profile {
				repositories = id.Repositories
			}
		}
		// fmt.Println(repositories)
		if len(repositories) == 0 {
			err = fmt.Errorf("grabber: no repositories cloned yet with " + profile + " profile.")
			common.CheckAndReturnError(err)
		}
		sort.Strings(repositories)
		tbl := table.New("Repositories")
		headerFmt := color.New(color.Underline).SprintfFunc()
		tbl.WithHeaderFormatter(headerFmt)
		for _, id := range repositories {
			tbl.AddRow(id)
		}
		tbl.Print()

	},
}

func init() {
	common.RootCmd.AddCommand(list)
	list.AddCommand(listProfiles)
	list.AddCommand(listRepositories)
	listRepositories.Flags().StringP("profile", "p", "", "grabber profile name")
	listRepositories.MarkFlagRequired("profile")
}

func getProfiles() {

	basicProfiles := common.ReadBasicProfilesFile()
	sshProfiles := common.ReadSshProfilesFile()
	tokenProfiles := common.ReadTokenProfilesFile()

	profilesMap := make(map[string]string)
	profilesList := make([]string, 0)

	for _, profile := range basicProfiles.SectionStrings() {
		if profile != "DEFAULT" {
			profilesList = append(profilesList, profile)
			profilesMap[profile] = "basic"
		}
	}

	for _, profile := range sshProfiles.SectionStrings() {
		if profile != "DEFAULT" {
			profilesList = append(profilesList, profile)
			profilesMap[profile] = "ssh"
		}
	}

	for _, profile := range tokenProfiles.SectionStrings() {
		if profile != "DEFAULT" {
			profilesList = append(profilesList, profile)
			profilesMap[profile] = "token"
		}
	}

	if len(profilesList) == 0 {
		err := fmt.Errorf("grabber: no profiles detected. Run `grabber new-profile profile [flags]` to create one.")
		common.CheckAndReturnError(err)
	}

	sort.Strings(profilesList)
	tbl := table.New("Profile", "AuthMethod")
	headerFmt := color.New(color.Underline).SprintfFunc()
	tbl.WithHeaderFormatter(headerFmt)
	for _, profile := range profilesList {
		tbl.AddRow(profile, profilesMap[profile])
	}
	tbl.Print()
}
