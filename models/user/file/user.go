package file

import (
	"encoding/json"
	"io/ioutil"

	"errors"

	log "github.com/meifamily/logrus"

	"reflect"

	"github.com/meifamily/ptt-alertor/models/user"
	"github.com/meifamily/ptt-alertor/myutil"
)

type User struct {
	user.User
}

var usersDir string = myutil.StoragePath() + "/users/"

func (u User) All() []*User {
	files, err := ioutil.ReadDir(usersDir)
	if err != nil {
		log.WithFields(log.Fields{
			"directory": usersDir,
			"runtime":   myutil.BasicRuntimeInfo(),
		}).WithError(err).Error("Read Directory Error")
	}
	us := make([]*User, 0)
	for _, file := range files {
		_, ok := myutil.JsonFile(file)
		if !ok {
			continue
		}
		userFile := usersDir + file.Name()
		userJSON, err := ioutil.ReadFile(userFile)
		if err != nil {
			log.WithFields(log.Fields{
				"file":    userFile,
				"runtime": myutil.BasicRuntimeInfo(),
			}).WithError(err).Error("Read File Error")
		}
		var user User
		err = json.Unmarshal(userJSON, &user)
		if err != nil {
			myutil.LogJSONDecode(err, userJSON)
		}
		us = append(us, &user)
	}
	return us
}

func (u User) Save() error {
	if u.Profile.Account == "" {
		return errors.New("account can not be empty")
	}

	if u.Profile.Email == "" {
		return errors.New("email can not be empty")
	}

	userFile := usersDir + u.Profile.Account + ".json"
	doc, err := ioutil.ReadFile(userFile)

	if reflect.TypeOf(err).String() != "*os.PathError" {
		log.WithFields(log.Fields{
			"file":    userFile,
			"runtime": myutil.BasicRuntimeInfo(),
		}).WithError(err).Error("Read File Error")
		return errors.New("Save User Unknow Error")
	}

	if doc != nil {
		return errors.New("user already exist")
	}

	uJSON, err := json.Marshal(u)
	if err != nil {
		myutil.LogJSONEncode(err, u)
		return err
	}

	err = ioutil.WriteFile(userFile, uJSON, 664)
	if err != nil {
		log.WithFields(log.Fields{
			"file":    userFile,
			"doc":     uJSON,
			"runtime": myutil.BasicRuntimeInfo(),
		}).WithError(err).Error("Write File Error")
		return err
	}
	return nil
}

func (u User) Update() error {
	userFile := usersDir + u.Profile.Account + ".json"
	val, err := ioutil.ReadFile(userFile)

	if val == nil {
		return errors.New("user not exist")
	}

	if err != nil {
		log.WithFields(log.Fields{
			"file":    userFile,
			"runtime": myutil.BasicRuntimeInfo(),
		}).WithError(err).Error("Read File Error")
	}

	if u.Profile.Account == "" {
		return errors.New("account can not be empty")
	}

	uJSON, err := json.Marshal(u)
	if err != nil {
		myutil.LogJSONEncode(err, u)
		return err
	}
	err = ioutil.WriteFile(userFile, uJSON, 664)
	if err != nil {
		log.WithFields(log.Fields{
			"file":    userFile,
			"doc":     uJSON,
			"runtime": myutil.BasicRuntimeInfo(),
		}).WithError(err).Error("Write File Error")
		return err
	}
	return nil
}

func (u User) Find(account string) User {
	userFile := usersDir + account + ".json"
	uJSON, err := ioutil.ReadFile(userFile)
	if err != nil {
		log.WithFields(log.Fields{
			"file":    userFile,
			"runtime": myutil.BasicRuntimeInfo(),
		}).WithError(err).Error("Read File Error")
		return u
	}
	err = json.Unmarshal(uJSON, &u)
	if err != nil {
		myutil.LogJSONDecode(err, uJSON)
	}
	return u
}
