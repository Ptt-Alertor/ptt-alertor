package redis

import (
	"encoding/json"
	"math"
	"strings"

	"github.com/garyburd/redigo/redis"
	log "github.com/meifamily/logrus"

	"github.com/meifamily/ptt-alertor/connections"
	"github.com/meifamily/ptt-alertor/crawler"
	"github.com/meifamily/ptt-alertor/models/ptt/article"
	"github.com/meifamily/ptt-alertor/models/ptt/board"
	"github.com/meifamily/ptt-alertor/myutil"
	"github.com/meifamily/ptt-alertor/myutil/maputil"
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

func (bd Board) All() (bds []*Board) {
	boards := bd.listName()
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

func (bd Board) GetArticles() (articles article.Articles) {
	conn := connections.Redis()
	defer conn.Close()

	key := prefix + bd.Name
	articlesJSON, err := redis.Bytes(conn.Do("GET", key))
	if err != nil && err != redis.ErrNil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}

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
	bd.NewArticles, bd.OnlineArticles = board.NewArticles(bd)
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
	conn.Send("WATCH", prefix+bd.Name)
	conn.Send("MULTI")
	conn.Send("SET", prefix+bd.Name, articlesJSON)
	_, err = conn.Do("EXEC")
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}

func (bd Board) Delete() error {
	conn := connections.Redis()
	defer conn.Close()
	_, err := conn.Do("DEL", prefix+bd.Name)
	_, err = conn.Do("SREM", "boards", bd.Name)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}

func (bd Board) SuggestBoardName() string {
	names := bd.listName()
	boardWeight := map[string]int{}
	chars := strings.Split(strings.ToLower(bd.Name), "")
	for _, name := range names {
		count := 0
		for _, char := range chars {
			if strings.Contains(name, char) {
				count++
			}
		}
		boardWeight[name] = count / (1 + int(math.Abs(float64(len(bd.Name)-len(name)))))
	}
	return maputil.FirstByValueInt(boardWeight)
}

func CheckBoardExist(boardName string) (bool, string) {
	bd := new(Board)
	bd.Name = boardName
	if bd.Exist() {
		return true, ""
	}
	if crawler.CheckBoardExist(boardName) {
		bd.Create()
		return true, ""
	}

	suggestBoard := bd.SuggestBoardName()
	log.WithFields(log.Fields{
		"inputBoard":   boardName,
		"suggestBoard": suggestBoard,
	}).Warning("Board Not Exist")
	return false, suggestBoard
}
