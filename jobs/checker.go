package jobs

import (
	"fmt"
	"strings"
	"sync"
	"time"

	log "github.com/meifamily/logrus"

	"github.com/meifamily/ptt-alertor/models/ptt/article"
	board "github.com/meifamily/ptt-alertor/models/ptt/board/redis"
	userProto "github.com/meifamily/ptt-alertor/models/user"
	user "github.com/meifamily/ptt-alertor/models/user/redis"
	"github.com/meifamily/ptt-alertor/myutil"
)

const checkBoardDuration = 200 * time.Millisecond
const checkHighBoardDuration = 1 * time.Second

var boardCh = make(chan *board.Board)
var ckerCh = make(chan Checker)
var highBoards []*board.Board

func init() {
	initHighBoards()
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

type Checker struct {
	board    string
	keyword  string
	author   string
	articles article.Articles
	subType  string
	word     string
	Profile  userProto.Profile
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
	var wgHigh sync.WaitGroup
	var wg sync.WaitGroup
	// step 1: check boards which one has new articles
	go func() {
		for {
			checkBoards(&wgHigh, highBoards, checkHighBoardDuration)
			wgHigh.Wait()
		}
	}()
	go func() {
		for {
			checkBoards(&wg, new(board.Board).All(), checkBoardDuration)
			wg.Wait()
		}
	}()

	for {
		select {
		//step 2: check user who subscribes board
		case bd := <-boardCh:
			go checkSubscriber(bd, cker)
		//step 3: send notification
		case cker := <-ckerCh:
			ckCh <- cker
		}
	}
}

func checkBoards(wg *sync.WaitGroup, bds []*board.Board, duration time.Duration) {
	wg.Add(len(bds))
	for _, bd := range bds {
		time.Sleep(duration)
		go checkNewArticle(wg, bd, boardCh)
	}
}

func checkNewArticle(wg *sync.WaitGroup, bd *board.Board, boardCh chan *board.Board) {
	defer wg.Done()
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

func checkSubscriber(bd *board.Board, cker Checker) {
	users := new(user.User).All()
	for _, user := range users {
		if user.Enable {
			cker.Profile = user.Profile
			go subscribeChecker(user, bd, cker)
		}
	}
}

func subscribeChecker(user *user.User, bd *board.Board, cker Checker) {
	for _, sub := range user.Subscribes {
		if bd.Name == sub.Board {
			cker.board = sub.Board
			for _, keyword := range sub.Keywords {
				go keywordChecker(keyword, bd, cker)
			}
			for _, author := range sub.Authors {
				go authorChecker(author, bd, cker)
			}
		}
	}
}

func keywordChecker(keyword string, bd *board.Board, cker Checker) {
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
		cker.subType = "keyword"
		cker.word = keyword
		ckerCh <- cker
	}
}

func authorChecker(author string, bd *board.Board, cker Checker) {
	authorArticles := make(article.Articles, 0)
	for _, newAtcl := range bd.NewArticles {
		if strings.EqualFold(newAtcl.Author, author) {
			authorArticles = append(authorArticles, newAtcl)
		}
	}
	if len(authorArticles) != 0 {
		cker.author = author
		cker.articles = authorArticles
		cker.subType = "author"
		cker.word = author
		ckerCh <- cker
	}
}
