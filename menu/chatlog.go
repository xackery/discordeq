package menu

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/xackery/eqemuconfig"
	"strings"
)

var theCommand = "REPLACE INTO `rule_values` (`ruleset_id`, `rule_name`, `rule_value`) VALUES (0, 'QueryServ:PlayerLogChat', 'true')"

func menuChatLog(config *eqemuconfig.Config) (err error) {
	fmt.Println("To enable chat logging, I will execute this command:")
	fmt.Println(theCommand)
	fmt.Println("Is this OK?")
	fmt.Println("Y) Yes")
	fmt.Println("N) No")
	option := ""
	fmt.Scan(&option)

	fmt.Println("You chose option:", option)
	option = strings.ToLower(option)
	if option == "y" || option == "yes" {
		err = enableChatLogging(config)
		if err != nil {
			fmt.Println("Error enabling chat logging:", err.Error())
			err = nil
			return
		}
		fmt.Println("Successfully enabled chat logging.")
	} else {
		fmt.Println("Cancelling")
	}
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
	err = db.Get(&count, "SELECT count(ruleset_id) FROM rule_values WHERE rule_name = 'QueryServ:PlayerLogChat' AND rule_value = 'true' && ruleset_id = 0 LIMIT 1")
	if err != nil {
		fmt.Println("Error initial", err.Error())
		return false
	}
	return count > 0
}

func connectDB(config *eqemuconfig.Config) (db *sqlx.DB, err error) {
	//Connect to DB
	db, err = sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", config.Database.Username, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Db))
	if err != nil {
		//fmt.Println("Error connecting to DB:", err.Error())
		return
	}
	return
}

func enableChatLogging(config *eqemuconfig.Config) (err error) {
	db, err := connectDB(config)
	if err != nil {
		return
	}
	var result sql.Result
	result, err = db.Exec(theCommand)
	if err != nil {
		return
	}
	var num int64
	num, err = result.RowsAffected()
	if err != nil {
		return
	}
	if num < 1 {
		err = fmt.Errorf("No rows affected")
	}
	return
}
