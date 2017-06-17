package subscription

import (
	"sort"
	"strings"

	"github.com/liam-lai/ptt-alertor/myutil"
)

type Subscription struct {
	Board    string             `json:"board"`
	Keywords myutil.StringSlice `json:"keywords"`
	Authors  myutil.StringSlice `json:"authors"`
}

func (s Subscription) String() string {
	if len(s.Keywords) == 0 {
		return ""
	}
	sort.Strings(s.Keywords)
	return s.Board + ": " + strings.Join(s.Keywords, ", ")
}

func (s Subscription) StringAuthor() string {
	if len(s.Authors) == 0 {
		return ""
	}
	sort.Strings(s.Authors)
	return s.Board + ": " + strings.Join(s.Authors, ", ")
}

func (s *Subscription) CleanUp() {
	s.Keywords.Clean()
	s.Keywords.RemoveStringsSpace()
	s.Authors.Clean()
	s.Authors.RemoveStringsSpace()
}
