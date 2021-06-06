package jobs

import (
	"sort"
	"strconv"

	log "github.com/meifamily/logrus"

	"strings"

	"github.com/meifamily/ptt-alertor/models"
	"github.com/meifamily/ptt-alertor/models/top"
)

type Top struct{}

func NewTop() *Top {
	return &Top{}
}

func (t Top) Run() {
	log.Info("Top List Generated")
	keywordMap := make(map[top.BoardWord]int)
	authorMap := make(map[top.BoardWord]int)
	pushSumMap := make(map[top.BoardWord]int)
	for _, u := range models.User().All() {
		for _, sub := range u.Subscribes {
			for _, keyword := range sub.Keywords {
				keyword = strings.ToLower(keyword)
				keywordMap[top.BoardWord{Board: sub.Board, Word: keyword}]++
			}
			for _, author := range sub.Authors {
				author = strings.ToLower(author)
				authorMap[top.BoardWord{Board: sub.Board, Word: author}]++
			}
			if sub.PushSum.Up != 0 {
				pushSumMap[top.BoardWord{Board: sub.Board, Word: strconv.Itoa(sub.PushSum.Up)}]++
			}
			if sub.PushSum.Down != 0 {
				pushSumMap[top.BoardWord{Board: sub.Board, Word: strconv.Itoa(sub.PushSum.Down * -1)}]++
			}
		}
	}
	topKeywords := rank(keywordMap)
	topKeywords.SaveKeywords()
	topAuthors := rank(authorMap)
	topAuthors.SaveAuthors()
	topPushSum := rank(pushSumMap)
	topPushSum.SavePushSum()
}

func rank(m map[top.BoardWord]int) (orderSlice top.WordOrders) {
	for key, count := range m {
		k := top.WordOrder{BoardWord: key, Count: count}
		orderSlice = append(orderSlice, k)
	}
	sort.Slice(orderSlice, func(i, j int) bool {
		return orderSlice[i].Count > orderSlice[j].Count
	})
	return orderSlice
}
