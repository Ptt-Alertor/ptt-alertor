package board

import (
	"reflect"

	"github.com/liam-lai/ptt-alertor/crawler"
	"github.com/liam-lai/ptt-alertor/models/ptt/article"
)

type Board struct {
	Name        string
	Articles    []article.Article
	NewArticles []article.Article
}

type BoardAction interface {
	OnlineArticles() []article.Article
	GetArticles() []article.Article
	WithArticles()
	Create() error
}

func (bd Board) OnlineArticles() []article.Article {
	bd.Articles = crawler.BuildArticles(bd.Name)
	return bd.Articles
}

func NewArticles(bd BoardAction) []article.Article {
	savedArticles := bd.GetArticles()
	onlineArticles := bd.OnlineArticles()
	newArticles := make([]article.Article, 0)
	for _, onlineArticle := range onlineArticles {
		for index, savedArticle := range savedArticles {
			if reflect.DeepEqual(onlineArticle, savedArticle) {
				break
			}
			if index == len(savedArticles)-1 {
				newArticles = append(newArticles, onlineArticle)
			}
		}
	}
	return newArticles
}
