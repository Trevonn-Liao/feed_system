package errs

import "net/http"

const (
	CodeBadRequest        = 40000
	CodeUnauthorized      = 40100
	CodeUserNotFound      = 40101
	CodePasswordIncorrect = 40102
	CodeRefreshInvalid    = 40103
	CodeTooManyRequests   = 42900
	CodeUsernameExists    = 40901
	CodeInternal          = 50000
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

func Unauthorized(code int, message string) *AppError {
	return New(code, message, http.StatusUnauthorized)
}

func UserNotFound() *AppError {
	return Unauthorized(CodeUserNotFound, "user not found")
}

func PasswordIncorrect() *AppError {
	return Unauthorized(CodePasswordIncorrect, "password incorrect")
}

func RefreshTokenInvalid() *AppError {
	return Unauthorized(CodeRefreshInvalid, "refresh token invalid")
}

func TooManyRequests() *AppError {
	return New(CodeTooManyRequests, "too many repeated requests", http.StatusTooManyRequests)
}

func UsernameExists() *AppError {
	return New(CodeUsernameExists, "username already exists", http.StatusConflict)
}

func Internal(message string) *AppError {
	return New(CodeInternal, message, http.StatusInternalServerError)
}
