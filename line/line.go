package line

import (
	"net/http"

	"strings"

	log "github.com/Sirupsen/logrus"
	user "github.com/liam-lai/ptt-alertor/models/user/redis"
	"github.com/liam-lai/ptt-alertor/myutil"
	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client
var err error

func init() {
	config := myutil.Config("line")
	bot, err = linebot.New(config["channelSecret"], config["channelAccessToken"])
	if err != nil {
		log.Fatal(err)
	}
}

func HandleRequest(r *http.Request) {
	events, err := bot.ParseRequest(r)
	if err != nil {
		log.Error(err)
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

/*
 * 清單
 * 追蹤 LoL 樂透
 * 取消 LoL 樂透
 * 指令
 */
func handleMessage(event *linebot.Event) {
	var responseText string
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		args := strings.Fields(message.Text)
		command := args[0]
		if strings.EqualFold(command, "清單") {
			responseText = "lol: 樂透"
		} else if strings.EqualFold(command, "指令") {
			responseText = "清單\n追蹤\n取消\n指令\n"
		} else if strings.EqualFold(command, "追蹤") {
			responseText = "追蹤成功"
		} else if strings.EqualFold(command, "取消") {
			responseText = "取消成功"
		} else {
			responseText = "無此指令，請打「指令」查看指令清單"
		}
	}
	_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(responseText)).Do()
	if err != nil {
		log.Error(err)
	}
}

func handleFollow(event *linebot.Event) {
	userID := event.Source.UserID
	profile, err := bot.GetProfile(userID).Do()
	if err != nil {
		log.Error(err)
	}
	log.Info("User Follow:" + profile.UserID)

	u := new(user.User).Find(profile.UserID)

	if u.Profile.Account != "" {
		u.Enable = true
		u.Update()
	} else {
		u.Profile.Account = profile.UserID
		u.Profile.Line = profile.UserID
		u.Enable = true
		err = u.Save()
		if err != nil {
			log.Error(err)
		}
	}

	_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(profile.DisplayName+" 歡迎使用 PTT Alertor")).Do()
	if err != nil {
		log.Error(err)
	}
}

func handleUnfollow(event *linebot.Event) {
	userID := event.Source.UserID
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

func broadcastTextMessage(ids []string, message string) {
	_, err := bot.Multicast(ids, linebot.NewTextMessage(message)).Do()
	if err != nil {
		panic(err)
	}
}
