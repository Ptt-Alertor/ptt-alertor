package line

import (
	"fmt"
	"net/http"
	"os"

	"strings"

	"regexp"

	"github.com/julienschmidt/httprouter"
	"github.com/line/line-bot-sdk-go/linebot"
	log "github.com/meifamily/logrus"
	"github.com/meifamily/ptt-alertor/command"
	"github.com/meifamily/ptt-alertor/models"
	"github.com/meifamily/ptt-alertor/myutil"
	"github.com/meifamily/ptt-alertor/shorturl"
)

const maxCharacters = 2000

var (
	bot                *linebot.Client
	err                error
	channelSecret      = os.Getenv("LINE_CHANNEL_SECRET")
	channelAccessToken = os.Getenv("LINE_CHANNEL_ACCESSTOKEN")
)

func init() {
	bot, err = linebot.New(channelSecret, channelAccessToken)
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

func handleMessage(event *linebot.Event) {
	var responseText string
	var lineMsg []linebot.Message
	userID := event.Source.UserID
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		text := strings.TrimSpace(message.Text)
		if strings.EqualFold(text, "notify") {
			lineMsg = append(lineMsg, linebot.NewTextMessage(shorturl.Gen(getAuthorizeURL(userID))))
			replyMessage(event.ReplyToken, lineMsg...)
			return
		}
		if !checkLineAccessTokenExist(userID) {
			lineMsg = append(lineMsg, linebot.NewTextMessage(getLineNotifyConnectMessage("", userID)))
			replyMessage(event.ReplyToken, lineMsg...)
			return
		}
		if match, _ := regexp.MatchString("^(刪除|刪除作者)+\\s.*\\*+", text); match {
			replyMessage(event.ReplyToken, genConfirmMessage(text))
			return
		}
		responseText = command.HandleCommand(text, userID)
	}
	for _, msg := range myutil.SplitTextByLineBreak(responseText, maxCharacters) {
		lineMsg = append(lineMsg, linebot.NewTextMessage(msg))
	}
	replyMessage(event.ReplyToken, lineMsg...)
}

func handleFollow(event *linebot.Event) {
	userID := event.Source.UserID
	profile, err := bot.GetProfile(userID).Do()
	if err != nil {
		log.WithError(err).Error("Line Get Profile Failed")
	}
	id := profile.UserID
	err = command.HandleLineFollow(id)
	if err != nil {
		log.WithError(err).Error("Line Follow Failed")
	}

	text := getLineNotifyConnectMessage(profile.DisplayName+" ", id)
	_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(text)).Do()
	if err != nil {
		log.WithError(err).Error("Line Follow Replay Message Failed")
	}
}

func getLineNotifyConnectMessage(displayName, userID string) string {
	url := shorturl.Gen(getAuthorizeURL(userID))
	return fmt.Sprintf(`%s歡迎使用 Ptt Alertor。
請按以下步驟啟用 LINE Notify 以獲得最新文章通知。
1. 開啟下方網址
2. 選擇”群組“ 第一個「透過1對1聊天接收LINE Notify的通知
3. 點擊「同意並連動
%s

你將可以在 Ptt Alertor 設定看板、作者、關鍵字，並在 LINE Notify 得到最新文章。
觀看Demo:
https://media.giphy.com/media/l0Iy28oboQbSw6Cn6/giphy.gif`, displayName, url)
}

func handleUnfollow(event *linebot.Event) {
	userID := event.Source.UserID
	log.WithFields(log.Fields{
		"ID": userID,
	}).Info("Line Unfollow")
	u := models.User.Find(userID)
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

func replyMessage(token string, message ...linebot.Message) {
	_, err := bot.ReplyMessage(token, message...).Do()
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
