package listener

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/xackery/discordeq/discord"
	"github.com/xackery/eqemuconfig"
)

type NameConfig struct {
	Discord string
	Name string
}

var disco *discord.Discord

// ListenToDiscord listens for discord messages
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

	//Look for messages to be relayed to OOC in game.
	if m.ChannelID == config.Discord.ChannelID &&
		len(m.Message.Content) > 0 &&
		m.Message.Content[0:1] != "!" {
		messageCreate(s, m)
		return
	}

	//Look for any commands.
	if len(m.Message.Content) > 0 &&
		m.Message.Content[0:1] == "!" {
		commandParse(s, m)
	}

}

func commandParse(s *discordgo.Session, m *discordgo.MessageCreate) {
	//This feature is currently not supported.
	return

	//Verify user is allowed to send commands
	isAllowed := false
	for _, admin := range config.Discord.Admins {
		if m.Author.ID == admin.Id {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		if _, err := disco.SendMessage(m.ChannelID, fmt.Sprintf("Sorry %s, access denied.", m.Author.Username)); err != nil {
			fmt.Printf("[Discord] Failed to send discord message: %s\n", err.Error())
			return
		}
		return
	}
	//figure out command, remove the ! bang
	command := strings.ToLower(m.Message.Content[1:])

	switch command {
	case "help":
		if _, err := disco.SendMessage(m.ChannelID, fmt.Sprintf("%s: !help: Available commands:", m.Author.Username)); err != nil {
			fmt.Printf("[Discord] Failed to send discord help command: %s\n", err.Error())
			return
		}
	case "who":

	default:
		if _, err := disco.SendMessage(m.ChannelID, fmt.Sprintf("%s: Invalid command. Use !help to learn commands.", m.Author.Username)); err != nil {
			fmt.Printf("[Discord] Failed to send discord command message: %s\n", err.Error())
			return
		}
	}

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	members, err := s.GuildMembers(config.Discord.ServerID, "0", 1000)
	if err != nil {
		log.Printf("[Discord] Failed to get member: %s (Make sure you have set the bot permissions to see members)", err.Error())
		return
	}

	name_config_file, err := ioutil.ReadFile("discord_names.json")
	if err != nil {
		log.Printf("Failed to open name config: %s", err.Error());
		return
	}
	
	var names []NameConfig

	nameJson := json.Unmarshal(name_config_file, &names)
	if nameJson != nil {
		log.Printf("Failed to unmarshal name config: %s", nameJson.Error());
		return
	}
	discord_id := ""
	ingame_nickname := ""
	for _, member := range members {
		if ingame_nickname != "" {
			break
		}

		for name_index := range names {
			if ingame_nickname != "" {
				break
			}

			discord_id = names[name_index].Discord
			if member.User.ID == discord_id {
				ingame_nickname = strings.TrimSpace(names[name_index].Name)
			}
		}
	}

	if ingame_nickname == "" {
		return
	}
	message := m.ContentWithMentionsReplaced()
	//Maximum limit of 4k
	if len(message) > 4000 {
		message = message[0:4000]
	}

	if len(message) < 1 {
		return
	}

	ingame_nickname = sanitize(ingame_nickname)
	message = sanitize(message)

	//Send message.
	if err = Sendln(fmt.Sprintf("emote world 260 %s says from Discord, '%s'", ingame_nickname, message)); err != nil {
		log.Printf("[Discord] Error sending message to telnet (%s:%s): %s\n", ingame_nickname, message, err.Error())
		return
	}

	log.Printf("[Discord] %s: %s\n", ingame_nickname, message)
}

func sanitize(data string) (sData string) {
	sData = data
	sData = strings.Replace(sData, `%`, "&PCT;", -1)
	for emoji, ascii := range emojis {
		sData = strings.Replace(sData, emoji, ascii, -1)
	}
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
