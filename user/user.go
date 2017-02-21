package user

import (
	"encoding/json"
	"io/ioutil"

	"github.com/liam-lai/ptt-alertor/myutil"
)

type Users []*User

type User struct {
	Profile struct {
		Email string
	}
	Subscribes []Subscribe
}

type Subscribe struct {
	Board    string
	Keywords []string
}

var usersDir string = myutil.StoragePath() + "/users/"

func (us Users) All() Users {
	files, _ := ioutil.ReadDir(usersDir)
	for _, file := range files {
		_, ok := myutil.JsonFile(file)
		if !ok {
			continue
		}
		userJSON, _ := ioutil.ReadFile(usersDir + file.Name())
		var user User
		_ = json.Unmarshal(userJSON, &user)
		us = append(us, &user)
	}
	return us
}
