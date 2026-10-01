package document

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"auth_info/internal/config"
	"auth_info/internal/pkg/apperr"
)

func TestLocalResourcesConfineTemplatesAndObserveCancellation(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "sample.json"), []byte("fixture"), 0600); err != nil {
		t.Fatal(err)
	}
	resources, err := NewResources(config.DocumentConfig{TemplateDir: dir, ImageTimeout: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := resources.Close(); err != nil {
			t.Error(err)
		}
	})
	data, err := resources.ReadTemplate(context.Background(), "sample", ".json")
	if err != nil || string(data) != "fixture" {
		t.Fatalf("read: %s %v", data, err)
	}
	for _, name := range []string{"../sample", `..\sample`, "", ".."} {
		_, err := resources.ReadTemplate(context.Background(), name, ".json")
		if !apperr.IsCode(err, apperr.CodeInvalidArgument) {
			t.Fatalf("unsafe path accepted: %q %v", name, err)
		}
	}
	outside := filepath.Join(t.TempDir(), "outside.json")
	if err := os.WriteFile(outside, []byte("outside"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(dir, "escape.json")); err != nil {
		t.Fatal(err)
	}
	if _, err := resources.ReadTemplate(context.Background(), "escape", ".json"); err == nil {
		t.Fatal("symlink escaped root")
	}
	_, err = resources.ReadTemplate(context.Background(), "absent", ".json")
	if !apperr.IsCode(err, apperr.CodeNotFound) {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := resources.ReadTemplate(ctx, "sample", ".json"); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}

func TestImageSuccessFailureAndCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/missing" {
			w.WriteHeader(404)
			return
		}
		_, _ = w.Write([]byte("image"))
	}))
	defer srv.Close()
	r := &Resources{client: srv.Client()}
	data, err := r.FetchImage(context.Background(), srv.URL, 10)
	if err != nil || string(data) != "image" {
		t.Fatalf("%s %v", data, err)
	}
	if _, err := r.FetchImage(context.Background(), srv.URL+"/missing", 10); err == nil {
		t.Fatal("missing status error")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := r.FetchImage(ctx, srv.URL, 10); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}
