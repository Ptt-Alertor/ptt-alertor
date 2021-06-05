package telegram

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	log "github.com/meifamily/logrus"

	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/julienschmidt/httprouter"
	"github.com/meifamily/ptt-alertor/command"
	"github.com/meifamily/ptt-alertor/myutil"
)

var (
	bot   *tgbotapi.BotAPI
	err   error
	token = os.Getenv("TELEGRAM_TOKEN")
	host  = os.Getenv("APP_HOST")
)

func init() {
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.WithError(err).Fatal("Telegram Bot Initialize Failed")
	}
	// bot.Debug = true
	log.Info("Telegram Authorized on " + bot.Self.UserName)

	webhookConfig := tgbotapi.NewWebhook(host + "/telegram/" + token)
	webhookConfig.MaxConnections = 100
	_, err = bot.SetWebhook(webhookConfig)
	if err != nil {
		log.WithError(err).Fatal("Telegram Bot Set Webhook Failed")
	}
	log.Info("Telegram Bot Sets Webhook Success")
}

// HandleRequest handles request from webhook
func HandleRequest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithError(err).Error("Telegram Read Request Body Failed")
	}

	var update tgbotapi.Update
	json.Unmarshal(bytes, &update)

	if update.CallbackQuery != nil {
		handleCallbackQuery(update)
		return
	}

	if update.Message != nil {
		if update.Message.IsCommand() {
			handleCommand(update)
			return
		}
		if update.Message.Text != "" {
			handleText(update)
			return
		}
	}
}

func handleCallbackQuery(update tgbotapi.Update) {
	var responseText string
	userID := strconv.Itoa(update.CallbackQuery.From.ID)
	switch update.CallbackQuery.Data {
	case "CANCEL":
		responseText = "取消"
	default:
		responseText = command.HandleCommand(update.CallbackQuery.Data, userID)
	}
	SendTextMessage(update.CallbackQuery.Message.Chat.ID, responseText)
}

// help - 所有指令清單
// list - 設定清單
// ranking - 熱門關鍵字、作者、推文數
// add - 新增看板關鍵字、作者、推文數
// del - 刪除看板關鍵字、作者、推文數
// showkeyboard - 顯示快捷小鍵盤
// hidekeyboard - 隱藏快捷小鍵盤
func handleCommand(update tgbotapi.Update) {
	var responseText string
	userID := strconv.Itoa(update.Message.From.ID)
	chatID := update.Message.Chat.ID

	switch update.Message.Command() {
	case "add", "del":
		text := update.Message.Command() + " " + update.Message.CommandArguments()
		responseText = command.HandleCommand(text, userID)
	case "start":
		command.HandleTelegramFollow(userID, chatID)
		responseText = "歡迎使用 Ptt Alertor\n輸入「指令」查看相關功能。\n\n觀看Demo:\nhttps://media.giphy.com/media/3ohzdF6vidM6I49lQs/giphy.gif"
	case "help":
		responseText = command.HandleCommand("help", userID)
	case "list":
		responseText = command.HandleCommand("list", userID)
	case "ranking":
		responseText = command.HandleCommand("ranking", userID)
	case "showkeyboard":
		showReplyKeyboard(chatID)
		return
	case "hidekeyboard":
		hideReplyKeyboard(chatID)
		return
	default:
		responseText = "I don't know the command"
	}
	SendTextMessage(chatID, responseText)
}

func handleText(update tgbotapi.Update) {
	var responseText string
	userID := strconv.Itoa(update.Message.From.ID)
	chatID := update.Message.Chat.ID
	text := update.Message.Text
	if match, _ := regexp.MatchString("^(刪除|刪除作者)+\\s.*\\*+", text); match {
		sendConfirmation(chatID, text)
		return
	}
	responseText = command.HandleCommand(text, userID)
	SendTextMessage(chatID, responseText)
}

func sendConfirmation(chatID int64, cmd string) {
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("是", cmd),
			tgbotapi.NewInlineKeyboardButtonData("否", "CANCEL"),
		))
	msg := tgbotapi.NewMessage(chatID, "確定"+cmd+"？")
	msg.ReplyMarkup = markup
	_, err := bot.Send(msg)
	if err != nil {
		log.WithError(err).Error("Telegram Send Confirmation Failed")
	}
}

const maxCharacters = 4096

// SendTextMessage sends text message to chatID
func SendTextMessage(chatID int64, text string) {
	for _, msg := range myutil.SplitTextByLineBreak(text, maxCharacters) {
		sendTextMessage(chatID, msg)
	}
}

func sendTextMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.DisableWebPagePreview = true
	_, err := bot.Send(msg)
	if err != nil {
		log.WithError(err).Error("Telegram Send Message Failed")
	}
}

func showReplyKeyboard(chatID int64) {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("清單"),
			tgbotapi.NewKeyboardButton("推文清單"),
			tgbotapi.NewKeyboardButton("排行"),
			tgbotapi.NewKeyboardButton("指令"),
		))
	msg := tgbotapi.NewMessage(chatID, "顯示小鍵盤")
	msg.ReplyMarkup = keyboard
	_, err := bot.Send(msg)
	if err != nil {
		log.WithError(err).Error("Telegram Show Reply Keyboard Failed")
	}
}

func hideReplyKeyboard(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "隱藏小鍵盤")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	_, err := bot.Send(msg)
	if err != nil {
		log.WithError(err).Error("Telegram Hide Reply Keyboard Failed")
	}
}
