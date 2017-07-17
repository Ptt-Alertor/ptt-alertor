package controllers

import (
	"fmt"
	"net/http"

	"encoding/json"

	"github.com/julienschmidt/httprouter"
	log "github.com/meifamily/logrus"
	"github.com/meifamily/ptt-alertor/jobs"
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
