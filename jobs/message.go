package jobs

import (
	"time"

	log "github.com/Sirupsen/logrus"

	"strings"

	"github.com/liam-lai/ptt-alertor/line"
	"github.com/liam-lai/ptt-alertor/mail"
	"github.com/liam-lai/ptt-alertor/messenger"
	"github.com/liam-lai/ptt-alertor/models/ptt/article"
	board "github.com/liam-lai/ptt-alertor/models/ptt/board/redis"
	user "github.com/liam-lai/ptt-alertor/models/user/redis"
)

type Message struct {
	email     string
	line      string
	messenger string
	board     string
	keyword   string
	author    string
	articles  article.Articles
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
		if user.Enable {
			msg.email = user.Profile.Email
			msg.line = user.Profile.Line
			msg.messenger = user.Profile.Messenger
			log.WithField("user", user.Profile.Account).Info("Checking User Subscribes")
			go userChecker(user, bds, msg, msgCh)
		}
	}

	for {
		select {
		case m := <-msgCh:
			sendMessage(m)
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
			for _, author := range sub.Authors {
				go authorChecker(author, bd, msg, msgCh)
			}
		}
	}
}

func keywordChecker(keyword string, bd *board.Board, msg Message, msgCh chan Message) {
	keywordArticles := make(article.Articles, 0)
	for _, newAtcl := range bd.NewArticles {
		if newAtcl.ContainKeyword(keyword) {
			newAtcl.Author = ""
			keywordArticles = append(keywordArticles, newAtcl)
		}
	}
	if len(keywordArticles) != 0 {
		msg.keyword = keyword
		msg.articles = keywordArticles
		msgCh <- msg
	}
}

func authorChecker(author string, bd *board.Board, msg Message, msgCh chan Message) {
	authorArticles := make(article.Articles, 0)
	for _, newAtcl := range bd.NewArticles {
		if strings.EqualFold(newAtcl.Author, author) {
			authorArticles = append(authorArticles, newAtcl)
		}
	}
	if len(authorArticles) != 0 {
		msg.author = author
		msg.articles = authorArticles
		msgCh <- msg
	}

}

func sendMessage(msg Message) {
	var account string
	if msg.email != "" {
		account = msg.email
		sendMail(msg)
	}
	if msg.line != "" {
		account = msg.line
		sendLine(msg)
	}
	if msg.messenger != "" {
		account = msg.messenger
		sendMessenger(msg)
	}
	log.WithFields(log.Fields{
		"account": account,
		"board":   msg.board,
		"keyword": msg.keyword,
	}).Info("Message Sent")
}

func sendMail(msg Message) {
	m := new(mail.Mail)
	m.Title.BoardName = msg.board
	m.Title.Keyword = msg.keyword
	m.Body.Articles = msg.articles
	m.Receiver = msg.email

	m.Send()
}

func sendLine(msg Message) {
	line.PushTextMessage(msg.line, msg.articles.String())
}

func sendMessenger(msg Message) {
	m := messenger.New()
	m.SendTextMessage(msg.messenger, msg.articles.String())
}
