package myutil

import "unicode/utf8"

func SplitTextByLineBreak(text string, limit int) (texts []string) {
	var runeCount, index, lineBreak, start int
	copyText := text
	for len(copyText) > 0 {
		if utf8.RuneCountInString(text[start:]) <= limit {
			texts = append(texts, text[start:])
			break
		}
		r, size := utf8.DecodeRuneInString(copyText)
		runeCount++
		index += size

		if runeCount > limit {
			runeCount = utf8.RuneCountInString(text[start:index])
			texts = append(texts, text[start:lineBreak])
			start = lineBreak
		}

		if string(r) == "\n" {
			lineBreak = index
		}
		copyText = copyText[size:]
	}
	return texts
}
