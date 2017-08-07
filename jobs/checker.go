package jobs

import (
	"fmt"
	"strings"
	"sync"
	"time"

	log "github.com/meifamily/logrus"

	"github.com/meifamily/ptt-alertor/models/author"
	"github.com/meifamily/ptt-alertor/models/keyword"
	"github.com/meifamily/ptt-alertor/models/ptt/article"
	board "github.com/meifamily/ptt-alertor/models/ptt/board/redis"
	userProto "github.com/meifamily/ptt-alertor/models/user"
	user "github.com/meifamily/ptt-alertor/models/user/redis"
	"github.com/meifamily/ptt-alertor/myutil"
)

const checkHighBoardDuration = 1 * time.Second

var boardCh = make(chan *board.Board, 700)
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

var cker *Checker
var ckerOnce sync.Once

type Checker struct {
	board    string
	keyword  string
	author   string
	articles article.Articles
	subType  string
	word     string
	Profile  userProto.Profile
	done     chan struct{}
	ch       chan Checker
	duration time.Duration
}

// NewChecker gets a Checker instance
func NewChecker() *Checker {
	ckerOnce.Do(func() {
		cker = &Checker{
			duration: 200 * time.Millisecond,
		}
		cker.done = make(chan struct{})
		cker.ch = make(chan Checker)
	})
	return cker
}

func (c Checker) String() string {
	subType := "關鍵字"
	if c.author != "" {
		subType = "作者"
	}
	return fmt.Sprintf("%s@%s\r\n看板：%s；%s：%s%s", c.word, c.board, c.board, subType, c.word, c.articles.String())
}

// Self return Checker itself
func (c Checker) Self() Checker {
	return c
}

// Run is main in Job
func (c Checker) Run() {
	// step 1: check boards which one has new articles
	c.runCheckBoards()

	for {
		select {
		//step 2: check user who subscribes board
		case bd := <-boardCh:
			go checkKeywordSubscriber(bd, c)
			go checkAuthorSubscriber(bd, c)
		//step 3: send notification
		case cker := <-c.ch:
			ckCh <- cker
		case <-c.done:
			return
		}
	}
}

func (c Checker) runCheckBoards() {
	go func() {
		for {
			select {
			case <-c.done:
				return
			default:
				checkBoards(highBoards, checkHighBoardDuration)
			}
		}
	}()
	offPeakCh := make(chan bool)
	go func(offPeakCh <-chan bool) {
		var offPeak bool
		duration := c.duration
		for {
			select {
			case op := <-offPeakCh:
				if offPeak != op {
					if op {
						log.Info("Switch to Slow Mode")
						duration = c.duration * 2
					} else {
						log.Info("Switch to Normal Mode")
						duration = c.duration
					}
					offPeak = op
				}
			case <-c.done:
				return
			default:
				checkBoards(new(board.Board).All(), duration)
			}
		}
	}(offPeakCh)

	// check off peak
	go c.checkOffPeak(offPeakCh)
}

func (c Checker) checkOffPeak(offPeakCh chan<- bool) {
	loc := time.FixedZone("CST", 8*60*60)
	for {
		t := time.Now().In(loc)
		if t.Hour() >= 3 && t.Hour() < 7 {
			offPeakCh <- true
		} else {
			offPeakCh <- false
		}
		time.Sleep(10 * time.Minute)
	}
}

func (c Checker) Stop() {
	for i := 0; i < 3; i++ {
		c.done <- struct{}{}
	}
	log.Info("Checker Stop")
}

func checkBoards(bds []*board.Board, duration time.Duration) {
	for _, bd := range bds {
		time.Sleep(duration)
		go checkNewArticle(bd, boardCh)
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

func checkKeywordSubscriber(bd *board.Board, cker Checker) {
	u := new(user.User)
	accounts := keyword.Subscribers(bd.Name)
	for _, account := range accounts {
		user := u.Find(account)
		if user.Enable {
			cker.Profile = user.Profile
			go checkKeywordSubscription(user, bd, cker)
		}
	}
}

func checkKeywordSubscription(user user.User, bd *board.Board, cker Checker) {
	for _, sub := range user.Subscribes {
		if bd.Name == sub.Board {
			cker.board = sub.Board
			for _, keyword := range sub.Keywords {
				go checkKeyword(keyword, bd, cker)
			}
		}
	}
}

func checkKeyword(keyword string, bd *board.Board, cker Checker) {
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
		cker.ch <- cker
	}
}

func checkAuthorSubscriber(bd *board.Board, cker Checker) {
	u := new(user.User)
	accounts := author.Subscribers(bd.Name)
	for _, account := range accounts {
		user := u.Find(account)
		if user.Enable {
			cker.Profile = user.Profile
			go checkAuthorSubscription(user, bd, cker)
		}
	}
}

func checkAuthorSubscription(user user.User, bd *board.Board, cker Checker) {
	for _, sub := range user.Subscribes {
		if bd.Name == sub.Board {
			cker.board = sub.Board
			for _, author := range sub.Authors {
				go checkAuthor(author, bd, cker)
			}
		}
	}
}

func checkAuthor(author string, bd *board.Board, cker Checker) {
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
		cker.ch <- cker
	}
}
