package listener

import (
	"github.com/bwmarrin/discordgo"
	"github.com/xackery/discordeq/discord"
	"github.com/xackery/eqemuconfig"
	"log"
	"strings"
	//"time"
	_ "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func ListenToDiscord(config *eqemuconfig.Config, disco *discord.Discord) {

	var session *discordgo.Session
	var guild *discordgo.Guild
	var err error

	session = disco.GetSession()
	guild, err = session.Guild(config.Discord.ServerID)
	if err != nil {
		log.Println("[discord] Error getting srver", config.Discord.ServerID, "with error:", err.Error())
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
			log.Println("[discord] Error getting member:", err.Error())
			return
		}

		roles, err := session.GuildRoles(config.Discord.ServerID)
		if err != nil {
			log.Println("[discord] Error getting roles:", err.Error())
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
		msg := m.Content
		//Maximum limit of 4k
		if len(msg) > 4000 {
			msg = msg[0:4000]
		}

		if len(msg) < 1 {
			return
		}

		//Insert entry
		db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", config.Database.Username, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Db))
		if err != nil {
			return
		}
		_, err = db.NamedExec("INSERT INTO qs_player_speech (`from`, `to`, `message`,`type`) VALUES (:ign, '!discord', :msg, 5)",
			map[string]interface{}{
				"ign": ign,
				"msg": msg,
			})
		if err != nil {
			log.Println("[discord] Invalid insert:", err.Error())
			return
		}
		log.Printf("[discord] %s: %s\n", ign, msg)
	}

	select {}
}
