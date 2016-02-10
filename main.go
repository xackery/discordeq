package main

import (
	"fmt"
	"github.com/xackery/discordeq/listener"
	//"github.com/xackery/discordeq/menu"
	"github.com/xackery/discordeq/discord"
	"github.com/xackery/eqemuconfig"
	"os"
)

func main() {
	startService()
}

func startService() {
	var option string
	//Load config
	config, err := eqemuconfig.GetConfig()
	if err != nil {
		fmt.Println("Error while loading eqemu_config.xml to start:", err.Error())
		fmt.Println("Press any key to exit.")
		fmt.Scan(&option)
		os.Exit(1)
	}
	if config.Discord.RefreshRate == 0 {
		config.Discord.RefreshRate = 10
	}

	if config.Discord.Username == "" {
		fmt.Println("I don't see a username set in your <discord><username> section of eqemuconfig.xml, please adjust.")
		fmt.Println("Press any key to exit.")
		fmt.Scan(&option)
		os.Exit(1)
	}

	if config.Discord.Password == "" {
		fmt.Println("I don't see a password set in your <discord><password> section of eqemuconfig.xml, please adjust.")
		fmt.Println("Press any key to exit.")
		fmt.Scan(&option)
		os.Exit(1)
	}

	if config.Discord.ServerID == "" {
		fmt.Println("I don't see a serverid set in your <discord><serverid> section of eqemuconfig.xml, please adjust.")
		fmt.Println("Press any key to exit.")
		fmt.Scan(&option)
		os.Exit(1)
	}

	if config.Discord.ChannelID == "" {
		fmt.Println("I don't see a channelid set in your <discord><channelid> section of eqemuconfig.xml, please adjust.")
		fmt.Println("Press any key to exit.")
		fmt.Scan(&option)
		os.Exit(1)
	}
	disco := discord.Discord{}
	err = disco.Connect(config.Discord.Username, config.Discord.Password)
	if err != nil {
		fmt.Println("Error connecting to discord:", err.Error())
		fmt.Println("Press any key to exit.")
		fmt.Scan(&option)
		os.Exit(1)
	}
	go listener.ListenToDiscord(config, &disco)
	go listener.ListenToOOC(config, &disco)
	select {}
}
