package rss

import (
	"errors"
	"net/http"
	"time"

	"github.com/Ptt-Alertor/ptt-alertor/models/article"
	pttHttp "github.com/Ptt-Alertor/ptt-alertor/ptt/http"
	"github.com/mmcdole/gofeed"
)

var ErrTooManyRequests = errors.New("Too Many Requests")

// CheckBoardExist use for checking board exist or not
func CheckBoardExist(board string) bool {
	feed, err := parseURL("https://www.ptt.cc/atom/" + board + ".xml")
	if err != nil {
		return false
	}
	// category didn't have any items but has rss link as well
	if len(feed.Items) == 0 {
		return false
	}
	return true
}

func BuildArticles(board string) (articles article.Articles, err error) {
	feed, err := parseURL("https://www.ptt.cc/atom/" + board + ".xml")
	if err != nil {
		if herr, ok := err.(gofeed.HTTPError); ok && herr.StatusCode == http.StatusTooManyRequests {
			return nil, ErrTooManyRequests
		}
		return nil, err
	}
	for _, item := range feed.Items {
		article := article.Article{
			Title:  item.Title,
			Link:   item.GUID,
			Date:   item.Published,
			Author: item.Author.Name,
		}
		article.ID = article.ParseID(item.GUID)
		articles = append(articles, article)
	}
	return articles, nil
}

var fp = gofeed.NewParser()
var client = http.Client{
	Timeout: 30 * time.Second,
}

func parseURL(feedURL string) (feed *gofeed.Feed, err error) {
	req, err := pttHttp.HttpRequest(feedURL)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, gofeed.HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
		}
	}

	return fp.Parse(resp.Body)
}
