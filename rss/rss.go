package rss

import (
	"net/http"
	"time"

	"github.com/meifamily/ptt-alertor/models/article"
	"github.com/mmcdole/gofeed"
)

func BuildArticles(board string) (articles article.Articles, err error) {
	feed, err := parseURL("https://www.ptt.cc/atom/" + board + ".xml")
	if err != nil {
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
	resp, err := client.Get(feedURL)
	if err != nil {
		return nil, err
	}

	if resp != nil {
		defer func() {
			ce := resp.Body.Close()
			if ce != nil {
				err = ce
			}
		}()
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, gofeed.HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
		}
	}

	return fp.Parse(resp.Body)
}
