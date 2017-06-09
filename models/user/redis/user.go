package redis

import (
	"encoding/json"

	"errors"

	log "github.com/Sirupsen/logrus"

	"github.com/garyburd/redigo/redis"
	"github.com/liam-lai/ptt-alertor/connections"
	"github.com/liam-lai/ptt-alertor/models/user"
	"github.com/liam-lai/ptt-alertor/myutil"
)

type User struct {
	user.User
}

const prefix string = "user:"

var usersDir string = myutil.StoragePath() + "/users/"

func (u User) All() []*User {
	conn := connections.Redis()
	userKeys, err := redis.Strings(conn.Do("KEYS", "user:*"))
	if err != nil {
		log.WithField("runtime", myutil.BasicRuntimeInfo()).WithError(err).Error()
	}
	us := make([]*User, 0)
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

	if u.Profile.Email == "" && u.Profile.Line == "" && u.Profile.Messenger == "" {
		return errors.New("Email or Line or Messenger have to be complete")
	}

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

	if err != nil {
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
