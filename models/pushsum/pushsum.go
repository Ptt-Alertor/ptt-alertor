package pushsum

import (
	"strconv"
	"strings"

	log "github.com/Ptt-Alertor/logrus"
	"github.com/Ptt-Alertor/ptt-alertor/connections"
	"github.com/Ptt-Alertor/ptt-alertor/myutil"
	"github.com/garyburd/redigo/redis"
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

func ConvertPushCount(str string) int {
	for num, text := range NumTextMap {
		if strings.EqualFold(str, text) {
			return num
		}
	}
	cnt, err := strconv.Atoi(str)
	if err != nil {
		cnt = 0
	}
	return cnt
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

func Remove(board string) error {
	conn := connections.Redis()
	defer conn.Close()
	_, err := conn.Do("SREM", prefix+"boards", board)
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

func DiffList(account, board, kind string, ids ...int) []int {
	if len(ids) == 0 {
		return []int{}
	}
	nowKey := prefix + account + ":" + board + ":" + kind + ":now"
	baseKey := prefix + account + ":" + board + ":" + kind + ":base"
	benchKey := prefix + account + ":" + board + ":" + kind + ":bench"
	conn := connections.Redis()
	defer conn.Close()
	bl, err := redis.Bool(conn.Do("EXISTS", baseKey))
	conn.Send("MULTI")
	conn.Send("SADD", redis.Args{}.Add(nowKey).AddFlat(ids)...)
	conn.Send("SDIFF", nowKey, baseKey, benchKey)
	conn.Send("DEL", nowKey)
	r, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
		return []int{}
	}
	ids, err = redis.Ints(r[1], err)
	if len(ids) > 0 {
		_, err = conn.Do("SADD", redis.Args{}.Add(baseKey).AddFlat(ids)...)
	}
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	if !bl {
		return []int{}
	}
	return ids
}

func DelDiffList(account, board, kind string) error {
	preKeyTemplate := prefix + account + ":" + board + ":" + kind + ":*"
	conn := connections.Redis()
	defer conn.Close()
	preKeys, err := redis.Strings(conn.Do("KEYS", preKeyTemplate))
	if len(preKeys) > 0 {
		_, err = conn.Do("DEL", redis.Args{}.AddFlat(preKeys)...)
	}
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return err
}

func ReplaceBenchKeys() error {
	baseKeyTemplate := prefix + "*:*:*:base"
	conn := connections.Redis()
	defer conn.Close()
	baseKeys, err := redis.Strings(conn.Do("KEYS", baseKeyTemplate))
	for _, baseKey := range baseKeys {
		key := strings.TrimSuffix(baseKey, "base") + "bench"
		conn.Send("WATCH", key)
		conn.Send("MULTI")
		conn.Send("RENAME", baseKey, key)
		_, err = conn.Do("EXEC")
		if err != nil {
			log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
		}
	}
	return err
}

func RenameDiffListKeys(preBoard, postBoard string) error {
	keyTemplate := prefix + "*:" + preBoard + ":*"
	conn := connections.Redis()
	defer conn.Close()
	keys, err := redis.Strings(conn.Do("KEYS", keyTemplate))
	for _, key := range keys {
		if postBoard == "" {
			_, err = conn.Do("DEL", key)
			continue
		}
		newKey := strings.Replace(key, preBoard, postBoard, -1)
		bl, err := redis.Bool(conn.Do("EXISTS", newKey))
		if err == nil {
			if bl {
				_, err = conn.Do("DEL", key)
			} else {
				conn.Send("WATCH", key)
				conn.Send("MULTI")
				conn.Send("RENAME", key, newKey)
				_, err = conn.Do("EXEC")
			}
		}
		if err != nil {
			log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
		}
	}
	return err
}
