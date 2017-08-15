package command

import (
	"strconv"
	"strings"

	log "github.com/meifamily/logrus"

	"github.com/meifamily/ptt-alertor/models/article"
	"github.com/meifamily/ptt-alertor/models/author"
	"github.com/meifamily/ptt-alertor/models/keyword"
	"github.com/meifamily/ptt-alertor/models/pushsum"
	"github.com/meifamily/ptt-alertor/models/subscription"
	"github.com/meifamily/ptt-alertor/models/user"
	"github.com/meifamily/ptt-alertor/myutil"
)

type updateAction func(u *user.User, sub subscription.Subscription, inputs ...string) error

func addKeywords(u *user.User, sub subscription.Subscription, inputs ...string) error {
	sub.Keywords = inputs
	err := u.Subscribes.Add(sub)
	if err == nil {
		log.WithFields(log.Fields{
			"board":   sub.Board,
			"account": u.Profile.Account,
		}).Debug("Add Board Subscriber")
		err = keyword.AddSubscriber(sub.Board, u.Profile.Account)
	}
	return err
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
	err := u.Subscribes.Remove(sub)
	if err == nil {
		for _, uSub := range u.Subscribes {
			if strings.EqualFold(sub.Board, uSub.Board) && len(uSub.Keywords) > 0 {
				return nil
			}
		}
		err = keyword.RemoveSubscriber(sub.Board, u.Profile.Account)
	}
	return err
}

func addAuthors(u *user.User, sub subscription.Subscription, inputs ...string) error {
	sub.Authors = inputs
	err := u.Subscribes.Add(sub)
	if err == nil {
		log.WithFields(log.Fields{
			"board":   sub.Board,
			"account": u.Profile.Account,
		}).Debug("Add Board Subscriber")
		err = author.AddSubscriber(sub.Board, u.Profile.Account)
	}
	return err
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
	err := u.Subscribes.Remove(sub)
	if err == nil {
		for _, uSub := range u.Subscribes {
			if strings.EqualFold(sub.Board, uSub.Board) && len(uSub.Authors) > 0 {
				return nil
			}
		}
		err = author.RemoveSubscriber(sub.Board, u.Profile.Account)
	}
	return err
}

func updatePushUp(u *user.User, sub subscription.Subscription, inputs ...string) error {
	up, err := strconv.Atoi(inputs[0])
	if err != nil {
		return err
	}
	for _, s := range u.Subscribes {
		if strings.EqualFold(s.Board, sub.Board) {
			sub.PushSum.Down = s.PushSum.Down
		}
	}
	sub.PushSum.Up = up
	err = u.Subscribes.Update(sub)
	if err == nil {
		err = dealPushSum(u.Profile.Account, sub)
	}
	return err
}

func updatePushDown(u *user.User, sub subscription.Subscription, inputs ...string) error {
	down, err := strconv.Atoi(inputs[0])
	if err != nil {
		return err
	}
	for _, s := range u.Subscribes {
		if strings.EqualFold(s.Board, sub.Board) {
			sub.PushSum.Up = s.PushSum.Up
		}
	}
	sub.PushSum.Down = down
	err = u.Subscribes.Update(sub)
	if err == nil {
		err = dealPushSum(u.Profile.Account, sub)
	}
	return err
}

func dealPushSum(account string, sub subscription.Subscription) (err error) {
	if !pushsum.Exist(sub.Board) {
		err = pushsum.Add(sub.Board)
	}
	if sub.Up == 0 {
		err = pushsum.DelDiffList(account, sub.Board, "up")
	}
	if sub.Down == 0 {
		err = pushsum.DelDiffList(account, sub.Board, "down")
	}
	if sub.Up == 0 && sub.Down == 0 {
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
