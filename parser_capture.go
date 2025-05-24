package textparser

import "fmt"

// StartCapture starts collecting output, see Captured to retrieve it.
func (x *Parser) StartCapture() *Parser {
	x.captureStartPosition = x.CurrentIndex()
	return x
}

// Captured retrieves output since StartCapture was last called.
func (x *Parser) Captured() string {
	return x.MustExtract(x.captureStartPosition, x.CurrentIndex())
}

func (x *Parser) MustExtract(start, end int) string {
	v, err := x.Extract(start, end)
	if err != nil {
		panic(err)
	}
	return v
}

// Extract gets a substring of the input by indices.
func (x *Parser) Extract(start, end int) (string, error) {
	if start < 0 {
		return "", fmt.Errorf("start index must be positive, got %d", start)
	}
	if end <= start {
		return "", fmt.Errorf("end index must be greater than start index, got start %d, end %d", start, end)
	}
	if end > len(x.input) {
		return "", fmt.Errorf("end index must be less than input length, got %d, input length %d", end, len(x.input))
	}
	return string(x.input[start:end]), nil
}
