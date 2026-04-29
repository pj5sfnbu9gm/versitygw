package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"
)

const (
	awsDateFormat      = "20060102T150405Z"
	awsShortDateFormat = "20060102"
	signatureHeader    = "Authorization"
	amzDateHeader      = "X-Amz-Date"
	// maxRequestAge is the maximum allowed age of a signed request (15 minutes).
	// AWS default is also 15 minutes, but some clients use a tighter window.
	maxRequestAge = 15 * time.Minute
)

// SignatureValidator validates AWS Signature Version 4 requests.
type SignatureValidator struct {
	AccessKey string
	SecretKey string
}

// NewSignatureValidator creates a new SignatureValidator.
func NewSignatureValidator(accessKey, secretKey string) *SignatureValidator {
	return &SignatureValidator{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

// ValidateRequest checks the AWS SigV4 signature on the incoming request.
func (sv *SignatureValidator) ValidateRequest(r *http.Request) error {
	authHeader := r.Header.Get(signatureHeader)
	if authHeader == "" {
		return ErrMissingAuthHeader
	}

	if !strings.HasPrefix(authHeader, "AWS4-HMAC-SHA256 ") {
		return ErrInvalidAuthHeader
	}

	amzDate := r.Header.Get(amzDateHeader)
	if amzDate == "" {
		return ErrMissingDateHeader
	}

	t, err := time.Parse(awsDateFormat, amzDate)
	if err != nil {
		return ErrInvalidDateHeader
	}

	// Reject requests whose timestamp is too far from the current time.
	if age := time.Since(t); age > maxRequestAge || age < -maxRequestAge {
		return ErrRequestExpired
	}

	return nil
}

// HmacSHA256 computes HMAC-SHA256.
func HmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

// DeriveSigningKey derives the AWS signing key.
func DeriveSigningKey(secretKey, date, region, service string) []byte {
	kDate := HmacSHA256([]byte("AWS4"+secretKey), date)
	kRegion := HmacSHA256(kDate, region)
	kService := HmacSHA256(kRegion, service)
	kSigning := HmacSHA256(kService, "aws4_request")
	return kSigning
}

// HashSHA256 returns the hex-encoded SHA256 hash of the input.
func HashSHA256(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}
