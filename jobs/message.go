package jobs

import (
	"fmt"

	"time"

	"github.com/liam-lai/ptt-alertor/mail"
	"github.com/liam-lai/ptt-alertor/ptt/article"
	"github.com/liam-lai/ptt-alertor/ptt/board"
	"github.com/liam-lai/ptt-alertor/user"
)

type Message struct {
	email    string
	board    string
	keyword  string
	articles []article.Article
}

func (msg Message) Run() {
	bds := new(board.Boards).All().WithNewArticles(true)
	users := new(user.Users).All()
	msgCh := make(chan Message)
	for _, user := range users {
		msg.email = user.Profile.Email
		go userChecker(user, bds, msg, msgCh)
	}

	for {
		select {
		case m := <-msgCh:
			sendMail(m)
		case <-time.After(time.Second * 3):
			fmt.Println("message done")
			return
		}
	}
}

func userChecker(user *user.User, bds board.Boards, msg Message, msgCh chan Message) {
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
	keywordArticles := bd.NewArticles.ContainKeyword(keyword)
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
