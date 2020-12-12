package errors

type LatteError interface {
	Error() string
	CliMessage() string
	ErrorName() string
}
