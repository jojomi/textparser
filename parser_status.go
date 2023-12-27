package textparser

// CurrentIndex returns the current internal parser index aka how many runes have been processed already.
func (x *Parser) CurrentIndex() int {
	return x.position
}

func (x *Parser) IsExhausted() bool {
	return x.CurrentIndex() >= len(x.input)
}

func (x *Parser) HasMore() bool {
	return !x.IsExhausted()
}

func (x *Parser) Processed() string {
	return string(x.input[0:x.position])
}

func (x *Parser) Remaining() string {
	return string(x.input[x.position:])
}

func (x *Parser) RemainingRuneCount() int {
	return len(x.input) - x.position
}

func (x *Parser) String() string {
	return x.Processed() + "|>" + x.Remaining()
}

func (x *Parser) CurrentContext() string {
	const contextLength = 10
	start := max(0, x.position-contextLength)
	end := min(x.position+contextLength, len(x.input))
	result := string(x.input[start:x.position]) + "|>" + string(x.input[x.position:end])
	return result
}
