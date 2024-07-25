package grabber

import (
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

// func CheckSectionExists(file *ini.File, profile string, method string) error {
// 	if !file.HasSection(profile) {
// 		err := fmt.Errorf("grabber: profile already exists with " + method + " authentication")
// 		return err
// 	}
// 	return nil
// }

func ReadBasiceProfilesFile() *ini.File {
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

func ReadCheckExistsProfiles(profile string) (*ini.File, *ini.File, *ini.File) {
	basicProfiles := ReadBasiceProfilesFile()
	sshProfiles := ReadSshProfilesFile()
	tokenProfiles := ReadTokenProfilesFile()

	switch {
	case basicProfiles.HasSection(profile):
		err := fmt.Errorf("grabber: profile " + profile + " already exists with basic method")
		CheckAndReturnError(err)
	case sshProfiles.HasSection(profile):
		err := fmt.Errorf("grabber: profile " + profile + " already exists with token method")
		CheckAndReturnError(err)
	case tokenProfiles.HasSection(profile):
		err := fmt.Errorf("grabber: profile " + profile + " already exists with ssh method")
		CheckAndReturnError(err)
	}

	return basicProfiles, sshProfiles, tokenProfiles
}
