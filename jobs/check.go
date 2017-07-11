package jobs

import (
	log "github.com/meifamily/logrus"

	"github.com/meifamily/ptt-alertor/line"
	"github.com/meifamily/ptt-alertor/mail"
	"github.com/meifamily/ptt-alertor/messenger"
)

const workers = 250

var ckCh = make(chan check)

func init() {
	for i := 0; i < workers; i++ {
		go messageWorker(ckCh)
	}
}

func messageWorker(ckCh chan check) {
	for {
		ck := <-ckCh
		sendMessage(ck)
	}
}

type check interface {
	String() string
	Self() Checker
}

func sendMessage(c check) {
	cr := c.Self()
	var account string
	var platform string

	if cr.email != "" {
		account = cr.email
		platform = "mail"
		sendMail(c)
	}
	if cr.lineNotify != "" {
		account = cr.line
		platform = "line"
		sendLineNotify(c)
	}
	if cr.messenger != "" {
		account = cr.messenger
		platform = "messenger"
		sendMessenger(c)
	}
	log.WithFields(log.Fields{
		"account":  account,
		"platform": platform,
		"board":    cr.board,
		"type":     cr.subType,
		"word":     cr.word,
	}).Info("Message Sent")
}

func sendMail(c check) {
	cr := c.Self()
	m := new(mail.Mail)
	m.Title.BoardName = cr.board
	m.Title.Keyword = cr.keyword
	m.Body.Articles = cr.articles
	m.Receiver = cr.email
	m.Send()
}

func sendLine(c check) {
	cr := c.Self()
	line.PushTextMessage(cr.line, c.String())
}

func sendLineNotify(c check) {
	cr := c.Self()
	line.Notify(cr.lineNotify, c.String())
}

func sendMessenger(c check) {
	cr := c.Self()
	m := messenger.New()
	m.SendTextMessage(cr.messenger, c.String())
}
