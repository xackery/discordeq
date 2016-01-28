package main

import (
	"fmt"
	"github.com/xackery/eqemuconfig"
	"os"
	"strings"
)

func main() {
	isFirstRun := true
	for {
		err := showMenu()
		if isFirstRun && err == nil { //Keep looping menu until no errors found
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
		status = fmt.Sprintf("Bad (%s)", err.Error())
	} else {
		isConfigLoaded = true
		status = fmt.Sprintf("Good (%s)", config.Longame)
	}
	fmt.Printf("1) Reload eqemu_config.xml (Status: %s)\n", status)
	if isConfigLoaded {
		isEverythingGood := true

		status = "Good"
		if config.Discord.Username == "" || config.Discord.Password == "" {
			status = "Not configured"
			isEverythingGood = false
		}
		if config.Discord.ServerID == 0 && config.Discord.ChannelID == 0 {
			status = "Bad"
			isEverythingGood = false
		}
		fmt.Printf("2) Configure Discord settings inside eqemu_config.xml (Status: %s)\n", status)

		status = "Good"
		err := checkChatLog()
		if err != nil {
			isEverythingGood = false
			status = "Bad: " + err.Error()
		}
		fmt.Printf("3) Enable Chat Logging (Status: %s)\n", status)
		fmt.Printf("4) Quest File For Discord Chat (Status: %s)\n", status)
		if isEverythingGood {
			fmt.Println("5) Start DiscordEQ")
		}
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

func checkChatLog() (err error) {
	return
}
