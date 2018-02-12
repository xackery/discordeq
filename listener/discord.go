package listener

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/xackery/eqemuconfig"
	"github.com/xackery/discordeq/discord"
)

var disco *discord.Discord

func ListenToDiscord(config *eqemuconfig.Config, disc *discord.Discord) (err error) {
	var session *discordgo.Session
	disco = disc
	if session, err = disco.GetSession(); err != nil {
		log.Printf("[Discord] Failed to get instance %s: %s (Make sure bot is part of server)", config.Discord.ServerID, err.Error())
		return
	}

	session.StateEnabled = true
	session.AddHandler(onMessageEvent)
	log.Printf("[Discord] Connected\n")
	if err = session.Open(); err != nil {
		log.Printf("[Discord] Session closed: %s", err.Error())
		return
	}
	select {}
}

func onMessageEvent(s *discordgo.Session, m *discordgo.MessageCreate) {

	//fmt.Printf("Debug: %s, %s\n", m.ChannelID, m.Message.Content)
	//Look for messages to be relayed to OOC in game.
	if m.ChannelID == config.Discord.ChannelID &&
		len(m.Message.Content) > 0 &&
		m.Message.Content[0:1] != "!" {
		messageCreate(s, m)
		return
	}

	//Look for any commands.
	if m.ChannelID == config.Discord.CommandChannelID &&
		len(m.Message.Content) > 0 &&
		m.Message.Content[0:1] == "!" {
		commandParse(s, m)
	}

}

//Commands are parsed on specific channels
func commandParse(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(m.Message.Content) < 1 {
		if _, err := disco.SendMessage(m.ChannelID, fmt.Sprintf("%s: !help for valid commands", m.Author.Username)); err != nil {
			fmt.Printf("[Discord] Failed to send discord help command: %s\n", err.Error())
			return
		}
		return
	}

	allowedCommands := []string{"unlock", "who", "lock", "setidentity", "worldshutdown"}
	//figure out command, remove the ! bang
	commandSplit := strings.Split(m.Message.Content[1:], " ")
	parameters := commandSplit[1:]
	command := commandSplit[0]
	command = strings.ToLower(command)
	if command == "help" {
		if _, err := disco.SendMessage(m.ChannelID, fmt.Sprintf("%s: !help: Available commands: %s", m.Author.Username, strings.Join(allowedCommands[:], ", "))); err != nil {
			fmt.Printf("[Discord] Failed to send discord help command: %s\n", err.Error())
			return
		}
	}
	for _, cmd := range allowedCommands {
		if strings.Index(command, cmd) != 0 {
			continue
		}

		if err := SendCommand(m.Author.Username, command, parameters); err != nil {
			if _, derr := disco.SendMessage(m.ChannelID, fmt.Sprintf("%s: %s: %s", m.Author.Username, command, err.Error())); derr != nil {
				fmt.Printf("[Discord] Failed to send discord command message: %s\n", err.Error())
				return
			}
		}
		return
	}

	if _, err := disco.SendMessage(m.ChannelID, fmt.Sprintf("%s: %s is an invalid command. Use !help to learn commands.", m.Author.Username, command)); err != nil {
		fmt.Printf("[Discord] Failed to send discord command message: %s\n", err.Error())
		return
	}

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	ign := ""

	member, err := s.GuildMember(config.Discord.ServerID, m.Author.ID)
	if err != nil {
		log.Printf("[Discord] Failed to get member: %s (Make sure you have set the bot permissions to see members) ServerID: %s, AuthorID: %s", err.Error(), config.Discord.ServerID, m.Author.ID)
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
	sendNATSMessage(ign, msg)

	//if err = Sendln(fmt.Sprintf("emote world 260 %s says from discord, '%s'", ign, msg)); err != nil {
	//	log.Printf("[Discord] Error sending message to telnet (%s:%s): %s\n", ign, msg, err.Error())
	//	return
	//}

	log.Printf("[Discord] %s: %s\n", ign, msg)
}

func sanitize(data string) (sData string) {
	sData = data
	sData = strings.Replace(sData, `%`, "&PCT;", -1)
	re := regexp.MustCompile("[^\x00-\x7F]+")
	sData = re.ReplaceAllString(sData, "")
	return
}

func alphanumeric(data string) (sData string) {
	sData = data
	re := regexp.MustCompile("[^a-zA-Z0-9_]+")
	sData = re.ReplaceAllString(sData, "")
	return
}
