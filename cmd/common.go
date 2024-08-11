package grabber

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"

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
	Profile      string       `json:"profile"`
	Type         string       `json:"type"`
	Repositories []Repository `json:"repositories"`
}

type Repository struct {
	Path string `json:"path"`
	Name string `json:"name"`
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
	checkMappingsFileExists(homeDir + MAPPINGS_FILE_PATH)
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

func checkMappingsFileExists(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if _, err := os.Create(filePath); os.IsPermission(err) {
			err = fmt.Errorf("grabber: unable to create file on .grabber directory. Permissions error.")
			CheckAndReturnError(err)
		}
		data, err := json.Marshal(Profiles{})
		CheckAndReturnError(err)
		os.WriteFile(filePath, data, 0700)
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

func ReadCheckExistsProfiles(profile string, validation bool) (bool, map[string]bool) {
	basicProfiles := ReadBasicProfilesFile()
	sshProfiles := ReadSshProfilesFile()
	tokenProfiles := ReadTokenProfilesFile()

	profileCheck := make(map[string]bool, 3)

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
		return true, profileCheck
	}

	switch {
	case basicProfiles.HasSection(profile):
		profileCheck["basic"] = true
		profileCheck["ssh"] = false
		profileCheck["token"] = false
		return true, profileCheck
	case sshProfiles.HasSection(profile):
		profileCheck["basic"] = false
		profileCheck["ssh"] = true
		profileCheck["token"] = false
		return true, profileCheck
	case tokenProfiles.HasSection(profile):
		profileCheck["basic"] = false
		profileCheck["ssh"] = false
		profileCheck["token"] = true
		return true, profileCheck
	}
	return false, profileCheck
}

func GetProfileByRepository(repository string) (string, string, error) {
	config := ReadGrabberConfig()
	sshRegex := regexp.MustCompile(`^git@`)
	httpsRegex := regexp.MustCompile(`^https://`)

	var authMethod string

	switch {
	case httpsRegex.MatchString(repository):
		authMethod = "token"
	case sshRegex.MatchString(repository):
		authMethod = "ssh"
	default:
		err := errors.New("grabber: not a valid origin")
		CheckAndReturnError(err)
	}

	var grabberProfile string
	for _, profile := range config.Profiles {
		if profile.Type == authMethod {
			for _, r := range profile.Repositories {
				if r.Name == repository {
					// fmt.Println("repository found")
					grabberProfile = profile.Profile
					return grabberProfile, authMethod, nil
				}
			}
		} else {
			continue
		}
	}

	err := errors.New("grabber: repository is outside grabber configuration.")
	return grabberProfile, authMethod, err
}
