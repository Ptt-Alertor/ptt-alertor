package board

import "github.com/liam-lai/ptt-alertor/models/ptt/article"

type Board struct {
	Name        string
	Articles    []article.Article
	NewArticles []article.Article
}

type BoardAction interface {
	OnlineArticles() []article.Article
	All() []*Board
	WithArticles()
	WithNewArticles()
	Create() error
}
