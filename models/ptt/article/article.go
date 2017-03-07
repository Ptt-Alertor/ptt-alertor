package article

import "strings"

type Article struct {
	Title  string
	Link   string
	Date   string
	Author string
}

type ArticleAction interface {
	ContainKeyword(keyword string) bool
}

func (a Article) ContainKeyword(keyword string) bool {
	return strings.Contains(a.Title, keyword)
}
