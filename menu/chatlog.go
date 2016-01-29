package menu

import (
	"fmt"
	"github.com/xackery/eqemuconfig"

	_ "database/sql"
	"github.com/jmoiron/sqlx"
)

func connectDB(config *eqemuconfig.Config) (db *sqlx.DB, err error) {
	//Connect to DB
	db, err = sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", config.Database.Username, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Db))
	if err != nil {
		//fmt.Println("Error connecting to DB:", err.Error())
		return
	}
	return
}

func menuChatLog(config *eqemuconfig.Config) (err error) {
	fmt.Println("ChatLog")
	return
}

func isChatLoggingEnabled(config *eqemuconfig.Config) bool {
	db, err := connectDB(config)
	if err != nil {
		return false
	}
	//QueryServ, PlayerLogChat
	//REPLACE INTO `rule_values` (`ruleset_id`, `rule_name`, `rule_value`, `notes`) VALUES (0, 'QueryServ:PlayerLogChat', 'false', '');
	count := 0
	err = db.Get(&count, "SELECT * FROM rule_values WHERE rule_name = 'QueryServ:PlayerLogChat' AND rule_value = 'true' && ruleset_id = 0")
	if err != nil {
		//fmt.Println("Error initial", err.Error())
		return false
	}
	return count > 0
}

func enableChatLogging(config *eqemuconfig.Config) (err error) {
	db, err := connectDB(config)
	if err != nil {
		return
	}
	db.Exec("REPLACE INTO `rule_values` (`ruleset_id`, `rule_name`, `rule_value`) VALUES (0, 'QueryServ:PlayerLogChat', 'true')")

	return
}
