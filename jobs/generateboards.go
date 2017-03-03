package jobs

import (
	"fmt"

	"github.com/liam-lai/ptt-alertor/models/file/ptt/board"
	"github.com/liam-lai/ptt-alertor/models/file/user"
)

type GenBoards struct {
}

func (gb GenBoards) Run() {
	usrs := new(user.User).All()
	bds := new(board.Board).All()
	boardNameBool := make(map[string]bool)
	for _, bd := range bds {
		boardNameBool[bd.Name] = true
	}

	for _, usr := range usrs {
		for _, sub := range usr.Subscribes {
			if !boardNameBool[sub.Board] {
				bd := new(board.Board)
				bd.Name = sub.Board
				fmt.Println(bd.Name)
				bd.Create()
			}
		}
	}
}
