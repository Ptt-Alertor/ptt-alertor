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

func (bds Boards) All() Boards {
	files, _ := ioutil.ReadDir(articlesDir)
	for _, file := range files {
		name, ok := myutil.JsonFile(file)
		if !ok {
			continue
		}
		bd := new(Board)
		bd.Name = name
		bds = append(bds, bd)
	}
	return bds
}

func (bds Boards) WithArticles() Boards {
	for _, bd := range bds {
		articles, err := ioutil.ReadFile(articlesDir + bd.Name + ".json")
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal(articles, &bd.Articles)
	}
	return bds
}

func (bds Boards) WithNewArticles(clean bool) Boards {
	bds = bds.WithArticles()
	for _, b := range bds {
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
		bds = bds.deleteNonNewArticleBoard()
	}

	return bds
}

func (bds Boards) deleteNonNewArticleBoard() Boards {

	for index, bd := range bds {
		if len(bd.NewArticles) == 0 {
			if index < len(bds)-1 {
				bds = append(bds[:index], bds[index+1:]...)
			} else {
				bds = bds[:index]
			}
		}
	}
	return bds
}

func (bd Board) Create() error {
	err := ioutil.WriteFile(articlesDir+bd.Name+".json", []byte("[]"), 664)
	return err
}
