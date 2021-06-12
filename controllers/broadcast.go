package controllers

import (
	"fmt"
	"net/http"

	"encoding/json"

	log "github.com/Ptt-Alertor/logrus"
	"github.com/Ptt-Alertor/ptt-alertor/jobs"
	"github.com/julienschmidt/httprouter"
)

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
	bc := new(jobs.Broadcaster)
	bc.Msg = body.Content
	err = bc.Send(body.Platforms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, "OK")
}
