package listener

import (
	_ "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/xackery/discordeq/discord"
	"github.com/xackery/eqemuconfig"
	"log"
	"strings"
	"time"
)

var lastId int
var db *sqlx.DB
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

func ListenToOOC(eqconfig *eqemuconfig.Config, disco *discord.Discord) {
	config = eqconfig
	var err error
	channelID = config.Discord.ChannelID
	log.Println("[ooc] Listening to OOC")
	for {
		db, err = connectDB(config)
		if err != nil {
			log.Println("[ooc] error while getting DB connection:", err.Error())
			time.Sleep((time.Duration(config.Discord.RefreshRate) + 10) * time.Second)

			continue
		}

		err = checkForMessages(db, disco)
		if err != nil {
			log.Println("[ooc] error while checking for messages:", err.Error())
			db.Close()
			time.Sleep((time.Duration(config.Discord.RefreshRate) + 10) * time.Second)
			continue
		}
		db.Close()
		time.Sleep(time.Duration(config.Discord.RefreshRate) * time.Second)
	}
}

func connectDB(config *eqemuconfig.Config) (db *sqlx.DB, err error) {
	//Connect to DB
	db, err = sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", config.Database.Username, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Db))
	if err != nil {
		return
	}
	return
}

func checkForMessages(db *sqlx.DB, disco *discord.Discord) (err error) {
	err = checkForOOCMessages(db, disco)
	if err != nil {
		return
	}
	checkForChannelMessages(db, disco)
	if err != nil {
		return
	}
	//We've parsed the entire file.
	err = db.Get(&lastId, "SELECT id from qs_player_speech ORDER BY id DESC limit 1")
	if err != nil {
		return
	}
	return
}

func checkForOOCMessages(db *sqlx.DB, disco *discord.Discord) (err error) {
	userMessages = nil
	//if lastID is set
	if lastId != 0 {
		//grab new ids if they match criteria
		err = db.Select(&userMessages, "SELECT `from`, `to`, message, type, timerecorded FROM qs_player_speech WHERE id > ? AND `type` = 5 AND `to` != '!discord' LIMIT 50", lastId)
	}
	if err != nil {
		return
	}

	//Iterate any results
	for _, msg := range userMessages {
		_, err := disco.SendMessage(channelID, fmt.Sprintf("**%s OOC**: %s", msg.From, msg.Message))
		if err != nil {
			log.Printf("[ooc] Error sending message (%s: %s) %s", msg.From, msg.Message, err.Error())
		} else {
			log.Printf("[ooc] %s: %s\n", msg.From, msg.Message)
		}
	}

	if len(userMessages) > 0 { //if results match, grab the last element's Id as our lastID
		lastId = userMessages[len(userMessages)-1].Id
		return
	}

	return
}

func checkForChannelMessages(db *sqlx.DB, disco *discord.Discord) (err error) {
	userMessages = nil
	//if lastID is set
	if lastId != 0 {
		//grab new ids if they match criteria
		err = db.Select(&userMessages, "SELECT `from`, `to`, message, type, timerecorded FROM qs_player_speech WHERE id > ? AND `type` = 6 AND `from` = '!eq2discord' LIMIT 50", lastId)
	}
	if err != nil {
		return
	}

	sendChannelID := ""
	//Iterate any results
	for _, msg := range userMessages {
		fmt.Println("Processing message:", msg.Message)
		for _, channel := range config.Discord.Channels {
			if strings.ToLower(msg.To) == strings.ToLower(channel.ChannelName) {
				sendChannelID = channel.ChannelID
				break
			}
		}

		if sendChannelID == "" {
			//Don't send the messgae if it's invalid
			log.Printf("[ooc] Error finding channel %s for id %d\n", msg.To, msg.Id)
			continue
		}
		fmt.Println("Sending message from sendChannelID: %s", sendChannelID)
		_, err := disco.SendMessage(sendChannelID, msg.Message)
		if err != nil {
			log.Printf("[ooc] Error sending message (%s: %s) %s\n", msg.From, msg.Message, err.Error())
		} else {
			log.Printf("[ooc] %s: %s\n", msg.From, msg.Message)
		}
	}

	return
}
