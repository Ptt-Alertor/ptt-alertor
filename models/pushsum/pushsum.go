package pushsum

import (
	"github.com/garyburd/redigo/redis"
	log "github.com/meifamily/logrus"
	"github.com/meifamily/ptt-alertor/connections"
	"github.com/meifamily/ptt-alertor/myutil"
)

const prefix string = "pushsum:"

func List() []string {
	conn := connections.Redis()
	defer conn.Close()
	boards, err := redis.Strings(conn.Do("SMEMBERS", prefix+"boards"))
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return boards
}

func Add(board string) error {
	conn := connections.Redis()
	defer conn.Close()
	_, err := conn.Do("SADD", prefix+"boards", board)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}

func AddSubscriber(board, account string) error {
	conn := connections.Redis()
	defer conn.Close()
	_, err := conn.Do("SADD", prefix+board+":subs", account)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}
