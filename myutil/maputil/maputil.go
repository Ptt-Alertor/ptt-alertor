package maputil

func FirstByValueInt(strs map[string]int) string {
	var first string
	var max int
	for str, cnt := range strs {
		if cnt > max {
			first, max = str, cnt
		}
	}
	return first
}
