package utils

type InternalError struct {
	Code    int
	Message string
}

func (ie *InternalError) Error() string {
	return ie.Message
}

func NewInternalError(code int, message string) error {
	return &InternalError{
		Code:    code,
		Message: message,
	}
}
