package controllers

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/julienschmidt/httprouter"
	"github.com/meifamily/ptt-alertor/models"
	"github.com/meifamily/ptt-alertor/models/author"
)

func AuthorBoards(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	type authorCount struct {
		board string
		count int
	}
	authorCounts := make([]authorCount, 0)
	boards := models.Board.List()
	for _, name := range boards {
		cnt := len(author.Subscribers(name))
		ac := authorCount{board: name, count: cnt}
		authorCounts = append(authorCounts, ac)
	}
	sort.Slice(authorCounts, func(i, j int) bool {
		return authorCounts[i].count > authorCounts[j].count
	})
	for _, ac := range authorCounts {
		fmt.Fprintf(w, "%s: %d\n", ac.board, ac.count)
	}
}
