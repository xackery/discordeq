package discord

import (
	"github.com/bwmarrin/discordgo"
)

type Discord struct {
	instance     *discordgo.Session
	lastUsername string
	lastPassword string
}

func (d *Discord) Connect(username string, password string) (err error) {
	if d.instance, err = discordgo.New(username, password); err != nil {
		return
	}
	return
}

func (d *Discord) GetName() (name string) {
	if d.instance == nil {
		err := d.Connect(d.lastUsername, d.lastPassword)
		if err != nil {
			name = "Unknown"
			return
		}
	}

	user, err := d.instance.User("@me")
	if err != nil {
		name = "Unknown"
		return
	}
	return user.Username
}

func (d *Discord) GetGuilds() (guilds []*discordgo.UserGuild, err error) {
	if d.instance == nil {
		err = d.Connect(d.lastUsername, d.lastPassword)
		if err != nil {
			return
		}
	}

	guilds, err = d.instance.UserGuilds(0, "", "")
	if err != nil {
		return
	}
	return
}

func (d *Discord) GetSession() (session *discordgo.Session, err error) {
	if d.instance == nil {
		if err = d.Connect(d.lastUsername, d.lastPassword); err != nil {
			return
		}
	}
	session = d.instance
	return
}

func (d *Discord) GetChannels(guildID string) (channels []*discordgo.Channel, err error) {
	if d.instance == nil {
		err = d.Connect(d.lastUsername, d.lastPassword)
		if err != nil {
			return
		}
	}

	channels, err = d.instance.GuildChannels(guildID)
	if err != nil {
		return
	}
	return
}

func (d *Discord) SendMessage(channelID string, message string) (msgReturn *discordgo.Message, err error) {
	if d.instance == nil {
		err = d.Connect(d.lastUsername, d.lastPassword)
		if err != nil {
			return
		}
	}
	msgReturn, err = d.instance.ChannelMessageSend(channelID, message)
	return
}
