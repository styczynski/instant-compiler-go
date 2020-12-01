package errors

type LatteError interface {
	Error() string
	CliMessage() string
}

type LatteSimpleError struct {
	message string
}

func (e *LatteSimpleError) Error() string {
	return e.message
}

func (e *LatteSimpleError) CliMessage() string {
	return e.message
}

func NewLatteSimpleError(err error) *LatteSimpleError {
	return &LatteSimpleError{
		message: err.Error(),
	}
}