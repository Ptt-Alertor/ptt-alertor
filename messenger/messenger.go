package messenger

import (
	"fmt"
	"net/http"

	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/liam-lai/ptt-alertor/command"
	"github.com/liam-lai/ptt-alertor/myutil"
)

const (
	SendAPIURL = "https://graph.facebook.com/v2.6/me/messages?access_token="
	ProfileURL = "https://graph.facebook.com/v2.6/me/messenger_profile?access_token="
)

type Messenger struct {
	VerifyToken string
	AccessToken string
}

func New() Messenger {
	msgConfig := myutil.Config("messenger")
	return Messenger{
		VerifyToken: msgConfig["verifyToken"],
		AccessToken: msgConfig["accessToken"],
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
						responseText := command.HandleCommand(text, id)
						m.SendTextMessage(id, responseText)
					}

				} else if messaging.Postback != nil {
					if messaging.Postback.Payload == "GET_STARTED_PAYLOAD" {
						fmt.Println(id)
						err := command.HandleMessengerFollow(id)
						if err != nil {
							log.WithError(err).Error("Messenger Follow Error")
						}
						m.SendTextMessage(id, "歡迎使用 PTT Alertor\n輸入「指令」查看相關功能。")
					}
				}
			}
		}
	}
}

func (m *Messenger) SendTextMessage(id string, message string) {
	body := Request{
		Recipient{id},
		Message{Text: message},
	}
	m.callSendAPI(body)
}

func (m *Messenger) callSendAPI(body Request) {
	url := SendAPIURL + m.AccessToken
	callAPI(url, body)
}
