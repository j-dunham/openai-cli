package util

import (
	"fmt"
	"strings"
	"time"
)

func LoadingAnimation(message string, done chan bool) {
	loadingChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	for {
		select {
		case <-done:
			// Clear the loading animation
			fmt.Print("\r\033[K")
			return
		default:
			fmt.Printf("\r\033[K %s %s ", message, loadingChars[i%len(loadingChars)])
			time.Sleep(100 * time.Millisecond)
			i++
		}
	}
}

func WrapText(text string, lineWidth int, prefix string) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	wrapped := prefix
	lineLen := 0
	for _, word := range words {
		if lineLen+len(word) > lineWidth {

			wrapped = strings.TrimSpace(wrapped) + "\n" + prefix
			lineLen = 0
		}
		wrapped += word + " "
		lineLen += len(word) + 1
	}
	return strings.TrimSpace(wrapped)
}
