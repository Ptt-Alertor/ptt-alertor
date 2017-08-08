package file

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/meifamily/logrus"

	"github.com/meifamily/ptt-alertor/myutil"
)

type User struct {
}

var usersDir string = myutil.StoragePath() + "/users/"

func (u User) List() (accounts []string) {
	files, err := ioutil.ReadDir(usersDir)
	if err != nil {
		log.WithFields(log.Fields{
			"directory": usersDir,
			"runtime":   myutil.BasicRuntimeInfo(),
		}).WithError(err).Error("Read Directory Error")
	}
	for _, file := range files {
		_, ok := myutil.JsonFile(file)
		if !ok {
			continue
		}
		accounts = append(accounts, file.Name())
	}
	return accounts
}

func (u User) Exist(account string) bool {
	userFile := usersDir + account + ".json"
	_, err := ioutil.ReadFile(userFile)
	if err != nil {
		return false
	}
	return true
}

func (u User) Save(account string, user interface{}) error {
	userFile := usersDir + account + ".json"
	uJSON, err := json.Marshal(user)
	if err != nil {
		myutil.LogJSONEncode(err, user)
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

func (u User) Update(account string, user interface{}) error {
	userFile := usersDir + account + ".json"
	uJSON, err := json.Marshal(user)
	if err != nil {
		myutil.LogJSONEncode(err, user)
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

func (u User) Find(account string, user interface{}) {
	userFile := usersDir + account + ".json"
	uJSON, err := ioutil.ReadFile(userFile)
	if err != nil {
		log.WithFields(log.Fields{
			"file":    userFile,
			"runtime": myutil.BasicRuntimeInfo(),
		}).WithError(err).Error("Read File Error")
	}
	err = json.Unmarshal(uJSON, &user)
	if err != nil {
		myutil.LogJSONDecode(err, uJSON)
	}
}
