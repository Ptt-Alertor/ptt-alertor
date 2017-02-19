package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/liam-lai/ptt-alertor/mail"
	"github.com/liam-lai/ptt-alertor/myutil"
	"github.com/liam-lai/ptt-alertor/pttboard"
)

type User struct {
	Profile struct {
		Email string
	}
	Subscribes
}

type Subscribes []Subscribe

type Subscribe struct {
	Board    string
	Keywords []string
}

type articles mail.Articles

var storageDir string = myutil.ProjectRootPath() + "/storage"

func main() {

	boards := boardsWithNewArticles()

	usersDir := storageDir + "/users/"
	files, _ := ioutil.ReadDir(usersDir)
	for _, file := range files {
		_, ok := jsonFile(file)
		if !ok {
			continue
		}
		userJSON, _ := ioutil.ReadFile(usersDir + file.Name())
		var user User
		_ = json.Unmarshal(userJSON, &user)
		for _, subscribe := range user.Subscribes {
			if articles, ok := boards[subscribe.Board]; ok {
				for _, keyword := range subscribe.Keywords {
					keywordArticles := articlesHaveKeyword(keyword, articles)
					if len(keywordArticles) != 0 {
						fmt.Println(user.Profile.Email + ":" + keyword + " in " + subscribe.Board)
						sendMail(user, subscribe.Board, keyword, articles)
					}
				}
			}
		}
	}

}

func articlesHaveKeyword(keyword string, newArticles mail.Articles) mail.Articles {
	articles := make(mail.Articles, 0)
	for _, article := range newArticles {
		if strings.Contains(article.Title, keyword) {
			articles = append(articles, article)
		}
	}
	return articles
}

func sendMail(user User, board string, keyword string, articles mail.Articles) {
	m := new(mail.Mail)
	m.Title.BoardName = board
	m.Title.Keyword = keyword
	m.Body.Articles = articles
	m.Receiver = user.Profile.Email

	m.Send()
}

func boardsWithNewArticles() (boards map[string]mail.Articles) {
	articlesDir := storageDir + "/articles/"
	files, _ := ioutil.ReadDir(articlesDir)
	boards = make(map[string]mail.Articles)
	for _, file := range files {
		boardName, ok := jsonFile(file)
		if !ok {
			continue
		}
		newArticles := newArticles(articlesDir, boardName)
		var articles mail.Articles
		json.Unmarshal(newArticles, &articles)
		boards[boardName] = append(boards[boardName], articles...)
	}
	return boards
}

func newArticles(dir string, BoardName string) []byte {

	oldArticlesJSON, err := ioutil.ReadFile(dir + BoardName + ".json")
	if err != nil {
		log.Fatal(err)
	}
	nowArticlesJSON := pttboard.Index(BoardName)

	newArticlesJSON := myutil.DifferenceJSON(oldArticlesJSON, nowArticlesJSON)

	return newArticlesJSON
}

func jsonFile(file os.FileInfo) (string, bool) {
	if file.IsDir() {
		return "", false
	}

	fileName, extension := myutil.FileNameAndExtension(file.Name())
	if extension != "json" {
		return "", false
	}
	return fileName, true
}
