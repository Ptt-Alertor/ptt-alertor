package article

import (
	"strings"
)

type Article struct {
	Title  string
	Link   string
	Date   string
	Author string
}

func (a Article) ContainKeyword(keyword string) bool {
	return strings.Contains(a.Title, keyword)
}
