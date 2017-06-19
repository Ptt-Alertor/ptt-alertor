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
	Name        string
	Articles    article.Articles
	NewArticles article.Articles
}

type BoardAction interface {
	OnlineArticles() article.Articles
	GetArticles() article.Articles
	WithArticles()
	Create() error
}

func (bd Board) OnlineArticles() article.Articles {
	bd.Articles = rss.BuildArticles(bd.Name)
	return bd.Articles
}

func NewArticles(bd BoardAction) (newArticles, onlineArticles article.Articles) {
	savedArticles := bd.GetArticles()
	onlineArticles = bd.OnlineArticles()
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
