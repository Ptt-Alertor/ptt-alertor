package command

import (
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	boardproto "github.com/liam-lai/ptt-alertor/models/ptt/board"
	"github.com/liam-lai/ptt-alertor/models/subscription"
	user "github.com/liam-lai/ptt-alertor/models/user/redis"
)

var Commands = map[string]string{
	"指令": "可使用的指令清單",
	"清單": "目前追蹤的看板與關鍵字",
	"新增": "新增看板關鍵字。範例：\n\t\t新增 nba 樂透\n\t\t新增 nba,lol 樂透 \n\t\t新增 nba,lol 樂透,情報",
	"刪除": "刪除看板關鍵字。範例：\n\t\t刪除 nba 樂透\n\t\t刪除 nba,lol 樂透 \n\t\t刪除 nba,lol 樂透,情報",
}

func HandleCommand(text string, userID string) string {
	command := strings.Fields(strings.TrimSpace(text))[0]
	log.WithFields(log.Fields{
		"account": userID,
		"command": command,
	}).Info("Command Request")
	switch command {
	case "清單":
		rspText := new(user.User).Find(userID).Subscribes.String()
		if rspText == "" {
			rspText = "尚未建立清單。請打「指令」查看新增方法。"
		}
		return rspText
	case "指令":
		return stringCommands()
	case "新增", "刪除":
		re := regexp.MustCompile("^(新增|刪除)\\s+([^,，][\\w\\d-_,，]+[^,，])\\s+(.+)")
		matched := re.MatchString(text)
		if !matched {
			return "指令格式錯誤。\n1.板名欄位開頭與結尾不可有逗號\n2.板名欄位間不允許空白字元。\n正確範例：" + command + " gossiping,lol 問卦,爆卦"
		}
		args := re.FindStringSubmatch(text)
		boardNames := splitParamString(args[2])
		keywords := splitParamString(args[3])
		if command == "新增" {
			for _, boardName := range boardNames {
				err := subscribe(userID, boardName, keywords)
				if bErr, ok := err.(boardproto.BoardNotExistError); ok {
					return "版名錯誤，請確認拼字。可能版名：\n" + bErr.Suggestion
				}
				if err != nil {
					return "新增失敗，請等待修復。"
				}
			}
			return "新增成功"
		}
		if command == "刪除" {
			for _, boardName := range boardNames {
				err := unsubscribe(userID, boardName, keywords)
				if bErr, ok := err.(boardproto.BoardNotExistError); ok {
					return "版名錯誤，請確認拼字。可能版名：\n" + bErr.Suggestion
				}
				if err != nil {
					return "刪除失敗，請等待修復。"
				}
			}
			return "刪除成功"
		}
	}
	return "無此指令，請打「指令」查看指令清單"
}

func stringCommands() string {
	str := ""
	for key, val := range Commands {
		str += key + "：" + val + "\n"
	}
	return str
}

func splitParamString(paramString string) (params []string) {

	paramString = strings.Trim(paramString, ",，")

	if !strings.ContainsAny(paramString, ",，") {
		return []string{paramString}
	}

	if strings.Contains(paramString, ",") {
		params = strings.Split(paramString, ",")
	} else {
		params = []string{paramString}
	}

	for i := 0; i < len(params); i++ {
		if strings.Contains(params[i], "，") {
			params = append(params[:i], append(strings.Split(params[i], "，"), params[i+1:]...)...)
			i--
		}
	}

	for i, param := range params {
		params[i] = strings.TrimSpace(param)
	}

	return params
}

func subscribe(account string, boardname string, keywords []string) error {
	u := new(user.User).Find(account)
	sub := subscription.Subscription{
		Board:    boardname,
		Keywords: keywords,
	}
	err := u.Subscribes.Add(sub)
	if err != nil {
		return err
	}
	err = u.Update()
	if err != nil {
		log.WithError(err).Error("Line Subscribe Update Error")
	}
	return err
}

func unsubscribe(account string, board string, keywords []string) error {
	u := new(user.User).Find(account)
	sub := subscription.Subscription{
		Board:    board,
		Keywords: keywords,
	}
	err := u.Subscribes.Remove(sub)
	if err != nil {
		return err
	}
	err = u.Update()
	if err != nil {
		log.WithError(err).Error("Line UnSubscribe Update Error")
	}
	return err
}

func HandleLineFollow(id string) error {
	u := new(user.User).Find(id)
	u.Profile.Line = id
	return handleFollow(u)
}

func HandleMessengerFollow(id string) error {
	u := new(user.User).Find(id)
	u.Profile.Messenger = id
	return handleFollow(u)
}

func handleFollow(u user.User) error {
	if u.Profile.Account != "" {
		u.Enable = true
		u.Update()
	} else {
		if u.Profile.Messenger != "" {
			u.Profile.Account = u.Profile.Messenger
		} else {
			u.Profile.Account = u.Profile.Line
		}
		u.Enable = true
		err := u.Save()
		if err != nil {
			return err
		}
	}
	return nil
}
