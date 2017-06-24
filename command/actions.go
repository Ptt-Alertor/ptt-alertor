package command

import (
	"strings"

	"github.com/liam-lai/ptt-alertor/models/subscription"
	user "github.com/liam-lai/ptt-alertor/models/user/redis"
	"github.com/liam-lai/ptt-alertor/myutil"
)

type updateAction func(u *user.User, sub subscription.Subscription, inputs []string) error

func addKeywords(u *user.User, sub subscription.Subscription, inputs []string) error {
	sub.Keywords = inputs
	return u.Subscribes.Add(sub)
}

func addAuthors(u *user.User, sub subscription.Subscription, inputs []string) error {
	sub.Authors = inputs
	return u.Subscribes.Add(sub)
}

func removeKeywords(u *user.User, sub subscription.Subscription, inputs []string) error {
	sub.Keywords = inputs
	if inputs[0] == "*" {
		for _, uSub := range u.Subscribes {
			if strings.EqualFold(uSub.Board, sub.Board) {
				sub.Keywords = make(myutil.StringSlice, len(uSub.Keywords))
				copy(sub.Keywords, uSub.Keywords)
			}
		}
	}
	return u.Subscribes.Remove(sub)
}

func removeAuthors(u *user.User, sub subscription.Subscription, inputs []string) error {
	sub.Authors = inputs
	if inputs[0] == "*" {
		for _, uSub := range u.Subscribes {
			if strings.EqualFold(uSub.Board, sub.Board) {
				sub.Authors = make(myutil.StringSlice, len(uSub.Authors))
				copy(sub.Authors, uSub.Authors)
			}
		}
	}
	return u.Subscribes.Remove(sub)
}
