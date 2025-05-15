package textparser

import (
	"fmt"
)

// TODO ReadToMatchingString(open, close string) (string, error) {

func (x *Parser) ReadToMatching(open, close rune) (string, error) {
	var (
		openCount          = 1
		runeCount          = 0
		differentRunes     = open != close
		r                  rune
		remainingRuneCount = x.RemainingRuneCount()
	)

	for {
		if openCount == 0 {
			break
		}

		if differentRunes {
			if openCount < 0 {
				return "", fmt.Errorf("could not find closing rune %v, because they were used unbalanced, closed before opened (rune %v)", close, open)
			}
			if runeCount > remainingRuneCount {
				return "", fmt.Errorf("could not find closing rune %v, input exhausted", close)
			}
		}

		r = x.input[x.position+runeCount]

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

func (x *Parser) MustReadToMatching(open, close rune) string {
	result, err := x.ReadToMatching(open, close)
	if err != nil {
		panic(err.Error())
	}
	return result
}

func (x *Parser) ReadToMatchingSkipDelims(open, close rune) (string, error) {
	if !x.LookingAtRune(open) {
		return "", fmt.Errorf("could not find opening rune %s, %s", string(open), x.CurrentContext())
	}
	err := x.Skip(1)
	if err != nil {
		return "", fmt.Errorf("could not skip opening rune %s, %s", string(open), x.CurrentContext())
	}
	content, err := x.ReadToMatching(open, close)
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

func (x *Parser) MustReadToMatchingSkipDelims(open, close rune) string {
	result, err := x.ReadToMatchingSkipDelims(open, close)
	if err != nil {
		panic(err.Error())
	}
	return result
}

func (x *Parser) ReadToMatchingEscaped(open, close, escape rune) (string, error) {
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

func (x *Parser) MustReadToMatchingEscaped(open, close, escape rune) string {
	result, err := x.ReadToMatchingEscaped(open, close, escape)
	if err != nil {
		panic(err.Error())
	}
	return result
}

func (x *Parser) ReadToMatchingEscapedSkipDelims(open, close, escape rune) (string, error) {
	if !x.LookingAtRune(open) {
		return "", fmt.Errorf("could not find opening rune %s, %s", string(open), x.CurrentContext())
	}
	err := x.Skip(1)
	if err != nil {
		return "", fmt.Errorf("could not skip opening rune %s, %s", string(open), x.CurrentContext())
	}
	content, err := x.ReadToMatchingEscaped(open, close, escape)
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

func (x *Parser) MustReadToMatchingEscapedSkipDelims(open, close, escape rune) string {
	result, err := x.ReadToMatchingEscapedSkipDelims(open, close, escape)
	if err != nil {
		panic(err.Error())
	}
	return result
}
