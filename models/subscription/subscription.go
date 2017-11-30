package subscription

import (
	"fmt"
	"sort"
	"strings"

	"github.com/meifamily/ptt-alertor/myutil"
)

var EmptyPushSum = PushSum{}

// Subscription is a struct show User Subscription
type Subscription struct {
	Board    string             `json:"board"`
	Keywords myutil.StringSlice `json:"keywords"`
	Authors  myutil.StringSlice `json:"authors"`
	Articles myutil.StringSlice `json:"articles"`
	PushSum  `json:"pushSum"`
}

type PushSum struct {
	Up   int `json:"up"`
	Down int `json:"down"`
}

func (s Subscription) String() string {
	if len(s.Keywords) == 0 {
		return ""
	}
	sort.Strings(s.Keywords)
	return s.Board + ": " + strings.Join(s.Keywords, ", ")
}

func (s Subscription) StringAuthor() string {
	if len(s.Authors) == 0 {
		return ""
	}
	sort.Strings(s.Authors)
	return s.Board + ": " + strings.Join(s.Authors, ", ")
}

func (s Subscription) StringPushSum() string {
	emptyPushSum := PushSum{}
	if s.PushSum == emptyPushSum {
		return ""
	}
	return fmt.Sprintf("%s: 推 %d; 噓 %d", s.Board, s.PushSum.Up, s.PushSum.Down)
}

func (s Subscription) StringArticle() string {
	if len(s.Articles) == 0 {
		return ""
	}
	sort.Strings(s.Articles)
	aURLs := make([]string, 0)
	for _, a := range s.Articles {
		aURLs = append(aURLs, buildArticleURL(s.Board, a))
	}
	return s.Board + ":\n" + strings.Join(aURLs, "\n")
}

func buildArticleURL(board, code string) string {
	return fmt.Sprintf("https://www.ptt.cc/bbs/%s/%s.html", board, code)
}

func (s *Subscription) CleanUp() {
	s.Keywords.Clean()
	s.Authors.Clean()
	s.Authors.RemoveStringsSpace()
}

func (s *Subscription) DeleteKeywords(keywords myutil.StringSlice) {
	s.Keywords.Delete(keywords, false)
}

func (s *Subscription) DeleteAuthors(authors myutil.StringSlice) {
	s.Authors.Delete(authors, false)
}

func (s *Subscription) DeleteArticles(articles myutil.StringSlice) {
	s.Articles.Delete(articles, false)
}
