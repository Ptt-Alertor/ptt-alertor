package jobs

import (
	"fmt"
	"time"

	"github.com/meifamily/ptt-alertor/crawler"
	"github.com/meifamily/ptt-alertor/models/ptt/article"
	"github.com/meifamily/ptt-alertor/models/pushsum"
)

const stopHour = 48 * time.Hour

type pushSumChecker struct {
	Checker
}

func NewPushSumChecker() *pushSumChecker {
	return &pushSumChecker{}
}

func (h pushSumChecker) Run() {
	pushSumArticles := make(article.Articles, 0)
	boards := pushsum.List()
	for _, board := range boards {
		currentPage, err := crawler.CurrentPage(board)
		if err != nil {
			panic(err)
		}

	Page:
		for i := currentPage - 1; ; i-- {
			articles, _ := crawler.BuildArticles(board, i)
			for _, a := range articles {
				if a.ID == 0 {
					continue
				}
				t, err := time.Parse("1/02", a.Date)
				now := time.Now()
				t = t.AddDate(now.Year(), 0, 0)
				if err != nil {
					panic(err)
				}
				if t.Before(now.Truncate(stopHour)) {
					break Page
				}
				if a.PushCount > 10 {
					pushSumArticles = append(pushSumArticles, a)
				}
			}
		}

		for _, a := range pushSumArticles {
			fmt.Println(a)
		}
	}
}
