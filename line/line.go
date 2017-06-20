package line

import (
	"net/http"

	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/liam-lai/ptt-alertor/command"
	board "github.com/liam-lai/ptt-alertor/models/ptt/board/redis"
	user "github.com/liam-lai/ptt-alertor/models/user/redis"
	"github.com/liam-lai/ptt-alertor/myutil"
	"github.com/liam-lai/ptt-alertor/shorturl"
	"github.com/line/line-bot-sdk-go/linebot"
)

var bd board.Board
var bot *linebot.Client
var err error
var config map[string]string

func init() {
	config = myutil.Config("line")
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
		text := strings.TrimSpace(message.Text)
		if strings.EqualFold(text, "notify") {
			responseText = shorturl.Gen(getAuthorizeURL(userID))
		} else {
			responseText = command.HandleCommand(text, userID)
		}
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
	url := shorturl.Gen(getAuthorizeURL(id))
	text := " 歡迎使用 PTT Alertor。\n請按以下步驟啟用 LINE Notify 以獲得最新文章通知。\n1. 開啟下方網址\n2. 選擇”群組“ 第一個「透過1對1聊天接收LINE Notify的通知\n3. 點擊「同意並連動」\n"
	_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(profile.DisplayName+text+url+"\n\n你將可以在 PTT Alertor 設定看板、作者、關鍵字，並在 LINE Notify 得到最新文章。")).Do()
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
