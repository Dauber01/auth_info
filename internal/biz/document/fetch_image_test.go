package document

import (
	"bytes"
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"auth_info/internal/apperr"
)

func TestFetchImageBytes_URLContentLengthTooLarge(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", "10485761")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	uc := NewUseCase()
	uc.httpClient = srv.Client()

	_, err := uc.fetchImageBytes(context.Background(), ImageValue{ImageURL: srv.URL})
	if !apperr.IsCode(err, apperr.CodeInvalidArgument) {
		t.Fatalf("expected invalid argument error, got: %v", err)
	}
}

func TestFetchImageBytes_URLBodyTooLargeWithoutContentLength(t *testing.T) {
	largeBody := bytes.Repeat([]byte("a"), maxImageBytes+1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		_, _ = w.Write(largeBody)
	}))
	defer srv.Close()

	uc := NewUseCase()
	uc.httpClient = srv.Client()

	_, err := uc.fetchImageBytes(context.Background(), ImageValue{ImageURL: srv.URL})
	if !apperr.IsCode(err, apperr.CodeInvalidArgument) {
		t.Fatalf("expected invalid argument error, got: %v", err)
	}
}

func TestFetchImageBytes_Base64TooLarge(t *testing.T) {
	encoded := base64.StdEncoding.EncodeToString(bytes.Repeat([]byte("a"), maxImageBytes+1))
	imageURL := "data:image/png;base64," + encoded

	uc := NewUseCase()
	_, err := uc.fetchImageBytes(context.Background(), ImageValue{ImageURL: imageURL})
	if !apperr.IsCode(err, apperr.CodeInvalidArgument) {
		t.Fatalf("expected invalid argument error, got: %v", err)
	}
}
