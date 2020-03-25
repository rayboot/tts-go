package errors

// APIError is the error that is returned from API calls.
type APIError struct {
	// Id is the request ID for which this error occurred.
	ID string `json:"id,omitempty"`

	// Status is the HTTP Status code for this error
	Status int `json:"status,omitempty"`

	// Code is the specific error code (for debugging purposes)
	Code string `json:"code,omitempty"`

	// Type is the type (BadRequest, NotFound, etc) of error
	Type string `json:"type,omitempty"`

	// Message is a descriptive message of the error, why it occurred, how to resolve, etc.
	Message string `json:"message,omitempty"`

	// Info is an optional field describing in detail the error for debugging purposes.
	Info string `json:"-"`
}

// Error converts the error to a human-readable, string format.
func (e APIError) Error() string {
	return e.Message
}
