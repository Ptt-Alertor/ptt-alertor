package top

import (
	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/liam-lai/ptt-alertor/connections"
	"github.com/liam-lai/ptt-alertor/myutil"
)

const prefix string = "top:"

type BoardWord struct {
	Board, Word string
}

type WordOrder struct {
	BoardWord
	Count int
}

type WordOrders []WordOrder

func (wos WordOrders) SaveKeywords() error {
	return wos.save("keywords")
}

func (wos WordOrders) SaveAuthors() error {
	return wos.save("authors")
}

func (wos WordOrders) save(kind string) error {
	conn := connections.Redis()
	for _, wo := range wos {
		if _, err := conn.Do("ZADD", prefix+kind, wo.Count, wo.String()); err != nil {
			log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
			return err
		}
	}
	return nil
}

func GetAuthorList(num int) []string {
	return getList("authors", num)
}

func GetKeywordList(num int) []string {
	return getList("keywords", num)
}

func getList(kind string, num int) []string {
	conn := connections.Redis()
	lists, err := redis.Strings(conn.Do("ZREVRANGE", prefix+kind, 0, num))
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return lists
}

func (wo WordOrder) String() string {
	return wo.Board + ": " + wo.Word
}
