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

func (a Article) MatchKeyword(keyword string) bool {
	if strings.Contains(keyword, "&") {
		keywords := strings.Split(keyword, "&")
		for _, keyword := range keywords {
			if !matchKeyword(a.Title, keyword) {
				return false
			}
		}
		return true
	}
	return matchKeyword(a.Title, keyword)
}

func matchKeyword(title string, keyword string) bool {
	if strings.HasPrefix(keyword, "!") {
		excludeKeyword := strings.Trim(keyword, "!")
		return !containKeyword(title, excludeKeyword)
	}
	return containKeyword(title, keyword)
}

func containKeyword(title string, keyword string) bool {
	return strings.Contains(strings.ToLower(title), strings.ToLower(keyword))
}

type Articles []Article

func (as Articles) String() string {
	var content string
	for _, article := range as {
		link := "https://www.ptt.cc" + article.Link
		content += "\r\n" + article.Title + "\r\n"
		if article.Author != "" {
			content += "作者: " + article.Author + "\r\n"
		}
		content += link + "\r\n"
	}
	return content
}
