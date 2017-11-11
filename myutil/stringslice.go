package myutil

import "strings"

// StringSlice is a type []string
type StringSlice []string

// Clean Removes "" and "*" in Slice
func (ss *StringSlice) Clean() {
	if *ss != nil {
		for i := 0; i < len(*ss); i++ {
			if (*ss)[i] == "" || (*ss)[i] == "*" {
				*ss = append((*ss)[:i], (*ss)[i+1:]...)
				i--
			}
		}
	}
}

func (ss *StringSlice) RemoveStringsSpace() {
	if *ss != nil {
		*ss = strings.Split(strings.Replace(strings.Join(*ss, ","), " ", "", -1), ",")
	}
}

func (ss *StringSlice) AppendNonRepeatElement(str string, caseSensitive bool) {
	if ss.Index(str, caseSensitive) == -1 {
		*ss = append(*ss, str)
	}
}

func (ss *StringSlice) AppendNonRepeat(objectStrs []string, caseSensitive bool) {
	for _, oStr := range objectStrs {
		if ss.Index(oStr, caseSensitive) == -1 {
			*ss = append(*ss, oStr)
		}
	}
}

func (ss *StringSlice) DeleteElement(s string, caseSensitive bool) {
	if i := ss.Index(s, caseSensitive); i != -1 {
		*ss = append((*ss)[:i], (*ss)[i+1:]...)
	}
}

func (ss *StringSlice) Delete(sDels []string, caseSensitive bool) {
	for _, v := range sDels {
		ss.DeleteElement(v, caseSensitive)
	}
}

func (ss StringSlice) Index(value string, caseSensitive bool) int {
	for i, v := range ss {
		if equalString(v, value, caseSensitive) {
			return i
		}
	}
	return -1
}

func equalString(a string, b string, caseSensitive bool) bool {
	if caseSensitive {
		return a == b
	}

	return strings.EqualFold(a, b)
}
