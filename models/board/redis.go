package board

import (
	"encoding/json"

	"github.com/garyburd/redigo/redis"
	log "github.com/meifamily/logrus"

	"github.com/meifamily/ptt-alertor/connections"
	"github.com/meifamily/ptt-alertor/models/article"
	"github.com/meifamily/ptt-alertor/myutil"
)

const prefix string = "board:"

type Redis struct {
}

func (Redis) List() []string {
	conn := connections.Redis()
	defer conn.Close()
	boards, err := redis.Strings(conn.Do("SMEMBERS", "boards"))
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return boards
}

func (Redis) Exist(boardName string) bool {
	conn := connections.Redis()
	defer conn.Close()
	bl, err := redis.Bool(conn.Do("SISMEMBER", "boards", boardName))
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return bl
}

func (Redis) GetArticles(boardName string) (articles article.Articles) {
	conn := connections.Redis()
	defer conn.Close()

	key := prefix + boardName
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

func (Redis) Create(boardName string) error {
	conn := connections.Redis()
	defer conn.Close()
	_, err := conn.Do("SADD", "boards", boardName)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}

func (Redis) Save(boardName string, articles article.Articles) error {
	conn := connections.Redis()
	defer conn.Close()

	articlesJSON, err := json.Marshal(articles)
	if err != nil {
		myutil.LogJSONEncode(err, articles)
		return err
	}
	conn.Send("WATCH", prefix+boardName)
	conn.Send("MULTI")
	conn.Send("SET", prefix+boardName, articlesJSON)
	_, err = conn.Do("EXEC")
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}

func (Redis) Delete(boardName string) error {
	conn := connections.Redis()
	defer conn.Close()
	_, err := conn.Do("DEL", prefix+boardName)
	_, err = conn.Do("SREM", "boards", boardName)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}
