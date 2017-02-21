package article

import (
	"strings"
)

type Articles []Article

type Article struct {
	Title  string
	Link   string
	Date   string
	Author string
}

func (as Articles) ContainKeyword(keyword string) Articles {
	keyworkAs := make(Articles, 0)
	for _, article := range as {
		if article.ContainKeyword(keyword) {
			keyworkAs = append(keyworkAs, article)
		}
	}
	return keyworkAs
}

func (a Article) ContainKeyword(keyword string) bool {
	return strings.Contains(a.Title, keyword)
}
