package document

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"auth_info/internal/pkg/apperr"
)

const maxImageBytes = 10 * 1024 * 1024

func TestFetchImageBytes_URLContentLengthTooLarge(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", "10485761")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	resources := &Resources{client: srv.Client()}

	_, err := resources.FetchImage(context.Background(), srv.URL, maxImageBytes)
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

	resources := &Resources{client: srv.Client()}

	_, err := resources.FetchImage(context.Background(), srv.URL, maxImageBytes)
	if !apperr.IsCode(err, apperr.CodeInvalidArgument) {
		t.Fatalf("expected invalid argument error, got: %v", err)
	}
}
