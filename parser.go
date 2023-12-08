package textparser

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Parser struct {
	input    []rune
	position int
}

func NewParser(input string) *Parser {
	return &Parser{
		input:    []rune(input),
		position: 0,
	}
}

func (x *Parser) Skip(runeCount int) error {
	if x.position+runeCount > len(x.input) {
		return fmt.Errorf(
			"can't advance another %d runes, the input contains %d runes and current pointer is at %d",
			runeCount,
			len(x.input),
			x.position,
		)
	}
	x.position += runeCount
	return nil
}

func (x *Parser) SkipToString(expected string) error {
	for {
		if x.IsExhausted() {

		}
	}
}

func (x *Parser) SkipRestOfLine() error {
	return nil
}

func (x *Parser) LookingAtString(expected string) bool {
	return x.GetLookingAtString(utf8.RuneCountInString(expected)) == expected
}

func (x *Parser) GetLookingAtString(runeCount int) string {
	return string(x.input[x.position : x.position+runeCount])
}

func (x *Parser) SkipString(expected string) error {
	if !x.LookingAtString(expected) {
		return fmt.Errorf(
			"expected string '%s' at position %d, but looking at '%s' instead",
			expected,
			x.position,
			x.GetLookingAtString(utf8.RuneCountInString(expected)),
		)
	}
	return x.Skip(utf8.RuneCountInString(expected))
}

func (x *Parser) MustSkipString(expected string) *Parser {
	err := x.SkipString(expected)
	if err != nil {
		panic(err)
	}
	return x
}

func (x *Parser) SkipSpaces() error {
	if !x.LookingAtString(" ") {
		return fmt.Errorf(
			"expected space(s) at position %d, but looking at '%s' instead",
			x.position,
			x.GetLookingAtString(10),
		)
	}

	var err error
	for {
		err = x.Skip(1)
		if err != nil {
			return err
		}

		if !x.LookingAtString(" ") {
			break
		}
	}
	return nil
}

func (x *Parser) MustSkipSpaces() *Parser {
	err := x.SkipSpaces()
	if err != nil {
		panic(err)
	}
	return x
}

func (x *Parser) IsExhausted() bool {
	return x.position == len(x.input)
}

// TODO allow for negative numbers
func (x *Parser) ReadInt() (int, error) {
	var (
		num    strings.Builder
		numPos = x.position
		r      rune
	)

	for {
		r = x.input[numPos]

		if r != '0' && r != '1' && r != '2' && r != '3' && r != '4' && r != '5' && r != '6' && r != '7' && r != '8' && r != '9' {
			break
		}

		num.WriteRune(r)
		numPos++
	}

	if num.Len() == 0 {
		return 0, fmt.Errorf(
			"expected to find an integer at position %d, but looking at '%s'",
			x.position,
			x.GetLookingAtString(10),
		)
	}

	// try to convert (should only fail due to the controlled buildup above when the number is too big)
	v, err := strconv.Atoi(num.String())
	if err != nil {
		return 0, err
	}

	_ = x.Skip(num.Len())

	return v, nil
}

func (x *Parser) String() string {
	return string(x.input[0:x.position]) + "|>" + string(x.input[x.position:])
}
