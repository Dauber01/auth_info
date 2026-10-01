// Package document implements local and HTTP document resource access.
package document

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"auth_info/internal/config"
	"auth_info/internal/pkg/apperr"
)

// Resources owns a confined template directory and a dedicated HTTP transport.
type Resources struct {
	root      *os.Root
	font      []byte
	client    *http.Client
	transport *http.Transport
}

// NewResources opens configured resources; callers own Close on success.
func NewResources(cfg config.DocumentConfig) (*Resources, error) {
	if cfg.ImageTimeout <= 0 {
		return nil, fmt.Errorf("image timeout must be positive")
	}
	var font []byte
	var err error
	if cfg.FontPath != "" {
		font, err = os.ReadFile(cfg.FontPath)
		if err != nil {
			return nil, fmt.Errorf("read configured font: %w", err)
		}
	}
	root, err := os.OpenRoot(cfg.TemplateDir)
	if err != nil {
		return nil, fmt.Errorf("open template directory: %w", err)
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	return &Resources{
		root: root, font: font, transport: transport,
		client: &http.Client{Transport: transport, Timeout: cfg.ImageTimeout},
	}, nil
}

// Close releases directory handles and idle connections after request processing ends.
func (r *Resources) Close() error { r.transport.CloseIdleConnections(); return r.root.Close() }

// ReadTemplate permits a single template name, confined even across symbolic links.
func (r *Resources) ReadTemplate(ctx context.Context, name, extension string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if name == "" || name == "." || name == ".." || strings.ContainsAny(name, "/\\\x00") ||
		(extension != ".json" && extension != ".docx") {
		return nil, apperr.New(apperr.CodeInvalidArgument, "invalid template name")
	}
	file, err := r.root.Open(name + extension)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, apperr.New(apperr.CodeNotFound, "template not found: "+name)
		}
		return nil, apperr.Wrap(apperr.CodeInternal, "failed to open template", err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, "failed to read template", err)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return data, nil
}

// ReadFont returns the configured font, or nil to use the renderer's built-in font.
func (r *Resources) ReadFont(ctx context.Context) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return r.font, nil
}

// FetchImage propagates request cancellation and bounds both known and chunked responses.
func (r *Resources) FetchImage(ctx context.Context, url string, limit int64) ([]byte, error) {
	if limit <= 0 {
		return nil, apperr.New(apperr.CodeInvalidArgument, "invalid image size limit")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeInvalidArgument, "invalid image URL", err)
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, "failed to fetch image", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("image response status: %d", resp.StatusCode)
	}
	if resp.ContentLength > limit {
		return nil, apperr.New(apperr.CodeInvalidArgument, "image size exceeds limit")
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, limit+1))
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeInternal, "failed to read image", err)
	}
	if int64(len(data)) > limit {
		return nil, apperr.New(apperr.CodeInvalidArgument, "image size exceeds limit")
	}
	return data, nil
}
