package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	board "github.com/liam-lai/ptt-alertor/models/ptt/board/redis"
	"github.com/liam-lai/ptt-alertor/myutil"
)

func BoardArticleIndex(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	bd := new(board.Board)
	bd.Name = strings.ToUpper(params.ByName("boardName"))
	articles := bd.OnlineArticles()
	articlesJSON, err := json.Marshal(articles)
	if err != nil {
		myutil.LogJSONEncode(err, articles)
	}
	fmt.Fprintf(w, "%s", articlesJSON)
}

func BoardIndex(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	bds := new(board.Board).All()
	fmt.Fprintf(w, "追蹤看板總數：%d", len(bds))
	for _, bd := range bds {
		fmt.Fprintf(w, "\n%s", bd.Name)
	}
}
