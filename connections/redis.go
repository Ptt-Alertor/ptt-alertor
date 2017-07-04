package connections

import (
	"time"

	"github.com/garyburd/redigo/redis"
	log "github.com/meifamily/logrus"
	"github.com/meifamily/ptt-alertor/myutil"
)

var config map[string]string
var pool = newPool()

func init() {
	config = myutil.Config("redis")
}

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", config["host"]+":"+config["port"])
			if err != nil {
				log.Fatal(err)
			}
			return conn, err
		},
	}
}

// Redis get redis connection
func Redis() redis.Conn {
	return pool.Get()
}
