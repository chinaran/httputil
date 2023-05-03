package httputil

import (
	"errors"
	"fmt"
)

// RequestError httputil request error
type RequestError struct {
	// http status code
	Code int
	// error message
	Message string
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("Status Code %d, Message: %s", e.Code, e.Message)
}

// NewRequestError new httputil error
func NewRequestError(statusCode int, message string) error {
	return &RequestError{Code: statusCode, Message: message}
}

// GetErrorCode get http request's status code
func GetErrorCode(err error) (int, bool) {
	var tmp *RequestError
	if errors.As(err, &tmp) {
		return tmp.Code, true
	}
	return 0, false
}

// GetErrorMessage get http request's error message
func GetErrorMessage(err error) string {
	var tmp *RequestError
	if errors.As(err, &tmp) {
		return tmp.Message
	}
	return err.Error()
}

// IsErrorCode juede http request's stauts code
func IsErrorCode(err error, code int) bool {
	var tmp *RequestError
	if errors.As(err, &tmp) {
		if tmp.Code == code {
			return true
		}
	}
	return false
}
