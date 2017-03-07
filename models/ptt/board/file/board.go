package file

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"reflect"

	"github.com/liam-lai/ptt-alertor/crawler"
	"github.com/liam-lai/ptt-alertor/models/ptt/article"
	"github.com/liam-lai/ptt-alertor/models/ptt/board"
	"github.com/liam-lai/ptt-alertor/myutil"
)

type Board struct {
	board.Board
	Articles    []article.Article
	NewArticles []article.Article
}

var articlesDir string = myutil.StoragePath() + "/articles/"

func (bd Board) OnlineArticles() []article.Article {
	bd.Articles = crawler.BuildArticles(bd.Name)
	return bd.Articles
}

func (bd Board) All() []*Board {
	files, _ := ioutil.ReadDir(articlesDir)
	bds := make([]*Board, 0)
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

func (bd *Board) WithArticles() {
	articles, err := ioutil.ReadFile(articlesDir + bd.Name + ".json")
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(articles, &bd.Articles)
}

func (bd *Board) WithNewArticles() {
	bd.WithArticles()
	savedArticles := bd.Articles
	onlineArticles := bd.OnlineArticles()
	for _, onlineArticle := range onlineArticles {
		for index, savedArticle := range savedArticles {
			if reflect.DeepEqual(onlineArticle, savedArticle) {
				break
			}
			if index == len(savedArticles)-1 {
				bd.NewArticles = append(bd.NewArticles, onlineArticle)
			}
		}
	}
}

func (bd Board) Create() error {
	err := ioutil.WriteFile(articlesDir+bd.Name+".json", []byte("[]"), 664)
	return err
}
