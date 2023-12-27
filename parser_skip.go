package textparser

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

// Skip skips runeCount runes without returning them.
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

// MustSkip (chainable).
func (x *Parser) MustSkip(runeCount int) *Parser {
	err := x.Skip(runeCount)
	if err != nil {
		panic(err)
	}
	return x
}

// SkipTo skips to index newIndex without returning the intermediate runes (forward only).
func (x *Parser) SkipTo(newIndex int) error {
	if x.position > newIndex {
		return fmt.Errorf("can't skip backwards, current %d, target position %d", x.position, newIndex)
	}
	return x.Skip(newIndex - x.position)
}

func (x *Parser) SkipString(expected string) error {
	if !x.LookingAtString(expected) {
		return fmt.Errorf(
			"expected string '%s' at position %d, but looking at '%s' instead",
			expected,
			x.position,
			x.GetNextMax(utf8.RuneCountInString(expected)),
		)
	}
	return x.Skip(utf8.RuneCountInString(expected))
}

func (x *Parser) MustSkipAnyWhitespaces() *Parser {
	err := x.SkipAny([]rune{'\n', '\t', ' '})
	if err != nil {
		panic(err)
	}
	return x
}

func (x *Parser) MustSkipString(expected string) *Parser {
	err := x.SkipString(expected)
	if err != nil {
		panic(err)
	}
	return x
}

func (x *Parser) MustSkipToString(expected string) *Parser {
	err := x.SkipToString(expected)
	if err != nil {
		panic(err)
	}
	return x
}

func (x *Parser) SkipToString(expected string) error {
	var err error
	for {
		if x.IsExhausted() {
			return EndOfInputError{}
		}

		if x.LookingAtString(expected) {
			break
		}

		err = x.Skip(1)
		if err != nil {
			return err
		}
	}

	return nil
}

func (x *Parser) SkipSpaces() error {
	return x.SkipAny([]rune{' '})
}

func (x *Parser) SkipAnyWhitespaces() error {
	return x.SkipAny([]rune{'\n', '\t', ' '})
}

func (x *Parser) SkipAny(runes []rune) error {
	var (
		goOn bool
		err  error
	)
	for {
		goOn = false
		for _, r := range runes {
			if x.LookingAtRune(r) {
				goOn = true
				break
			}
		}
		if !goOn {
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
	return x.SkipAny([]rune{'\n'})
}

func (x *Parser) MustSkipNewlines() *Parser {
	err := x.SkipNewlines()
	if err != nil {
		panic(err)
	}
	return x
}

func (x *Parser) getToIndex(endIndex int) (string, error) {
	if x.position > endIndex {
		return "", fmt.Errorf("can't read backwards, current %d, target position %d", x.position, endIndex)
	}
	if endIndex > len(x.input) {
		return "", fmt.Errorf("can't read that far, current index %d, target position %d, input length %d", x.position, endIndex, len(x.input))
	}

	return x.MustGetNext(endIndex - x.position), nil
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
