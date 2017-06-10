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
	command := strings.Fields(strings.TrimSpace(text))[0]
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
		re := regexp.MustCompile("^(新增|刪除)\\s+([^,，][\\w\\d-_,，\\s]+[^,，])\\s+([^,，].*[^,，]$)")
		matched := re.MatchString(text)
		if !matched {
			return "指令格式錯誤。看板與關鍵字欄位開始與最後不可有逗號。範例：\n" + command + " gossiping,lol 問卦,爆卦"
		}
		args := re.FindStringSubmatch(text)
		boardName := strings.Replace(args[2], " ", "", -1)
		keywords := splitKeywords(args[3])
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

func splitKeywords(keywordText string) (keywords []string) {

	if !strings.ContainsAny(keywordText, ",，") {
		return []string{keywordText}
	}

	if strings.Contains(keywordText, ",") {
		keywords = strings.Split(keywordText, ",")
	} else {
		keywords = []string{keywordText}
	}

	for i := 0; i < len(keywords); i++ {
		if strings.Contains(keywords[i], "，") {
			keywords = append(keywords[:i], append(strings.Split(keywords[i], "，"), keywords[i+1:]...)...)
			i--
		}
	}

	for i, keyword := range keywords {
		keywords[i] = strings.TrimSpace(keyword)
	}

	return keywords
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
