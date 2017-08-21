package jobs

import (
	log "github.com/meifamily/logrus"

	"github.com/meifamily/ptt-alertor/models/article"
	"github.com/meifamily/ptt-alertor/models/author"
	"github.com/meifamily/ptt-alertor/models/board"
	"github.com/meifamily/ptt-alertor/models/keyword"
	"github.com/meifamily/ptt-alertor/models/pushsum"
	"github.com/meifamily/ptt-alertor/models/subscription"
	"github.com/meifamily/ptt-alertor/models/user"
)

type Generator struct {
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (gb Generator) Run() {
	users := user.NewUser().All()
	bds := board.NewBoard().All()
	boardNameBool := make(map[string]bool)
	for _, bd := range bds {
		boardNameBool[bd.Name] = true
	}

	emptyPushSum := subscription.PushSum{}
	for _, u := range users {
		for _, sub := range u.Subscribes {
			if !boardNameBool[sub.Board] {
				addBoard(sub.Board)
			}
			if sub.PushSum != emptyPushSum {
				addPushsumSub(u.Profile.Account, sub.Board)
			}
			if len(sub.Keywords) > 0 {
				addKeywordSub(u.Profile.Account, sub.Board)
			}
			if len(sub.Authors) > 0 {
				addAuthorSub(u.Profile.Account, sub.Board)
			}
			if len(sub.Articles) > 0 {
				for _, a := range sub.Articles {
					addArticleSub(u.Profile.Account, a)
				}
			}
		}
	}
	log.Info("Generated Done")
}

func addBoard(boardName string) {
	bd := board.NewBoard()
	bd.Name = boardName
	bd.Create()
	log.WithField("board", bd.Name).Info("Added Board")
}

func addPushsumSub(account, board string) {
	pushsum.Add(board)
	pushsum.AddSubscriber(board, account)
	log.WithField("board", board).Info("Added PushSum Board and Subscriber")
}

func addKeywordSub(account, board string) {
	keyword.AddSubscriber(board, account)
	log.WithField("board", board).Info("Added Keyword Subscriber")
}

func addAuthorSub(account, board string) {
	author.AddSubscriber(board, account)
	log.WithField("board", board).Info("Added Author Subscriber")
}

func addArticleSub(account, articleID string) {
	a := article.Article{
		Code: articleID,
	}
	a.AddSubscriber(account)
	log.WithField("article", articleID).Info("Added Article Subscriber")
}
