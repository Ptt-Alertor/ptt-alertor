package redis

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"

	"github.com/liam-lai/ptt-alertor/connections"
	"github.com/liam-lai/ptt-alertor/models/ptt/article"
	"github.com/liam-lai/ptt-alertor/models/ptt/board"
	"github.com/liam-lai/ptt-alertor/myutil"
)

const prefix string = "board:"

type Board struct {
	board.Board
}

func (bd Board) Exist() bool {
	names := bd.listName()
	for _, name := range names {
		if bd.Name == name {
			return true
		}
	}
	return false
}

func (bd Board) All() []*Board {
	boards := bd.listName()
	bds := make([]*Board, 0)
	for _, board := range boards {
		bd := new(Board)
		bd.Name = board
		bds = append(bds, bd)
	}
	return bds
}

func (bd Board) listName() []string {
	conn := connections.Redis()
	defer conn.Close()
	boards, err := redis.Strings(conn.Do("SMEMBERS", "boards"))
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return boards
}

func (bd Board) GetArticles() []article.Article {
	conn := connections.Redis()
	defer conn.Close()

	key := prefix + bd.Name
	articlesJSON, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}

	articles := make([]article.Article, 0)
	if articlesJSON != nil {
		err = json.Unmarshal(articlesJSON, &articles)
		if err != nil {
			myutil.LogJSONDecode(err, articlesJSON)
		}
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
	conn := connections.Redis()
	defer conn.Close()
	_, err := conn.Do("SADD", "boards", bd.Name)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}

func (bd Board) Save() error {
	conn := connections.Redis()
	defer conn.Close()

	articlesJSON, err := json.Marshal(bd.Articles)
	if err != nil {
		myutil.LogJSONEncode(err, bd.Articles)
		return err
	}
	_, err = conn.Do("SET", prefix+bd.Name, articlesJSON)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}
