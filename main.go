package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/xackery/discordeq/applog"
	"github.com/xackery/discordeq/discord"
	"github.com/xackery/discordeq/listener"
	"github.com/xackery/eqemuconfig"
)

// Version will be set during build
var Version = "0.52.0"

func main() {
	applog.StartupInteractive()
	log.SetOutput(applog.DefaultOutput)
	err := run()
	if err != nil {
		var option string
		applog.Error.Println(err)
		fmt.Println("press a key then enter to exit.")
		fmt.Scan(&option)
		os.Exit(1)
	}
	os.Exit(0)
}

func run() (err error) {
	log.Println("Starting DiscordEQ", Version)

	//Load config
	config, err := eqemuconfig.GetConfig()
	if err != nil {
		err = errors.Wrap(err, "error while loading eqemu_config to start")
		return
	}
	if config.Discord.RefreshRate == 0 {
		config.Discord.RefreshRate = 10
	}

	if !isNewTelnetConfig(config) {
		err = fmt.Errorf("telnet must be enabled for this tool to work. Check your eqemu_config, and please adjust")
		return
	}

	if config.Discord.Username == "" {
		err = fmt.Errorf("username not set in your discord > username section of eqemu_config, please adjust")
		return
	}

	if config.Discord.Password == "" && config.Discord.ClientID == "" {
		err = fmt.Errorf("password not set in your discord > password section of eqemu_config, as well as no client id, please adjust")
		return
	}

	if config.Discord.ServerID == "" {
		err = fmt.Errorf("serverid not set in your discord > serverid section of eqemuconfig, please adjust")
		return
	}

	if config.Discord.ChannelID == "" {
		err = fmt.Errorf("channelid not set in your  discord > channelid section of eqemuconfig.xml, please adjust")
		return
	}
	disco := discord.Discord{}
	err = disco.Connect(config.Discord.Username, config.Discord.Password)
	if err != nil {
		err = errors.Wrap(err, "discord")
		return
	}
	go listenToDiscord(config, &disco)
	go listenToOOC(config, &disco)
	select {}
}

func isNewTelnetConfig(config *eqemuconfig.Config) bool {
	if strings.ToLower(config.World.Telnet.Enabled) == "true" {
		return true
	}
	if strings.ToLower(config.World.Tcp.Telnet) == "enabled" {
		return true
	}
	return false
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

func listenToOOC(config *eqemuconfig.Config, disco *discord.Discord) (err error) {
	for {
		listener.ListenToOOC(config, disco)
		applog.Info.Println("[OOC] Reconnecting in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}
