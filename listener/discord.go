package listener

import (
	"github.com/xackery/discordeq/discord"
	"github.com/xackery/eqemuconfig"
	//"log"
	"time"
)

func ListenToDiscord(config *eqemuconfig.Config, disco *discord.Discord) {
	//var err error
	for {

		time.Sleep(10 * time.Second)
		/*
			err = disco.Connect(config.Discord.Username, config.Discord.Password)
			if err != nil {
				log.Println("[discord] Failed to connect to discord:", err.Error())
				time.Sleep(10 * time.Second)
				continue
			}
		*/
	}
}
