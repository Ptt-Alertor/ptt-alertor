package jobs

import (
	"sync"
	"time"

	log "github.com/meifamily/logrus"

	"fmt"

	"github.com/meifamily/ptt-alertor/crawler"
	"github.com/meifamily/ptt-alertor/models/article"
	"github.com/meifamily/ptt-alertor/models/user"
)

// TODO: rename to commentList

var plcker *pushListChecker
var plcOnce sync.Once

// pushListChecker embedding Checker for checking pushlist
type pushListChecker struct {
	Checker
	Article article.Article
	ch      chan pushListChecker
}

// NewPushListChecker return Empty PushChecker pointer
func NewPushListChecker() *pushListChecker {
	plcOnce.Do(func() {
		plcker = &pushListChecker{}
		plcker.duration = 1 * time.Second
		plcker.done = make(chan struct{})
		plcker.ch = make(chan pushListChecker)
	})
	return plcker
}

func (plc pushListChecker) String() string {
	return fmt.Sprintf("推文@%s\n\n%s\n%s\n%s", plc.Article.Board, plc.Article.Title, plc.Article.Link, plc.Article.PushList.String())
}

func (plc pushListChecker) Stop() {
	plc.done <- struct{}{}
	plc.done <- struct{}{}
	log.Info("Pushlist Checker Stop")
}

// Run start job
func (plc pushListChecker) Run() {

	ach := make(chan article.Article)

	go func() {
		for {
			select {
			case <-plc.done:
				return
			default:
				codes := new(article.Articles).List()
				for _, code := range codes {
					time.Sleep(plc.duration)
					go plc.checkPushList(code, ach)
				}
			}
		}
	}()

	for {
		select {
		case a := <-ach:
			plc.Article = a
			plc.checkSubscribers()
		case pc := <-plc.ch:
			ckCh <- pc
		case <-plc.done:
			return
		}
	}
}

func (plc pushListChecker) checkPushList(code string, c chan article.Article) {
	a := new(article.Article).Find(code)
	new, err := crawler.BuildArticle(a.Board, a.Code)
	if _, ok := err.(crawler.URLNotFoundError); ok {
		plc.destroyPushList(a)
	}
	if subs, _ := a.Subscribers(); len(subs) == 0 {
		plc.destroyPushList(a)
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

func (plc pushListChecker) destroyPushList(a article.Article) {
	a.Destroy()
	log.WithFields(log.Fields{
		"board": a.Board,
		"code":  a.Code,
	}).Info("Destroy PushList")
}

func (plc pushListChecker) checkSubscribers() {
	subs, err := plc.Article.Subscribers()
	if err != nil {
		log.WithError(err).Error("Get Subscribers Failed")
	}

	for _, account := range subs {
		go send(account, plc, plc.ch)
	}
}

func send(account string, plc pushListChecker, pch chan pushListChecker) {
	u := user.NewUser().Find(account)
	plc.board = plc.Article.Board
	plc.subType = "push"
	plc.word = plc.Article.Code
	plc.Profile = u.Profile
	pch <- plc
}
