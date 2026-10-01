package document

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"testing"

	"auth_info/internal/pkg/apperr"
)

type memoryResources struct {
	templates map[string][]byte
	err       error
}

func (r *memoryResources) ReadTemplate(ctx context.Context, name, ext string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if r.err != nil {
		return nil, r.err
	}
	if data, ok := r.templates[name+ext]; ok {
		return data, nil
	}
	return nil, apperr.New(apperr.CodeNotFound, "missing template")
}
func (*memoryResources) ReadFont(context.Context) ([]byte, error) { return nil, nil }
func (r *memoryResources) FetchImage(context.Context, string, int64) ([]byte, error) {
	return nil, r.err
}

func testResources(t *testing.T) *memoryResources {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, content := range map[string]string{
		"[Content_Types].xml":          `<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="png" ContentType="image/png"/></Types>`,
		"word/document.xml":            `<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body><w:p><w:r><w:t>{Name}</w:t></w:r></w:p><w:p><w:r><w:t>{Signature}</w:t></w:r></w:p></w:body></w:document>`,
		"word/_rels/document.xml.rels": `<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"></Relationships>`,
	} {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return &memoryResources{templates: map[string][]byte{
		"word_template_test.docx": buf.Bytes(),
		"example_template.json": []byte(`{"title":"Invoice","sections":[
			{"type":"paragraph","content":"Hello"},{"type":"image","data":"{{.SignatureBase64}}"}]}`),
	}}
}

func TestResourcesFailureMissingAndCancellation(t *testing.T) {
	failure := errors.New("resource failed")
	uc := NewUseCase(&memoryResources{err: failure})
	if _, err := uc.GeneratePDF(context.Background(), "demo", nil); !errors.Is(err, failure) {
		t.Fatal(err)
	}
	uc = NewUseCase(testResources(t))
	_, err := uc.GenerateWord(context.Background(), "absent", WordTemplateData{})
	if !apperr.IsCode(err, apperr.CodeNotFound) {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := uc.GeneratePDF(ctx, "example_template", nil); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	if _, err := uc.GenerateWord(ctx, "word_template_test", WordTemplateData{}); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}
