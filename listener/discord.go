package listener

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/xackery/discordeq/discord"
	"github.com/xackery/eqemuconfig"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"log"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

func ListenToDiscord(config *eqemuconfig.Config, disco *discord.Discord) (err error) {
	var session *discordgo.Session
	var guild *discordgo.Guild
	//log.Println("Listen to discord..")
	if session, err = disco.GetSession(); err != nil {
		log.Printf("[Discord] Failed to get instance %s: %s (Make sure bot is part of server)", config.Discord.ServerID, err.Error())
		return
	}

	if guild, err = session.Guild(config.Discord.ServerID); err != nil {
		log.Printf("[Discord] Failed to get server %s: %s (Make sure bot is part of server)", config.Discord.ServerID, err.Error())
		return
	}

	isNotAvail := true
	if guild.Unavailable == &isNotAvail {
		log.Printf("[Discord] Failed to get server %s: Server unavailable (Make sure bot is part of server, and has permission)", config.Discord.ServerID, err.Error())
		return
	}

	session.StateEnabled = true
	session.AddHandler(messageCreate)
	if err = session.Open(); err != nil {
		log.Printf("[Discord] Session closed: %s", err.Error())
		return
	}
	select {}
	return
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID != config.Discord.ChannelID {
		return
	}

	ign := ""
	member, err := s.State.Member(config.Discord.ServerID, m.Author.ID)
	if err != nil {
		log.Printf("[Discord] Failed to get member: %s (Make sure you have set the bot permissions to see members)", err.Error())
		return
	}

	roles, err := s.GuildRoles(config.Discord.ServerID)
	if err != nil {
		log.Printf("[Discord] Failed to get roles: %s (Make sure you have set the bot permissions to see roles)", err.Error())
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
	msg := m.ContentWithMentionsReplaced()
	//Maximum limit of 4k
	if len(msg) > 4000 {
		msg = msg[0:4000]
	}

	if len(msg) < 1 {
		return
	}

	ign = sanitize(ign)
	msg = sanitize(msg)

	//Send message.
	if err = Sendln(fmt.Sprintf("emote world 260 %s says from discord, '%s'", ign, msg)); err != nil {
		log.Printf("[Discord] Error sending message to telnet (%s:%s): %s\n", ign, msg, err.Error())
		return
	}

	log.Printf("[Discord] %s: %s\n", ign, msg)
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func sanitize(data string) (sData string) {
	re := regexp.MustCompile("/^[\x00-\x7F]+")
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	sData, _, _ = transform.String(t, data)
	sData = strings.Replace(sData, "�", "", -1) //emotes
	if !utf8.ValidString(sData) {
		v := make([]rune, 0, len(sData))
		for i, r := range sData {
			if r == utf8.RuneError {
				_, size := utf8.DecodeRuneInString(sData[i:])
				if size == 1 {
					continue
				}
			}
			v = append(v, r)
		}
		sData = string(v)
	}

	//finally, regex strip
	sData = re.ReplaceAllString(sData, "")
	return
}
