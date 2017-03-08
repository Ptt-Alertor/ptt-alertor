package redis

import (
	"encoding/json"

	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/liam-lai/ptt-alertor/connections"
	"github.com/liam-lai/ptt-alertor/models/ptt/article"
	"github.com/liam-lai/ptt-alertor/models/ptt/board"
)

const prefix string = "board:"

type Board struct {
	board.Board
}

func (bd Board) All() []*Board {
	conn := connections.Redis()
	defer conn.Close()
	boards, err := redis.Strings(conn.Do("SMEMBERS", "boards"))
	if err != nil {
		log.Fatal(err)
	}
	bds := make([]*Board, 0)
	for _, board := range boards {
		bd := new(Board)
		bd.Name = board
		bds = append(bds, bd)
	}
	return bds
}

func (bd Board) GetArticles() []article.Article {
	conn := connections.Redis()
	defer conn.Close()

	key := prefix + bd.Name
	articlesJSON, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		log.Fatal(err)
	}

	articles := make([]article.Article, 0)
	if articlesJSON != nil {
		json.Unmarshal(articlesJSON, &articles)
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
	return err
}

func (bd Board) Save() error {
	conn := connections.Redis()
	defer conn.Close()

	articlesJSON, err := json.Marshal(bd.Articles)
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Do("SET", prefix+bd.Name, articlesJSON)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
