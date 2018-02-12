package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/xackery/discordeq/applog"
	"github.com/xackery/discordeq/discord"
	"github.com/xackery/discordeq/listener"
	"github.com/xackery/eqemuconfig"
)

func main() {
	applog.StartupInteractive()
	log.SetOutput(applog.DefaultOutput)
	startService()
}

func startService() {
	log.Println("Starting DiscordEQ (NATS edition)")
	var option string
	//Load config
	config, err := eqemuconfig.GetConfig()
	if err != nil {
		applog.Error.Println("Error while loading eqemu_config.xml to start:", err.Error())
		fmt.Println("press a key then enter to exit.")

		fmt.Scan(&option)
		os.Exit(1)
	}
	if config.Discord.RefreshRate == 0 {
		config.Discord.RefreshRate = 10
	}

	if config.Discord.Username == "" {
		applog.Error.Println("I don't see a username set in your <discord><username> section of eqemu_config.xml, please adjust.")
		fmt.Println("press a key then enter to exit.")
		fmt.Scan(&option)
		os.Exit(1)
	}

	if config.Discord.Password == "" && config.Discord.ClientID == "" {
		applog.Error.Println("I don't see a password set in your discord > password section of eqemu_config, as well as no client id, please adjust.")
		fmt.Println("press a key then enter to exit.")
		fmt.Scan(&option)
		os.Exit(1)
	}

	if config.Discord.Password == "" && config.Discord.ClientID == "" {
		applog.Error.Println("I don't see a password set in your discord > password section of eqemu_config, as well as no client id, please adjust.")
		fmt.Println("press a key then enter to exit.")
		fmt.Scan(&option)
		os.Exit(1)
	}

	if config.Discord.ServerID == "" {
		applog.Error.Println("I don't see a serverid set in your <discord><serverid> section of eqemuconfig.xml, please adjust.")
		fmt.Println("press a key then enter to exit.")
		fmt.Scan(&option)
		os.Exit(1)
	}

	if config.Discord.ChannelID == "" {
		applog.Error.Println("I don't see a channelid set in your <discord><channelid> section of eqemuconfig.xml, please adjust.")
		fmt.Println("press a key then enter to exit.")
		fmt.Scan(&option)
		os.Exit(1)
	}
	disco := discord.Discord{}
	err = disco.Connect(config.Discord.Username, config.Discord.Password)
	if err != nil {
		applog.Error.Println("Error connecting to discord:", err.Error())
		fmt.Println("press a key then enter to exit.")
		fmt.Scan(&option)
		os.Exit(1)
	}
	go listenToDiscord(config, &disco)
	go listenToNATS(config, &disco)
	select {}
}

func listenToDiscord(config *eqemuconfig.Config, disco *discord.Discord) (err error) {
	for {
		if len(config.Discord.Password) > 0 { //don't show username if it's token based
			applog.Info.Println("[Discord] Connecting as", config.Discord.Username, "...")
		} else {
			applog.Info.Println("[Discord] Connecting...")
		}
		if err = listener.ListenToDiscord(config, disco); err != nil {
			if strings.Contains(err.Error(), "Unauthorized") {
				applog.Info.Printf("Your bot is not authorized to access this server.\nClick this link and give the bot access: https://discordapp.com/oauth2/authorize?&client_id=%s&scope=bot&permissions=268446736", config.Discord.ClientID)
				return
			}
			applog.Error.Println("[Discord] Disconnected with error:", err.Error())
		}

		applog.Info.Println("[Discord] Reconnecting in 5 seconds...")
		time.Sleep(5 * time.Second)
		err = disco.Connect(config.Discord.Username, config.Discord.Password)
		if err != nil {
			applog.Error.Println("[Discord] Error connecting to discord:", err.Error())
		}
	}
}

func listenToNATS(config *eqemuconfig.Config, disco *discord.Discord) (err error) {
	for {
		listener.ListenToNATS(config, disco)
		applog.Info.Println("[NATS] Reconnecting in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}
