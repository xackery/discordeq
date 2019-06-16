package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

//Discord wraps discord objects
type Discord struct {
	instance     *discordgo.Session
	lastUsername string
	lastPassword string
}

//Connect to a discord endpoint
func (d *Discord) Connect(username string, password string) (err error) {
	if len(password) > 0 {
		if d.instance, err = discordgo.New(username, password); err != nil {
			err = errors.Wrap(err, "failed to connect with user/pass")
			return
		}
	}
	d.instance, err = discordgo.New("Bot " + username)
	if err != nil {
		err = errors.Wrap(err, "failed to connect with token")
		return
	}
	return
}

// GetName returns user's name
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

// GetGuilds gets a list of guilds
func (d *Discord) GetGuilds() (guilds []*discordgo.UserGuild, err error) {
	if d.instance == nil {
		err = d.Connect(d.lastUsername, d.lastPassword)
		if err != nil {
			err = errors.Wrap(err, "getguilds: failed to connect")
			return
		}
	}

	guilds, err = d.instance.UserGuilds(0, "", "")
	if err != nil {
		err = errors.Wrap(err, "failed to user guilds")
		return
	}
	return
}

// GetSession returns a session object
func (d *Discord) GetSession() (session *discordgo.Session, err error) {
	if d.instance == nil {
		if err = d.Connect(d.lastUsername, d.lastPassword); err != nil {
			err = errors.Wrap(err, "GetSession.Connect")
			return
		}
	}
	session = d.instance
	return
}

// GetChannels returns a channel
func (d *Discord) GetChannels(guildID string) (channels []*discordgo.Channel, err error) {
	if d.instance == nil {
		err = d.Connect(d.lastUsername, d.lastPassword)
		if err != nil {
			err = errors.Wrap(err, "GetChannels.Connect")
			return
		}
	}

	channels, err = d.instance.GuildChannels(guildID)
	if err != nil {
		err = errors.Wrap(err, "GetChannels.guildchannels")
		return
	}
	return
}

// SendMessage sends a message through discord
func (d *Discord) SendMessage(channelID string, message string) (msgReturn *discordgo.Message, err error) {
	if d.instance == nil {
		err = d.Connect(d.lastUsername, d.lastPassword)
		if err != nil {
			err = errors.Wrap(err, "Sendmessage.Connect")
			return
		}
	}
	msgReturn, err = d.instance.ChannelMessageSend(channelID, message)
	return
}
