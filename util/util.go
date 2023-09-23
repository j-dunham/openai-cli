package util

import (
	"strings"
)

func WrapText(text string, lineWidth int, prefix string) string {
	words := strings.Fields(text)
	if len(words) == 0 {
			return text
	}

	wrapped := ""
	lineLen := 0
	for _, word := range words {
			if lineLen+len(word) > lineWidth {
					wrapped += "\n" + prefix
					lineLen = 0
			}
			wrapped += word + " "
			lineLen += len(word) + 1
	}

	return wrapped
}