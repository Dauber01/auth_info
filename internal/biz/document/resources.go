package document

import "context"

// Resources is the document use case's external resource boundary.
// Implementations must observe ctx and enforce the requested download size limit.
type Resources interface {
	ReadTemplate(ctx context.Context, name, extension string) ([]byte, error)
	ReadFont(ctx context.Context) ([]byte, error)
	FetchImage(ctx context.Context, url string, limit int64) ([]byte, error)
}
