package user

import (
	"errors"
	"time"

	"github.com/meifamily/ptt-alertor/models/subscription"
	"github.com/meifamily/ptt-alertor/models/user/redis"
)

type User struct {
	Enable     bool      `json:"enable"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
	Profile    `json:"Profile"`
	Subscribes subscription.Subscriptions
	driver     Driver
}

type Profile struct {
	Account         string `json:"account"`
	Email           string `json:"email"`
	Line            string `json:"line"`
	LineAccessToken string `json:"lineAccessToken"`
	Messenger       string `json:"messenger"`
	Telegram        string `json:"telegram"`
	TelegramChat    int64  `json:"telegramChat"`
}

type Driver interface {
	List() (accounts []string)
	Exist(account string) bool
	Save(account string, user interface{}) error
	Update(account string, user interface{}) error
	Find(account string, user interface{})
}

var driver = new(redis.User)

func NewUser() *User {
	return &User{
		driver: driver,
	}
}

func (u User) All() (us []*User) {
	accounts := u.driver.List()
	for _, account := range accounts {
		user := u.Find(account)
		us = append(us, &user)
	}
	return us
}

func (u User) Save() error {

	if u.driver.Exist(u.Profile.Account) {
		return errors.New("user already exist")
	}

	if u.Profile.Account == "" {
		return errors.New("account can not be empty")
	}

	if u.Profile.Email == "" && u.Profile.Line == "" && u.Profile.Messenger == "" && u.Profile.Telegram == "" {
		return errors.New("one of Email, Line, Messenger and Telegram have to be filled")
	}
	u.CreateTime = time.Now()
	u.UpdateTime = time.Now()

	return u.driver.Save(u.Profile.Account, u)
}

func (u User) Update() error {

	if !u.driver.Exist(u.Profile.Account) {
		return errors.New("user not exist")
	}

	if u.Profile.Account == "" {
		return errors.New("account can not be empty")
	}

	u.UpdateTime = time.Now()
	return u.driver.Update(u.Profile.Account, u)
}

func (u User) Find(account string) User {
	u.driver.Find(account, &u)
	return u
}
