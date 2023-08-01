package customerror

type CustomError struct {
	Field   string
	Message string
}

func (e CustomError) Error() string {
	return e.Message
}
