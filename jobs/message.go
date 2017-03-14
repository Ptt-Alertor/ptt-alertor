package jobs

import (
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/liam-lai/ptt-alertor/mail"
	"github.com/liam-lai/ptt-alertor/models/ptt/article"
	board "github.com/liam-lai/ptt-alertor/models/ptt/board/redis"
	user "github.com/liam-lai/ptt-alertor/models/user/redis"
)

type Message struct {
	email    string
	board    string
	keyword  string
	articles []article.Article
}

func (msg Message) Run() {
	bds := new(board.Board).All()
	for _, bd := range bds {
		bd.WithNewArticles()
	}
	bds = deleteNonNewArticleBoard(bds)
	users := new(user.User).All()
	msgCh := make(chan Message)
	for _, user := range users {
		msg.email = user.Profile.Email
		log.WithField("user", user.Profile.Account).Info("Checking User Subscribes")
		go userChecker(user, bds, msg, msgCh)
	}

	for {
		select {
		case m := <-msgCh:
			sendMail(m)
		case <-time.After(time.Second * 3):
			log.Info("Message Done")
			return
		}
	}
}

func deleteNonNewArticleBoard(bds []*board.Board) []*board.Board {
	for index, bd := range bds {
		if len(bd.NewArticles) == 0 {
			if index < len(bds)-1 {
				bds = append(bds[:index], bds[index+1:]...)
			} else {
				bds = bds[:index]
			}
		}
	}
	return bds
}

func userChecker(user *user.User, bds []*board.Board, msg Message, msgCh chan Message) {
	for _, bd := range bds {
		go subscribeChecker(user, bd, msg, msgCh)
	}
}

func subscribeChecker(user *user.User, bd *board.Board, msg Message, msgCh chan Message) {
	for _, sub := range user.Subscribes {
		if bd.Name == sub.Board {
			msg.board = sub.Board
			for _, keyword := range sub.Keywords {
				go keywordChecker(keyword, bd, msg, msgCh)
			}
		}
	}
}

func keywordChecker(keyword string, bd *board.Board, msg Message, msgCh chan Message) {
	keywordArticles := make([]article.Article, 0)
	for _, newAtcl := range bd.NewArticles {
		if newAtcl.ContainKeyword(keyword) {
			keywordArticles = append(keywordArticles, newAtcl)
		}
	}
	if len(keywordArticles) != 0 {
		msg.keyword = keyword
		msg.articles = keywordArticles
		msgCh <- msg
	}
}

func sendMail(msg Message) {
	m := new(mail.Mail)
	m.Title.BoardName = msg.board
	m.Title.Keyword = msg.keyword
	m.Body.Articles = msg.articles
	m.Receiver = msg.email

	m.Send()
}
