package jobs

import (
	log "github.com/meifamily/logrus"

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
				createBoard(sub.Board)
			}
			if sub.PushSum != emptyPushSum {
				createPushSumKeys(u.Profile.Account, sub.Board)
			}
			if len(sub.Keywords) > 0 {
				createKeyword(u.Profile.Account, sub.Board)
			}
			if len(sub.Authors) > 0 {
				createAuthor(u.Profile.Account, sub.Board)
			}
		}
	}
	log.Info("Generated Done")
}

func createBoard(boardName string) {
	bd := board.NewBoard()
	bd.Name = boardName
	bd.Create()
	log.WithField("board", bd.Name).Info("Added Board")
}

func createPushSumKeys(account, board string) {
	pushsum.Add(board)
	pushsum.AddSubscriber(board, account)
	log.WithField("board", board).Info("Added PushSum Board and Subscriber")
}

func createKeyword(account, board string) {
	keyword.AddSubscriber(board, account)
	log.WithField("board", board).Info("Added Keyword Subscriber")
}

func createAuthor(account, board string) {
	author.AddSubscriber(board, account)
	log.WithField("board", board).Info("Added Author Subscriber")
}
