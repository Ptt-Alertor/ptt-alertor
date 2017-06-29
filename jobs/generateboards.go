package jobs

import (
	log "github.com/meifamily/logrus"

	board "github.com/liam-lai/ptt-alertor/models/ptt/board/redis"
	"github.com/liam-lai/ptt-alertor/models/subscription"
	user "github.com/liam-lai/ptt-alertor/models/user/redis"
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

	for _, usr := range usrs {
		for _, sub := range usr.Subscribes {
			if !boardNameBool[sub.Board] {
				createBoard(sub)
			}
		}
	}
	log.Info("Boards Generated")
}

func createBoard(sub subscription.Subscription) {
	bd := new(board.Board)
	bd.Name = sub.Board
	bd.Create()
	log.WithField("board", bd.Name).Info("Added Board")
}
