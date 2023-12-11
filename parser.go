package textparser

import (
	"errors"
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

func (x *Parser) SkipTo(newIndex int) error {
	if x.position > newIndex {
		return fmt.Errorf("can't skip backwards, current %d, target position %d", x.position, newIndex)
	}
	return x.Skip(newIndex - x.position)
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

func (x *Parser) MustSkip(runeCount int) *Parser {
	err := x.Skip(runeCount)
	if err != nil {
		panic(err)
	}
	return x
}

func (x *Parser) SkipToString(expected string) error {
	for {
		if x.IsExhausted() {

		}
	}
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
	var err error
	for {
		if !x.LookingAtString(" ") {
			break
		}

		err = x.Skip(1)
		if err != nil {
			return err
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

func (x *Parser) MustSkipRestOfLine() *Parser {
	err := x.SkipRestOfLine()
	if err != nil {
		panic(err)
	}
	return x
}

func (x *Parser) SkipRestOfLine() error {
	end, err := x.findNext("\n")

	// end of input?
	if errors.Is(err, EndOfInputError{}) {
		x.SkipToEnd()
		return nil
	}

	// something else went wrong
	if err != nil {
		return err
	}

	return x.SkipTo(end)
}

func (x *Parser) SkipToEnd() *Parser {
	x.position = len(x.input)
	return x
}

func (x *Parser) SkipNewlines() error {
	var err error
	for {
		if !x.LookingAtString("\n") {
			break
		}

		err = x.Skip(1)
		if err != nil {
			return err
		}
	}
	return nil
}

func (x *Parser) MustSkipNewlines() *Parser {
	err := x.SkipNewlines()
	if err != nil {
		panic(err)
	}
	return x
}

func (x *Parser) IsExhausted() bool {
	return x.position >= len(x.input)
}

func (x *Parser) MustReadRune() rune {
	value, err := x.ReadRune()
	if err != nil {
		panic(err)
	}
	return value
}

func (x *Parser) ReadRune() (rune, error) {
	if x.IsExhausted() {
		var empty rune
		return empty, fmt.Errorf("could not read a single rune, reader is exhausted")
	}
	value := x.input[x.position]
	err := x.Skip(1)
	if err != nil {
		var empty rune
		return empty, err
	}
	return value, nil
}

func (x *Parser) MustReadToRune(stopRune rune) string {
	value, err := x.ReadToRune(stopRune)
	if err != nil {
		panic(err)
	}
	return value
}

func (x *Parser) ReadToRune(stopRune rune) (string, error) {
	if x.IsExhausted() {
		return "", fmt.Errorf("could not read on, reader is exhausted")
	}
	pos := x.position + 1
	for {
		// TODO check input limit
		value := x.input[pos]
		if value == stopRune {
			break
		}
		pos++
	}
	return x.readRunes(pos - x.position)
}

func (x *Parser) readRunes(runeCount int) (string, error) {
	// TODO check length
	value := x.input[x.position : x.position+runeCount]
	err := x.Skip(len(value))
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func (x *Parser) ReadRestOfLine() (string, error) {
	end, err := x.findNext("\n")

	// end of input?
	if errors.Is(err, EndOfInputError{}) {
		value, err := x.getToIndex(len(x.input))
		if err != nil {
			return "", err
		}

		err = x.Skip(utf8.RuneCountInString(value))
		if err != nil {
			return "", err
		}

		return value, nil
	}

	// something else went wrong
	if err != nil {
		return "", err
	}

	value, err := x.getToIndex(end)
	if err != nil {
		return "", err
	}

	err = x.SkipTo(end)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (x *Parser) MustReadRestOfLine() string {
	result, err := x.ReadRestOfLine()
	if err != nil {
		panic(err)
	}
	return result
}

// TODO allow for negative numbers
func (x *Parser) ReadInt() (int, error) {
	var (
		num    strings.Builder
		numPos = x.position
		r      rune
	)

	for {
		// exhausted?
		if numPos >= len(x.input) {
			break
		}

		r = x.input[numPos]

		if r != '0' && r != '1' && r != '2' && r != '3' && r != '4' && r != '5' && r != '6' && r != '7' && r != '8' && r != '9' && (numPos > x.position || r != '-') {
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

func (x *Parser) MustReadInt() int {
	value, err := x.ReadInt()
	if err != nil {
		panic(err)
	}
	return value
}

func (x *Parser) getToIndex(endIndex int) (string, error) {
	if x.position > endIndex {
		return "", fmt.Errorf("can't read backwards, current %d, target position %d", x.position, endIndex)
	}

	return x.GetLookingAtString(endIndex - x.position), nil
}

func (x *Parser) findNext(r string) (int, error) {
	return x.findAnyNext([]string{r})
}

func (x *Parser) findAnyNext(r []string) (int, error) {
	pos := x.position

	for {
		// exhausted?
		if pos >= len(x.input) {
			return 0, EndOfInputError{}
		}
		for _, s := range r {
			if string(x.input[pos:pos+len(s)]) == s {
				return pos, nil
			}
		}
		pos++
	}
}

func (x *Parser) String() string {
	return string(x.input[0:x.position]) + "|>" + string(x.input[x.position:])
}

func (x *Parser) CurrentContext() string {
	const contextLength = 10
	start := max(0, x.position-contextLength)
	end := min(x.position+contextLength, len(x.input))
	result := string(x.input[start:x.position]) + "|>" + string(x.input[x.position:end])
	return result
}
