package textparser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParser_Basic(t *testing.T) {
	a := assert.New(t)

	p := NewParser(`äöü`)
	a.Equal(`|>äöü`, p.String())
	err := p.Skip(1)
	a.Nil(err)
	a.Equal(`ä|>öü`, p.String())

	a.Equal("öü", p.MustGetNext(2))
	a.True(p.LookingAtString("öü"))
	a.False(p.LookingAtString("ü"))
	err = p.SkipString("ö")
	a.Nil(err)
	a.Equal(`äö|>ü`, p.String())
	a.False(p.IsExhausted())

	_ = p.Skip(1)
	a.True(p.IsExhausted())
}

func TestParser_Types(t *testing.T) {
	a := assert.New(t)

	p := NewParser(`Game 22: 13,  14, 216, 90121. Irrelevant Info.`)
	err := p.SkipString("Game ")
	a.Nil(err)
	gameId, err := p.ReadInt()
	a.Nil(err)
	a.Equal(22, gameId)
	err = p.SkipString(": ")
	a.Nil(err)

	numbers := make([]int, 0)
	for {
		if p.IsExhausted() {
			break
		}
		n, err := p.ReadInt()
		a.Nil(err)
		numbers = append(numbers, n)

		if p.LookingAtString(".") {
			break
		}

		p.MustSkipString(",").MustSkipSpaces()
	}
	a.Len(numbers, 4)
}

func TestParser_SkipRestOfLine(t *testing.T) {
	a := assert.New(t)

	p := NewParser("abc\ndef")
	err := p.SkipRestOfLine()
	a.Nil(err)
	a.True(p.LookingAtString("\ndef"), p.CurrentContext())

	p.MustSkipNewlines()
	a.True(p.LookingAtString(`def`), p.CurrentContext())
	v, err := p.ReadRestOfLine()
	a.Nil(err, p.CurrentContext())
	a.Equal("def", v, p.CurrentContext())
	a.True(p.IsExhausted(), p.CurrentContext())
}

func TestParser_MustReadInt(t *testing.T) {
	a := assert.New(t)

	p := NewParser("-15 16")
	a.Equal(-15, p.MustReadInt())
	p.MustSkipString(" ")

	a.Equal(16, p.MustReadInt())
}

func TestParser_MustReadToRune(t *testing.T) {
	a := assert.New(t)

	p := NewParser("15 16")
	a.Equal("15", p.MustReadToRune(' '), p.CurrentContext())
	a.True(p.LookingAtString(" 16"), p.CurrentContext())
}

func TestParser_MustReadToAnyString(t *testing.T) {
	a := assert.New(t)

	p := NewParser("a > b = 8")
	a.Equal("a ", p.MustReadToAnyString([]string{">", "= "}), p.CurrentContext())
	a.True(p.LookingAtString("> b"), p.CurrentContext())

	p = NewParser("a = b > 8")
	a.Equal("a ", p.MustReadToAnyString([]string{">", "= "}), p.CurrentContext())
	a.True(p.LookingAtString("= b"), p.CurrentContext())
}

func TestParser_MustReadToPositionString(t *testing.T) {
	a := assert.New(t)

	p := NewParser("a or b")
	a.Equal("a o", p.MustReadToPositionString(3), p.CurrentContext())
	a.True(p.LookingAtString("r b"), p.CurrentContext())
}

func TestParser_CurrentContext(t *testing.T) {
	a := assert.New(t)

	p := NewParser(`abcdefghijklmnopqrstuvwxyz`)
	p.MustSkip(13)
	a.Equal(`defghijklm|>nopqrstuvw`, p.CurrentContext())
}
