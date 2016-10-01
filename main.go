package main

import (
	"fmt"
	"github.com/xackery/discordeq/discord"
	"github.com/xackery/discordeq/listener"
	"github.com/xackery/eqemuconfig"
	"log"
	"os"
	"time"
)

func main() {
	startService()
}

func startService() {
	log.Println("Starting DiscordEQ v0.4")
	var option string
	//Load config
	config, err := eqemuconfig.GetConfig()
	if err != nil {
		log.Println("Error while loading eqemu_config.xml to start:", err.Error())
		log.Println("Press any key to exit.")
		fmt.Scan(&option)
		os.Exit(1)
	}
	if config.Discord.RefreshRate == 0 {
		config.Discord.RefreshRate = 10
	}

	if config.Discord.Username == "" {
		log.Println("I don't see a username set in your <discord><username> section of eqemuconfig.xml, please adjust.")
		log.Println("Press any key to exit.")
		fmt.Scan(&option)
		os.Exit(1)
	}

	if config.Discord.Password == "" {
		log.Println("I don't see a password set in your <discord><password> section of eqemuconfig.xml, please adjust.")
		log.Println("Press any key to exit.")
		fmt.Scan(&option)
		os.Exit(1)
	}

	if config.Discord.ServerID == "" {
		log.Println("I don't see a serverid set in your <discord><serverid> section of eqemuconfig.xml, please adjust.")
		log.Println("Press any key to exit.")
		fmt.Scan(&option)
		os.Exit(1)
	}

	if config.Discord.ChannelID == "" {
		log.Println("I don't see a channelid set in your <discord><channelid> section of eqemuconfig.xml, please adjust.")
		log.Println("Press any key to exit.")
		fmt.Scan(&option)
		os.Exit(1)
	}
	disco := discord.Discord{}
	err = disco.Connect(config.Discord.Username, config.Discord.Password)
	if err != nil {
		log.Println("Error connecting to discord:", err.Error())
		log.Println("Press any key to exit.")
		fmt.Scan(&option)
		os.Exit(1)
	}
	go listenToDiscord(config, &disco)
	go listenToOOC(config, &disco)
	select {}
}

func listenToDiscord(config *eqemuconfig.Config, disco *discord.Discord) (err error) {
	for {
		log.Println("[Discord] Connecting as", config.Discord.Username, "...")
		err = listener.ListenToDiscord(config, disco)
		if err != nil {
			log.Println("[Discord] Disconnected with error:", err.Error())
		}

		log.Println("[Discord] Reconnecting in 5 seconds...")
		time.Sleep(5 * time.Second)
		err = disco.Connect(config.Discord.Username, config.Discord.Password)
		if err != nil {
			log.Println("[Discord] Error connecting to discord:", err.Error())
		}
	}
}

func listenToOOC(config *eqemuconfig.Config, disco *discord.Discord) (err error) {
	for {
		log.Println("[OOC] Connecting to ", config.Database.Host, "...")
		listener.ListenToOOC(config, disco)
		log.Println("[OOC] Reconnecting in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}
