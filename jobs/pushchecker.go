package jobs

import (
	log "github.com/meifamily/logrus"

	"fmt"

	"github.com/liam-lai/ptt-alertor/crawler"
	"github.com/liam-lai/ptt-alertor/models/ptt/article"

	"time"

	user "github.com/liam-lai/ptt-alertor/models/user/redis"
)

const checkPushDuration = 5

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
	return fmt.Sprintf("推文@%s\n%s\n%s\n%s", pc.Article.Board, pc.Article.Title, pc.Article.Link, pc.Article.PushList.String())
}

// Run start run job
func (pc PushChecker) Run() {

	ach := make(chan article.Article)
	pch := make(chan PushChecker)

	codes := new(article.Articles).List()
	for _, code := range codes {
		go checkPushList(code, ach)
		time.Sleep(checkBoardDuration * time.Second)
	}

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
	new := crawler.BuildArticle(a.Board, a.Code)
	newPushList := make([]article.Push, 0)
	if new.LastPushDateTime.After(a.LastPushDateTime) {
		for _, push := range new.PushList {
			if push.DateTime.After(a.LastPushDateTime) {
				newPushList = append(newPushList, push)
			}
		}
		a.LastPushDateTime = new.LastPushDateTime
		a.Save()
		a.PushList = newPushList
		c <- a
	}
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
