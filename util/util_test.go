package util

import (
	"testing"
)

func TestWrapText(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		lineWidth int
		prefix    string
		want      string
	}{
		{
			name:      "WrapText wraps text correctly",
			text:      "Hello, world!",
			lineWidth: 8,
			prefix:    "",
			want:      "Hello, \nworld! ",
		},
		{
			name:      "WrapText wraps text correctly with prefix",
			text:      "Hello, world!",
			lineWidth: 8,
			prefix:    "|",
			want:      "|Hello, \n|world! ",
		},
		{
			name:      "WrapText wraps text correctly for a single line",
			text:      "Hello, world!",
			lineWidth: 80,
			prefix:    "|",
			want:      "|Hello, world! ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WrapText(tt.text, tt.lineWidth, tt.prefix); got != tt.want {
				t.Errorf("WrapText() = %v\n%v", got, tt.want)
			}
		})
	}
}
