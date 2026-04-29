package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateRequest_MissingAuthHeader(t *testing.T) {
	v := NewSignatureValidator("testkey", "testsecret")
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	err := v.ValidateRequest(req)
	if err != ErrMissingAuthHeader {
		t.Errorf("expected ErrMissingAuthHeader, got %v", err)
	}
}

func TestValidateRequest_InvalidAuthHeader(t *testing.T) {
	v := NewSignatureValidator("testkey", "testsecret")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "InvalidScheme credentials")

	err := v.ValidateRequest(req)
	if err != ErrInvalidAuthHeader {
		t.Errorf("expected ErrInvalidAuthHeader, got %v", err)
	}
}

func TestValidateRequest_MissingDateHeader(t *testing.T) {
	v := NewSignatureValidator("testkey", "testsecret")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "AWS4-HMAC-SHA256 Credential=testkey/20240101/us-east-1/s3/aws4_request")

	err := v.ValidateRequest(req)
	if err != ErrMissingDateHeader {
		t.Errorf("expected ErrMissingDateHeader, got %v", err)
	}
}

func TestValidateRequest_InvalidDateHeader(t *testing.T) {
	v := NewSignatureValidator("testkey", "testsecret")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "AWS4-HMAC-SHA256 Credential=testkey/20240101/us-east-1/s3/aws4_request")
	req.Header.Set("X-Amz-Date", "not-a-valid-date")

	err := v.ValidateRequest(req)
	if err != ErrInvalidDateHeader {
		t.Errorf("expected ErrInvalidDateHeader, got %v", err)
	}
}

func TestValidateRequest_Valid(t *testing.T) {
	v := NewSignatureValidator("testkey", "testsecret")
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "AWS4-HMAC-SHA256 Credential=testkey/20240101/us-east-1/s3/aws4_request")
	req.Header.Set("X-Amz-Date", "20240101T120000Z")

	err := v.ValidateRequest(req)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestDeriveSigningKey(t *testing.T) {
	key := DeriveSigningKey("wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY", "20150830", "us-east-1", "iam")
	if len(key) != 32 {
		t.Errorf("expected 32 byte key, got %d", len(key))
	}
}

func TestAuthMiddleware_BlocksInvalidRequest(t *testing.T) {
	v := NewSignatureValidator("testkey", "testsecret")
	mw := AuthMiddleware(v)

	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	mw(next).ServeHTTP(rec, req)

	if called {
		t.Error("handler should not have been called with missing auth header")
	}
	if rec.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rec.Code)
	}
}
