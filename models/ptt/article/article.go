package article

import "strings"
import "regexp"

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
	if strings.HasPrefix(keyword, "regexp:") {
		return matchRegex(a.Title, keyword)
	}
	return matchKeyword(a.Title, keyword)
}

func matchRegex(title string, regex string) bool {
	pattern := strings.TrimPrefix(regex, "regexp:")
	b, err := regexp.MatchString(pattern, title)
	if err != nil {
		return false
	}
	return b
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

func (a Article) String() string {
	return a.Title + "\r\n" + a.Link
}

type Articles []Article

func (as Articles) String() string {
	var content string
	for _, a := range as {
		content += "\r\n\r\n" + a.String()
	}
	return content
}
