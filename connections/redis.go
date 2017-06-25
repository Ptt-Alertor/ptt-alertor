package connections

import (
	log "github.com/meifamily/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/liam-lai/ptt-alertor/myutil"
)

var config map[string]string

func init() {
	config = myutil.Config("redis")
}

func Redis() redis.Conn {
	conn, err := redis.Dial("tcp", config["host"]+":"+config["port"])
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
