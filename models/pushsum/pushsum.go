package pushsum

import (
	"github.com/garyburd/redigo/redis"
	log "github.com/meifamily/logrus"
	"github.com/meifamily/ptt-alertor/connections"
	"github.com/meifamily/ptt-alertor/myutil"
)

const prefix string = "pushsum:"

var NumTextMap = map[int]string{
	100:  "爆",
	-10:  "X1",
	-20:  "X2",
	-30:  "X3",
	-40:  "X4",
	-50:  "X5",
	-60:  "X6",
	-70:  "X7",
	-80:  "X8",
	-90:  "X9",
	-100: "XX",
}

func List() []string {
	conn := connections.Redis()
	defer conn.Close()
	boards, err := redis.Strings(conn.Do("SMEMBERS", prefix+"boards"))
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return boards
}

func Exist(board string) bool {
	conn := connections.Redis()
	defer conn.Close()
	bl, err := redis.Bool(conn.Do("SISMEMBER", prefix+"boards", board))
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return bl
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

func RemoveSubscriber(board, account string) error {
	conn := connections.Redis()
	defer conn.Close()
	_, err := conn.Do("SREM", prefix+board+":subs", account)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}

func ListSubscribers(board string) []string {
	conn := connections.Redis()
	defer conn.Close()
	subs, err := redis.Strings(conn.Do("SMEMBERS", prefix+board+":subs"))
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return subs
}

func DiffList(account, board, kind string, ids ...int) []int {
	nowKey := prefix + account + ":" + board + ":" + kind + ":now"
	preKey := prefix + account + ":" + board + ":" + kind + ":pre"
	conn := connections.Redis()
	defer conn.Close()
	bl, err := redis.Bool(conn.Do("EXISTS", preKey))
	conn.Send("MULTI")
	conn.Send("SADD", redis.Args{}.Add(nowKey).AddFlat(ids)...)
	conn.Send("SDIFF", nowKey, preKey)
	conn.Send("RENAME", nowKey, preKey)
	r, err := redis.Values(conn.Do("EXEC"))
	if !bl {
		ids = []int{}
	} else {
		ids, err = redis.Ints(r[1], err)
	}
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return ids
}

func DelDiffList(account, board, kind string) error {
	preKey := prefix + account + ":" + board + ":" + kind + ":pre"
	conn := connections.Redis()
	defer conn.Close()
	_, err := conn.Do("DEL", preKey)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}
