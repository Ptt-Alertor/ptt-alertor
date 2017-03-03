package board

import "github.com/liam-lai/ptt-alertor/ptt/article"

type Board interface {
	OnlineArticles() []article.Article
	All() []*Board
	WithArticles()
	WithNewArticles()
	Create() error
}
