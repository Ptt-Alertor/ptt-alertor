package jobs

import (
	"fmt"
	"strings"
	"time"

	log "github.com/meifamily/logrus"

	"github.com/liam-lai/ptt-alertor/models/ptt/article"
	board "github.com/liam-lai/ptt-alertor/models/ptt/board/redis"
	user "github.com/liam-lai/ptt-alertor/models/user/redis"
)

const checkBoardDuration = 150

type Checker struct {
	email      string
	line       string
	lineNotify string
	messenger  string
	board      string
	keyword    string
	author     string
	articles   article.Articles
	subType    string
	word       string
}

func (cker Checker) String() string {
	subType := "關鍵字"
	if cker.author != "" {
		subType = "作者"
	}
	return fmt.Sprintf("%s@%s\r\n看板：%s；%s：%s%s", cker.word, cker.board, cker.board, subType, cker.word, cker.articles.String())
}

func (cker Checker) Self() Checker {
	return cker
}

func (cker Checker) Run() {
	boardCh := make(chan *board.Board)
	go func() {
		for {
			bds := new(board.Board).All()
			for _, bd := range bds {
				time.Sleep(checkBoardDuration * time.Millisecond)
				go checkNewArticle(bd, boardCh)
			}
		}
	}()
	ckerCh := make(chan Checker)

	for {
		select {
		case bd := <-boardCh:
			checkSubscriber(bd, cker, ckerCh)
		case cker := <-ckerCh:
			cker.subType = "keyword"
			cker.word = cker.keyword
			if cker.author != "" {
				cker.subType = "author"
				cker.word = cker.author
			}
			go sendMessage(cker)
		}
	}
}

func checkNewArticle(bd *board.Board, boardCh chan *board.Board) {
	bd.WithNewArticles()
	if bd.NewArticles == nil {
		bd.Articles = bd.OnlineArticles
		log.WithField("board", bd.Name).Info("Created Articles")
		bd.Save()
	}
	if len(bd.NewArticles) != 0 {
		bd.Articles = bd.OnlineArticles
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
