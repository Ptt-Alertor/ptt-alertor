package jobs

import (
	"regexp"

	"github.com/garyburd/redigo/redis"
	log "github.com/meifamily/logrus"
	"github.com/meifamily/ptt-alertor/connections"
	"github.com/meifamily/ptt-alertor/myutil"
)

func NewCacheCleaner() *cacheCleaner {
	return &cacheCleaner{}
}

type cacheCleaner struct {
}

// TODO: clean keyword, author, pushsum, pushlist keys
func (c *cacheCleaner) Run() {
	conn := connections.Redis()
	defer conn.Close()
	// c.cleanBoardKeys(conn)
	c.cleanUpperCaseKeys(conn)
	log.Info("Clean Up Cache")
}

func (c *cacheCleaner) cleanBoardKeys(conn redis.Conn) {
	_, err := redis.Strings(conn.Do("DEL", "boards"))
	boards, err := redis.Strings(conn.Do("KEYS", "board:*"))
	for _, board := range boards {
		_, _ = conn.Do("DEL", board)
	}
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	log.Info("Clean Up Board Keys")
}

func (c *cacheCleaner) cleanUpperCaseKeys(conn redis.Conn) {
	keys, _ := redis.Strings(conn.Do("KEYS", "keyword:*:subs"))
	tmp, _ := redis.Strings(conn.Do("KEYS", "author:*:subs"))
	keys = append(keys, tmp...)
	tmp, _ = redis.Strings(conn.Do("KEYS", "pushsum:*:subs"))
	keys = append(keys, tmp...)
	for _, key := range keys {
		if bl, _ := regexp.MatchString("[A-Z]", key); bl {
			_, err := conn.Do("DEL", key)
			if err == nil {
				log.Info("Delete Key: ", key)
			}
		}
	}
	log.Info("Clean Up UpperCase Keys")
}
