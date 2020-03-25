// Package error is a collection of error-related functionality that aims to
// unify all of the errors across the SDK. Specifically, it abstracts away low
// level errors into higher-level, easier-to-digest errors of two types: SDK
// errors (represented by the `Error` class) and API errors (represented by the
// `APIError` class).
package errors

// Error is a generic error that is returned when something SDK-related
// goes wrong.
type Error struct {
	// Code is the specific error code (for debugging purposes)
	Code string `json:"code,omitempty"`

	// Message is a descriptive message of the error, why it occurred, how to resolve, etc.
	Message string `json:"message,omitempty"`

	// Info is an optional field describing in detail the error for debugging purposes.
	Info string `json:"-"`
}

// Error converts the error to a human-readable, string format.
func (e Error) Error() string {
	return e.Message
}

// NewError creates and `Error` object from the given information.
func NewError(code string, message string, info string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Info:    info,
	}
}

// NewFromErrorCode creates an `Error` object based on a predefined code and
// message.
func NewFromErrorCode(code ErrorCode) *Error {
	return &Error{
		Code:    string(code),
		Message: errorMessages[code],
	}
}

// NewFromErrorCodeInfo creates an `Error` object based on a predefined code
// and also includes some extra information about the error.
func NewFromErrorCodeInfo(code ErrorCode, info string) *Error {
	return &Error{
		Code:    string(code),
		Message: errorMessages[code],
		Info:    info,
	}
}
