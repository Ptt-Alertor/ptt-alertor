package command

import (
	"regexp"
	"strings"

	"fmt"

	log "github.com/meifamily/logrus"
	"github.com/meifamily/ptt-alertor/crawler"
	"github.com/meifamily/ptt-alertor/models/ptt/article"
	boardproto "github.com/meifamily/ptt-alertor/models/ptt/board"
	"github.com/meifamily/ptt-alertor/models/subscription"
	"github.com/meifamily/ptt-alertor/models/top"
	user "github.com/meifamily/ptt-alertor/models/user/redis"
)

const subArticlesLimit int = 25

var Commands = map[string]map[string]string{
	"一般": {
		"指令": "可使用的指令清單",
		"清單": "設定的看板、關鍵字、作者",
		"排行": "前五名追蹤的關鍵字、作者",
	},
	"關鍵字相關": {
		"新增 看板 關鍵字": "新增追蹤關鍵字",
		"刪除 看板 關鍵字": "取消追蹤關鍵字",
		"範例":        "新增 gossiping,movie 金城武,結衣",
	},
	"作者相關": {
		"新增作者 看板 作者": "新增追蹤作者",
		"刪除作者 看板 作者": "取消追蹤作者",
		"範例":         "新增作者 gossiping ffaarr,obov",
	},
	"推噓文數相關": {
		"新增(推/噓)文數 看板 總數": "通知推或噓文數",
		"範例":    "新增推文數 joke,beauty 10",
		"歸零即刪除": "新增噓文數 joke 0",
	},
	"推文相關": {
		"新增推文 網址": "新增推文追蹤",
		"刪除推文 網址": "刪除推文追蹤",
		"範例":      "新增推文 https://www.ptt.cc/bbs/EZsoft/M.1497363598.A.74E.html",
	},
	"進階應用": {
		"參考連結": "https://pttalertor.dinolai.com/docs",
	},
}

var commandActionMap = map[string]updateAction{
	"新增":    addKeywords,
	"刪除":    removeKeywords,
	"新增作者":  addAuthors,
	"刪除作者":  removeAuthors,
	"新增推文":  addArticles,
	"刪除推文":  removeArticles,
	"新增推文數": updatePushUp,
	"新增噓文數": updatePushDown,
}

func HandleCommand(text string, userID string) string {
	command := strings.Fields(strings.TrimSpace(text))[0]
	log.WithFields(log.Fields{
		"account": userID,
		"command": command,
	}).Info("Command Request")
	switch command {
	case "debug":
		return handleDebug(userID)
	case "清單":
		return handleList(userID)
	case "指令":
		return stringCommands()
	case "排行":
		return listTop()
	case "新增", "刪除":
		return handleKeyword(command, userID, text)
	case "新增作者", "刪除作者":
		return handleAuthor(command, userID, text)
	case "新增推文數", "新增噓文數":
		return handlePushSum(command, userID, text)
	case "新增推文", "刪除推文":
		return handlePush(command, userID, text)
	case "清理推文":
		return cleanPushList(userID)
	case "推文清單":
		return handlePushList(userID)
	}
	return "無此指令，請打「指令」查看指令清單"
}

func handleDebug(account string) string {
	profile := new(user.User).Find(account).Profile
	return profile.Account
}

func handleList(account string) string {
	subs := new(user.User).Find(account).Subscribes
	if len(subs) == 0 {
		return "尚未建立清單。請打「指令」查看新增方法。"
	}
	return subs.String()
}

func cleanPushList(account string) string {
	subs := new(user.User).Find(account).Subscribes
	var i int
	for _, sub := range subs {
		for _, code := range sub.Articles {
			a := article.Article{
				Code: code,
			}
			bl, err := a.Exist()
			if err != nil {
				return "清理推文失敗，請洽至粉絲團或 LINE 首頁留言。"
			}
			if !bl {
				update(removeArticles, account, []string{sub.Board}, code)
				i++
			}
		}
	}
	return fmt.Sprintf("清理 %d 則推文", i)
}

func handlePushList(account string) string {
	subs := new(user.User).Find(account).Subscribes
	if len(subs) == 0 {
		return "尚未建立清單。請打「指令」查看新增方法。"
	}
	return "推文追蹤清單，上限 25 篇：\n" + subs.StringPushList() + "\n輸入「清理推文」，可刪除無效連結。"
}

func stringCommands() string {
	str := ""
	for cat, cmds := range Commands {
		str += "[" + cat + "]\n"
		for cmd, doc := range cmds {
			str += cmd
			if doc != "" {
				str += "：" + doc
			}
			str += "\n"
		}
		str += "\n"
	}
	return strings.TrimSpace(str)
}

