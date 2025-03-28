package textparser

// StartCapture starts collecting output, see Captured to retrieve it.
func (x *Parser) StartCapture() *Parser {
	x.captureStartPosition = x.CurrentIndex()
	return x
}

// Captured retrieves output since StartCapture was last called.
func (x *Parser) Captured() string {
	return string(x.input[x.captureStartPosition:x.CurrentIndex()])
}
