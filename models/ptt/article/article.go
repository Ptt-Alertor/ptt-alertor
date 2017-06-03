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
	return strings.Contains(strings.ToLower(a.Title), strings.ToLower(keyword))
}

type Articles []Article

func (as Articles) String() string {
	var content string
	for _, article := range as {
		content += article.Title + "\r\n" +
			"https://www.ptt.cc" + article.Link + "\r\n" +
			"\r\n"
	}
	return content
}
