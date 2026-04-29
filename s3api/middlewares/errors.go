package middlewares

import "errors"

// Sentinel errors for authentication middleware.
var (
	ErrMissingAuthHeader  = errors.New("missing Authorization header")
	ErrInvalidAuthHeader  = errors.New("invalid Authorization header format")
	ErrMissingDateHeader  = errors.New("missing X-Amz-Date header")
	ErrInvalidDateHeader  = errors.New("invalid X-Amz-Date header format")
	ErrSignatureMismatch  = errors.New("signature does not match")
	ErrExpiredRequest     = errors.New("request has expired")
	ErrInvalidAccessKey   = errors.New("invalid access key")
)

// AuthError wraps an authentication error with an HTTP status code.
type AuthError struct {
	Code    int
	Message string
	Err     error
}

func (e *AuthError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AuthError) Unwrap() error {
	return e.Err
}

// NewAuthError creates a new AuthError.
func NewAuthError(code int, message string, err error) *AuthError {
	return &AuthError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
