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
	d.instance, err = discordgo.New(username, password)
	if err != nil {
		return
	}
	return
}

func (d *Discord) GetGuilds() (guilds []*discordgo.Guild, err error) {
	if d.instance == nil {
		err = d.Connect(d.lastUsername, d.lastPassword)
		if err != nil {
			return
		}
	}

	guilds, err = d.instance.UserGuilds()
	if err != nil {
		return
	}
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
