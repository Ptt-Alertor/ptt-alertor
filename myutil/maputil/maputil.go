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

func FirstByValueFloat(strs map[string]float64) string {
	var first string
	var max float64
	for str, result := range strs {
		if result > max {
			first, max = str, result
		}
	}
	return first
}
