package user

import (
	"encoding/json"
	"io/ioutil"

	"errors"

	"log"

	"github.com/liam-lai/ptt-alertor/myutil"
)

type Users []*User

type User struct {
	Profile struct {
		Account string `json:"account"`
		Email   string `json:"email"`
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

func (u User) Save() error {
	_, err := ioutil.ReadFile(usersDir + u.Profile.Account + ".json")
	if err == nil {
		return errors.New("user already exist")
	}

	if u.Profile.Account == "" {
		return errors.New("account can not be empty")
	}

	if u.Profile.Email == "" {
		return errors.New("email can not be empty")
	}

	uJSON, err := json.Marshal(u)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(usersDir+u.Profile.Account+".json", uJSON, 664)
	if err != nil {
		return err
	}
	return nil
}

func (u User) Find(account string) User {
	uJSON, err := ioutil.ReadFile(usersDir + account + ".json")
	if err != nil {
		return u
	}
	err = json.Unmarshal(uJSON, &u)
	if err != nil {
		log.Fatal(err)
	}
	return u
}
