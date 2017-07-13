package jobs

import (
	"fmt"
	"time"

	"strings"

	"strconv"

	log "github.com/meifamily/logrus"

	"github.com/meifamily/ptt-alertor/crawler"
	"github.com/meifamily/ptt-alertor/models/ptt/article"
	"github.com/meifamily/ptt-alertor/models/pushsum"
	"github.com/meifamily/ptt-alertor/models/subscription"
	user "github.com/meifamily/ptt-alertor/models/user/redis"
)

const stopHour = 48 * time.Hour
const checkPushSumDuration = 500 * time.Millisecond
const pauseCheckPushSum = 3 * time.Minute

var psckerCh = make(chan pushSumChecker)
var boardFinish = make(map[string]bool)

type pushSumChecker struct {
	Checker
	account string
}

func NewPushSumChecker() *pushSumChecker {
	return &pushSumChecker{}
}

func (psc pushSumChecker) String() string {
	textMap := map[string]string{
		"pushup":   "推文數",
		"pushdown": "噓文數",
	}
	subType := textMap[psc.subType]
	return fmt.Sprintf("%s@%s\r\n看板：%s；%s：%s%s", psc.word, psc.board, psc.board, subType, psc.word, psc.articles.StringWithPushSum())
}

type BoardArticles struct {
	board    string
	articles article.Articles
}

func (psc pushSumChecker) Run() {
	baCh := make(chan BoardArticles)

	go func() {
		for {
			boards := pushsum.List()
			for _, board := range boards {
				bl, ok := boardFinish[board]
				if !ok {
					boardFinish[board] = true
					ok, bl = true, true
				}
				if bl && ok {
					ba := BoardArticles{board: board}
					boardFinish[board] = false
					go psc.crawlArticles(ba, baCh)
				}
				time.Sleep(checkPushSumDuration)
			}
			time.Sleep(pauseCheckPushSum)
		}
	}()

	for {
		select {
		case ba := <-baCh:
			psc.board = ba.board
			boardFinish[ba.board] = true
			if len(ba.articles) > 0 {
				go psc.checkSubscribers(ba)
			}
		case pscker := <-psckerCh:
			ckCh <- pscker
		}
	}
}

func (psc pushSumChecker) crawlArticles(ba BoardArticles, baCh chan BoardArticles) {
	currentPage, err := crawler.CurrentPage(ba.board)
	if err != nil {
		log.WithFields(log.Fields{
			"board": ba.board,
		}).WithError(err).Error("Get CurrentPage Failed")
		baCh <- ba
	}

Page:
	for page := currentPage; ; page-- {
		articles, _ := crawler.BuildArticles(ba.board, page)
		for i := len(articles) - 1; i > 0; i-- {
			a := articles[i]
			if a.ID == 0 {
				continue
			}
			loc := time.FixedZone("CST", 8*60*60)
			t, err := time.ParseInLocation("1/02", a.Date, loc)
			now := time.Now()
			nowDate := now.Truncate(24 * time.Hour)
			t = t.AddDate(now.Year(), 0, 0)
			if err != nil {
				log.WithFields(log.Fields{
					"board": ba.board,
					"page":  page,
				}).WithError(err).Error("Parse DateTime Error")
				continue
			}
			if nowDate.After(t.Add(stopHour)) {
				break Page
			}
			ba.articles = append(ba.articles, a)
		}
	}

	log.WithFields(log.Fields{
		"board": ba.board,
		"total": len(ba.articles),
	}).Info("PushSum Crawl Finish")

	baCh <- ba
}

func (psc pushSumChecker) checkSubscribers(ba BoardArticles) {
	subs := pushsum.ListSubscribers(ba.board)
	for _, account := range subs {
		u := new(user.User).Find(account)
		psc.account = u.Profile.Account
		psc.email = u.Profile.Email
		psc.line = u.Profile.Line
		psc.lineNotify = u.Profile.LineAccessToken
		psc.messenger = u.Profile.Messenger
		psc.telegram = u.Profile.Telegram
		psc.telegramChat = u.Profile.TelegramChat
		go psc.checkPushSum(u, ba, checkUp)
		go psc.checkPushSum(u, ba, checkDown)
	}
}

type checkPushSum func(*pushSumChecker, subscription.Subscription, article.Articles) (article.Articles, []int)

func checkUp(psc *pushSumChecker, sub subscription.Subscription, articles article.Articles) (upArticles article.Articles, ids []int) {
	psc.word = strconv.Itoa(sub.Up)
	psc.subType = "pushup"
	if sub.Up != 0 {
		for _, a := range articles {
			if a.PushSum >= sub.Up {
				upArticles = append(upArticles, a)
				ids = append(ids, a.ID)
			}
		}
	}
	return upArticles, ids
}

func checkDown(psc *pushSumChecker, sub subscription.Subscription, articles article.Articles) (downArticles article.Articles, ids []int) {
	down := sub.Down * -1
	psc.word = strconv.Itoa(down)
	psc.subType = "pushdown"
	if sub.Down != 0 {
		for _, a := range articles {
			if a.PushSum <= down {
				downArticles = append(downArticles, a)
				ids = append(ids, a.ID)
			}
		}
	}
	return downArticles, ids
}

func (psc pushSumChecker) checkPushSum(u user.User, ba BoardArticles, checkFn checkPushSum) {
	var articles article.Articles
	var ids []int
	for _, sub := range u.Subscribes {
		if strings.EqualFold(sub.Board, ba.board) {
			articles, ids = checkFn(&psc, sub, ba.articles)
		}
	}
	if len(articles) > 0 {
		psc.articles = psc.toSendArticles(ids, articles)
		if len(psc.articles) > 0 {
			psckerCh <- psc
		}
	}
}

func (psc pushSumChecker) toSendArticles(ids []int, articles article.Articles) article.Articles {
	kindMap := map[string]string{
		"pushup":   "up",
		"pushdown": "down",
	}
	date := strconv.Itoa(time.Now().Day())
	ids = pushsum.DiffList(psc.account, psc.board, date, kindMap[psc.subType], ids...)
	diffIds := make(map[int]bool)
	for _, id := range ids {
		diffIds[id] = true
	}
	sendArticles := make(article.Articles, 0)
	for _, a := range articles {
		if diffIds[a.ID] {
			sendArticles = append(sendArticles, a)
		}
	}
	return sendArticles
}
