package textparser

import (
	"errors"
	"unicode"
	"unicode/utf8"
)

// LookingAtString determines if from the current position, the next runes would equal the expected string provided.
func (x *Parser) LookingAtString(expected string) bool {
	result, err := x.GetNext(utf8.RuneCountInString(expected))
	if errors.Is(err, EndOfInputError{}) {
		return false
	}
	return result == expected
}

func (x *Parser) LookingAtRune(r rune) bool {
	result, err := x.GetNextRune()
	if errors.Is(err, EndOfInputError{}) {
		return false
	}
	return result == r
}

func (x *Parser) LookingAtWhitespace() bool {
	nextRune, err := x.GetNextRune()
	if errors.Is(err, EndOfInputError{}) {
		return false
	}
	if nextRune == ' ' || nextRune == '\t' || nextRune == '\r' || nextRune == '\n' {
		return true
	}
	return false
}

func (x *Parser) LookingAtDigit() bool {
	nextRune, err := x.GetNextRune()
	if errors.Is(err, EndOfInputError{}) {
		return false
	}
	if nextRune >= '0' && nextRune <= '9' {
		return true
	}
	return false
}

func (x *Parser) LookingAtLetter() bool {
	nextRune, err := x.GetNextRune()
	if errors.Is(err, EndOfInputError{}) {
		return false
	}
	return unicode.IsLetter(nextRune)
}

func (x *Parser) MustGetNext(runeCount int) string {
	result, err := x.GetNext(runeCount)
	if err != nil {
		panic(err)
	}
	return result
}

// GetNext returns the next runeCount runes as a string.
func (x *Parser) GetNext(runeCount int) (string, error) {
	if x.RemainingRuneCount() < runeCount {
		return "", EndOfInputError{}
	}
	return string(x.input[x.position : x.position+runeCount]), nil
}

func (x *Parser) GetNextMax(runeCount int) string {
	if runeCount > x.RemainingRuneCount() {
		runeCount = x.RemainingRuneCount()
	}
	return string(x.input[x.position : x.position+runeCount])
}

func (x *Parser) MustGetNextRune() rune {
	result, err := x.GetNextRune()
	if err != nil {
		panic(err)
	}
	return result
}

func (x *Parser) GetNextRune() (rune, error) {
	if x.IsExhausted() {
		var empty rune
		return empty, EndOfInputError{}
	}
	return x.input[x.position], nil
}

func (x *Parser) MustGetNextWord() string {
	// save position
	pos := x.position

	result := x.MustReadWord()

	// reset position
	x.position = pos

	return result
}
