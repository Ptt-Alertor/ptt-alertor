package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/meifamily/ptt-alertor/models"
	"github.com/meifamily/ptt-alertor/myutil"
)

func BoardArticleIndex(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	bd := models.Board
	bd.Name = strings.ToUpper(params.ByName("boardName"))
	articles := bd.FetchArticles()
	articlesJSON, err := json.Marshal(articles)
	if err != nil {
		myutil.LogJSONEncode(err, articles)
	}
	fmt.Fprintf(w, "%s", articlesJSON)
}

func BoardArticle(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	code := params.ByName("code")
	a := models.Article.Find(code)
	aJSON, err := json.Marshal(a)
	if err != nil {
		myutil.LogJSONEncode(err, a)
	}
	fmt.Fprintf(w, "%s", aJSON)
}

func BoardIndex(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	bds := models.Board.All()
	fmt.Fprintf(w, "追蹤看板總數：%d", len(bds))
	for _, bd := range bds {
		fmt.Fprintf(w, "\n%s", bd.Name)
	}
}
