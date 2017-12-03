package shorturl

import (
	"crypto/md5"
	"fmt"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
	log "github.com/meifamily/logrus"

	"strconv"

	"github.com/meifamily/ptt-alertor/connections"
	"github.com/meifamily/ptt-alertor/myutil"
)

const redisPrefix = "sum:"

var url = os.Getenv("APP_HOST") + "/redirect/"

func Gen(longURL string) string {
	data := []byte(longURL)
	sum := fmt.Sprintf("%x", md5.Sum(data))
	sum += strconv.FormatInt(time.Now().Unix(), 10)
	conn := connections.Redis()
	defer conn.Close()
	_, err := conn.Do("SET", redisPrefix+sum, longURL, "EX", 600)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	shortURL := url + sum
	return shortURL
}

func Original(sum string) string {
	conn := connections.Redis()
	defer conn.Close()
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
