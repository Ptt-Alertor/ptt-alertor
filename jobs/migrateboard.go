package jobs

import (
	"strings"

	log "github.com/meifamily/logrus"
	"github.com/meifamily/ptt-alertor/models"
	"github.com/meifamily/ptt-alertor/models/article"
	"github.com/meifamily/ptt-alertor/models/author"
	"github.com/meifamily/ptt-alertor/models/keyword"
	"github.com/meifamily/ptt-alertor/models/pushsum"
)

type migrateBoard struct {
	boardMap map[string]string
}

func NewMigrateBoard(boardMap map[string]string) *migrateBoard {
	return &migrateBoard{boardMap: boardMap}
}

func (m migrateBoard) Run() {
	for pre, post := range m.boardMap {
		log.WithField("board", pre).Info("Board Migrating")
		m.RunSingle(pre, post)
	}
	log.Info("All Board Migrated")
}

func (migrateBoard) RunSingle(preBoard string, postBoard string) {
	// board list
	if postBoard != "" {
		addBoard(postBoard)
	}
	bd := models.Board()
	bd.Name = preBoard
	bd.Delete()
	log.WithField("board", preBoard).Info("Board List Migrated")

	// keyword
	if postBoard != "" {
		subs := keyword.Subscribers(preBoard)
		for _, sub := range subs {
			keyword.AddSubscriber(postBoard, sub)
		}
	}
	keyword.Destroy(preBoard)
	log.WithField("board", preBoard).Info("Keyword Migrated")

	// author
	if postBoard != "" {
		subs := author.Subscribers(preBoard)
		for _, sub := range subs {
			author.AddSubscriber(postBoard, sub)
		}
	}
	author.Destroy(preBoard)
	log.WithField("board", preBoard).Info("Author Migrated")

	// pushsum
	if postBoard != "" {
		pushsum.Add(postBoard)
		subs := pushsum.ListSubscribers(preBoard)
		for _, sub := range subs {
			pushsum.AddSubscriber(postBoard, sub)
		}
	}
	pushsum.Remove(preBoard)
	pushsum.Destroy(preBoard)
	pushsum.RenameDiffListKeys(preBoard, postBoard)
	log.WithField("board", preBoard).Info("Pushsum Migrated")

	// articles
	codes := new(article.Articles).List()
	for _, code := range codes {
		a := models.Article().Find(code)
		if strings.EqualFold(a.Board, preBoard) {
			if postBoard == "" {
				a.Destroy()
			} else {
				a.Board = postBoard
				a.Save()
			}
			log.WithFields(log.Fields{
				"board": preBoard,
				"code":  code,
			}).Info("Article Migrated")
		}
	}
	log.Info("Articles Migrated")

	// user
	for _, u := range models.User.All() {
		for _, sub := range u.Subscribes {
			if strings.EqualFold(sub.Board, preBoard) {
				u.Subscribes.Delete(sub)
				if postBoard != "" {
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
				}
				u.Update()
				log.WithFields(log.Fields{
					"account": u.Account,
					"board":   preBoard,
				}).Info("User Subscription Migrated")
			}
		}
	}
}
