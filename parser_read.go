package textparser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

func (x *Parser) MustReadRune() rune {
	return x.MustReadRunes(1)[0]
}

func (x *Parser) MustReadRunes(runeCount int) []rune {
	value, err := x.ReadRunes(runeCount)
	if err != nil {
		panic(err)
	}
	return value
}

func (x *Parser) ReadRune() (rune, error) {
	runes, err := x.ReadRunes(1)
	return runes[0], err
}

func (x *Parser) ReadRunes(runeCount int) ([]rune, error) {
	if x.RemainingRuneCount() < runeCount {
		var empty []rune
		return empty, fmt.Errorf("could not read %d runes, reader is exhausted", runeCount)
	}
	value := x.input[x.position : x.position+runeCount]
	err := x.Skip(runeCount)
	if err != nil {
		var empty []rune
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

func (x *Parser) MustReadWord() string {
	value, err := x.ReadWord()
	if err != nil {
		panic(err)
	}
	return value
}

// ReadWord reads exactly one word of input. Stops at space, newline and end of input. Error if the input end is already reached.
func (x *Parser) ReadWord() (string, error) {
	if x.IsExhausted() {
		return "", fmt.Errorf("could not read on, reader is exhausted")
	}
	pos := x.position + 1
	for {
		if len(x.input) <= pos {
			break
		}
		value := x.input[pos]
		if value == ' ' {
			break
		}
		if value == '\n' {
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
			x.GetNextMax(10),
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

func (x *Parser) ReadToPositionString(newPosition int) (string, error) {
	if newPosition >= len(x.input) {
		return "", EndOfInputError{}
	}
	content := x.input[x.position:newPosition]
	x.position = newPosition
	return string(content), nil
}

func (x *Parser) MustReadToPositionString(newPosition int) string {
	value, err := x.ReadToPositionString(newPosition)
	if err != nil {
		panic(err)
	}
	return value
}

func (x *Parser) ReadToAnyString(limitStrings []string) (string, error) {
	newPos, err := x.findAnyNext(limitStrings)
	if err != nil {
		return "", err
	}
	content, err := x.ReadToPositionString(newPos)
	if err != nil {
		return "", err
	}
	return content, nil
}

func (x *Parser) MustReadToAnyString(limitStrings []string) string {
	content, err := x.ReadToAnyString(limitStrings)
	if err != nil {
		panic(err)
	}
	return content
}
