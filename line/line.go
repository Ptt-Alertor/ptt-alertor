/**
 * [x]remove multiple keywords at once
 * [x]add multiple keywords at once
 **/

package line

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/liam-lai/ptt-alertor/command"
	board "github.com/liam-lai/ptt-alertor/models/ptt/board/redis"
	user "github.com/liam-lai/ptt-alertor/models/user/redis"
	"github.com/liam-lai/ptt-alertor/myutil"
	"github.com/line/line-bot-sdk-go/linebot"
)

var bd board.Board
var bot *linebot.Client
var err error

func init() {
	config := myutil.Config("line")
	bot, err = linebot.New(config["channelSecret"], config["channelAccessToken"])
	if err != nil {
		log.Fatal(err)
	}
}

func HandleRequest(_ http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	events, err := bot.ParseRequest(r)
	if err != nil {
		log.WithError(err).Error("Line ParseRequest Error")
	}
	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeMessage:
			handleMessage(event)
		case linebot.EventTypeFollow:
			handleFollow(event)
		case linebot.EventTypeUnfollow:
			handleUnfollow(event)
		}
	}
}

/**
 * TODO: check board exist or not
 * 1. hotboard
 * 2. allboard
 **/
func handleMessage(event *linebot.Event) {
	var responseText string
	userID := event.Source.UserID
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		responseText = command.HandleCommand(message.Text, userID)
	}
	_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(responseText)).Do()
	if err != nil {
		log.WithError(err).Error("Line Reply Message Failed")
	}
}

func handleFollow(event *linebot.Event) {
	userID := event.Source.UserID
	profile, err := bot.GetProfile(userID).Do()
	if err != nil {
		log.WithError(err).Error("")
	}

	id := profile.UserID

	log.WithFields(log.Fields{
		"ID": id,
	}).Info("Line Follow")

	err = command.HandleLineFollow(id)
	if err != nil {
		log.WithError(err).Error("Line Follow Error")
	}

	_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(profile.DisplayName+" 歡迎使用 PTT Alertor\n輸入「指令」查看相關功能。")).Do()
	if err != nil {
		log.WithError(err).Error("Line Follow Replay Message Failed")
	}
}

func handleUnfollow(event *linebot.Event) {
	userID := event.Source.UserID
	log.WithFields(log.Fields{
		"ID": userID,
	}).Info("Line Unfollow")
	u := new(user.User).Find(userID)
	u.Enable = false
	u.Update()
}

// useless
func handleLeave() {

}

// useless
func handleJoin() {

}

// useless
func handlePostback() {

}

// useless
func handleBeacon() {

}

func PushTextMessage(id string, message string) {
	_, err := bot.PushMessage(id, linebot.NewTextMessage(message)).Do()
	if err != nil {
		log.WithError(err).Error("Line Push Message Failed")
	} else {
		log.WithFields(log.Fields{
			"ID": id,
		}).Info("Line Push Message")
	}
}

func BroadcastTextMessage(ids []string, message string) {
	_, err := bot.Multicast(ids, linebot.NewTextMessage(message)).Do()
	if err != nil {
		log.WithError(err).Error("Line Broadcast Message Failed")
	} else {
		log.WithFields(log.Fields{
			"IDs": ids,
		}).Info("Line BroadCast Message")
	}
}
