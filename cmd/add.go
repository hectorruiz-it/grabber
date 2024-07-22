/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

type Profile struct {
	privateKey string `toml:"privateKey"`
	publicKey  string `toml:"publicKey"`
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")
		addKey()

	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func checkPathAbsoluteExists(keyPath string) error {

	if !filepath.IsAbs(keyPath) {
		return fmt.Errorf("grabber: path is not absolute.")
	}

	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return fmt.Errorf("grabber: ssh key does not exist on the given path.")
	}
	return nil
}

func addKey() {
	var id string
	fmt.Print("> Enter profile ID: ")
	fmt.Scanln(&id)

	// var privateSshKeyPath string
	// fmt.Print("> Enter SSH private key absolute path: ")
	// fmt.Scanln(&privateSshKeyPath)

	// if err := checkPathAbsoluteExists(privateSshKeyPath); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// if err := checkPathAbsoluteExists(privateSshKeyPath + ".pub"); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	file, err := os.ReadFile("/home/hruiz/.grabber/grabber-profiles.toml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// doc := `
	// [github]
	// privateKey = '/home/hruiz/.ssh/github'
	// publicKey = '/home/hruiz/.ssh/github.pub'
	// `

	v := make(map[string]interface{})
	v[id] = map[string]string{}
	// data[id] = map[string]string{
	// 	"privateKey": privateSshKeyPath,
	// 	"publicKey":  privateSshKeyPath + ".pub",
	// }

	if err := toml.Unmarshal([]byte(file), &v); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	profileMap, ok := v[id].(map[string]string)

	if !ok {
		fmt.Println("Type assertion failed")
		os.Exit(1)
	}

	fmt.Println(profileMap["privateKey"])

	// fmt.Println("Public key: ", v["github"].p)
	// fmt.Println("Private key: ", v[id].publicKey)

	// data := map[string]interface{}{
	// 	id: map[string]string{
	// 		"privateKey": privateSshKeyPath,
	// 		"publicKey":  privateSshKeyPath + ".pub",
	// 	},
	// }

	// b, _ := toml.Marshal(data)
	// fmt.Println(string(b))

	// // Now you can use the profile variable
	// fmt.Printf("Profile created: %s\n", id)

}
