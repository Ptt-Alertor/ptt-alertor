package user

import (
	"time"

	"github.com/meifamily/ptt-alertor/models/subscription"
)

type User struct {
	Enable     bool      `json:"enable"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
	Profile
	Subscribes subscription.Subscriptions
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

type UserAction interface {
	All() []*User
	Save() error
	Update() error
	Find() User
}
