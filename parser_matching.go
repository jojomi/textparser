package textparser

import (
	"fmt"
)

func (x *Parser) ReadToMatchingString(open, close string) (string, error) {
	var (
		oldIndex           = x.position
		openCount          = 1
		runeCount          = 0
		differentDelims    = open != close
		remainingRuneCount = x.RemainingRuneCount()
		err                error
	)

	for {
		if differentDelims {
			if openCount < 0 {
				x.position = oldIndex
				return "", fmt.Errorf("could not find closing rune %s, because they were used unbalanced, closed before opened (rune %s)", string(close), string(open))
			}
			if runeCount > remainingRuneCount {
				x.position = oldIndex
				return "", fmt.Errorf("could not find closing rune %s, input exhausted", string(close))
			}
		}

		if x.LookingAtString(close) {
			openCount--
			if openCount == 0 {
				break
			}

			err = x.SkipString(close)
			if err != nil {
				x.position = oldIndex
				return "", err
			}
			continue
		}

		if x.LookingAtString(open) {
			err = x.SkipString(open)
			if err != nil {
				x.position = oldIndex
				return "", err
			}
			openCount++
			continue
		}
		if x.IsExhausted() {
			x.position = oldIndex
			return "", fmt.Errorf("could not find closing rune %s, input exhausted", string(close))
		}

		x.MustSkip(1)
	}
	return x.MustExtract(oldIndex, x.position), nil
}

func (x *Parser) MustReadToMatchingString(open, close string) string {
	result, err := x.ReadToMatchingString(open, close)
	if err != nil {
		panic(err.Error())
	}
	return result
}

func (x *Parser) ReadToMatchingStringSkipDelims(open, close string) (string, error) {
	if !x.LookingAtString(open) {
		return "", fmt.Errorf("could not find opening string %s, %s", open, x.CurrentContext())
	}
	err := x.SkipString(open)
	if err != nil {
		return "", fmt.Errorf("could not skip opening string %s, %s", open, x.CurrentContext())
	}
	content, err := x.ReadToMatchingString(open, close)
	if err != nil {
		return "", err
	}

	if !x.LookingAtString(close) {
		return "", fmt.Errorf("could not find opening rune %s, %s", close, x.CurrentContext())
	}
	err = x.SkipString(close)
	if err != nil {
		return "", fmt.Errorf("could not skip opening rune %s, %s", close, x.CurrentContext())
	}

	return content, nil
}

func (x *Parser) MustReadToMatchingStringSkipDelims(open, close string) string {
	result, err := x.ReadToMatchingStringSkipDelims(open, close)
	if err != nil {
		panic(err.Error())
	}
	return result
}

func (x *Parser) ReadToMatchingRune(open, close rune) (string, error) {
	var (
		openCount          = 1
		runeCount          = 0
		differentRunes     = open != close
		r                  rune
		remainingRuneCount = x.RemainingRuneCount()
		newPos             int
	)

	for {
		if openCount == 0 {
			break
		}

		if differentRunes {
			if openCount < 0 {
				return "", fmt.Errorf("could not find closing rune %s, because they were used unbalanced, closed before opened (rune %s)", string(close), string(open))
			}
			if runeCount > remainingRuneCount {
				return "", fmt.Errorf("could not find closing rune %s, input exhausted", string(close))
			}
		}

		newPos = x.position + runeCount
		if len(x.input) <= newPos {
			return "", fmt.Errorf("could not find closing rune %s, input exhausted", string(close))
		}
		r = x.input[newPos]

		if differentRunes {
			switch r {
			// order is important here because close == open is possible and need precedence then
			case close:
				openCount--
			case open:
				openCount++
			}
		} else {
			if r == close {
				openCount--
			}
		}

		runeCount++
	}
	return string(x.MustReadRunes(max(0, runeCount-1))), nil
}

