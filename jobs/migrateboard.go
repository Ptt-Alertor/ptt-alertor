package jobs

import (
	"github.com/meifamily/ptt-alertor/models/author"
	"github.com/meifamily/ptt-alertor/models/keyword"
	"github.com/meifamily/ptt-alertor/models/pushsum"
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

	// TODO: delete preboard's key
	// TODO: change user profile's board name
}
