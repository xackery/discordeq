package listener

import (
	"github.com/bwmarrin/discordgo"
	"github.com/xackery/discordeq/discord"
	"github.com/xackery/eqemuconfig"
	"log"
	"strings"
	//"time"
)

func ListenToDiscord(config *eqemuconfig.Config, disco *discord.Discord) {

	var session *discordgo.Session
	var guild *discordgo.Guild
	var err error

	session = disco.GetSession()
	guild, err = session.Guild(config.Discord.ServerID)
	log.Println(guild.Name)
	log.Println(guild.Members)
	if err != nil {
		log.Println("Error getting srver", config.Discord.ServerID, "with error:", err.Error())
		return
	}
	//var err error
	for {

		session.OnMessageCreate = func(s *discordgo.Session, m *discordgo.Message) {
			log.Println(m.ChannelID, config.Discord.ChannelID)
			if m.ChannelID != config.Discord.ChannelID {
				return
			}

			ign := ""
			for _, member := range guild.Members {
				log.Println(member.User.Username)
				for _, role := range member.Roles {
					log.Println(role)
					if strings.Contains(strings.ToLower(role), "IGN:") {
						split := strings.SplitAfter(role, "IGN:")
						log.Println(split)
					}
				}
			}
			if ign == "" {
				return
			}
			log.Println(ign, m.Content)

		}

		select {}
	}
}
