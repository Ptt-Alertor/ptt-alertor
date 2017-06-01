package line

import (
	"net/http"

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
	"清單": "正在追蹤的看板與關鍵字",
	"新增": "新增訂閱看板關鍵字。ex: 新增 lol 樂透",
	"刪除": "取消訂閱看板關鍵字。ex: 刪除 lol 樂透",
	"指令": "可使用的指令清單",
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

/**
 * check board exist or not
 * 1. hotboard
 * 2. allboard
 **/
func handleMessage(event *linebot.Event) {
	var responseText string
	userID := event.Source.UserID
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		args := strings.Fields(message.Text)
		command := args[0]
		if strings.EqualFold(command, "清單") {
			responseText = new(user.User).Find(userID).Subscribes.String()
		} else if strings.EqualFold(command, "指令") {
			responseText = stringCommands()
		} else if strings.EqualFold(command, "新增") {
			if len(args) < 3 {
				responseText = "指令格式錯誤。\n範例:\n新增 lol 樂透\n新增 lol 樂透,電競"
			} else {
				err := subscribe(userID, args[1], strings.Split(args[2], ","))
				if err != nil {
					responseText = "新增失敗，請等待修復。"
				} else {
					responseText = "新增成功"
				}
			}
		} else if strings.EqualFold(command, "刪除") {
			if len(args) < 3 {
				responseText = "指令格式錯誤。\n範例:\n刪除 lol 樂透\n刪除 lol 樂透,電競"
			} else {
				err := unsubscribe(userID, args[1], strings.Split(args[2], ","))
				if err != nil {
					responseText = "刪除失敗，請等待修復。"
				} else {
					responseText = "刪除成功"
				}
			}
		} else {
			responseText = "無此指令，請打「指令」查看指令清單"
		}
	}
	_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(responseText)).Do()
	if err != nil {
		log.Error(err)
	}
}

func stringCommands() string {
	str := ""
	for key, val := range commands {
		str += key + ": " + val + "\n"
	}
	return str
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
		log.Error(err)
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
		log.Error(err)
	}
	return err
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
