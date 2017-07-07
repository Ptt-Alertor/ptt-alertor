package jobs

import (
	"fmt"
	"time"

	"github.com/meifamily/ptt-alertor/crawler"
	"github.com/meifamily/ptt-alertor/models/ptt/article"
)

type hotChecker struct {
}

func NewHotChecker() *hotChecker {
	return &hotChecker{}
}

func (h hotChecker) Run() {
	hotArticles := make(article.Articles, 0)
Page:
	for i := 25475; i > 25270; i-- {
		articles, _ := crawler.BuildArticles("gossiping", i)
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
			if t.Before(now.Truncate(24 * time.Hour)) {
				break Page
			}
			if a.PushCount == 100 {
				hotArticles = append(hotArticles, a)
			}
		}
	}

	for _, a := range hotArticles {
		fmt.Println(a)
	}

}
