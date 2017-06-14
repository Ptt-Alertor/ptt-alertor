package myutil

import "strings"

// StringSlice is a type []string
type StringSlice []string

// ToStringSlice converts type []string to type StringSlice
func ToStringSlice(strs []string) StringSlice {
	ss := make(StringSlice, len(strs))
	for i, s := range strs {
		ss[i] = s
	}
	return ss
}

// Clean Removes "" in Slice
func (ss *StringSlice) Clean() {
	if *ss != nil {
		for i := 0; i < len(*ss); i++ {
			if (*ss)[i] == "" {
				*ss = append((*ss)[:i], (*ss)[i+1:]...)
			}
		}
	}
}

func (ss *StringSlice) RemoveStringsSpace() {
	if *ss != nil {
		*ss = strings.Split(strings.Replace(strings.Join(*ss, ","), " ", "", -1), ",")
	}
}

func (ss *StringSlice) AppendNonRepeat(objectStrs []string, caseSensitive bool) {
	for _, oStr := range objectStrs {
		if ss.Index(oStr, caseSensitive) == -1 {
			*ss = append(*ss, oStr)
		}
	}
}

func (ss *StringSlice) DeleteSlice(sDels []string, caseSensitive bool) {
	for _, v := range sDels {
		ss.Delete(v, caseSensitive)
	}
}

func (ss *StringSlice) Delete(s string, caseSensitive bool) {
	if i := ss.Index(s, caseSensitive); i != -1 {
		*ss = append((*ss)[:i], (*ss)[i+1:]...)
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
