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
	session.StateEnabled = true

	session.OnMessageCreate = func(s *discordgo.Session, m *discordgo.Message) {
		if m.ChannelID != config.Discord.ChannelID {
			return
		}

		ign := ""
		member, err := session.State.Member(config.Discord.ServerID, m.Author.ID)
		if err != nil {
			log.Println("Error getting member:", err.Error())
			return
		}

		roles, err := session.GuildRoles(config.Discord.ServerID)
		if err != nil {
			log.Println("Error getting roles:", err.Error())
			return
		}
		for _, role := range member.Roles {
			if ign != "" {
				break
			}
			for _, gRole := range roles {
				if ign != "" {
					break
				}
				if strings.TrimSpace(gRole.ID) == strings.TrimSpace(role) {
					if strings.Contains(gRole.Name, "IGN:") {
						splitStr := strings.Split(gRole.Name, "IGN:")
						if len(splitStr) > 1 {
							ign = strings.TrimSpace(splitStr[1])
						}
					}
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
