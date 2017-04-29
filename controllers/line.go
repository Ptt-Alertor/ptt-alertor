package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/liam-lai/ptt-alertor/line"
)

func LineCallback(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	line.HandleRequest(r)
}
