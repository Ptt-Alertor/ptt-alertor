package article

import (
	"encoding/json"

	log "github.com/Ptt-Alertor/logrus"
	"github.com/Ptt-Alertor/ptt-alertor/connections"
	"github.com/Ptt-Alertor/ptt-alertor/myutil"
	"github.com/garyburd/redigo/redis"
)

type Redis struct{}

var connectRedis = connections.Redis

const detailSuffix = ":detail"

func (Redis) Find(code string, a *Article) {
	conn := connectRedis()
	defer conn.Close()

	aMap, err := redis.StringMap(conn.Do("HGETALL", prefix+code+detailSuffix))
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	a.Board = aMap["board"]

	if err = json.Unmarshal([]byte(aMap["content"]), &a); err != nil {
		log.WithField("code", code).Error("Article Content Unmarshal Failed")
		myutil.LogJSONDecode(err, aMap["content"])
	}
}

func (Redis) Save(a Article) error {
	conn := connectRedis()
	defer conn.Close()

	articleJSON, err := json.Marshal(a)
	if err != nil {
		myutil.LogJSONEncode(err, a)
		return err
	}
	_, err = conn.Do("HMSET", prefix+a.Code+detailSuffix, "board", a.Board, "content", articleJSON)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}

func (Redis) Delete(articleCode string) error {
	conn := connectRedis()
	defer conn.Close()

	_, err := conn.Do("DEL", prefix+articleCode+subsSuffix)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}
