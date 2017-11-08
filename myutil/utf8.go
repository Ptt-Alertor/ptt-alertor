package myutil

import "unicode/utf8"

func SplitTextByLineBreak(text string, limit int) (texts []string) {
	var runeCount, index, lineBreak, start int
	copyText := text
	for len(copyText) > 0 {
		r, size := utf8.DecodeRuneInString(copyText)
		index = index + size
		runeCount++

		if runeCount > limit {
			runeCount = 0
			texts = append(texts, text[start:lineBreak])
			start = lineBreak
		}
		if string(r) == "\n" {
			lineBreak = index
		}
		copyText = copyText[size:]
	}
	texts = append(texts, text[start:])
	return texts
}
