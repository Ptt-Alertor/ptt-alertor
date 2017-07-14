package jobs

import (
	"errors"

	user "github.com/meifamily/ptt-alertor/models/user/redis"
)

var platforms = map[string]bool{
	"email":     true,
	"line":      true,
	"messenger": true,
}

type Broadcast struct {
	Checker
	Msg string
}

func (bc Broadcast) String() string {
	return bc.Msg
}

func (bc Broadcast) Send(plfms []string) error {
	var platformBl = make(map[string]bool)
	for _, plfm := range plfms {
		if _, ok := platforms[plfm]; !ok {
			return errors.New("platform " + plfm + "is not in broadcast list")
		}
		platformBl[plfm] = true
	}

	users := new(user.User).All()
	for _, u := range users {
		bc.subType = "broadcast"
		if platformBl["line"] {
			go bc.sendLine(u)
		}
		if platformBl["messenger"] {
			go bc.sendMessenger(u)
		}
		if platformBl["email"] {
			go bc.sendEmail(u)
		}
	}
	return nil
}

func (bc Broadcast) sendEmail(u *user.User) {
	bc.Profile.Email = u.Profile.Email
	ckCh <- bc
}

func (bc Broadcast) sendLine(u *user.User) {
	bc.Profile.Line = u.Profile.Line
	bc.Profile.LineAccessToken = u.Profile.LineAccessToken
	ckCh <- bc
}

func (bc Broadcast) sendMessenger(u *user.User) {
	bc.Profile.Messenger = u.Profile.Messenger
	ckCh <- bc
}
