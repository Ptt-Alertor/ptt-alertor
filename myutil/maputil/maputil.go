package maputil

func MaxIntKey(strs map[string]int) (first string) {
	var max int
	for str, cnt := range strs {
		if cnt > max {
			first, max = str, cnt
		}
	}
	return first
}

func MaxFloatKey(strs map[string]float64) (first string) {
	var max float64
	for str, result := range strs {
		if result > max {
			first, max = str, result
		}
	}
	return first
}
