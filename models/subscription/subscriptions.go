package subscription

import (
	"strings"

	"fmt"

	"github.com/liam-lai/ptt-alertor/crawler"
	boardProto "github.com/liam-lai/ptt-alertor/models/ptt/board"
	board "github.com/liam-lai/ptt-alertor/models/ptt/board/redis"
)

type Subscriptions []Subscription

func (ss Subscriptions) String() string {
	str := "關鍵字\n"
	for _, sub := range ss {
		if sub.String() != "" {
			str += sub.String() + "\n"
		}
	}
	str += "----\n作者\n"
	for _, sub := range ss {
		if sub.StringAuthor() != "" {
			str += sub.StringAuthor() + "\n"
		}
	}
	return str
}

func (ss *Subscriptions) Add(sub Subscription) error {
	sub.Board = strings.ToLower(sub.Board)
	if ok, suggestion := checkBoardExist(sub.Board); !ok {
		return boardProto.BoardNotExistError{suggestion}
	}
	sub.CleanUp()
	for i, s := range *ss {
		if strings.EqualFold(s.Board, sub.Board) {
			s.Keywords.AppendNonRepeat(sub.Keywords, false)
			s.Authors.AppendNonRepeat(sub.Authors, false)
			(*ss)[i] = s
			fmt.Printf("%+v", ss)
			return nil
		}
	}
	*ss = append(*ss, sub)

	return nil
}

func (ss *Subscriptions) Remove(sub Subscription) error {
	sub.Board = strings.ToLower(sub.Board)
	if ok, suggestion := checkBoardExist(sub.Board); !ok {
		return boardProto.BoardNotExistError{suggestion}
	}
	sub.CleanUp()
	for i := 0; i < len(*ss); i++ {
		s := (*ss)[i]
		if strings.EqualFold(s.Board, sub.Board) {
			s.Keywords.DeleteSlice(sub.Keywords, false)
			s.Authors.DeleteSlice(sub.Authors, false)
			(*ss)[i] = s
		}
		if len((*ss)[i].Keywords) == 0 {
			*ss = append((*ss)[:i], (*ss)[i+1:]...)
			i--
		}
	}
	return nil
}

func checkBoardExist(boardName string) (bool, string) {
	bd := new(board.Board)
	bd.Name = boardName
	if bd.Exist() {
		return true, ""
	}
	if crawler.CheckBoardExist(boardName) {
		bd.Create()
		return true, ""
	}

	return false, bd.SuggestBoardName()
}
