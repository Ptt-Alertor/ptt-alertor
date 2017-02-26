package main

import (
	"fmt"

	"time"

	"github.com/liam-lai/ptt-alertor/mail"
	"github.com/liam-lai/ptt-alertor/ptt/article"
	"github.com/liam-lai/ptt-alertor/ptt/board"
	"github.com/liam-lai/ptt-alertor/user"
)

type Message struct {
	user     *user.User
	board    string
	keyword  string
	articles []article.Article
}

func main() {
	bs := new(board.Boards).All().WithNewArticles(true)
	users := new(user.Users).All()
	msgCh := make(chan Message)
	for _, user := range users {
		go userChecker(user, bs, msgCh)
	}

	for {
		select {
		case m := <-msgCh:
			sendMail(m)
		case <-time.After(time.Second * 3):
			fmt.Println("time out")
			return
		}
	}
}

func userChecker(user *user.User, bds board.Boards, msgCh chan Message) {
	msg := Message{
		user: user,
	}
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
		fmt.Printf("%+v", msg)
		msgCh <- msg
	}
}

func sendMail(msg Message) {
	m := new(mail.Mail)
	m.Title.BoardName = msg.board
	m.Title.Keyword = msg.keyword
	m.Body.Articles = msg.articles
	m.Receiver = msg.user.Profile.Email

	m.Send()
}
