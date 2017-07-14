package redis

import (
	"encoding/json"
	"time"

	"errors"

	log "github.com/meifamily/logrus"

	"github.com/garyburd/redigo/redis"
	"github.com/meifamily/ptt-alertor/connections"
	"github.com/meifamily/ptt-alertor/models/user"
	"github.com/meifamily/ptt-alertor/myutil"
)

type User struct {
	user.User
}

const prefix string = "user:"

func (u User) All() (us []*User) {
	conn := connections.Redis()
	defer conn.Close()
	userKeys, err := redis.Strings(conn.Do("KEYS", "user:*"))
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	for _, uKey := range userKeys {
		uJSON, _ := redis.Bytes(conn.Do("GET", uKey))
		var user User
		err = json.Unmarshal(uJSON, &user)
		if err != nil {
			myutil.LogJSONDecode(err, uJSON)
		}
		us = append(us, &user)
	}
	return us
}

func (u User) Save() error {
	conn := connections.Redis()
	defer conn.Close()
	key := prefix + u.Profile.Account
	val, err := conn.Do("GET", key)

	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}

	if val != nil {
		return errors.New("user already exist")
	}

	if u.Profile.Account == "" {
		return errors.New("account can not be empty")
	}

	if u.Profile.Email == "" && u.Profile.Line == "" && u.Profile.Messenger == "" && u.Profile.Telegram == "" {
		return errors.New("one of Email, Line, Messenger and Telegram have to be complete")
	}

	u.CreateTime = time.Now()
	u.UpdateTime = time.Now()
	uJSON, err := json.Marshal(u)
	if err != nil {
		myutil.LogJSONEncode(err, u)
		return err
	}

	_, err = conn.Do("SET", key, uJSON, "NX")

	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
		return err
	}
	return nil
}

func (u User) Update() error {

	conn := connections.Redis()
	defer conn.Close()

	key := prefix + u.Profile.Account
	val, err := conn.Do("GET", key)

	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}

	if val == nil {
		return errors.New("user not exist")
	}

	if u.Profile.Account == "" {
		return errors.New("account can not be empty")
	}

	u.UpdateTime = time.Now()
	uJSON, err := json.Marshal(u)
	if err != nil {
		myutil.LogJSONEncode(err, u)
		return err
	}

	_, err = conn.Do("SET", key, uJSON, "XX")

	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
		return err
	}
	return nil
}

func (u User) Find(account string) User {
	conn := connections.Redis()
	defer conn.Close()

	key := prefix + account
	uJSON, err := redis.Bytes(conn.Do("GET", key))

	if err != nil && err != redis.ErrNil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}

	if uJSON != nil {
		err = json.Unmarshal(uJSON, &u)
		if err != nil {
			myutil.LogJSONDecode(err, uJSON)
		}
	}

	return u
}
