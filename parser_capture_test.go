package textparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_Capture(t *testing.T) {
	a := assert.New(t)

	p := NewParser("a > b = 8")
	a.NotNil(p)

	p.MustSkip(2)
	p.StartCapture()
	p.MustSkipToString("=")
	p.MustSkip(1)

	a.Equal("> b =", p.Captured())
}

func TestParser_Extract(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		start   int
		end     int
		want    string
		wantErr bool
	}{
		{"basic", "a > b = 8", 0, 3, "a >", false},
		{"tail", "a > b = 8", 6, 9, "= 8", false},
		{"single rune", "a > b = 8", 2, 3, ">", false},
		{"nothing", "a > b = 8", 3, 3, "", true},
		{"bad start", "a > b = 8", -1, 3, "", true},
		{"bad end", "a > b = 8", 0, 300, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := NewParser(tt.input).Extract(tt.start, tt.end)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, res)
			}
		})
	}
}
