package jobs

import (
	"sort"

	log "github.com/meifamily/logrus"

	"github.com/liam-lai/ptt-alertor/models/top"
	user "github.com/liam-lai/ptt-alertor/models/user/redis"
)

type Top struct{}

func NewTop() *Top {
	return &Top{}
}

func (t Top) Run() {
	log.Info("Top List Generated")
	keywordMap := make(map[top.BoardWord]int)
	authorMap := make(map[top.BoardWord]int)
	us := new(user.User).All()
	for _, u := range us {
		for _, sub := range u.Subscribes {
			for _, keyword := range sub.Keywords {
				keywordMap[top.BoardWord{sub.Board, keyword}]++
			}
			for _, author := range sub.Authors {
				authorMap[top.BoardWord{sub.Board, author}]++
			}
		}
	}
	topKeywords := rank(keywordMap)
	topKeywords.SaveKeywords()
	topAuthors := rank(authorMap)
	topAuthors.SaveAuthors()
}

func rank(m map[top.BoardWord]int) (orderSlice top.WordOrders) {
	for key, count := range m {
		k := top.WordOrder{key, count}
		orderSlice = append(orderSlice, k)
	}
	sort.Slice(orderSlice, func(i, j int) bool {
		return orderSlice[i].Count > orderSlice[j].Count
	})
	return orderSlice
}
