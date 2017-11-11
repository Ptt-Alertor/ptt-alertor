package messenger

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"encoding/json"

	"github.com/julienschmidt/httprouter"
	log "github.com/meifamily/logrus"
	"github.com/meifamily/ptt-alertor/command"
	"github.com/meifamily/ptt-alertor/myutil"
)

const (
	SendAPIURL    = "https://graph.facebook.com/v2.6/me/messages?access_token="
	ProfileURL    = "https://graph.facebook.com/v2.6/me/messenger_profile?access_token="
	maxCharacters = 640
)

type Messenger struct {
	VerifyToken string
	AccessToken string
}

var config map[string]string

func init() {
	config = myutil.Config("messenger")
}

func New() Messenger {
	return Messenger{
		VerifyToken: config["verifyToken"],
		AccessToken: config["accessToken"],
	}
}

func (m *Messenger) Verify(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.FormValue("hub.mode") == "subscribe" && r.FormValue("hub.verify_token") == m.VerifyToken {
		log.Info("Validating webhook")
		resStr := r.FormValue("hub.challenge")
		fmt.Fprintln(w, resStr)
	} else {
		log.Info("Failed validation. Make sure the validation tokens match.")
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	}
}

func (m *Messenger) Received(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := Webhook{}
	json.NewDecoder(r.Body).Decode(&data)
	if data.Object == "page" {
		for _, entry := range data.Entry {
			for _, messaging := range entry.Messaging {
				id := messaging.Sender.ID
				if messaging.Message != nil {
					text := messaging.Message.Text
					if text != "" {
						if match, _ := regexp.MatchString("^(刪除|刪除作者)+\\s.*\\*+", text); match {
							m.SendConfirmation(id, text)
							return
						}
						responseText := command.HandleCommand(text, id)
						m.SendTextMessage(id, responseText)
					}
				} else if messaging.Postback != nil {
					payload := messaging.Postback.Payload
					log.WithField("payload", payload).Info("Messenger Postback")
					m.handlePostback(id, payload)
				}
			}
		}
	}
}

func (m *Messenger) handlePostback(id string, payload string) {
	var responseText string
	switch payload {
	case "GET_STARTED_PAYLOAD":
		err := command.HandleMessengerFollow(id)
		if err != nil {
			log.WithError(err).Error("Messenger Follow Error")
		}
		responseText = "歡迎使用 Ptt Alertor\n輸入「指令」查看相關功能。\n\n觀看Demo:\nhttps://media.giphy.com/media/3ohzdF6vidM6I49lQs/giphy.gif"
	case "COMMANDS_PAYLOAD":
		// responseText = command.HandleCommand("指令", id)
		var str string
		commands := make(map[string]string)
		for cat, cmds := range command.Commands {
			if strings.EqualFold(cat, "進階應用") || strings.EqualFold(cat, "一般") {
				continue
			}
			for cmd, doc := range cmds {
				str += cmd
				if doc != "" {
					str += "：" + doc
				}
				str += "\n"
			}
			commands[cat] = str
			str = ""
		}
		m.SendListMessage(id, commands)
	case "SUBSCRIPTIONS_PAYLOAD":
		responseText = command.HandleCommand("清單", id)
	case "TOP_PAYLOAD":
		responseText = command.HandleCommand("排行", id)
	case "PUSHLIST_PAYLOAD":
		responseText = command.HandleCommand("推文清單", id)
	case "CANCEL":
		responseText = "取消"
	default:
		responseText = command.HandleCommand(payload, id)
	}
	m.SendTextMessage(id, responseText)
}

func (m *Messenger) SendTextMessage(id string, message string) {
	for _, msg := range myutil.SplitTextByLineBreak(message, maxCharacters) {
		body := Request{
			Recipient{id},
			Message{Text: msg},
		}
		m.callSendAPI(body)
	}
}

func (m *Messenger) SendConfirmation(id string, cmd string) {
	attachment := Attachment{
		Type: "template",
		Payload: ButtonPayload{
			TemplateType: "button",
			Text:         "確認" + cmd,
			Buttons: Buttons{
				Button{"postback", "是", cmd},
				Button{"postback", "否", "CANCEL"},
			},
		},
	}
	body := Request{}
	body.Recipient.ID = id
	body.Message.Attachment = &attachment
	m.callSendAPI(body)
}

func (m *Messenger) SendQuickReplies(id string, payload string) {
	qrs := QuickReplies{
		QuickReply{"text", "是", payload},
		QuickReply{"text", "否", "CANCEL"},
	}
	body := Request{
		Recipient{id},
		Message{Text: "確定" + payload, QuickReplies: &qrs},
	}
	m.callSendAPI(body)
}

func (m *Messenger) SendListMessage(id string, StringMap map[string]string) {
	elements := []Element{}
	for key, str := range StringMap {
		elements = append(elements, Element{
			Title:    key,
			Subtitle: str,
		})
	}
	attachment := Attachment{
		Type: "template",
		Payload: ListPayload{
			TemplateType:    "list",
			TopElementStyle: "compact",
			Elements:        elements,
		},
	}
	body := Request{}
	body.Recipient.ID = id
	body.Message.Attachment = &attachment
	m.callSendAPI(body)
}

func (m *Messenger) callSendAPI(body Request) {
	url := SendAPIURL + m.AccessToken
	err := callAPI(url, body)
	if err != nil {
		log.WithError(err).Error("Call Send API Error")
	}
}
