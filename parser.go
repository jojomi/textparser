package textparser

import (
	"io"
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
func NewParserFromReader(input io.Reader) (*Parser, error) {
	content, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	return &Parser{
		input:    []rune(string(content)),
		position: 0,
	}, nil
}
