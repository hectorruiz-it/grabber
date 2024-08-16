// Migrate from JSON file and from DynamoDB
// Add preview before starting cloning repositories <- ta difisi

// recursive and simple mode
package git_actions

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	common "github.com/hectorruiz-it/grabber/cmd"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

type report struct {
	Repository string
	Status     string
}

// pullCmd represents the pull command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Clones all repositories from a JSON grabber config file respecting paths.",
	Long:  `Clones all repositories from a JSON grabber config file respecting paths.`,
	// Args:  cobra.ExactArgs(1),
	Example: `  grabber migrate --all
  grabber migrate --profile github`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("pull called")

		// profile, err := cmd.Flags().GetString("profile")
		// common.CheckAndReturnError(err)
		// fmt.Println(profile)

		migrate()
	},
}

func init() {
	common.RootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().BoolP("all", "A", false, "migrates repositories from all profiles")
	migrateCmd.Flags().StringP("profile", "p", "", "profile repositories to migrate")
	// migrateCmd.MarkFlagsOneRequired("all", "profile")
	// migrateCmd.MarkFlagsMutuallyExclusive("all", "profile")
}

var (
	WarningLog *log.Logger
	InfoLog    *log.Logger
	ErrorLog   *log.Logger
)

func migrate() {
	dt := time.Now()
	logFilePath := common.GetHomeDirectory() + common.ROOT_DIR + "/migration-" + dt.Format("01-02-2006") + "_" + dt.Format("15:04:05")
	_, err := os.Create(logFilePath)
	common.CheckAndReturnError(err)
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	common.CheckAndReturnError(err)
	defer logFile.Close()

	InfoLog = log.New(logFile, "INFO: ", log.Ldate|log.Ltime)
	WarningLog = log.New(logFile, "WARNING: ", log.Ldate|log.Ltime)
	ErrorLog = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime)

	config := common.ReadGrabberConfig()
	currentProfiles := common.GetProfiles()
	// var configProfiles []string

	// var profiles common.Profiles
	checkProfiles := make(map[string]bool, len(config.Profiles))

	for _, configProfile := range config.Profiles {
		checkProfiles[configProfile.Profile] = false

		for _, currentProfile := range currentProfiles {
			if configProfile.Profile == currentProfile {
				checkProfiles[configProfile.Profile] = true
				break
			}
		}
		if !checkProfiles[configProfile.Profile] {
			fmt.Println("Profile `" + configProfile.Profile + "` not found. What do you want to do?")
			fmt.Println("1: Migrate repositories to another profile")
			fmt.Println("2: Create a new profile")
			fmt.Println("3: Omit profile")
			var option string
		loop:
			for {
				switch option {
				case "1":
					migrateRepositoriesToAnotherProfile()
					delete(checkProfiles, configProfile.Profile)
					break loop
				case "2":
					addProfile(configProfile.Profile, configProfile.Type)
					checkProfiles[configProfile.Profile] = true
					break loop
				case "3":
					delete(checkProfiles, configProfile.Profile)
					break loop
				default:
					fmt.Println("grabber: not a valid option")
				}
			}
		}
	}

	// fmt.Println(checkProfiles)
	var cloneReport report
	var cloned int
	var ommited int
	var failed int

	for _, configProfile := range config.Profiles {
		reportChannel := make(chan report, len(configProfile.Repositories))

		if checkProfiles[configProfile.Profile] {
			for _, repository := range configProfile.Repositories {
				go migrateClone(repository.Name, repository.Path, configProfile.Profile, configProfile.Type, reportChannel)
			}

			for range configProfile.Repositories {
				cloneReport = <-reportChannel
				switch cloneReport.Status {
				case "cloned":
					cloned += 1
				case "ommited":
					ommited += 1
				case "failed":
					failed += 1
				}
			}
			// fmt.Println(status)
			close(reportChannel)
		} else {
			fmt.Println("grabber: profile `" + configProfile.Profile + "` ommited. Profile doesn´t exist.")
		}
	}
	fmt.Println("Migration complete!", cloned, "cloned,", ommited, "ommited and", failed, "failed")

	// get config and existing profiles.
	// Profiles validation procedure:
	// If exists continue
	// If not (tui selector):
	// 1. Map profile to an existing one.
	// This will require to move all repositories to the existing profile if they are of the same type.
	// 2. Create a new profile based on config profile type:
	// Ask for name rebranding if wanted (prompt modification)?
	// Add the profile to the auth method file but not add it to json.
	// Modify Profile name if it has been a rebranding.
	// 3. Omit profile.
	// Start cloning repositories in parallel.
	// Is there a way to limit parallelism? -> Yes, limiting channel size
	// Handle err repository already exists with continue, not error.
	// Show a report of cloned and ommited repositories.
	// Use bubbletea packamanager tui.
}

func migrateClone(repository string, path string, profile string, profileType string, reportChannel chan report) {
	sshRegex := regexp.MustCompile(`^git@`)
	httpsRegex := regexp.MustCompile(`^https://`)
	service := "grabber"

	fmt.Println("Cloning:", repository)
	switch {
	case httpsRegex.MatchString(repository):
		password, err := keyring.Get(service, profile+"-profile")
		common.CheckAndReturnError(err)

		_, err = git.PlainClone(path, false, &git.CloneOptions{
			Auth: &http.BasicAuth{
				Username: "git",
				Password: password,
			},
			URL: repository,
			// Progress: os.Stdout,
		})

		switch err {
		case nil:
			InfoLog.Println("grabber: succesfully cloned repository `" + repository + "`.")
			reportChannel <- report{Repository: repository, Status: "cloned"}
		case git.ErrRepositoryAlreadyExists:
			WarningLog.Println("grabber: repository `" + repository + "` already exists.")
			reportChannel <- report{Repository: repository, Status: "ommited"}
		default:
			ErrorLog.Println("grabber:", err)
			reportChannel <- report{Repository: repository, Status: "failed"}
		}

	case sshRegex.MatchString(repository):
		if profileType == AUTH_METHODS[1] {
			sshProfiles := common.ReadSshProfilesFile()
			privateKey, err := sshProfiles.Section(profile).GetKey("private_key")
			common.CheckAndReturnError(err)

			password, err := keyring.Get(service, profile+"-profile")
			common.CheckAndReturnError(err)

			publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKey.Value(), password)
			common.CheckAndReturnError(err)

			_, err = git.PlainClone(path, false, &git.CloneOptions{
				Auth: publicKeys,
				// Progress: os.Stdout,
				URL: repository,
			})

			switch err {
			case nil:
				InfoLog.Println("grabber: succesfully cloned repository `" + repository + "`.")
				reportChannel <- report{Repository: repository, Status: "cloned"}
			case git.ErrRepositoryAlreadyExists:
				WarningLog.Println("grabber: repository `" + repository + "` already exists.")
				reportChannel <- report{Repository: repository, Status: "ommited"}
			default:
				ErrorLog.Println("grabber:", err)
				reportChannel <- report{Repository: repository, Status: "failed"}
			}
		} else {
			err := errors.New("grabber: profile `" + profile + "` is not an ssh profile.")
			common.CheckAndReturnError(err)
		}
	}
}

func migrateRepositoriesToAnotherProfile()          {}
func addProfile(profile string, profileType string) {}
