package main

import (
	"fmt"
	"github.com/xackery/eqemuconfig"
	"os"
	"strings"
)

func main() {
	for {
		err := showMenu()
		if err == nil { //Keep looping menu until no errors found
			break
		}
	}
}

func showMenu() (err error) {
	var option string
	var isConfigLoaded bool
	fmt.Println("\n===DiscordEQ Plugin===")

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
		status = "Good"
		if config.Discord.Username == "" || config.Discord.Password == "" {
			status = "Bad"
		}
		fmt.Printf("2) Discord settings inside eqemu_config.xml (Status: %s)\n", status)
	}
	fmt.Println("Q) Quit")

	fmt.Scan(&option)
	fmt.Println("You chose option:", option)
	option = strings.ToLower(option)
	if option == "q" || option == "exit" || option == "quit" {
		fmt.Println("Quitting")
		os.Exit(0)
	} else {
		fmt.Println("Invalid option")
		err = fmt.Errorf("No option chosen")
	}
	return
}
