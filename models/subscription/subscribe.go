package subscription

import "strings"

type Subscribe struct {
	Board    string
	Keywords []string
}

func (s Subscribe) String() string {
	return s.Board + ": " + strings.Join(s.Keywords, ", ")
}

type Subscribes []Subscribe

func (ss Subscribes) String() string {
	str := ""
	for _, sub := range ss {
		str += sub.String() + "\n"
	}
	return str
}

//TODO: check keyword exist or not
func (ss *Subscribes) Add(sub Subscribe) {
	sub.Keywords = removeStringsSpace(sub.Keywords)
	for i, s := range *ss {
		if s.Board == sub.Board {
			(*ss)[i].Keywords = append((*ss)[i].Keywords, sub.Keywords...)
			return
		}
	}
	*ss = append(*ss, sub)
}

func (ss *Subscribes) Remove(sub Subscribe) {
	sub.Keywords = removeStringsSpace(sub.Keywords)
	for i, s := range *ss {
		if s.Board == sub.Board {
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
		}
	}
}

func removeStringsSpace(strs []string) []string {
	return strings.Split(strings.Replace(strings.Join(strs, ","), " ", "", -1), ",")
}
