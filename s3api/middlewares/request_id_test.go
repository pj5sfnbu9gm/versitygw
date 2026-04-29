// Copyright 2023 Versity Software
// This file is licensed under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific
// language governing permissions and limitations under the License.

package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerateRequestID(t *testing.T) {
	id1, err := generateRequestID()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(id1) != 32 {
		t.Errorf("expected length 32, got %d", len(id1))
	}

	id2, err := generateRequestID()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id1 == id2 {
		t.Error("expected unique IDs, got duplicates")
	}
}

func TestGetRequestID_Missing(t *testing.T) {
	ctx := context.Background()
	if id := GetRequestID(ctx); id != "" {
		t.Errorf("expected empty string, got %q", id)
	}
}

func TestRequestIDMiddleware_SetsHeader(t *testing.T) {
	handler := RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	requestID := rec.Header().Get(RequestIDHeader)
	if requestID == "" {
		t.Error("expected request ID header to be set")
	}
	if len(requestID) != 32 {
		t.Errorf("expected request ID length 32, got %d", len(requestID))
	}
}

func TestRequestIDMiddleware_ContextValue(t *testing.T) {
	var capturedID string
	handler := RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID = GetRequestID(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	headerID := rec.Header().Get(RequestIDHeader)
	if capturedID == "" {
		t.Error("expected request ID in context")
	}
	if capturedID != headerID {
		t.Errorf("context ID %q does not match header ID %q", capturedID, headerID)
	}
}
