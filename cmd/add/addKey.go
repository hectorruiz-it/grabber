/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package add

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addKeyCmd represents the addKey command
var newProfile = &cobra.Command{
	Use:   "new-profile",
	Short: "Creates a new grabber profile",
	Long: `Creates a new grabber profile. A profile is an ID that has mapped a pair of
	private an public keys.A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("addKey called")
	},
}

func init() {
	addCmd.AddCommand(newProfile)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addKeyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addKeyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