func (x *Parser) MustReadToMatchingRune(open, close rune) string {
	result, err := x.ReadToMatchingRune(open, close)
	if err != nil {
		panic(err.Error())
	}
	return result
}

func (x *Parser) ReadToMatchingRuneSkipDelims(open, close rune) (string, error) {
	if !x.LookingAtRune(open) {
		return "", fmt.Errorf("could not find opening rune %s, %s", string(open), x.CurrentContext())
	}
	err := x.Skip(1)
	if err != nil {
		return "", fmt.Errorf("could not skip opening rune %s, %s", string(open), x.CurrentContext())
	}
	content, err := x.ReadToMatchingRune(open, close)
	if err != nil {
		return "", err
	}

	if !x.LookingAtRune(close) {
		return "", fmt.Errorf("could not find opening rune %s, %s", string(close), x.CurrentContext())
	}
	err = x.Skip(1)
	if err != nil {
		return "", fmt.Errorf("could not skip opening rune %s, %s", string(close), x.CurrentContext())
	}

	return content, nil
}

func (x *Parser) MustReadToMatchingRuneSkipDelims(open, close rune) string {
	result, err := x.ReadToMatchingRuneSkipDelims(open, close)
	if err != nil {
		panic(err.Error())
	}
	return result
}

func (x *Parser) ReadToMatchingRuneEscaped(open, close, escape rune) (string, error) {
	var (
		openCount          = 1
		runeCount          = 0
		differentRunes     = open != close
		r                  rune
		rBefore            rune
		remainingRuneCount = x.RemainingRuneCount()
		result             = ""
	)

	for {
		if differentRunes {
			if openCount < 0 {
				return "", fmt.Errorf("could not find closing rune %v, because they were used unbalanced, closed before opened (rune %v)", close, open)
			}
			if runeCount > remainingRuneCount {
				return "", fmt.Errorf("could not find closing rune %v, input exhausted", close)
			}
		}

		r = x.input[x.position+runeCount]

		// check the escaping
		rBefore = x.input[x.position+runeCount-1]
		if rBefore == escape && (r == close || r == open) {
			result = result[0:len(result)-1] + string(r)
			runeCount++
			continue
		}

		if differentRunes {
			switch r {
			// order is important here because close == open is possible and need precedence then
			case close:
				openCount--
			case open:
				openCount++
			}
		} else {
			if r == close {
				openCount--
			}
		}

		if openCount == 0 {
			break
		}

		result += string(r)
		runeCount++
	}

	x.position += runeCount
	return result, nil
}

func (x *Parser) MustReadToMatchingRuneEscaped(open, close, escape rune) string {
	result, err := x.ReadToMatchingRuneEscaped(open, close, escape)
	if err != nil {
		panic(err.Error())
	}
	return result
}

func (x *Parser) ReadToMatchingRuneEscapedSkipDelims(open, close, escape rune) (string, error) {
	if !x.LookingAtRune(open) {
		return "", fmt.Errorf("could not find opening rune %s, %s", string(open), x.CurrentContext())
	}
	err := x.Skip(1)
	if err != nil {
		return "", fmt.Errorf("could not skip opening rune %s, %s", string(open), x.CurrentContext())
	}
	content, err := x.ReadToMatchingRuneEscaped(open, close, escape)
	if err != nil {
		return "", err
	}

	if !x.LookingAtRune(close) {
		return "", fmt.Errorf("could not find closing rune %s, %s", string(close), x.CurrentContext())
	}
	err = x.Skip(1)
	if err != nil {
		return "", fmt.Errorf("could not skip closing rune %s, %s", string(close), x.CurrentContext())
	}

	return content, nil
}

func (x *Parser) MustReadToMatchingRuneEscapedSkipDelims(open, close, escape rune) string {
	result, err := x.ReadToMatchingRuneEscapedSkipDelims(open, close, escape)
	if err != nil {
		panic(err.Error())
	}
	return result
}
