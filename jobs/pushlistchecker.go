package jobs

import (
	"sync"
	"time"

	log "github.com/meifamily/logrus"

	"fmt"

	"github.com/meifamily/ptt-alertor/crawler"
	"github.com/meifamily/ptt-alertor/models/ptt/article"

	user "github.com/meifamily/ptt-alertor/models/user/redis"
)

const checkPushListDuration = 1 * time.Second

var plcker *PushListChecker
var plcOnce sync.Once

// PushListChecker embedding Checker for checking pushlist
type PushListChecker struct {
	Checker
	Article article.Article
}

// NewPushListChecker return Empty PushChecker pointer
func NewPushListChecker() *PushListChecker {
	plcOnce.Do(func() {
		plcker = &PushListChecker{}
		plcker.done = make(chan struct{})
	})
	return plcker
}

func (plc PushListChecker) String() string {
	return fmt.Sprintf("推文@%s\n\n%s\n%s\n%s", plc.Article.Board, plc.Article.Title, plc.Article.Link, plc.Article.PushList.String())
}

func (plc PushListChecker) Stop() {
	plc.done <- struct{}{}
	plc.done <- struct{}{}
	log.Info("Pushlist Checker Stop")
}

// Run start job
func (plc PushListChecker) Run() {

	ach := make(chan article.Article)
	pch := make(chan PushListChecker)

	go func() {
		for {
			select {
			case <-plc.done:
				return
			default:
				codes := new(article.Articles).List()
				for _, code := range codes {
					time.Sleep(checkPushListDuration)
					go checkPushList(code, ach)
				}
			}
		}
	}()

	for {
		select {
		case a := <-ach:
			plc.Article = a
			plc.checkSubscribers(pch)
		case pc := <-pch:
			ckCh <- pc
		case <-plc.done:
			return
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

func (plc PushListChecker) checkSubscribers(pch chan PushListChecker) {
	subs, err := plc.Article.Subscribers()
	if err != nil {
		log.WithError(err).Error("Get Subscribers Failed")
	}

	for _, account := range subs {
		go send(account, plc, pch)
	}
}

func send(account string, plc PushListChecker, pch chan PushListChecker) {
	u := user.User{}
	u = u.Find(account)
	plc.board = plc.Article.Board
	plc.subType = "push"
	plc.word = plc.Article.Code
	plc.Profile = u.Profile
	pch <- plc
}
