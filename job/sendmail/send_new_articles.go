package main

import (
	"fmt"

	"github.com/liam-lai/ptt-alertor/mail"
	"github.com/liam-lai/ptt-alertor/myutil"
	"github.com/liam-lai/ptt-alertor/ptt/article"
	"github.com/liam-lai/ptt-alertor/ptt/board"
	"github.com/liam-lai/ptt-alertor/user"
)

type articles []article.Article

var storageDir string = myutil.StoragePath()

func main() {
	bs := new(board.Boards).All().WithNewArticles(true)
	users := new(user.Users).All()
	for _, bd := range bs {
		for _, user := range users {
			for _, subscribe := range user.Subscribes {
				if bd.Name == subscribe.Board {
					for _, keyword := range subscribe.Keywords {
						keywordArticles := bd.NewArticles.ContainKeyword(keyword)
						if len(keywordArticles) != 0 {
							fmt.Println(user.Profile.Email + ":" + keyword + " in " + subscribe.Board)
							sendMail(user, subscribe.Board, keyword, keywordArticles)
						}
					}
				}
			}
		}
	}
}

func sendMail(user *user.User, board string, keyword string, articles []article.Article) {
	m := new(mail.Mail)
	m.Title.BoardName = board
	m.Title.Keyword = keyword
	m.Body.Articles = articles
	m.Receiver = user.Profile.Email

	m.Send()
}
