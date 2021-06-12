//+build !test

package user

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/Ptt-Alertor/logrus"

	"github.com/Ptt-Alertor/ptt-alertor/myutil"
)

type File struct {
}

var usersDir string = myutil.StoragePath() + "/users/"

func (File) List() (accounts []string) {
	files, err := ioutil.ReadDir(usersDir)
	if err != nil {
		log.WithFields(log.Fields{
			"directory": usersDir,
			"runtime":   myutil.BasicRuntimeInfo(),
		}).WithError(err).Error("Read Directory Error")
	}
	for _, file := range files {
		_, ok := myutil.JSONFile(file)
		if !ok {
			continue
		}
		accounts = append(accounts, file.Name())
	}
	return accounts
}

func (File) Exist(account string) bool {
	userFile := usersDir + account + ".json"
	if _, err := ioutil.ReadFile(userFile); err != nil {
		return false
	}
	return true
}

func (File) Save(account string, user interface{}) error {
	userFile := usersDir + account + ".json"
	uJSON, err := json.Marshal(user)
	if err != nil {
		myutil.LogJSONEncode(err, user)
		return err
	}

	if err := ioutil.WriteFile(userFile, uJSON, 664); err != nil {
		log.WithFields(log.Fields{
			"file":    userFile,
			"doc":     uJSON,
			"runtime": myutil.BasicRuntimeInfo(),
		}).WithError(err).Error("Write File Error")
		return err
	}
	return nil
}

func (File) Update(account string, user interface{}) error {
	userFile := usersDir + account + ".json"
	uJSON, err := json.Marshal(user)
	if err != nil {
		myutil.LogJSONEncode(err, user)
		return err
	}
	if err := ioutil.WriteFile(userFile, uJSON, 664); err != nil {
		log.WithFields(log.Fields{
			"file":    userFile,
			"doc":     uJSON,
			"runtime": myutil.BasicRuntimeInfo(),
		}).WithError(err).Error("Write File Error")
		return err
	}
	return nil
}

func (File) Find(account string, user *User) {
	userFile := usersDir + account + ".json"
	uJSON, err := ioutil.ReadFile(userFile)
	if err != nil {
		log.WithFields(log.Fields{
			"file":    userFile,
			"runtime": myutil.BasicRuntimeInfo(),
		}).WithError(err).Error("Read File Error")
	}
	if err := json.Unmarshal(uJSON, &user); err != nil {
		myutil.LogJSONDecode(err, uJSON)
	}
}
