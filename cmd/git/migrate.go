// Migrate from JSON file and from DynamoDB
// Add preview before starting cloning repositories <- ta difisi

// recursive and simple mode
package git_actions

import (
	"os"

	common "github.com/hectorruiz-it/grabber/cmd"
	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var trackCmd = &cobra.Command{
	Use:   "comodoro",
	Short: "Adds a git repository to grabber configuration",
	Long:  `Adds a git repository to grabber configuration`,
	Args:  cobra.ExactArgs(1),
	Example: `  grabber track ./ --profile github
  grabber track -r ./* --profile github`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("pull called")
		cmd.Flags().GetBool("recursive")
		cmd.Flags().GetString("profile")

		track()
	},
}

func init() {
	common.RootCmd.AddCommand(trackCmd)
	trackCmd.Flags().StringP("profile", "p", "", "grabber profile to use")
	trackCmd.Flags().BoolP("recursive", "r", false, "add repositories recursively")
	trackCmd.MarkFlagRequired("profile")
}

func track() {
	os.DirFS("./")
}
