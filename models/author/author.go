package author

import (
	log "github.com/meifamily/logrus"

	"github.com/garyburd/redigo/redis"
	"github.com/meifamily/ptt-alertor/connections"
	"github.com/meifamily/ptt-alertor/myutil"
)

const prefix string = "author:"

func Subscribers(board string) []string {
	key := prefix + board + ":subs"
	conn := connections.Redis()
	defer conn.Close()
	accounts, err := redis.Strings(conn.Do("SMEMBERS", key))
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return accounts
}

func AddSubscriber(board, account string) error {
	key := prefix + board + ":subs"
	conn := connections.Redis()
	defer conn.Close()
	_, err := conn.Do("SADD", key, account)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}

func RemoveSubscriber(board, account string) error {
	key := prefix + board + ":subs"
	conn := connections.Redis()
	defer conn.Close()
	_, err := conn.Do("SREM", key, account)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}

func Destroy(board string) error {
	key := prefix + board + ":subs"
	conn := connections.Redis()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}
