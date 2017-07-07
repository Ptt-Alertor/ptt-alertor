package jobs

import (
	"fmt"
	"strings"
	"time"

	log "github.com/meifamily/logrus"

	"github.com/meifamily/ptt-alertor/models/ptt/article"
	board "github.com/meifamily/ptt-alertor/models/ptt/board/redis"
	user "github.com/meifamily/ptt-alertor/models/user/redis"
	"github.com/meifamily/ptt-alertor/myutil"
)

const checkBoardDuration = 200 * time.Millisecond
const checkHighBoardDuration = 1 * time.Second
const workers = 250

var highBoards []*board.Board
var boardCh = make(chan *board.Board)

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

// Self return Checker itself
func (cker Checker) Self() Checker {
	return cker
}

// Run is main in Job
func (cker Checker) Run() {
	initHighBoards()
	go func() {
		for {
			checkBoards(highBoards, checkHighBoardDuration)
		}
	}()
	go func() {
		for {
			checkBoards(new(board.Board).All(), checkBoardDuration)
		}
	}()
	ckerCh := make(chan Checker)

	for i := 0; i < workers; i++ {
		go messageWorker(ckerCh)
	}

	for {
		select {
		case bd := <-boardCh:
			checkSubscriber(bd, cker, ckerCh)
		}
	}
}

func initHighBoards() {
	boardcfg := myutil.Config("board")
	highBoardNames := strings.Split(boardcfg["high"], ",")
	for _, name := range highBoardNames {
		bd := new(board.Board)
		bd.Name = name
		highBoards = append(highBoards, bd)
	}
}

func checkBoards(bds []*board.Board, duration time.Duration) {
	for _, bd := range bds {
		time.Sleep(duration)
		go checkNewArticle(bd, boardCh)
	}
}

func messageWorker(ckerCh chan Checker) {
	for {
		cker := <-ckerCh
		cker.subType = "keyword"
		cker.word = cker.keyword
		if cker.author != "" {
			cker.subType = "author"
			cker.word = cker.author
		}
		sendMessage(cker)
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
		err := bd.Save()
		if err == nil {
			boardCh <- bd
		}
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
