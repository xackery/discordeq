package menu

import (
	"fmt"
	"github.com/xackery/discordeq/discord"
	"github.com/xackery/eqemuconfig"
	"strconv"
	"strings"
)

func menuDiscord(config *eqemuconfig.Config) (err error) {
	disco := discord.Discord{}
	disco.Connect(config.Discord.Username, config.Discord.Password)
	fmt.Println("You are logged in as:", disco.GetName())
	guilds, err := disco.GetGuilds()
	if err != nil {
		fmt.Println("There is an error with your discord settings:", err.Error())
		err = nil
		return
	}
	if len(guilds) < 1 {
		fmt.Println("No guilds found, please join or create one prior to running this.")
		err = nil
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
		err = nil
		return
	} else if optionVal <= int64(len(guilds)) && optionVal >= 0 {
		serverID := guilds[optionVal].ID

		fmt.Printf("You chose server %s (%s)\n", guilds[optionVal].Name, serverID)

		channels, newErr := disco.GetChannels(serverID)
		if newErr != nil {
			fmt.Println("There is an error with your discord settings:", err.Error())
			err = nil
			return
		}
		if len(channels) < 1 {
			err = fmt.Errorf("No channels found, please join or create one prior to running this.")
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
				fmt.Errorf("Channel is %s, please select a text channel.", channels[optionVal].Type)
				return
			}
			fmt.Println("Inside your eqemu_config.xml, please update two lines as follows:")
			fmt.Printf("<serverid>%s</serverid>\n<channelid>%s</channelid>\n", serverID, channels[optionVal].ID)
		} else {
			fmt.Errorf("Invalid option")
			return
		}
	} else {
		fmt.Errorf("Invalid Option")
		return
	}
	return
}
