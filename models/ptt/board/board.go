package board

import (
	"strings"

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
		articles, err = crawler.BuildArticles(bd.Name)
		if err != nil {
			log.WithField("board", bd.Name).WithError(err).Error("HTML Parse Failed")
		}
	}
	if strings.EqualFold(bd.Name, "allpost") {
		fixLink(&articles)
	}
	return articles
}

func fixLink(articles *article.Articles) {
	for i, a := range *articles {
		preParenthesesIndex := strings.LastIndex(a.Title, "(")
		backParenthesesIndex := strings.LastIndex(a.Title, ")")
		realBoard := a.Title[preParenthesesIndex+1 : backParenthesesIndex]
		a.Link = strings.Replace(a.Link, "ALLPOST", realBoard, -1)
		(*articles)[i] = a
	}
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
