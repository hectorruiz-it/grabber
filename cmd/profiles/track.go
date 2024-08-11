// recursive and simple mode
package profiles

import (
	"errors"
	"os"
	"path"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	common "github.com/hectorruiz-it/grabber/cmd"
	track_tui "github.com/hectorruiz-it/grabber/tui/track"

	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var trackCmd = &cobra.Command{
	Use:   "track",
	Short: "Adds a git repository to grabber configuration",
	Long:  `Adds a git repository to grabber configuration`,
	Example: `  grabber track ./ --profile github
  grabber track -r ./* --profile github`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("pull called")
		profile, err := cmd.Flags().GetString("profile")
		common.CheckAndReturnError(err)

		var repositories []common.Repository

		for _, directory := range args {
			switch path.IsAbs(directory) {
			case true:
				r := getRepositoriesAbsolutePath(directory)
				for _, id := range r {
					repositories = append(repositories, id)
				}
			case false:
				absolutePathDirectory, err := filepath.Abs(directory)
				common.CheckAndReturnError(err)
				r := getRepositoriesAbsolutePath(absolutePathDirectory)
				for _, id := range r {
					repositories = append(repositories, id)
				}
			}
		}

		// fmt.Println("Full repository list: ", repositories)
		if len(repositories) > 0 {
			track_tui.Track(repositories, profile)
		}
	},
}

func init() {
	common.RootCmd.AddCommand(trackCmd)
	trackCmd.Flags().StringP("profile", "p", "", "grabber profile to use")
	trackCmd.MarkFlagRequired("profile")
}

func getRepositoriesAbsolutePath(directory string) []common.Repository {
	var repositories []common.Repository

	// fmt.Println("Received path: ", directory)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err = errors.New("grabber: No such directory `" + directory + "`")
		common.CheckAndReturnError(err)
	}

	// Check if received path is a repository
	r, err := git.PlainOpen(directory)
	if err == nil {
		remotes, err := r.Remotes()

		if err == nil && len(remotes) > 0 {
			remoteURLs := remotes[0].Config().URLs
			repositories = append(repositories, common.Repository{
				Path: directory,
				Name: remoteURLs[0],
			})
			return repositories
		}
	}

	files, err := os.ReadDir(directory)
	common.CheckAndReturnError(err)

	var directories []string

	for _, file := range files {
		if file.IsDir() {
			// fmt.Println("Readed path: ", file.Name())
			path := path.Join(directory, file.Name())
			directories = append(directories, path)
		}
	}

	// fmt.Println("Full directory list:", directories)

	switch {
	case len(directories) == 0:
		err := errors.New("grabber: no directories found on the given path.")
		common.CheckAndReturnError(err)
	default:
		for _, directory := range directories {
			// fmt.Println(directory)
			r, err := git.PlainOpen(directory)
			common.CheckAndReturnError(err)

			remotes, err := r.Remotes()
			common.CheckAndReturnError(err)

			if err == nil && len(remotes) > 0 {
				remoteURLs := remotes[0].Config().URLs
				repositories = append(repositories, common.Repository{
					Path: directory,
					Name: remoteURLs[0],
				})
			} else {
				continue
				// fmt.Println("grabber: directory `" + directory + "` is not a repository. Skipping it")
			}
		}
	}
	return repositories
}
