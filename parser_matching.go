package textparser

import "fmt"

// TODO ReadToMatchingString(open, close string) (string, error) {

func (x *Parser) ReadToMatching(open, close rune) (string, error) {
	var (
		openCount          = 1
		runeCount          = 0
		r                  rune
		remainingRuneCount = x.RemainingRuneCount()
	)

	for {
		if openCount < 0 {
			return "", fmt.Errorf("could not find closing rune %v, because they were used unbalanced, closed before opened (rune %v)", close, open)
		}
		if openCount == 0 {
			break
		}
		if runeCount > remainingRuneCount {
			return "", fmt.Errorf("could not find closing rune %v, input exhausted", close)
		}

		r = x.input[x.position+runeCount]
		switch r {
		// order is important here because close == open is possible and need precedence then
		case close:
			openCount--
		case open:
			openCount++
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
