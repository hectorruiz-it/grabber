/*
Copyright © 2024 HÉCTOR <EMAIL ADDRESS>
*/
package main

import (
	grabber "github.com/hectorruiz-it/grabber/cmd"
	_ "github.com/hectorruiz-it/grabber/cmd/git"
	_ "github.com/hectorruiz-it/grabber/cmd/profiles"
)

func main() {
	grabber.Execute()
}
