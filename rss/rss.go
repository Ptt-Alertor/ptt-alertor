package rss

import (
	log "github.com/Sirupsen/logrus"
	"github.com/liam-lai/ptt-alertor/crawler"
	"github.com/liam-lai/ptt-alertor/models/ptt/article"
	"github.com/mmcdole/gofeed"
)

func BuildArticles(board string) article.Articles {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://www.ptt.cc/atom/" + board + ".xml")
	if err != nil {
		log.WithField("board", board).WithError(err).Error("RSS Parse Failed, Switch to HTML Crawler")
		return crawler.BuildArticles(board)
	}
	articles := make(article.Articles, 0)
	for _, item := range feed.Items {
		article := article.Article{
			Title:  item.Title,
			Link:   item.GUID,
			Date:   item.Published,
			Author: item.Author.Name,
		}
		articles = append(articles, article)
	}
	return articles
}
