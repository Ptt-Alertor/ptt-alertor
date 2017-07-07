package board

import (
	log "github.com/meifamily/logrus"
	"github.com/meifamily/ptt-alertor/crawler"
	"github.com/meifamily/ptt-alertor/models/ptt/article"
	"github.com/meifamily/ptt-alertor/rss"
)

type BoardNotExistError struct {
	Suggestion string
}

func (e BoardNotExistError) Error() string {
	return "board is not exist"
}

type Board struct {
	Name           string
	Articles       article.Articles
	OnlineArticles article.Articles
	NewArticles    article.Articles
}

type BoardAction interface {
	FetchArticles() article.Articles
	GetArticles() article.Articles
	WithArticles()
	Create() error
}

func (bd Board) FetchArticles() (articles article.Articles) {
	articles, err := rss.BuildArticles(bd.Name)
	if err != nil {
		log.WithField("board", bd.Name).WithError(err).Error("RSS Parse Failed, Switch to HTML Crawler")
		articles, err = crawler.BuildArticles(bd.Name, -1)
		if err != nil {
			log.WithField("board", bd.Name).WithError(err).Error("HTML Parse Failed")
		}
	}
	return articles
}

func NewArticles(bd BoardAction) (newArticles, onlineArticles article.Articles) {
	newArticles = make(article.Articles, 0)
	savedArticles := bd.GetArticles()
	onlineArticles = bd.FetchArticles()
	if len(savedArticles) == 0 {
		return nil, onlineArticles
	}
	for _, onlineArticle := range onlineArticles {
		for index, savedArticle := range savedArticles {
			if onlineArticle.ID <= savedArticle.ID {
				break
			}
			if index == len(savedArticles)-1 {
				newArticles = append(newArticles, onlineArticle)
			}
		}
	}
	return newArticles, onlineArticles
}
