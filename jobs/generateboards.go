package jobs

import (
	log "github.com/meifamily/logrus"

	board "github.com/meifamily/ptt-alertor/models/ptt/board/redis"
	"github.com/meifamily/ptt-alertor/models/pushsum"
	"github.com/meifamily/ptt-alertor/models/subscription"
	user "github.com/meifamily/ptt-alertor/models/user/redis"
)

type Generator struct {
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (gb Generator) Run() {
	usrs := new(user.User).All()
	bds := new(board.Board).All()
	boardNameBool := make(map[string]bool)
	for _, bd := range bds {
		boardNameBool[bd.Name] = true
	}

	emptyPushSum := subscription.PushSum{}
	for _, usr := range usrs {
		for _, sub := range usr.Subscribes {
			if !boardNameBool[sub.Board] {
				createBoard(sub.Board)
			}
			if sub.PushSum != emptyPushSum {
				createPushSumKeys(usr.Profile.Account, sub.Board)
			}
		}
	}
	log.Info("Boards Generated")
}

func createBoard(boardName string) {
	bd := new(board.Board)
	bd.Name = boardName
	bd.Create()
	log.WithField("board", bd.Name).Info("Added Board")
}

func createPushSumKeys(account, board string) {
	pushsum.Add(board)
	pushsum.AddSubscriber(board, account)
	log.WithField("board", board).Info("Added PushSum Board")
}
