package middlewares

import (
	"encoding/xml"
	"net/http"
)

// ErrorResponse represents an S3-compatible XML error response.
type ErrorResponse struct {
	XMLName   xml.Name `xml:"Error"`
	Code      string   `xml:"Code"`
	Message   string   `xml:"Message"`
	RequestID string   `xml:"RequestId"`
}

// AuthMiddleware returns an HTTP middleware that validates AWS SigV4 signatures.
func AuthMiddleware(validator *SignatureValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := validator.ValidateRequest(r); err != nil {
				writeAuthError(w, err)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func writeAuthError(w http.ResponseWriter, err error) {
	var code int
	var errCode, message string

	switch err {
	case ErrMissingAuthHeader, ErrInvalidAuthHeader:
		code = http.StatusForbidden
		errCode = "InvalidSecurity"
		message = "The provided security credentials are not valid."
	case ErrMissingDateHeader, ErrInvalidDateHeader:
		code = http.StatusBadRequest
		errCode = "InvalidArgument"
		message = "Invalid or missing X-Amz-Date header."
	case ErrSignatureMismatch:
		code = http.StatusForbidden
		errCode = "SignatureDoesNotMatch"
		message = "The request signature we calculated does not match the signature you provided."
	case ErrExpiredRequest:
		code = http.StatusForbidden
		errCode = "RequestExpired"
		message = "Request has expired."
	default:
		code = http.StatusInternalServerError
		errCode = "InternalError"
		message = "An internal error occurred."
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(code)
	xml.NewEncoder(w).Encode(ErrorResponse{
		Code:    errCode,
		Message: message,
	})
}
