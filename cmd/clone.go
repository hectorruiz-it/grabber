/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package grabber

import (
	"fmt"

	"github.com/spf13/cobra"
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clones a git repository using your credentials configuration",
	Long:  `Clones a git repository using your credentials configuration`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("clone called")
	},
}

/*
grabber clone --profile
*/

func init() {
	RootCmd.AddCommand(cloneCmd)
}

func clone(repository string) {
	fmt.Println(repository)
	// sshRegex := regexp.MustCompile(`^git@`)
	// httpsRegex := regexp.MustCompile(`^https://`)

	// switch {
	// case sshRegex.MatchString(repository):
	// 	git.PlainClone(".", false, &git.CloneOptions{
	// 		URL:  repository,
	// 		Auth: transport.AuthMethod,
	// 	})
	// }
}
