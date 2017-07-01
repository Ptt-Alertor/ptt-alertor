package user

import (
	"time"

	"github.com/meifamily/ptt-alertor/models/subscription"
)

type User struct {
	Enable     bool      `json:"enable"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
	Profile    struct {
		Account         string `json:"account"`
		Email           string `json:"email"`
		Line            string `json:"line"`
		LineAccessToken string `json:"lineAccessToken"`
		Messenger       string `json:"messenger"`
	}
	Subscribes subscription.Subscriptions
}

type UserAction interface {
	All() []*User
	Save() error
	Update() error
	Find() User
}
