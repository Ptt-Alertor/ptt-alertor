package jobs

import (
	"strings"

	"github.com/meifamily/ptt-alertor/models/author"
	"github.com/meifamily/ptt-alertor/models/board"
	"github.com/meifamily/ptt-alertor/models/keyword"
	"github.com/meifamily/ptt-alertor/models/pushsum"
	"github.com/meifamily/ptt-alertor/models/user"
)

const preBoard = "iphone"
const postBoard = "ios"

type migrateBoard struct{}

func NewMigrateBoard() *migrateBoard {
	return &migrateBoard{}
}

func (migrateBoard) Run() {
	addBoard(postBoard)
	pushsum.Add(postBoard)

	subs := pushsum.ListSubscribers(preBoard)
	for _, sub := range subs {
		pushsum.AddSubscriber(postBoard, sub)
	}

	subs = author.Subscribers(preBoard)
	for _, sub := range subs {
		author.AddSubscriber(postBoard, sub)
	}

	subs = keyword.Subscribers(preBoard)
	for _, sub := range subs {
		keyword.AddSubscriber(postBoard, sub)
	}

	users := user.NewUser().All()
	for _, u := range users {
		for _, sub := range u.Subscribes {
			if strings.EqualFold(sub.Board, preBoard) {
				u.Subscribes.Remove(sub)
				sub.Board = postBoard
				u.Subscribes.Add(sub)
				u.Update()
			}
		}
	}

	keyword.Destroy(preBoard)
	author.Destroy(preBoard)
	pushsum.Remove(preBoard)
	pushsum.Destroy(preBoard)

	// delete board from board list and board content
	bd := board.NewBoard()
	bd.Name = preBoard
	bd.Delete()

	// rename article's content board name
}
