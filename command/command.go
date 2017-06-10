package command

import (
	"fmt"
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
	"新增": "新增看板關鍵字。範例：\n\t\t新增 gossiping 爆卦\n\t\t新增 gossiping 爆卦,問卦\n\t\t新增 gossiping 爆卦，問卦",
	"刪除": "刪除看板關鍵字。範例：\n\t\t刪除 gossiping 爆卦\n\t\t刪除 gossiping 爆卦,問卦\n\t\t刪除 gossiping 爆卦，問卦",
}

func HandleCommand(text string, userID string) string {
	args := strings.Fields(text)
	command, args := args[0], args[1:]
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
		matched, err := regexp.MatchString("^(新增|刪除)(\\s+)([\\w\\d-_]+)(\\s+)([^\\s]+)$", text)
		if err != nil {
			log.WithError(err).Error("Line Check Command Failed")
		}
		if !matched {
			return "指令格式錯誤。關鍵字與逗號間不可有空格。範例：\n" + command + " gossiping 問卦,爆卦"
		}
		boardName := args[0]
		keywords := splitKeywords(args[1])
		if command == "新增" {
			err := subscribe(userID, boardName, keywords)
			if bErr, ok := err.(boardproto.BoardNotExistError); ok {
				return "版名錯誤，請確認拼字。可能版名：\n" + bErr.Suggestion
			}
			if err != nil {
				return "新增失敗，請等待修復。"
			}
			return "新增成功"
		}
		if command == "刪除" {
			err := unsubscribe(userID, boardName, keywords)
			if bErr, ok := err.(boardproto.BoardNotExistError); ok {
				return "版名錯誤，請確認拼字。可能版名：\n" + bErr.Suggestion
			}
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
	for key, val := range Commands {
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

func subscribe(account string, boardname string, keywords []string) error {
	u := new(user.User).Find(account)
	sub := subscription.Subscribe{
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
	sub := subscription.Subscribe{
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
	fmt.Printf("%+v", u)
	return nil
}
