package main

import (
	"fmt"
	"github.com/xackery/eqemuconfig"
	"os"
	"strings"
)

func main() {
	for {
		showMenu()
	}
}

func showMenu() {
	var err error
	var option string
	var isConfigLoaded bool
	fmt.Println("\n\n===DiscordEQ Plugin===")

	config, err := eqemuconfig.LoadConfig()
	status := "Good"
	if err != nil {
		isConfigLoaded = false
		status = fmt.Sprintf("Bad (%s)", err.Error())
	} else {
		isConfigLoaded = true
		status = fmt.Sprintf("Good (%s)", config.Longame)
	}
	fmt.Printf("1) Reload eqemu_config.xml (Status: %s)\n", status)

	if isConfigLoaded {
		fmt.Println("More options")
	}
	fmt.Println("Q) Quit")

	fmt.Scan(&option)
	fmt.Println("You chose option:", option)
	option = strings.ToLower(option)
	if option == "q" || option == "exit" || option == "quit" {
		fmt.Println("Quitting")
		os.Exit(0)
	}
}
