package controllers

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/meifamily/ptt-alertor/models/pushsum"
)

func PushSumBoards(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	boards := pushsum.List()
	for _, board := range boards {
		fmt.Fprintf(w, "%s\n", board)
	}
}
