package subscription

import (
	"sort"
	"strings"

	"github.com/liam-lai/ptt-alertor/myutil"
)

type Subscription struct {
	Board    string             `json:"board"`
	Keywords myutil.StringSlice `json:"keywords"`
	Authors  myutil.StringSlice `json:"authors"`
	Articles myutil.StringSlice `json:"articles"`
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

func (s Subscription) StringArticle() string {
	if len(s.Articles) == 0 {
		return ""
	}
	sort.Strings(s.Articles)
	return s.Board + ": " + strings.Join(s.Articles, ", ")
}

func (s *Subscription) CleanUp() {
	s.Keywords.Clean()
	s.Authors.Clean()
	s.Authors.RemoveStringsSpace()
}

func (s *Subscription) DeleteKeywords(keywords myutil.StringSlice) {
	s.Keywords.DeleteSlice(keywords, false)
}

func (s *Subscription) DeleteAuthors(authors myutil.StringSlice) {
	s.Authors.DeleteSlice(authors, false)
}

func (s *Subscription) DeleteArticles(articles myutil.StringSlice) {
	s.Articles.DeleteSlice(articles, false)
}
