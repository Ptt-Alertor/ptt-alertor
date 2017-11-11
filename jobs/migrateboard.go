package jobs

import (
	"strings"

	log "github.com/meifamily/logrus"
	"github.com/meifamily/ptt-alertor/models/article"
	"github.com/meifamily/ptt-alertor/models/author"
	"github.com/meifamily/ptt-alertor/models/board"
	"github.com/meifamily/ptt-alertor/models/keyword"
	"github.com/meifamily/ptt-alertor/models/pushsum"
	"github.com/meifamily/ptt-alertor/models/user"
)

const preBoard = "iphone"
const postBoard = "ios"

type migrateBoard struct{}

func NewMigrateBoard() *migrateBoard {
	return &migrateBoard{}
}

func (migrateBoard) Run() {
	// board list
	addBoard(postBoard)
	bd := board.NewBoard()
	bd.Name = preBoard
	bd.Delete()
	log.Info("Board List Migrated")

	// keyword
	subs := keyword.Subscribers(preBoard)
	for _, sub := range subs {
		keyword.AddSubscriber(postBoard, sub)
	}
	keyword.Destroy(preBoard)
	log.Info("Keyword Migrated")

	// author
	subs = author.Subscribers(preBoard)
	for _, sub := range subs {
		author.AddSubscriber(postBoard, sub)
	}
	author.Destroy(preBoard)
	log.Info("Author Migrated")

	// pushsum
	pushsum.Add(postBoard)
	subs = pushsum.ListSubscribers(preBoard)
	for _, sub := range subs {
		pushsum.AddSubscriber(postBoard, sub)
	}
	pushsum.Remove(preBoard)
	pushsum.Destroy(preBoard)
	pushsum.RenameDiffListKeys(preBoard, postBoard)
	log.Info("Pushsum Migrated")

	// articles
	codes := new(article.Articles).List()
	for _, code := range codes {
		a := new(article.Article).Find(code)
		if strings.EqualFold(a.Board, preBoard) {
			a.Board = postBoard
			a.Save()
			log.WithField("code", code).Info("Article Migrated")
		}
	}
	log.Info("Articles Migrated")

	// user
	users := user.NewUser().All()
	for _, u := range users {
		for _, sub := range u.Subscribes {
			if strings.EqualFold(sub.Board, preBoard) {
				u.Subscribes.Delete(sub)
				for _, postSub := range u.Subscribes {
					if strings.EqualFold(postSub.Board, postBoard) {
						if postSub.PushSum.Up == 0 && postSub.PushSum.Down == 0 {
							postSub.PushSum.Up, postSub.PushSum.Down = sub.PushSum.Up, sub.PushSum.Down
							u.Subscribes.Update(postSub)
						}
					}
				}
				sub.Board = postBoard
				u.Subscribes.Add(sub)
				u.Update()
				log.WithField("account", u.Account).Info("User Subscription Migrated")
			}
		}
	}
	log.Info("Board Migrated")
}
