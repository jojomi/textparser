package textparser

type EndOfInputError struct{}

func (x EndOfInputError) Error() string {
	return "end of input"
}
