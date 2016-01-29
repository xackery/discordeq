package listener

import (
	_ "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/xackery/discordeq/discord"
	"github.com/xackery/eqemuconfig"
	"log"
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

func ListenToOOC(config *eqemuconfig.Config, disco *discord.Discord) {
	var err error
	channelID = config.Discord.ChannelID
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
	userMessages = nil
	//if lastID is set
	if lastId != 0 {
		//grab new ids if they match criteria
		err = db.Select(&userMessages, "SELECT `from`, `to`, message, type, timerecorded FROM qs_player_speech WHERE id > ? AND `type` = 5 LIMIT 50", lastId)
	}
	if err != nil {
		return
	}

	//Iterate any resluts
	for _, msg := range userMessages {
		log.Println(msg.From, msg.Message, msg.CreateDate, "vs", time.Now().UTC())
		disco.SendMessage(channelID, fmt.Sprintf("**%s OOC**: %s", msg.From, msg.Message))
	}

	if len(userMessages) > 0 { //if results match, grab the last element's Id as our lastID
		lastId = userMessages[len(userMessages)-1].Id
		return
	}

	//We've parsed the entire file.
	err = db.Get(&lastId, "SELECT id from qs_player_speech ORDER BY id DESC limit 1")
	if err != nil {
		return
	}

	return
}
