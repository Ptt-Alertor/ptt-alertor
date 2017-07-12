package rss

import (
	"github.com/meifamily/ptt-alertor/models/ptt/article"
	"github.com/mmcdole/gofeed"
)

func BuildArticles(board string) (articles article.Articles, err error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://www.ptt.cc/atom/" + board + ".xml")
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
