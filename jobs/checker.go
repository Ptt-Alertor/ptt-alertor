package jobs

import (
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/liam-lai/ptt-alertor/line"
	"github.com/liam-lai/ptt-alertor/mail"
	"github.com/liam-lai/ptt-alertor/messenger"
	"github.com/liam-lai/ptt-alertor/models/ptt/article"
	board "github.com/liam-lai/ptt-alertor/models/ptt/board/redis"
	user "github.com/liam-lai/ptt-alertor/models/user/redis"
)

const checkBoardDuration = 30

type Checker struct {
	email      string
	line       string
	lineNotify string
	messenger  string
	board      string
	keyword    string
	author     string
	articles   article.Articles
}

func (cker Checker) String() string {
	subType := "關鍵字"
	subText := cker.keyword
	if cker.author != "" {
		subType = "作者"
		subText = cker.author
	}
	return fmt.Sprintf("\r\n看版：%s；%s：%s%s", cker.board, subType, subText, cker.articles.String())
}

func (cker Checker) Run() {
	boardCh := make(chan *board.Board)
	go func() {
		for {
			bds := new(board.Board).All()
			for _, bd := range bds[:len(bds)/2] {
				go checkNewArticle(bd, boardCh)
			}
			time.Sleep(checkBoardDuration * time.Second)
			for _, bd := range bds[len(bds)/2:] {
				go checkNewArticle(bd, boardCh)
			}
			time.Sleep(checkBoardDuration * time.Second)
		}
	}()
	ckerCh := make(chan Checker)

	for {
		select {
		case bd := <-boardCh:
			checkSubscriber(bd, cker, ckerCh)
		case m := <-ckerCh:
			sendMessage(m)
		}
	}
}

func checkNewArticle(bd *board.Board, boardCh chan *board.Board) {
	bd.WithNewArticles()
	if len(bd.NewArticles) != 0 {
		bd.Articles = bd.OnlineArticles()
		log.WithField("board", bd.Name).Info("Updated Articles")
		bd.Save()
		boardCh <- bd
	}
}

func checkSubscriber(bd *board.Board, cker Checker, ckerCh chan Checker) {
	users := new(user.User).All()
	for _, user := range users {
		if user.Enable {
			cker.email = user.Profile.Email
			cker.line = user.Profile.Line
			cker.lineNotify = user.Profile.LineAccessToken
			cker.messenger = user.Profile.Messenger
			go subscribeChecker(user, bd, cker, ckerCh)
		}
	}
}

func subscribeChecker(user *user.User, bd *board.Board, cker Checker, ckerCh chan Checker) {
	for _, sub := range user.Subscribes {
		if bd.Name == sub.Board {
			cker.board = sub.Board
			for _, keyword := range sub.Keywords {
				go keywordChecker(keyword, bd, cker, ckerCh)
			}
			for _, author := range sub.Authors {
				go authorChecker(author, bd, cker, ckerCh)
			}
		}
	}
}

func keywordChecker(keyword string, bd *board.Board, cker Checker, ckerCh chan Checker) {
	keywordArticles := make(article.Articles, 0)
	for _, newAtcl := range bd.NewArticles {
		if newAtcl.MatchKeyword(keyword) {
			newAtcl.Author = ""
			keywordArticles = append(keywordArticles, newAtcl)
		}
	}
	if len(keywordArticles) != 0 {
		cker.keyword = keyword
		cker.articles = keywordArticles
		ckerCh <- cker
	}
}

func authorChecker(author string, bd *board.Board, cker Checker, ckerCh chan Checker) {
	authorArticles := make(article.Articles, 0)
	for _, newAtcl := range bd.NewArticles {
		if strings.EqualFold(newAtcl.Author, author) {
			authorArticles = append(authorArticles, newAtcl)
		}
	}
	if len(authorArticles) != 0 {
		cker.author = author
		cker.articles = authorArticles
		ckerCh <- cker
	}

}

func sendMessage(cker Checker) {
	var account string
	var platform string
	subType := "keyword"
	word := cker.keyword

	if cker.author != "" {
		subType = "author"
		word = cker.author
	}

	if cker.email != "" {
		account = cker.email
		platform = "mail"
		sendMail(cker)
	}
	if cker.lineNotify != "" {
		account = cker.line
		platform = "line"
		sendLineNotify(cker)
	}
	if cker.messenger != "" {
		account = cker.messenger
		platform = "messenger"
		sendMessenger(cker)
	}
	log.WithFields(log.Fields{
		"account":  account,
		"platform": platform,
		"board":    cker.board,
		"type":     subType,
		"word":     word,
	}).Info("Message Sent")
}

func sendMail(cker Checker) {
	m := new(mail.Mail)
	m.Title.BoardName = cker.board
	m.Title.Keyword = cker.keyword
	m.Body.Articles = cker.articles
	m.Receiver = cker.email

	m.Send()
}

func sendLine(cker Checker) {
	line.PushTextMessage(cker.line, cker.String())
}

func sendLineNotify(cker Checker) {
	line.Notify(cker.lineNotify, cker.String())
}

func sendMessenger(cker Checker) {
	m := messenger.New()
	m.SendTextMessage(cker.messenger, cker.String())
}
