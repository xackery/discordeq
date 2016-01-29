package main

import (
	"fmt"
	//"github.com/bwmarrin/discordgo"
	"github.com/xackery/discordeq/discord"
	"github.com/xackery/eqemuconfig"
	"os"
	"strconv"
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
	} else if option == "2" {
		configureDiscord(&config)
		err = fmt.Errorf("Configuring discord")
	} else {
		fmt.Println("Invalid option")
		err = fmt.Errorf("No option chosen")
	}
	return
}

func checkChatLog() (err error) {
	return
}

func configureDiscord(config *eqemuconfig.Config) (err error) {
	disco := discord.Discord{}
	disco.Connect(config.Discord.Username, config.Discord.Password)

	guilds, err := disco.GetGuilds()
	if err != nil {
		fmt.Println("There is an error with your discord settings:", err.Error())
		return
	}
	if len(guilds) < 1 {
		fmt.Println("No guilds found, please join or create one prior to running this.")
		return
	}
	fmt.Println("Select which server you would like to use.")
	for i, guild := range guilds {
		fmt.Printf("%d) %s\n", i, guild.Name)
	}
	fmt.Println("C) Cancel")
	option := ""
	fmt.Scan(&option)

	fmt.Println("You chose option:", option)
	option = strings.ToLower(option)
	optionVal, _ := strconv.ParseInt(option, 10, 64)
	if option == "c" || option == "cancel" || option == "exit" || option == "quit" {
		fmt.Println("Cancelling")
		return
	} else if optionVal <= int64(len(guilds)) && optionVal >= 0 {
		serverID := guilds[optionVal].ID

		fmt.Printf("You chose server %s (%s)\n", guilds[optionVal].Name, serverID)

		channels, newErr := disco.GetChannels(serverID)
		if newErr != nil {
			fmt.Println("There is an error with your discord settings:", err.Error())
			err = newErr
			return
		}
		if len(channels) < 1 {
			fmt.Println("No channels found, please join or create one prior to running this.")
			return
		}
		fmt.Println("Select which text channel you would like to use.")
		for i, channel := range channels {
			if channel.Type != "text" {
				continue
			}
			fmt.Printf("%d) %s (%s)\n", i, channel.Name, channel.Type)
		}
		fmt.Println("C) Cancel")

		fmt.Scan(&option)

		fmt.Println("You chose option:", option)
		option = strings.ToLower(option)

		optionVal, _ := strconv.ParseInt(option, 10, 64)
		if option == "c" || option == "cancel" || option == "exit" || option == "quit" {
			fmt.Println("Cancelling")
			return
		} else if optionVal <= int64(len(channels)) && optionVal >= 0 {
			if channels[optionVal].Type != "text" {
				fmt.Printf("Channel is %s, please select a text channel.", channels[optionVal].Type, ", please select text.")
				return
			}
			fmt.Println("Inside your eqemu_config.xml, please update two lines as follows:")
			fmt.Printf("<serverid>%s</serverid>\n<channelid>%s</channelid>\n", serverID, channels[optionVal].ID)
		} else {
			fmt.Println("Invalid option")
			return
		}

	} else {
		fmt.Println("Invalid option")
		return
	}
	return
}
