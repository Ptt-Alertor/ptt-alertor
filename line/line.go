package line

import (
	"net/http"

	"strings"

	"regexp"

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
		case linebot.EventTypePostback:
			handlePostback(event)
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
			if match, _ := regexp.MatchString("^(刪除|刪除作者)+\\s.*\\*+", text); match {
				replyMessage(event.ReplyToken, genConfirmMessage(text))
				return
			}
			responseText = command.HandleCommand(text, userID)
		}
	}
	replyMessage(event.ReplyToken, linebot.NewTextMessage(responseText))
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
	_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(profile.DisplayName+text+url+"\n\n你將可以在 PTT Alertor 設定看板、作者、關鍵字，並在 LINE Notify 得到最新文章。\n\n觀看Demo:\nhttps://media.giphy.com/media/l0Iy28oboQbSw6Cn6/giphy.gif")).Do()
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

func handlePostback(event *linebot.Event) {
	data := event.Postback.Data
	if data == "cancel" {
		replyMessage(event.ReplyToken, linebot.NewTextMessage("取消"))
		return
	}
	responseText := command.HandleCommand(data, event.Source.UserID)
	replyMessage(event.ReplyToken, linebot.NewTextMessage(responseText))
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

func genConfirmMessage(command string) *linebot.TemplateMessage {
	leftBtn := linebot.NewPostbackTemplateAction("是", command, "")
	rightBtn := linebot.NewPostbackTemplateAction("否", "cancel", "")

	template := linebot.NewConfirmTemplate("確定"+command+"？", leftBtn, rightBtn)
	message := linebot.NewTemplateMessage("批次刪除", template)
	return message
}

func replyMessage(token string, message linebot.Message) {
	_, err := bot.ReplyMessage(token, message).Do()
	if err != nil {
		log.WithError(err).Error("Line Reply Message Failed")
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
