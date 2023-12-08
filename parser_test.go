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

	a.Equal("öü", p.GetLookingAtString(2))
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
