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
	/*
		isFirstRun := true
		for {
			err := menu.ShowMenu()
			if isFirstRun && err == nil { //Keep looping menu until no errors found
				break
			}
		}*/

	startService()
}

func startService() {
	//Load config
	config, err := eqemuconfig.LoadConfig()
	if err != nil {
		fmt.Println("Error while loading config to start:", err.Error())
		os.Exit(1)
	}
	if config.Discord.RefreshRate == 0 {
		config.Discord.RefreshRate = 60
	}
	disco := discord.Discord{}
	err = disco.Connect(config.Discord.Username, config.Discord.Password)
	go listener.ListenToDiscord(&config, &disco)
	go listener.ListenToOOC(&config, &disco)
	fmt.Println("Listening")
	select {}
}
