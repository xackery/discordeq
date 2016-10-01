package listener

import (
	"fmt"
	"github.com/xackery/discordeq/discord"
	"github.com/xackery/eqemuconfig"
	"github.com/ziutek/telnet"
	"log"
	"strings"
	"time"
)

var lastId int
var channelID string

type UserMessage struct {
	Id         int       `db:"id"`
	From       string    `db:"from"`
	To         string    `db:"to"`
	Message    string    `db:"message"`
	Type       int       `db:"type"`
	CreateDate time.Time `db:"timerecorded"`
}

var userMessages []UserMessage
var config *eqemuconfig.Config

var t *telnet.Conn

func GetTelnet() (conn *telnet.Conn) {
	conn = t
	return
}

func ListenToOOC(eqconfig *eqemuconfig.Config, disco *discord.Discord) {
	var err error
	config = eqconfig
	channelID = config.Discord.ChannelID
	var t *telnet.Conn
	for {
		err = connectTelnet(config)
		if err != nil {
			log.Println("[OOC] Warning while getting telnet connection:", err.Error())
			time.Sleep((time.Duration(config.Discord.RefreshRate) + 10) * time.Second)
			continue
		}

		err = checkForMessages(t, disco)
		if err != nil {
			log.Println("[OOC] Warning while checking for messages:", err.Error())
			t.Close()
			time.Sleep((time.Duration(config.Discord.RefreshRate) + 10) * time.Second)
			continue
		}
		t.Close()
		time.Sleep(time.Duration(config.Discord.RefreshRate) * time.Second)
	}
	log.Println("Connected.")
}

func connectTelnet(config *eqemuconfig.Config) (err error) {
	if t, err = telnet.Dial("tcp", fmt.Sprintf("%s:%s", config.World.Tcp.Ip, config.World.Tcp.Port)); err != nil {
		return
	}
	t.SetReadDeadline(time.Now().Add(10 * time.Second))
	t.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err = t.SkipUntil("Username:"); err != nil {
		return
	}
	if err = Sendln(config.Discord.TelnetUsername); err != nil {
		return
	}

	if err = t.SkipUntil("Password:"); err != nil {
		return
	}
	if err = Sendln(config.Discord.TelnetPassword); err != nil {
		return
	}

	if err = Sendln("echo off"); err != nil {
		return
	}

	if err = Sendln("acceptmessages on"); err != nil {
		return
	}

	t.SetReadDeadline(time.Time{})
	return
}

func Sendln(s string) (err error) {
	buf := make([]byte, len(s)+1)
	copy(buf, s)
	buf[len(s)] = '\n'
	_, err = t.Write(buf)
	return
}

func checkForMessages(t *telnet.Conn, disco *discord.Discord) (err error) {
	data := []byte{}
	message := ""
	for {
		if data, err = t.ReadUntil("\n"); err != nil {
			err = fmt.Errorf("Error reading", err.Error())
			return
		}
		message = string(data)
		if len(message) < 3 { //ignore small messages
			continue
		}
		if !strings.Contains(message, "says ooc,") { //ignore non-ooc
			continue
		}
		if strings.Index(message, ">") < 1 { //ignore prompts
			continue
		}
		if message[0:1] == "*" { //ignore echo backs
			continue
		}

		message = message[strings.Index(message, ">")+2 : len(message)-1]
		sender := message[0:strings.Index(message, " says ooc,")]
		message = message[strings.Index(message, "says ooc, '")+11 : len(message)-2]

		if _, err = disco.SendMessage(channelID, fmt.Sprintf("**%s OOC**: %s", sender, message)); err != nil {
			log.Printf("[OOC] Error sending message (%s: %s) %s", sender, message, err.Error())
			continue
		}
		log.Printf("[OOC] %s: %s\n", sender, message)
	}
}
