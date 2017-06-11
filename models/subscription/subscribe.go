package subscription

import (
	"strings"

	"github.com/liam-lai/ptt-alertor/crawler"
	boardProto "github.com/liam-lai/ptt-alertor/models/ptt/board"
	board "github.com/liam-lai/ptt-alertor/models/ptt/board/redis"
	"github.com/liam-lai/ptt-alertor/myutil/collection"
)

type Subscription struct {
	Board    string   `json:"board"`
	Keywords []string `json:"keywords"`
}

func (s Subscription) String() string {
	return s.Board + ": " + strings.Join(s.Keywords, ", ")
}

type Subscriptions []Subscription

func (ss Subscriptions) String() string {
	str := ""
	for _, sub := range ss {
		str += sub.String() + "\n"
	}
	return str
}

func (ss *Subscriptions) Add(sub Subscription) error {
	sub.Board = strings.ToLower(sub.Board)
	if ok, suggestion := checkBoardExist(sub.Board); !ok {
		return boardProto.BoardNotExistError{suggestion}
	}
	sub.Keywords = removeStringsSpace(sub.Keywords)
	for i, s := range *ss {
		if strings.EqualFold(s.Board, sub.Board) {
			for _, keyword := range sub.Keywords {
				if !collection.In(s.Keywords, keyword) {
					(*ss)[i].Keywords = append((*ss)[i].Keywords, keyword)
				}
			}
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
	sub.Keywords = removeStringsSpace(sub.Keywords)
	for i := 0; i < len(*ss); i++ {
		s := (*ss)[i]
		if strings.EqualFold(s.Board, sub.Board) {
			for _, subKeyword := range sub.Keywords {
				for j := 0; j < len(s.Keywords); j++ {
					if s.Keywords[j] == subKeyword {
						s.Keywords = append(s.Keywords[:j], s.Keywords[j+1:]...)
						j--
					}
				}
				(*ss)[i].Keywords = s.Keywords
			}
		}
		if len((*ss)[i].Keywords) == 0 {
			*ss = append((*ss)[:i], (*ss)[i+1:]...)
			i--
		}
	}
	return nil
}

func removeStringsSpace(strs []string) []string {
	return strings.Split(strings.Replace(strings.Join(strs, ","), " ", "", -1), ",")
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
