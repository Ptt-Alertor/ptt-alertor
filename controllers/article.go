package controllers

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/meifamily/ptt-alertor/models/article"
)

// ArticleIndex show all subscribed article codes
func ArticleIndex(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	codes := new(article.Articles).List()
	fmt.Fprintf(w, "推文追蹤文章總數：%d\n", len(codes))
	for _, code := range codes {
		fmt.Fprintf(w, "%s\n", code)
	}
}
