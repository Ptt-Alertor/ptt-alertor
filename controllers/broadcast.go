package controllers

import (
	"fmt"
	"net/http"

	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/liam-lai/ptt-alertor/line"
	"github.com/liam-lai/ptt-alertor/messenger"
	user "github.com/liam-lai/ptt-alertor/models/user/redis"
)

type platformFunc map[string]func([]*user.User, string)

var ptFunc = platformFunc{
	"line":      broadcastLine,
	"messenger": broadcastMessenger,
}

func Broadcast(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type requestBody struct {
		Platforms []string `json:"platforms"`
		Content   string   `json:"content"`
	}
	body := requestBody{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.WithError(err).Error("Decode Notify Body Failed")
	}
	users := new(user.User).All()
	for _, platform := range body.Platforms {
		if f, ok := ptFunc[platform]; ok {
			go f(users, body.Content)
		} else {
			http.Error(w, "platform "+platform+" is not valid", http.StatusBadRequest)
			return
		}
	}
	fmt.Fprintln(w, "OK")
}

func broadcastLine(users []*user.User, text string) {
	for _, user := range users {
		if user.Profile.LineAccessToken != "" {
			go line.Notify(user.Profile.LineAccessToken, text)
		}
	}
}

func broadcastMessenger(users []*user.User, text string) {
	m := messenger.New()
	for _, user := range users {
		if user.Profile.Messenger != "" {
			go m.SendTextMessage(user.Profile.Messenger, text)
		}
	}
}
