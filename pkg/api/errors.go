// Package api owns the HTTP handlers and the JSON conventions: the {data}
// and {error} envelopes, the closed error code set, cursor pagination, and
// the auth endpoints with their session middleware.
package api

import (
	"fmt"
	"net/http"
)

// Code is one of the closed set of API error codes. The HTTP status always
// matches the code, so clients switch on the code alone.
type Code string

// The full code set; handlers never invent new ones.
const (
	CodeBadRequest   Code = "bad_request"
	CodeUnauthorized Code = "unauthorized"
	CodeForbidden    Code = "forbidden"
	CodeNotFound     Code = "not_found"
	CodeConflict     Code = "conflict"
	CodeRateLimited  Code = "rate_limited"
	CodeInternal     Code = "internal"
)

func (c Code) status() int {
	switch c {
	case CodeBadRequest:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeConflict:
		return http.StatusConflict
	case CodeRateLimited:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

// apiError is what handlers return internally. Message renders to the
// client; Err is the wrapped cause and stays server-side, for the logs only.
type apiError struct {
	Code    Code
	Message string
	Status  int // optional override, e.g. 415 with code bad_request
	Err     error
}

func (e *apiError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *apiError) Unwrap() error {
	return e.Err
}

// errf builds a client-visible error with a formatted message.
func errf(code Code, format string, args ...any) *apiError {
	return &apiError{Code: code, Message: fmt.Sprintf(format, args...)}
}

// internalErr wraps a server-side failure; the client sees only a generic
// message, the cause goes to the log.
func internalErr(err error) *apiError {
	return &apiError{Code: CodeInternal, Message: "internal error", Err: err}
}
