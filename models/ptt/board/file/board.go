package file

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/meifamily/logrus"

	"github.com/liam-lai/ptt-alertor/models/ptt/article"
	"github.com/liam-lai/ptt-alertor/models/ptt/board"
	"github.com/liam-lai/ptt-alertor/myutil"
)

type Board struct {
	board.Board
}

var articlesDir string = myutil.StoragePath() + "/articles/"

func (bd Board) All() []*Board {
	files, err := ioutil.ReadDir(articlesDir)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
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

func (bd Board) GetArticles() article.Articles {
	file := articlesDir + bd.Name + ".json"
	articlesJSON, err := ioutil.ReadFile(file)
	if err != nil {
		log.WithFields(log.Fields{
			"file":    file,
			"runtime": myutil.BasicRuntimeInfo(),
		}).WithError(err).Error("Read File Error")
	}
	articles := make(article.Articles, 0)
	err = json.Unmarshal(articlesJSON, &articles)
	if err != nil {
		myutil.LogJSONDecode(err, articlesJSON)
	}
	return articles
}

func (bd *Board) WithArticles() {
	bd.Articles = bd.GetArticles()
}

func (bd *Board) WithNewArticles() {
	bd.NewArticles = board.NewArticles(bd)
}

func (bd Board) Create() error {
	err := ioutil.WriteFile(articlesDir+bd.Name+".json", []byte("[]"), 664)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}

func (bd Board) Save() error {
	articlesJSON, err := json.Marshal(bd.Articles)
	if err != nil {
		myutil.LogJSONEncode(err, bd.Articles)
	}
	err = ioutil.WriteFile(articlesDir+bd.Name+".json", articlesJSON, 0644)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}
