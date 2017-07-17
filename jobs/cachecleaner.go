package jobs

import (
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
	_, err := redis.Strings(conn.Do("DEL", "boards"))
	boards, err := redis.Strings(conn.Do("KEYS", "board:*"))
	for _, board := range boards {
		_, _ = conn.Do("DEL", board)
	}
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	log.Info("Clean Up Cache")
}
