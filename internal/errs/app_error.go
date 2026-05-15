package errs

import "net/http"

const (
	CodeBadRequest = 40000
	CodeInternal   = 50000
)

type AppError struct {
	Code       int
	Message    string
	HTTPStatus int
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}

	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code int, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

func Wrap(err error, code int, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Err:        err,
	}
}

func BadRequest(message string) *AppError {
	return New(CodeBadRequest, message, http.StatusBadRequest)
}

func Internal(message string) *AppError {
	return New(CodeInternal, message, http.StatusInternalServerError)
}
