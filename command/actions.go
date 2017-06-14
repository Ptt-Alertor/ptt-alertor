package command

import (
	"github.com/liam-lai/ptt-alertor/models/subscription"
	user "github.com/liam-lai/ptt-alertor/models/user/redis"
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
	return u.Subscribes.Remove(sub)
}

func removeAuthors(u *user.User, sub subscription.Subscription, inputs []string) error {
	sub.Authors = inputs
	return u.Subscribes.Remove(sub)
}
