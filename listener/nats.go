package listener

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/eqemu/server/protobuf/go/eqproto"
	"github.com/go-yaml/yaml"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
	"github.com/robfig/cron"
	"github.com/xackery/discordeq/discord"
	"github.com/xackery/eqemuconfig"
)

var (
	newNATS     bool
	nc          *nats.Conn
	isCronSet   bool
	chanType    string
	ok          bool
	dailyReport DailyReport
	chans       = map[int]string{
		4:   "auctions:",
		5:   "OOC:",
		13:  "", //encounter
		15:  "", //system wide message
		256: "says:",
		260: "OOC:",
		261: "auctions:",
		263: "announcement:",
		269: "", //broadcast
	}
)

type DailyReport struct {
	DailyGains map[int32]*DailyGain
}

type DailyGain struct {
	CharacterId int32
	Identity    string
	Exp         int32
	Lvl         int32
	Money       int32
}

func GetNATS() (conn *nats.Conn) {
	conn = nc
	return
}

func ListenToNATS(eqconfig *eqemuconfig.Config, disco *discord.Discord) {
	var err error
	config = eqconfig
	channelID = config.Discord.ChannelID

	if err = connectNATS(config); err != nil {
		log.Println("[NATS] Warning while getting NATS connection:", err.Error())
		return
	}

	if err = checkForNATSMessages(nc, disco); err != nil {
		log.Println("[NATS] Warning while checking for messages:", err.Error())
	}
	nc.Close()
	nc = nil
	return
}

