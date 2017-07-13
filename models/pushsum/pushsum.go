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

func DiffList(account, board, day, kind string, ids ...int) []int {
	nowKey := prefix + account + ":" + board + ":" + kind + ":now"
	preKey := prefix + account + ":" + board + ":" + kind + ":pre"
	tmpKey := prefix + account + ":" + board + ":" + kind + ":tmp"
	dayKey := prefix + account + ":" + board + ":" + kind + ":day:" + day
	conn := connections.Redis()
	defer conn.Close()
	bl, err := redis.Bool(conn.Do("EXISTS", preKey))
	conn.Send("MULTI")
	conn.Send("SADD", redis.Args{}.Add(nowKey).AddFlat(ids)...)
	conn.Send("SDIFF", nowKey, preKey)
	conn.Send("RENAME", nowKey, preKey)
	r, err := redis.Values(conn.Do("EXEC"))
	if !bl {
		return []int{}
	}
	ids, err = redis.Ints(r[1], err)
	if len(ids) > 0 {
		conn.Send("MULTI")
		conn.Send("SADD", redis.Args{}.Add(tmpKey).AddFlat(ids)...)
		conn.Send("SDIFF", tmpKey, dayKey)
		conn.Send("SADD", redis.Args{}.Add(dayKey).AddFlat(ids)...)
		conn.Send("DEL", tmpKey)
		r, err = redis.Values(conn.Do("EXEC"))
		ids, err = redis.Ints(r[1], err)
	}
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return ids
}

func DelDiffList(account, board, kind string) error {
	preKeyTemplate := prefix + account + ":" + board + ":" + kind + ":*"
	conn := connections.Redis()
	defer conn.Close()
	preKeys, err := redis.Strings(conn.Do("KEYS", preKeyTemplate))
	_, err = conn.Do("DEL", redis.Args{}.AddFlat(preKeys)...)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}

func DelDayKeys(day string) error {
	keyTemplate := prefix + "*:*:*:day:" + day
	conn := connections.Redis()
	defer conn.Close()
	preKeys, err := redis.Strings(conn.Do("KEYS", keyTemplate))
	_, err = conn.Do("DEL", redis.Args{}.AddFlat(preKeys)...)
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}
