package controllers

import (
	"fmt"
	"net/http"

	"github.com/Ptt-Alertor/ptt-alertor/models/pushsum"
	"github.com/julienschmidt/httprouter"
)

func PushSumBoards(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	boards := pushsum.List()
	fmt.Fprintf(w, "推噓文數看板總數：%d\n", len(boards))
	for _, board := range boards {
		fmt.Fprintf(w, "%s\n", board)
	}
}
