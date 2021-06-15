package line

import (
	"fmt"
	"net/http"
	"os"

	"strings"

	"regexp"

	log "github.com/Ptt-Alertor/logrus"
	"github.com/Ptt-Alertor/ptt-alertor/command"
	"github.com/Ptt-Alertor/ptt-alertor/models"
	"github.com/Ptt-Alertor/ptt-alertor/myutil"
	"github.com/Ptt-Alertor/ptt-alertor/shorturl"
	"github.com/julienschmidt/httprouter"
	"github.com/line/line-bot-sdk-go/linebot"
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
		case linebot.EventTypeFollow, linebot.EventTypeJoin:
			handleFollowAndJoin(event)
		case linebot.EventTypeUnfollow, linebot.EventTypeLeave:
			handleUnfollowAndLeave(event)
		}
	}
}

func handleMessage(event *linebot.Event) {
	var responseText string
	var lineMsg []linebot.SendingMessage
	accountID, accountType := getAccountIDAndType(event)

	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		text := strings.TrimSpace(message.Text)
		if strings.EqualFold(text, "notify") {
			lineMsg = append(lineMsg, linebot.NewTextMessage(shorturl.Gen(getAuthorizeURL(accountID))))
			replyMessage(event.ReplyToken, lineMsg...)
			return
		}
		if !checkLineAccessTokenExist(accountID) {
			lineMsg = append(lineMsg, linebot.NewTextMessage(getLineNotifyConnectMessage(accountID, accountType)))
			replyMessage(event.ReplyToken, lineMsg...)
			return
		}
		if match, _ := regexp.MatchString("^(刪除|刪除作者)+\\s.*\\*+", text); match {
			replyMessage(event.ReplyToken, genConfirmMessage(text))
			return
		}
		responseText = command.HandleCommand(text, accountID, accountType == accountTypeUser)
	}

	if responseText == "" {
		return
	}
	for _, msg := range myutil.SplitTextByLineBreak(responseText, maxCharacters) {
		lineMsg = append(lineMsg, linebot.NewTextMessage(msg))
	}
	replyMessage(event.ReplyToken, lineMsg...)
}

func handleFollowAndJoin(event *linebot.Event) {
	// TODO: make all user naming change to account
	// account will include group, room and user and make accountType as enum
	accountID, accountType := getAccountIDAndType(event)
	if err = command.HandleLineFollow(accountID, accountType); err != nil {
		log.WithError(err).Error("Line Follow Failed")
	}

	text := getLineNotifyConnectMessage(accountID, accountType)
	_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(text)).Do()
	if err != nil {
		log.WithError(err).Error("Line Follow Reply Message Failed")
	}
}

func getLineNotifyConnectMessage(accountID, accountType string) string {
	url := shorturl.Gen(getAuthorizeURL(accountID))

	groupName := "透過1對1聊天接收LINE Notify的通知"
	demoLink := "https://media.giphy.com/media/l0Iy28oboQbSw6Cn6/giphy.gif"
	if accountType == accountTypeGroup {
		groupName = "此群組名稱"
		demoLink = "https://www.facebook.com/PttAlertor/posts/975327009676889"
	}

	return fmt.Sprintf(`歡迎使用 Ptt Alertor。
請按以下步驟啟用 LINE Notify 以獲得最新文章通知。
1. 開啟下方網址
2. 選擇「%s」
3. 點擊「同意並連動」
%s

你將可以在 Ptt Alertor 設定看板、作者、關鍵字，並在 LINE Notify 得到最新文章。
觀看Demo:
%s`, groupName, url, demoLink)
}

func handleUnfollowAndLeave(event *linebot.Event) {
	accountID, _ := getAccountIDAndType(event)
	log.WithFields(log.Fields{
		"ID": accountID,
	}).Info("Line Unfollow")
	u := models.User().Find(accountID)
	u.Enable = false
	u.Update()
}

func handlePostback(event *linebot.Event) {
	data := event.Postback.Data
	if data == "cancel" {
		replyMessage(event.ReplyToken, linebot.NewTextMessage("取消"))
		return
	}
	accountID, accountType := getAccountIDAndType(event)
	responseText := command.HandleCommand(data, accountID, accountType == accountTypeUser)
	if responseText == "" {
		return
	}
	replyMessage(event.ReplyToken, linebot.NewTextMessage(responseText))
}

// useless
func handleBeacon() {

}

const (
	accountTypeUser  = "user"
	accountTypeGroup = "group"
	accountTypeRoom  = "room"
)

func getAccountIDAndType(event *linebot.Event) (id, accountType string) {
	switch event.Source.Type {
	case linebot.EventSourceTypeUser:
		return event.Source.UserID, accountTypeUser
	case linebot.EventSourceTypeGroup:
		return event.Source.GroupID, accountTypeGroup
	case linebot.EventSourceTypeRoom:
		return event.Source.RoomID, accountTypeRoom
	}
	return
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
	leftBtn := linebot.NewPostbackAction("是", command, "", "")
	rightBtn := linebot.NewPostbackAction("否", "cancel", "", "")

	template := linebot.NewConfirmTemplate("確定"+command+"？", leftBtn, rightBtn)
	message := linebot.NewTemplateMessage("批次刪除", template)
	return message
}

func replyMessage(token string, message ...linebot.SendingMessage) {
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