func listTop() string {
	content := "關鍵字"
	for i, keyword := range top.ListKeywords(5) {
		content += fmt.Sprintf("\n%d. %s", i+1, keyword)
	}
	content += "\n----\n作者"
	for i, author := range top.ListAuthors(5) {
		content += fmt.Sprintf("\n%d. %s", i+1, author)
	}
	content += "\n----\n推噓文"
	for i, pushSum := range top.ListPushSum(5) {
		content += fmt.Sprintf("\n%d. %s", i+1, pushSum)
	}
	content += "\n\nTOP 100:\nhttp://pttalertor.dinolai.com/top"
	return content
}

func handleKeyword(command, userID, text string) string {
	re := regexp.MustCompile("^(新增|刪除)\\s+([^,，][\\w-_,，\\.]*[^,，:\\s]):?\\s+(\\*|.*[^\\s])")
	matched := re.MatchString(text)
	if !matched {
		return inputErrorTips() + "\n\n正確範例：\n" + command + " gossiping,lol 問卦,爆卦"
	}
	args := re.FindStringSubmatch(text)
	boardNames := splitParamString(args[2])
	input := args[3]
	var inputs []string
	if strings.HasPrefix(input, "regexp:") {
		if !checkRegexp(input) {
			return "正規表示式錯誤，請檢查規則。"
		}
		inputs = []string{args[3]}
	} else {
		inputs = splitParamString(args[3])
	}
	log.WithFields(log.Fields{
		"id":      userID,
		"command": command,
		"boards":  boardNames,
		"words":   inputs,
	}).Info("Keyword Command")
	err := update(commandActionMap[command], userID, boardNames, inputs...)
	if msg, ok := checkBoardError(err); ok {
		return msg
	}
	if err != nil {
		log.WithError(err).Error("Keyword Command Failed")
		return command + "失敗，請嘗試封鎖再解封鎖，並重新執行註冊步驟。\n若問題未解決，請至粉絲團或 LINE 首頁留言。"
	}
	return command + "成功"

}

func handleAuthor(command, userID, text string) string {
	re := regexp.MustCompile("^(新增作者|刪除作者)\\s+([^,，][\\w-_,，\\.]*[^,，:\\s]):?\\s+(\\*|[\\s,\\w]+)")
	matched := re.MatchString(text)
	if !matched {
		return inputErrorTips() + "\n4. 作者為半形英文與數字組成。\n\n正確範例：\n" + command + " gossiping,lol ffaarr,obov"
	}
	args := re.FindStringSubmatch(text)
	boardNames := splitParamString(args[2])
	inputs := splitParamString(args[3])
	log.WithFields(log.Fields{
		"id":      userID,
		"command": command,
		"boards":  boardNames,
		"words":   inputs,
	}).Info("Author Command")
	err := update(commandActionMap[command], userID, boardNames, inputs...)
	if msg, ok := checkBoardError(err); ok {
		return msg
	}
	if err != nil {
		log.WithError(err).Error("Author Command Failed")
		return command + "失敗，請嘗試封鎖再解封鎖，並重新執行註冊步驟。\n若問題未解決，請至粉絲團或 LINE 首頁留言。"
	}
	return command + "成功"
}

func handlePush(command, userID, text string) string {
	re := regexp.MustCompile("^(新增推文|刪除推文)\\s+https?://www.ptt.cc/bbs/([\\w-_]*)/(M\\.\\d+.A.\\w*)\\.html$")
	matched := re.MatchString(text)
	if !matched {
		return "指令格式錯誤。\n1. 網址與指令需至少一個空白。\n2. 網址錯誤格式。\n正確範例：\n" + command + " https://www.ptt.cc/bbs/EZsoft/M.1497363598.A.74E.html"
	}
	args := re.FindStringSubmatch(text)
	boardName := args[2]
	articleCode := args[3]
	log.WithFields(log.Fields{
		"id":      userID,
		"command": command,
		"boards":  boardName,
		"words":   articleCode,
	}).Info("Push Command")
	if !checkArticleExist(boardName, articleCode) {
		return "文章不存在"
	}
	if strings.EqualFold("新增推文", command) && countUserArticles(userID) > subArticlesLimit {
		return "推文追蹤最多 25 篇，輸入「推文清單」，整理追蹤列表。"
	}
	err := update(commandActionMap[command], userID, []string{boardName}, articleCode)
	if err != nil {
		log.WithError(err).Error("Pushlist Command Failed")
		return command + "失敗，請嘗試封鎖再解封鎖，並重新執行註冊步驟。\n若問題未解決，請至粉絲團或 LINE 首頁留言。"
	}
	return command + "成功"
}

