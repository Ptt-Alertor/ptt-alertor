package command

import (
	"strconv"
	"strings"

	"github.com/meifamily/ptt-alertor/models/ptt/article"
	"github.com/meifamily/ptt-alertor/models/pushsum"
	"github.com/meifamily/ptt-alertor/models/subscription"
	user "github.com/meifamily/ptt-alertor/models/user/redis"
	"github.com/meifamily/ptt-alertor/myutil"
)

type updateAction func(u *user.User, sub subscription.Subscription, inputs ...string) error

func addKeywords(u *user.User, sub subscription.Subscription, inputs ...string) error {
	sub.Keywords = inputs
	return u.Subscribes.Add(sub)
}

func removeKeywords(u *user.User, sub subscription.Subscription, inputs ...string) error {
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

func addAuthors(u *user.User, sub subscription.Subscription, inputs ...string) error {
	sub.Authors = inputs
	return u.Subscribes.Add(sub)
}

func removeAuthors(u *user.User, sub subscription.Subscription, inputs ...string) error {
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

func updatePushMax(u *user.User, sub subscription.Subscription, inputs ...string) error {
	max, err := strconv.Atoi(inputs[0])
	if err != nil {
		return err
	}
	for _, s := range u.Subscribes {
		if strings.EqualFold(s.Board, sub.Board) {
			sub.PushSum.Min = s.PushSum.Min
		}
	}
	sub.PushSum.Max = max
	err = u.Subscribes.Update(sub)
	if err == nil {
		dealPushSum(u.Profile.Account, sub)
	}
	return err
}

func updatePushMin(u *user.User, sub subscription.Subscription, inputs ...string) error {
	min, err := strconv.Atoi(inputs[0])
	if err != nil {
		return err
	}
	for _, s := range u.Subscribes {
		if strings.EqualFold(s.Board, sub.Board) {
			sub.PushSum.Max = s.PushSum.Max
		}
	}
	sub.PushSum.Min = min
	err = u.Subscribes.Update(sub)
	if err == nil {
		dealPushSum(u.Profile.Account, sub)
	}
	return err
}

func dealPushSum(account string, sub subscription.Subscription) (err error) {
	if !pushsum.Exist(sub.Board) {
		err = pushsum.Add(sub.Board)
	}
	if sub.Max == 0 && sub.Min == 0 {
		err = pushsum.RemoveSubscriber(sub.Board, account)
	} else {
		err = pushsum.AddSubscriber(sub.Board, account)
	}
	return err
}

func addArticles(u *user.User, sub subscription.Subscription, inputs ...string) error {
	sub.Articles = inputs
	a := article.Article{
		Code: inputs[0],
	}
	a.AddSubscriber(u.Profile.Account)
	return u.Subscribes.Add(sub)
}

func removeArticles(u *user.User, sub subscription.Subscription, inputs ...string) error {
	sub.Articles = inputs
	a := article.Article{
		Code: inputs[0],
	}
	a.RemoveSubscriber(u.Profile.Account)
	return u.Subscribes.Remove(sub)
}
