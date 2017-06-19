package board

import (
	"github.com/liam-lai/ptt-alertor/models/ptt/article"
	"github.com/liam-lai/ptt-alertor/rss"
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

func (bd Board) FetchArticles() article.Articles {
	return rss.BuildArticles(bd.Name)
}

func NewArticles(bd BoardAction) (newArticles, onlineArticles article.Articles) {
	savedArticles := bd.GetArticles()
	onlineArticles = bd.FetchArticles()
	if len(savedArticles) == 0 {
		return onlineArticles, onlineArticles
	}
	for _, onlineArticle := range onlineArticles {
		for index, savedArticle := range savedArticles {
			if onlineArticle.Link == savedArticle.Link {
				break
			}
			if index == len(savedArticles)-1 {
				newArticles = append(newArticles, onlineArticle)
			}
		}
	}
	return newArticles, onlineArticles
}
