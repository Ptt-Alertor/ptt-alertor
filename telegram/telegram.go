package telegram

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"

	log "github.com/meifamily/logrus"

	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/julienschmidt/httprouter"
	"github.com/meifamily/ptt-alertor/command"
	"github.com/meifamily/ptt-alertor/myutil"
)

var bot *tgbotapi.BotAPI
var err error

// Token is telegram authentication token
var Token string

func init() {
	Token = myutil.Config("telegram")["token"]
	host := myutil.Config("app")["host"]
	bot, err = tgbotapi.NewBotAPI(Token)
	if err != nil {
		log.WithError(err).Fatal("Telegram Bot Initialize Failed")
	}
	// bot.Debug = true
	log.Info("Telegram Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(host + "/telegram/" + Token))
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

func handleCommand(update tgbotapi.Update) {
	var responseText string
	userID := strconv.Itoa(update.Message.From.ID)
	chatID := update.Message.Chat.ID

	switch update.Message.Command() {
	case "start":
		command.HandleTelegramFollow(userID, chatID)
		responseText = "歡迎使用 Ptt Alertor\n輸入「指令」查看相關功能。\n\n觀看Demo:\nhttps://media.giphy.com/media/3ohzdF6vidM6I49lQs/giphy.gif"
	case "help":
		responseText = command.HandleCommand("指令", userID)
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

// SendTextMessage sends text message to chatID
func SendTextMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.DisableWebPagePreview = true
	_, err := bot.Send(msg)
	if err != nil {
		log.WithError(err).Error("Telegram Send Message Failed")
	}
}