func handlePushSum(command, account, text string) string {
	re := regexp.MustCompile("^(新增推文數|新增噓文數)\\s+([^,，][\\w-_,，\\.]*[^,，:\\s]):?\\s+(100|[1-9][0-9]|[0-9])$")
	matched := re.MatchString(text)
	if !matched {
		return inputErrorTips() + "\n4. 推噓文數需為介於 0-100 的數字 \n\n正確範例：\n" + command + " gossiping,beauty 100"
	}
	args := re.FindStringSubmatch(text)
	boardNames := splitParamString(args[2])
	inputs := args[3]
	log.WithFields(log.Fields{
		"id":      account,
		"command": command,
		"boards":  boardNames,
		"words":   inputs,
	}).Info("PushSum Command")
	for _, boardName := range boardNames {
		if strings.EqualFold(boardName, "allpost") {
			return "推文數通知不支持 ALLPOST 板。"
		}
	}
	err := update(commandActionMap[command], account, boardNames, inputs)
	if msg, ok := checkBoardError(err); ok {
		return msg
	}
	if err != nil {
		log.WithError(err).Error("PushSum Command Failed")
		return command + "失敗，請嘗試封鎖再解封鎖，並重新執行註冊步驟。\n若問題未解決，請至粉絲團或 LINE 首頁留言。"
	}
	return command + "成功"
}

func countUserArticles(account string) (cnt int) {
	u := new(user.User).Find(account)
	for _, sub := range u.Subscribes {
		cnt += len(sub.Articles)
	}
	return cnt
}

func checkArticleExist(boardName, articleCode string) bool {
	a := new(article.Article)
	a.Code = articleCode
	if bl, _ := a.Exist(); bl {
		return true
	}
	if crawler.CheckArticleExist(boardName, articleCode) {
		a.Board = boardName
		initialArticle(*a)
		return true
	}
	return false
}

func initialArticle(a article.Article) error {
	a, err := crawler.BuildArticle(a.Board, a.Code)
	if err != nil {
		return err
	}
	err = a.Save()
	return err
}

func checkBoardError(err error) (string, bool) {
	if bErr, ok := err.(boardproto.BoardNotExistError); ok {
		return "板名錯誤，請確認拼字。可能板名：\n" + bErr.Suggestion, true
	}
	return "", false
}

func inputErrorTips() string {
	return "指令格式錯誤。\n1. 需以空白分隔動作、板名、參數\n2. 板名欄位開頭與結尾不可有逗號\n3. 板名欄位間不允許空白字元。"
}

func checkRegexp(input string) bool {
	pattern := strings.Replace(strings.TrimPrefix(input, "regexp:"), "//", "////", -1)
	_, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return true
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

func update(action updateAction, account string, boardNames []string, inputs ...string) error {
	u := new(user.User).Find(account)
	if boardNames[0] == "**" {
		boardNames = nil
		for _, uSub := range u.Subscribes {
			boardNames = append(boardNames, uSub.Board)
		}
	}
	for _, boardName := range boardNames {
		sub := subscription.Subscription{
			Board: boardName,
		}
		err := action(&u, sub, inputs...)
		if err != nil {
			return err
		}
		err = u.Update()
		if err != nil {
			log.WithError(err).Error("Subscription Update Error")
			return err
		}
	}
	return nil
}

func HandleLineFollow(id string) error {
	u := new(user.User).Find(id)
	u.Profile.Line = id
	log.WithFields(log.Fields{
		"id":       id,
		"platform": "line",
	}).Info("User Join")
	return handleFollow(u)
}

func HandleMessengerFollow(id string) error {
	u := new(user.User).Find(id)
	u.Profile.Messenger = id
	log.WithFields(log.Fields{
		"id":       id,
		"platform": "messenger",
	}).Info("User Join")
	return handleFollow(u)
}

func HandleTelegramFollow(id string, chatID int64) error {
	u := new(user.User).Find(id)
	u.Profile.Telegram = id
	u.Profile.TelegramChat = chatID
	log.WithFields(log.Fields{
		"id":       id,
		"platform": "telegram",
	}).Info("User Join")
	return handleFollow(u)
}

func handleFollow(u user.User) error {
	if u.Profile.Account != "" {
		u.Enable = true
		u.Update()
	} else {
		if u.Profile.Messenger != "" {
			u.Profile.Account = u.Profile.Messenger
		}
		if u.Profile.Line != "" {
			u.Profile.Account = u.Profile.Line
		}
		if u.Profile.Telegram != "" {
			u.Profile.Account = u.Profile.Telegram
		}
		u.Enable = true
		err := u.Save()
		if err != nil {
			return err
		}
	}
	return nil
}
