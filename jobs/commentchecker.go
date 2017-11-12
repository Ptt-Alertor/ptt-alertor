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

var cmtcker *commentChecker
var cmtOnce sync.Once

// commentChecker embedding Checker for checking comment
type commentChecker struct {
	Checker
	Article article.Article
	ch      chan commentChecker
}

// NewCommentChecker return Empty PushChecker pointer
func NewCommentChecker() *commentChecker {
	cmtOnce.Do(func() {
		cmtcker = &commentChecker{}
		cmtcker.duration = 1 * time.Second
		cmtcker.done = make(chan struct{})
		cmtcker.ch = make(chan commentChecker)
	})
	return cmtcker
}

func (cc commentChecker) String() string {
	return fmt.Sprintf("推文@%s\n\n%s\n%s\n%s", cc.Article.Board, cc.Article.Title, cc.Article.Link, cc.Article.Comments.String())
}

func (cc commentChecker) Stop() {
	cc.done <- struct{}{}
	cc.done <- struct{}{}
	log.Info("Comment Checker Stop")
}

// Run start job
func (cc commentChecker) Run() {

	ach := make(chan article.Article)

	go func() {
		for {
			select {
			case <-cc.done:
				return
			default:
				codes := new(article.Articles).List()
				for _, code := range codes {
					time.Sleep(cc.duration)
					go cc.checkComments(code, ach)
				}
			}
		}
	}()

	for {
		select {
		case a := <-ach:
			cc.Article = a
			cc.checkSubscribers()
		case pc := <-cc.ch:
			ckCh <- pc
		case <-cc.done:
			return
		}
	}
}

func (cc commentChecker) checkComments(code string, c chan article.Article) {
	a := new(article.Article).Find(code)
	new, err := crawler.BuildArticle(a.Board, a.Code)
	if _, ok := err.(crawler.URLNotFoundError); ok {
		cc.destroyComments(a)
	}
	if subs, _ := a.Subscribers(); len(subs) == 0 {
		cc.destroyComments(a)
	}
	newComments := make([]article.Comment, 0)
	if new.LastPushDateTime.After(a.LastPushDateTime) {
		for _, push := range new.Comments {
			if push.DateTime.After(a.LastPushDateTime) {
				newComments = append(newComments, push)
			}
		}
		a.LastPushDateTime = new.LastPushDateTime
		a.Comments = newComments
		a.Save()
		log.WithFields(log.Fields{
			"board": a.Board,
			"code":  a.Code,
		}).Info("Updated Comments")
		c <- a
	}
}

func (cc commentChecker) destroyComments(a article.Article) {
	a.Destroy()
	log.WithFields(log.Fields{
		"board": a.Board,
		"code":  a.Code,
	}).Info("Destroy Comments")
}

func (cc commentChecker) checkSubscribers() {
	subs, err := cc.Article.Subscribers()
	if err != nil {
		log.WithError(err).Error("Get Subscribers Failed")
	}

	for _, account := range subs {
		go send(account, cc, cc.ch)
	}
}

func send(account string, cc commentChecker, pch chan commentChecker) {
	u := user.NewUser().Find(account)
	cc.board = cc.Article.Board
	cc.subType = "push"
	cc.word = cc.Article.Code
	cc.Profile = u.Profile
	pch <- cc
}
