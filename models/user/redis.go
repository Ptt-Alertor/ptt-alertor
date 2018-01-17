package user

import (
	"encoding/json"

	log "github.com/meifamily/logrus"

	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/meifamily/ptt-alertor/connections"
	"github.com/meifamily/ptt-alertor/myutil"
)

type Redis struct {
}

var connectRedis = connections.Redis

const prefix string = "user:"

func (Redis) List() (accounts []string) {
	conn := connectRedis()
	defer conn.Close()
	userKeys, err := redis.Strings(conn.Do("KEYS", "user:*"))
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	for _, key := range userKeys {
		accounts = append(accounts, strings.TrimPrefix(key, "user:"))
	}
	return accounts
}

func (Redis) Exist(account string) bool {
	conn := connectRedis()
	defer conn.Close()
	key := prefix + account
	bl, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	return bl
}

func (Redis) Save(account string, data interface{}) error {
	conn := connectRedis()
	defer conn.Close()
	key := prefix + account
	uJSON, err := json.Marshal(data)
	if err != nil {
		myutil.LogJSONEncode(err, data)
		return err
	}

	_, err = conn.Do("SET", key, uJSON, "NX")
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
		return err
	}
	return nil
}

func (Redis) Update(account string, user interface{}) error {
	conn := connectRedis()
	defer conn.Close()
	key := prefix + account
	uJSON, err := json.Marshal(user)
	if err != nil {
		myutil.LogJSONEncode(err, user)
		return err
	}

	_, err = conn.Do("SET", key, uJSON, "XX")
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
		return err
	}
	return nil
}

func (Redis) Find(account string, user *User) {
	conn := connectRedis()
	defer conn.Close()

	key := prefix + account
	uJSON, err := redis.Bytes(conn.Do("GET", key))
	if err != nil && err != redis.ErrNil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}

	if uJSON != nil {
		err = json.Unmarshal(uJSON, &user)
		if err != nil {
			myutil.LogJSONDecode(err, uJSON)
		}
	}
}
