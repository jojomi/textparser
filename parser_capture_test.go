package textparser

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
