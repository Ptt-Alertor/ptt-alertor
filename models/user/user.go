package user

import (
	"time"

	"github.com/meifamily/ptt-alertor/models/subscription"
)

// TODO: fetch Profile to outside, for share with other package
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
		Telegram        string `json:"telegram"`
		TelegramChat    int64  `json:"telegramChat"`
	}
	Subscribes subscription.Subscriptions
}

type UserAction interface {
	All() []*User
	Save() error
	Update() error
	Find() User
}
