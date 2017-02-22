package main

import (
	"fmt"

	"sync"

	"github.com/liam-lai/ptt-alertor/mail"
	"github.com/liam-lai/ptt-alertor/ptt/article"
	"github.com/liam-lai/ptt-alertor/ptt/board"
	"github.com/liam-lai/ptt-alertor/user"
)

type articles []article.Article

var wg sync.WaitGroup

func main() {
	bs := new(board.Boards).All().WithNewArticles(true)
	users := new(user.Users).All()
	for _, bd := range bs {
		wg.Add(1)
		go checkUserSubscribeBoard(users, bd)
	}
	wg.Wait()
}

func checkUserSubscribeBoard(users user.Users, bd *board.Board) {
	defer wg.Done()
	for _, u := range users {
		wg.Add(1)
		go func(u *user.User) {
			defer wg.Done()
			for _, subscribe := range u.Subscribes {
				if bd.Name == subscribe.Board {
					for _, keyword := range subscribe.Keywords {
						keywordArticles := bd.NewArticles.ContainKeyword(keyword)
						if len(keywordArticles) != 0 {
							fmt.Println(u.Profile.Email + ":" + keyword + " in " + subscribe.Board)
							//sendMail(u, subscribe.Board, keyword, keywordArticles)
						}

					}
				}
			}
		}(u)
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
