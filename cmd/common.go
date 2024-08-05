package grabber

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

const (
	ROOT_DIR                 = "/.grabber"
	PROFILES_BASIC_FILE_PATH = "/.grabber/grabber-basic-profiles"
	PROFILES_SSH_FILE_PATH   = "/.grabber/grabber-ssh-profiles"
	PROFILES_TOKEN_FILE_PATH = "/.grabber/grabber-token-profiles"
	MAPPINGS_FILE_PATH       = "/.grabber/.grabber-config.json"
)

// User struct which contains a name
// a type and a list of social links
type Profiles struct {
	Profiles []Profile `json:"profiles"`
}

type Profile struct {
	Profile      string   `json:"profile"`
	Repositories []string `json:"repositories"`
}

func CheckAndReturnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func GetHomeDirectory() string {
	homeDirectory, err := os.UserHomeDir()
	CheckAndReturnError(err)
	return homeDirectory
}

func setup() {
	homeDir := GetHomeDirectory()

	if err := os.Mkdir(homeDir+ROOT_DIR, 0700); !os.IsExist(err) {
		CheckAndReturnError(err)
	}
	checkFileExists(homeDir + PROFILES_BASIC_FILE_PATH)
	checkFileExists(homeDir + PROFILES_SSH_FILE_PATH)
	checkFileExists(homeDir + PROFILES_TOKEN_FILE_PATH)
	checkFileExists(homeDir + MAPPINGS_FILE_PATH)
}

func checkFileExists(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if _, err := os.Create(filePath); os.IsPermission(err) {
			err = fmt.Errorf("grabber: unable to create file on .grabber directory. Permissions error.")
			CheckAndReturnError(err)
		}
	} else {
		CheckAndReturnError(err)
	}
}

func ReadGrabberConfig() Profiles {
	homeDir := GetHomeDirectory()
	jsonFile, err := os.ReadFile(homeDir + MAPPINGS_FILE_PATH)
	CheckAndReturnError(err)

	var config Profiles

	err = json.Unmarshal(jsonFile, &config)
	CheckAndReturnError(err)

	return config
}

func ReadBasicProfilesFile() *ini.File {
	homeDir := GetHomeDirectory()
	basicProfiles, err := ini.Load(homeDir + PROFILES_BASIC_FILE_PATH)
	CheckAndReturnError(err)
	return basicProfiles
}

func ReadSshProfilesFile() *ini.File {
	homeDir := GetHomeDirectory()
	sshProfiles, err := ini.Load(homeDir + PROFILES_SSH_FILE_PATH)
	CheckAndReturnError(err)
	return sshProfiles
}

func ReadTokenProfilesFile() *ini.File {
	homeDir := GetHomeDirectory()
	tokenProfiles, err := ini.Load(homeDir + PROFILES_TOKEN_FILE_PATH)
	CheckAndReturnError(err)
	return tokenProfiles
}

func ReadCheckExistsProfiles(profile string, validation bool) bool {
	basicProfiles := ReadBasicProfilesFile()
	sshProfiles := ReadSshProfilesFile()
	tokenProfiles := ReadTokenProfilesFile()

	var err error
	switch {
	case validation && basicProfiles.HasSection(profile):
		err = fmt.Errorf("grabber: profile " + profile + " already exists with basic method")
		CheckAndReturnError(err)

	case validation && sshProfiles.HasSection(profile):
		err = fmt.Errorf("grabber: profile " + profile + " already exists with ssh method")
		CheckAndReturnError(err)

	case validation && tokenProfiles.HasSection(profile):
		err = fmt.Errorf("grabber: profile " + profile + " already exists with token method")
		CheckAndReturnError(err)
	case validation:
		return true
	}

	switch {
	case basicProfiles.HasSection(profile):
		return true
	case sshProfiles.HasSection(profile):
		return true
	case tokenProfiles.HasSection(profile):
		return true
	}
	return false
}
