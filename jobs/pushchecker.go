package jobs

import (
	"time"

	log "github.com/meifamily/logrus"

	"fmt"

	"github.com/meifamily/ptt-alertor/crawler"
	"github.com/meifamily/ptt-alertor/models/ptt/article"

	user "github.com/meifamily/ptt-alertor/models/user/redis"
)

const checkPushDuration = 3

// PushChecker embedding Checker for checking pushlist
type PushChecker struct {
	Checker
	Article article.Article
}

// NewPushChecker return Empty PushChecker pointer
func NewPushChecker() *PushChecker {
	return &PushChecker{}
}

func (pc PushChecker) String() string {
	return fmt.Sprintf("推文@%s\n\n%s\n%s\n%s", pc.Article.Board, pc.Article.Title, pc.Article.Link, pc.Article.PushList.String())
}

// Run start job
func (pc PushChecker) Run() {

	ach := make(chan article.Article)
	pch := make(chan PushChecker)

	go func() {
		for {
			codes := new(article.Articles).List()
			for _, code := range codes {
				time.Sleep(checkPushDuration * time.Second)
				go checkPushList(code, ach)
			}
		}
	}()

	for {
		select {
		case a := <-ach:
			pc.Article = a
			pc.checkSubscribers(pch)
		case pc := <-pch:
			sendMessage(pc)
		}
	}

}

func checkPushList(code string, c chan article.Article) {
	a := new(article.Article).Find(code)
	new, err := crawler.BuildArticle(a.Board, a.Code)
	if _, ok := err.(crawler.URLNotFoundError); ok {
		destroyPushList(a)
	}
	if subs, _ := a.Subscribers(); len(subs) == 0 {
		destroyPushList(a)
	}
	newPushList := make([]article.Push, 0)
	if new.LastPushDateTime.After(a.LastPushDateTime) {
		for _, push := range new.PushList {
			if push.DateTime.After(a.LastPushDateTime) {
				newPushList = append(newPushList, push)
			}
		}
		a.LastPushDateTime = new.LastPushDateTime
		a.PushList = newPushList
		a.Save()
		log.WithFields(log.Fields{
			"board": a.Board,
			"code":  a.Code,
		}).Info("Updated PushList")
		c <- a
	}
}

func destroyPushList(a article.Article) {
	a.Destroy()
	log.WithFields(log.Fields{
		"board": a.Board,
		"code":  a.Code,
	}).Info("Destroy PushList")
}

func (pc PushChecker) checkSubscribers(pch chan PushChecker) {
	subs, err := pc.Article.Subscribers()
	if err != nil {
		log.WithError(err).Error("Get Subscribers Failed")
	}

	for _, account := range subs {
		go send(account, pc, pch)
	}
}

func send(account string, pc PushChecker, pch chan PushChecker) {
	u := user.User{}
	u = u.Find(account)
	pc.board = pc.Article.Board
	pc.subType = "push"
	pc.word = pc.Article.Code
	pc.email = u.Profile.Email
	pc.line = u.Profile.Line
	pc.lineNotify = u.Profile.LineAccessToken
	pc.messenger = u.Profile.Messenger
	pch <- pc
}