func connectNATS(config *eqemuconfig.Config) (err error) {
	if nc != nil {
		return
	}
	if config.NATS.Host != "" && config.NATS.Port != "" {
		if nc, err = nats.Connect(fmt.Sprintf("nats://%s:%s", config.NATS.Host, config.NATS.Port)); err != nil {
			log.Fatal(err)
		}
	} else {
		if nc, err = nats.Connect(nats.DefaultURL); err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("[NATS] Connected\n")
	return
}

func checkForNATSMessages(nc *nats.Conn, disco *discord.Discord) (err error) {
	if !isCronSet {
		isCronSet = true
		c := cron.New()
		c.AddFunc("@midnight", DoDailyReport)
		c.Start()
	}

	LoadDailyReport()

	nc.Subscribe("ChannelMessage", OnChannelMessage)
	nc.Subscribe("AdminMessage", OnAdminMessage)

	for {
		if nc.Status() != nats.CONNECTED {
			log.Println("[NATS] Disconnected, status is", nc.Status())
			break
		}
		time.Sleep(6000 * time.Second)
	}
	return
}

func SendCommand(author string, command string, parameters []string) (err error) {
	if nc == nil {
		err = fmt.Errorf("nats is not connected.")
		return
	}

	commandMessage := &eqproto.CommandMessage{
		Author:  author,
		Command: command,
		Params:  parameters,
	}
	log.Println(commandMessage)
	var cmd []byte
	if cmd, err = proto.Marshal(commandMessage); err != nil {
		err = fmt.Errorf("Failed to marshal command: %s", err.Error())
		return
	}

	var msg *nats.Msg
	if msg, err = nc.Request("CommandMessageWorld", cmd, 2*time.Second); err != nil {
		return
	}

	//now process reply
	if err = proto.Unmarshal(msg.Data, commandMessage); err != nil {
		err = fmt.Errorf("Failed to unmarshal response: %s", err.Error())
		return
	}

	if _, err = disco.SendMessage(config.Discord.CommandChannelID, fmt.Sprintf("**%s** %s: %s", commandMessage.Author, commandMessage.Command, commandMessage.Result)); err != nil {
		log.Printf("[NATS] Error sending message (%s: %s) %s", commandMessage.Author, commandMessage.Result, err.Error())
		err = nil
		return
	}
	return
}

func OnAdminMessage(nm *nats.Msg) {
	var err error
	channelMessage := &eqproto.ChannelMessage{}
	proto.Unmarshal(nm.Data, channelMessage)

	if _, err = disco.SendMessage(config.Discord.CommandChannelID, fmt.Sprintf("**Admin:** %s", channelMessage.Message)); err != nil {
		log.Printf("[NATS] Error sending admin message (%s) %s", channelMessage.Message, err.Error())
		return
	}

	log.Printf("[NATS] AdminMessage: %s\n", channelMessage.Message)
}

func OnChannelMessage(nm *nats.Msg) {
	var err error
	channelMessage := &eqproto.ChannelMessage{}
	proto.Unmarshal(nm.Data, channelMessage)

	if channelMessage.IsEmote {
		channelMessage.ChanNum = channelMessage.Type
	}

	if chanType, ok = chans[int(channelMessage.ChanNum)]; !ok {
		log.Printf("[NATS] Unknown channel: %d with message: %s", channelMessage.ChanNum, channelMessage.Message)
		return
	}

	channelMessage.From = strings.Replace(channelMessage.From, "_", " ", -1)

	if strings.Contains(channelMessage.From, " ") {
		channelMessage.From = fmt.Sprintf("%s [%s]", channelMessage.From[:strings.Index(channelMessage.From, " ")], channelMessage.From[strings.Index(channelMessage.From, " ")+1:])
	}
	channelMessage.From = alphanumeric(channelMessage.From) //purify name to be alphanumeric

	if strings.Contains(channelMessage.Message, "Summoning you to") { //GM messages are relaying to discord!
		return
	}

	//message = message[strings.Index(message, "says ooc, '")+11 : len(message)-padOffset]
	if channelMessage.ChanNum == 269 && strings.Contains(channelMessage.Message, "opened a box to find") {
		channelMessage.From = ":gift:"
		channelMessage.Message += " :gift:"
	}
	if channelMessage.ChanNum == 15 {
		channelMessage.From = ":loudspeaker:"
	}

	if channelMessage.ChanNum == 269 && strings.Contains(channelMessage.Message, "Welcome back to the server,") {
		channelMessage.From = ":hand_splayed::skin-tone-1:"
	}

	if channelMessage.ChanNum == 13 && strings.Contains(channelMessage.Message, "successfully stopped") {
		channelMessage.From = ":whale:"
		channelMessage.Message += " :crocodile:"
	}
	channelMessage.Message = convertLinks("", channelMessage.Message)

	if _, err = disco.SendMessage(channelID, fmt.Sprintf("**%s %s** %s", channelMessage.From, chanType, channelMessage.Message)); err != nil {
		log.Printf("[NATS] Error sending message (%s: %s) %s", channelMessage.From, channelMessage.Message, err.Error())
		return
	}

	log.Printf("[NATS] %d %s: %s\n", channelMessage.ChanNum, channelMessage.From, channelMessage.Message)
}

func sendNATSMessage(from string, message string) {
	if nc == nil {
		log.Println("[NATS] not connected?")
		return
	}
	channelMessage := &eqproto.ChannelMessage{
		//From:    from + " says from discord, '",
		IsEmote: true,
		Message: fmt.Sprintf("%s says from discord, '%s'", from, message),
		ChanNum: 260,
		Type:    260,
	}
	msg, err := proto.Marshal(channelMessage)
	if err != nil {
		log.Printf("[NATS] Error marshalling %s %s: %s", from, message, err.Error())
		return
	}
	err = nc.Publish("ChannelMessageWorld", msg)
	if err != nil {
		log.Printf("[NATS] Error publishing: %s", err.Error())
		return
	}
}

func DoDailyReport() {
	var err error
	topLvl := int32(-1)
	topExp := int32(-1)
	topMoney := int32(-1)
	for k, v := range dailyReport.DailyGains {
		if topExp < 0 || v.Exp > dailyReport.DailyGains[topExp].Exp {
			topExp = k
		}
		if topLvl < 0 || v.Lvl > dailyReport.DailyGains[topLvl].Lvl {
			topLvl = k
		}
		if topMoney < 0 || v.Money > dailyReport.DailyGains[topMoney].Money {
			topMoney = k
		}
	}
	if _, err = disco.SendMessage(channelID, "==== 24 Hour Report ===="); err != nil {
		log.Printf("[NATS] Failed to send 24 hour report: %s", err.Error())
		return
	}
	if topExp >= 0 {
		if _, err = disco.SendMessage(channelID, fmt.Sprintf("Top Experince Gains: %s with %0.2f bottles worth of experience!", dailyReport.DailyGains[topExp].Identity, float32(dailyReport.DailyGains[topExp].Exp/23976503))); err != nil {
			log.Printf("[NATS] Error sending message: %s", err.Error())
			return
		}
	}
	if topLvl >= 0 {
		if _, err = disco.SendMessage(channelID, fmt.Sprintf("Top Level Gains: %s with %d levels gained!", dailyReport.DailyGains[topLvl].Identity, int(dailyReport.DailyGains[topLvl].Lvl))); err != nil {
			log.Printf("[NATS] Error sending message: %s", err.Error())
			return
		}
	}
	if topExp >= 0 {
		if _, err = disco.SendMessage(channelID, fmt.Sprintf("Top Money Gains: %s with %0.2f platinum earned!", dailyReport.DailyGains[topMoney].Identity, float32(dailyReport.DailyGains[topMoney].Money/1000))); err != nil {
			log.Printf("[NATS] Error sending message: %s", err.Error())
			return
		}
	}
	//flush daily reports
	dailyReport.DailyGains = map[int32]*DailyGain{}
}

func SaveDailyReport() {
	var err error
	out, err := yaml.Marshal(&dailyReport)
	if err != nil {
		log.Fatal("Error marshalling daily report:", err.Error())
	}
	ioutil.WriteFile("dailyreport.yml", out, 0644)
}

func LoadDailyReport() {
	var err error
	if _, err := os.Stat("dailyreport.yml"); os.IsNotExist(err) {
		SaveDailyReport()
		return
	}
	inFile, err := ioutil.ReadFile("dailyreport.yml")
	if err != nil {
		log.Fatal("Failed to parse dailyreport.yml:", err.Error())
	}
	err = yaml.Unmarshal(inFile, &dailyReport)
	if err != nil {
		log.Fatal("Failed to unmarshal dailyreport.yml:", err.Error())
	}
}
