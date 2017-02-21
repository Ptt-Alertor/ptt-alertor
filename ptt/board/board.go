package board

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"reflect"

	"github.com/liam-lai/ptt-alertor/crawler"
	"github.com/liam-lai/ptt-alertor/myutil"
	"github.com/liam-lai/ptt-alertor/ptt/article"
)

type Boards []*Board

type Board struct {
	Name        string
	Articles    article.Articles
	NewArticles article.Articles
}

var articlesDir string = myutil.StoragePath() + "/articles/"

func (b Board) IndexJSON() []byte {
	articles := b.Index()
	articlesJSON, err := json.Marshal(articles)
	if err != nil {
		log.Fatal(err)
	}
	return articlesJSON
}

func (b Board) Index() article.Articles {
	b.Articles = crawler.BuildArticles(b.Name)
	return b.Articles
}

func (bs Boards) All() Boards {

	files, _ := ioutil.ReadDir(articlesDir)
	for _, file := range files {
		name, ok := myutil.JsonFile(file)
		if !ok {
			continue
		}
		b := new(Board)
		b.Name = name
		articles, err := ioutil.ReadFile(articlesDir + name + ".json")
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal(articles, &b.Articles)
		bs = append(bs, b)
	}
	return bs
}

func (bs Boards) WithNewArticles(clean bool) Boards {
	for _, b := range bs {
		savedArticles := b.Articles
		onlineArticles := b.Index()
		for _, onlineArticle := range onlineArticles {
			for index, savedArticle := range savedArticles {
				if reflect.DeepEqual(onlineArticle, savedArticle) {
					break
				}
				if index == len(savedArticles)-1 {
					b.NewArticles = append(b.NewArticles, onlineArticle)
				}
			}
		}

	}

	if clean {
		bs = bs.deleteNonNewArticleBoard()
	}

	return bs
}

func (bs Boards) deleteNonNewArticleBoard() Boards {

	for index, bd := range bs {
		if len(bd.NewArticles) == 0 {
			if index < len(bs)-1 {
				bs = append(bs[:index], bs[index+1:]...)
			} else {
				bs = bs[:index]
			}
		}
	}
	return bs
}
