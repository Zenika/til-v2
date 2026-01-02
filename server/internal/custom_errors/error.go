package custom_errors

import (
	"errors"
)

type CustomError struct {
	Content   string
	ErrorCode int
	BaseError error
}

func (e *CustomError) Error() string {
	return e.Content
}

func NewInvalidArgumentError(content string, base error) *CustomError {
	return &CustomError{
		Content:   content,
		BaseError: base,
		ErrorCode: 400,
	}
}

func NewPermissionDeniedError(content string, base error) *CustomError {
	return &CustomError{
		Content:   content,
		ErrorCode: 403,
		BaseError: base,
	}
}

func NewNotFoundError(content string, base error) *CustomError {
	return &CustomError{
		Content:   content,
		ErrorCode: 404,
		BaseError: base,
	}
}

func NewInternalServerError(content string, base error) *CustomError {
	return &CustomError{
		Content:   content,
		ErrorCode: 500,
		BaseError: base,
	}
}

func IsNotFoundError(err error) bool {
	var customError *CustomError
	ok := errors.As(err, &customError)
	return ok && customError.ErrorCode == 404
}
