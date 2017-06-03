/**
 * [x]remove multiple keywords at once
 * [x]add multiple keywords at once
 **/

package line

import (
	"net/http"
	"regexp"

	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/liam-lai/ptt-alertor/models/subscription"
	user "github.com/liam-lai/ptt-alertor/models/user/redis"
	"github.com/liam-lai/ptt-alertor/myutil"
	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client
var err error
var commands = map[string]string{
	"指令": "可使用的指令清單",
	"清單": "目前追蹤的看板與關鍵字",
	"新增": "新增看板關鍵字。範例：\n\t\t新增 gossiping 爆卦\n\t\t新增 gossiping 爆卦,問卦\n\t\t新增 gossiping 爆卦，問卦",
	"刪除": "刪除看板關鍵字。範例：\n\t\t刪除 gossiping 爆卦\n\t\t刪除 gossiping 爆卦,問卦\n\t\t刪除 gossiping 爆卦，問卦",
}

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
		responseText = handleCommand(message.Text, userID)
	}
	_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(responseText)).Do()
	if err != nil {
		log.WithError(err).Error("Line Reply Message Failed")
	}
}

func handleCommand(text string, userID string) string {
	args := strings.Fields(text)
	command, args := args[0], args[1:]
	switch command {
	case "清單":
		return new(user.User).Find(userID).Subscribes.String()
	case "指令":
		return stringCommands()
	case "新增", "刪除":
		matched, err := regexp.MatchString("^(新增|刪除)(\\s+)([\\w\\d-_]+)(\\s+)([^\\s]+)$", text)
		if err != nil {
			log.WithError(err).Error("Line Check Command Failed")
		}
		if !matched {
			return "指令格式錯誤。關鍵字與逗號間不可有空格。範例：\n" + command + " gossiping 問卦,爆卦"
		}
		board := args[0]
		keywords := splitKeywords(args[1])
		if command == "新增" {
			err := subscribe(userID, board, keywords)
			if err != nil {
				return "新增失敗，請等待修復。"
			}
			return "新增成功"
		}
		if command == "刪除" {
			err := unsubscribe(userID, board, keywords)
			if err != nil {
				return "刪除失敗，請等待修復。"
			}
			return "刪除成功"
		}
	}
	return "無此指令，請打「指令」查看指令清單"
}

func stringCommands() string {
	str := ""
	for key, val := range commands {
		str += key + "：" + val + "\n"
	}
	return str
}

func splitKeywords(keywords string) []string {
	if strings.Contains(keywords, ",") {
		return strings.Split(keywords, ",")
	}

	if strings.Contains(keywords, "，") {
		return strings.Split(keywords, "，")
	}

	return []string{keywords}
}

func subscribe(account string, board string, keywords []string) error {
	u := new(user.User).Find(account)
	sub := subscription.Subscribe{
		Board:    board,
		Keywords: keywords,
	}
	u.Subscribes.Add(sub)
	err := u.Update()
	if err != nil {
		log.WithError(err).Error("Line Subscribe Update Error")
	}
	return err
}

func unsubscribe(account string, board string, keywords []string) error {
	u := new(user.User).Find(account)
	sub := subscription.Subscribe{
		Board:    board,
		Keywords: keywords,
	}
	u.Subscribes.Remove(sub)
	err := u.Update()
	if err != nil {
		log.WithError(err).Error("Line UnSubscribe Update Error")
	}
	return err
}

func handleFollow(event *linebot.Event) {
	userID := event.Source.UserID
	profile, err := bot.GetProfile(userID).Do()
	if err != nil {
		log.WithError(err).Error("")
	}

	log.WithFields(log.Fields{
		"ID": profile.UserID,
	}).Info("Line Follow")

	u := new(user.User).Find(profile.UserID)

	if u.Profile.Account != "" {
		log.WithFields(log.Fields{
			"ID": profile.UserID,
		}).Info("Line ReFollow")
		u.Enable = true
		u.Update()
	} else {
		log.WithFields(log.Fields{
			"ID": profile.UserID,
		}).Info("Line Follow")
		u.Profile.Account = profile.UserID
		u.Profile.Line = profile.UserID
		u.Enable = true
		err = u.Save()
		if err != nil {
			log.WithError(err).Error("Line Follow Save User Failed")
		}
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
