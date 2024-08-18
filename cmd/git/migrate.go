// Migrate from JSON file and from DynamoDB
// Add preview before starting cloning repositories <- ta difisi

// recursive and simple mode
package git_actions

import (
	"fmt"
	"log"
	"os"
	"time"

	common "github.com/hectorruiz-it/grabber/cmd"
	migrate_tui "github.com/hectorruiz-it/grabber/tui/migrate/progress"
	"github.com/spf13/cobra"
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
		migrate()
	},
}

func init() {
	common.RootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().BoolP("all", "A", true, "migrates repositories from all profiles")
	migrateCmd.Flags().StringP("profile", "p", "", "profile repositories to migrate (WIP)")
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
			fmt.Println("  1: Migrate repositories to another profile")
			fmt.Println("  2: Create a new profile")
			fmt.Println("  3: Omit profile")
			var option string
			fmt.Printf("Option: ")
			fmt.Scanln(&option)
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

	var migrationRepositories []migrate_tui.Repository

	for _, configProfile := range config.Profiles {
		if checkProfiles[configProfile.Profile] {
			for _, repository := range configProfile.Repositories {
				migrationRepositories = append(migrationRepositories, migrate_tui.Repository{
					Path:        repository.Path,
					Profile:     configProfile.Profile,
					ProfileType: configProfile.Type,
					Repository:  repository.Name,
				})
			}
		}
	}

	migrate_tui.ProgressTui(migrationRepositories)
}

func migrateRepositoriesToAnotherProfile()          {}
func addProfile(profile string, profileType string) {}
