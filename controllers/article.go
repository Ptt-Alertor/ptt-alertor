package controllers

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/meifamily/ptt-alertor/models/ptt/article"
)

// ArticleIndex show all subscribed article codes
func ArticleIndex(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	codes := new(article.Articles).List()
	for _, code := range codes {
		fmt.Fprintf(w, "%s\n", code)
	}
}
