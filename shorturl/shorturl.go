package shorturl

import (
	"crypto/sha1"

	"fmt"

	"github.com/garyburd/redigo/redis"
	log "github.com/meifamily/logrus"

	"github.com/meifamily/ptt-alertor/connections"
	"github.com/meifamily/ptt-alertor/myutil"
)

const redisPrefix = "sum:"

var url = "https://pttalertor.dinolai.com/redirect/"

func init() {
	config := myutil.Config("app")
	url = config["host"] + "/redirect/"
}

func Gen(longURL string) string {
	data := []byte(longURL)
	sum := fmt.Sprintf("%x", sha1.Sum(data))
	conn := connections.Redis()
	_, err := conn.Do("SET", redisPrefix+sum, longURL, "EX", 600)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	shortURL := url + sum
	return shortURL
}

func Original(sum string) string {
	conn := connections.Redis()
	key := redisPrefix + sum
	conn.Send("MULTI")
	conn.Send("GET", key)
	conn.Send("DEL", key)
	result, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	if result[0] == nil {
		return ""
	}
	return string(result[0].([]byte))
}
