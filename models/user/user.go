package user

import "github.com/liam-lai/ptt-alertor/models/subscription"

type User struct {
	Enable  bool `json:"enable"`
	Profile struct {
		Account string `json:"account"`
		Email   string `json:"email"`
		Line    string `json:"line"`
	}
	Subscribes subscription.Subscribes
}

type UserAction interface {
	All() []*User
	Save() error
	Update() error
	Find() User
}
