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

var listProfiles = &cobra.Command{
	Use:   "list-profiles",
	Short: "List grabber configured authentication methods.",
	Long:  `List grabber configured authentication methods`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		getProfiles()
	},
}

func init() {
	common.RootCmd.AddCommand(listProfiles)
}

func getProfiles() {
	headerFmt := color.New(color.Underline).SprintfFunc()

	basicProfiles := common.ReadBasiceProfilesFile()
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
	tbl.WithHeaderFormatter(headerFmt)
	for _, profile := range profilesList {
		tbl.AddRow(profile, profilesMap[profile])
	}
	tbl.Print()
}
