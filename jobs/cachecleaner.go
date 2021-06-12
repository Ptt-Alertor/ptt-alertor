package jobs

import (
	"regexp"

	log "github.com/Ptt-Alertor/logrus"
	"github.com/Ptt-Alertor/ptt-alertor/connections"
	"github.com/Ptt-Alertor/ptt-alertor/myutil"
	"github.com/garyburd/redigo/redis"
)

func NewCacheCleaner() *cacheCleaner {
	return &cacheCleaner{}
}

type cacheCleaner struct {
}

func (c *cacheCleaner) Run() {
	conn := connections.Redis()
	defer conn.Close()
	c.cleanBoardKeys(conn)
	c.cleanKeywordKeys(conn)
	c.cleanAuthorKeys(conn)
	c.cleanPushsumKeys(conn)
	c.cleanCommentKeys(conn)
	// c.cleanUpperCaseKeys(conn)
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

func (c *cacheCleaner) cleanKeywordKeys(conn redis.Conn) {
	keys, _ := redis.Strings(conn.Do("KEYS", "keyword:*:subs"))
	for _, key := range keys {
		_, err := conn.Do("DEL", key)
		if err == nil {
			log.Info("Delete Key: ", key)
		} else {
			log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
		}
	}
	log.Info("Clean Up Keyword Keys")
}

func (c *cacheCleaner) cleanAuthorKeys(conn redis.Conn) {
	keys, _ := redis.Strings(conn.Do("KEYS", "author:*:subs"))
	for _, key := range keys {
		_, err := conn.Do("DEL", key)
		if err == nil {
			log.Info("Delete Key: ", key)
		} else {
			log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
		}
	}
	log.Info("Clean Up Author Keys")
}

func (c *cacheCleaner) cleanPushsumKeys(conn redis.Conn) {
	keys, _ := redis.Strings(conn.Do("KEYS", "keyword:*:subs"))
	for _, key := range keys {
		_, err := conn.Do("DEL", key)
		if err == nil {
			log.Info("Delete Key: ", key)
		} else {
			log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
		}
	}
	log.Info("Clean Up Pushsum Keys")
}

func (c *cacheCleaner) cleanCommentKeys(conn redis.Conn) {
	keys, _ := redis.Strings(conn.Do("KEYS", "article:*:subs"))
	for _, key := range keys {
		_, err := conn.Do("DEL", key)
		if err == nil {
			log.Info("Delete Key: ", key)
		} else {
			log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
		}
	}
	log.Info("Clean Up Comment Keys")
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
