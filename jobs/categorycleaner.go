package jobs

import (
	"strings"
	"time"

	log "github.com/Ptt-Alertor/logrus"
	"github.com/Ptt-Alertor/ptt-alertor/models"
	"github.com/Ptt-Alertor/ptt-alertor/models/author"
	"github.com/Ptt-Alertor/ptt-alertor/models/keyword"
	"github.com/Ptt-Alertor/ptt-alertor/models/pushsum"
	"github.com/Ptt-Alertor/ptt-alertor/myutil"
	"github.com/Ptt-Alertor/ptt-alertor/ptt/rss"
)

type categoryCleaner struct {
}

func NewCategoryCleaner() *categoryCleaner {
	return &categoryCleaner{}
}

func (cc categoryCleaner) Run() {
	boardNames := myutil.StringSlice(models.Board().List())
	boardNames.AppendNonRepeat(myutil.StringSlice(pushsum.List()), false)

	for _, boardName := range boardNames {
		time.Sleep(100 * time.Millisecond)
		if !rss.CheckBoardExist(boardName) {
			log.WithField("category", boardName).Info("Delete Category")
			cc.CleanAccountSetting(boardName)
			cc.CleanKeywordAuthorBoard(boardName)
			cc.CleanPushsumBoard(boardName)
		}
	}

}

func (cc categoryCleaner) CleanAccountSetting(boardName string) {
	subs := myutil.StringSlice(keyword.Subscribers(boardName))
	subs.AppendNonRepeat(author.Subscribers(boardName), false)
	subs.AppendNonRepeat(pushsum.ListSubscribers(boardName), false)

	for _, sub := range subs {
		u := models.User().Find(sub)
		for index, subBoard := range u.Subscribes {
			if strings.EqualFold(subBoard.Board, boardName) {
				u.Subscribes = append(u.Subscribes[:index], u.Subscribes[index+1:]...)
			}
		}
		if err := u.Save(); err != nil {
			log.WithField("category", boardName).WithError(err).Error("Remove Category in User Failed")
		}
	}
}

func (cc categoryCleaner) CleanKeywordAuthorBoard(boardName string) {
	board := models.Board()
	board.Name = boardName

	if err := board.Delete(); err != nil {
		log.WithField("category", boardName).WithError(err).Error("Delete Board Category Failed")
	}

	if err := keyword.Destroy(boardName); err != nil {
		log.WithField("category", boardName).WithError(err).Error("Delete Keyword Category Failed")
	}

	if err := author.Destroy(boardName); err != nil {
		log.WithField("category", boardName).WithError(err).Error("Delete Author Category Failed")
	}
}

func (cc categoryCleaner) CleanPushsumBoard(boardName string) {
	if err := pushsum.Remove(boardName); err != nil {
		log.WithField("category", boardName).WithError(err).Error("Remove Pushsum Category Failed")
	}
	if err := pushsum.Destroy(boardName); err != nil {
		log.WithField("category", boardName).WithError(err).Error("Delete Pushsum Category Failed")
	}
}
