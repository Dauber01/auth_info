package document

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"
)

// injectImagesToDocx 向 docx 模板字节流中注入图片，直接替换 {key} 占位符。
func (uc *UseCase) injectImagesToDocx(ctx context.Context, docxBytes []byte, imageData map[string]ImageValue) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(docxBytes), int64(len(docxBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to read docx zip: %w", err)
	}

	files := make(map[string][]byte, len(zr.File))
	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open zip entry %s: %w", f.Name, err)
		}
		var buf bytes.Buffer
		buf.ReadFrom(rc)
		rc.Close()
		files[f.Name] = buf.Bytes()
	}

	relsFile := "word/_rels/document.xml.rels"
	relsXML := string(files[relsFile])
	docXML := string(files["word/document.xml"])

	for _, p := range indexedImagePairs(imageData) {
		imgBytes, err := uc.fetchImageBytes(ctx, p.val)
		if err != nil {
			return nil, fmt.Errorf("image key %q: %w", p.key, err)
		}

		imgExt := detectImageExtension(imgBytes)
		mediaName := fmt.Sprintf("image%d.%s", p.idx+1, imgExt)
		files["word/media/"+mediaName] = imgBytes

		rId := fmt.Sprintf("rIdImg%d", p.idx+1)
		relEntry := fmt.Sprintf(
			`<Relationship Id="%s" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="media/%s"/>`,
			rId, mediaName,
		)
		relsXML = strings.Replace(relsXML, "</Relationships>", relEntry+"</Relationships>", 1)

		cx, cy := imageEMU(imgBytes, p.val)
		drawingXML := buildDrawingXML(rId, p.idx+1, cx, cy)

		marker := "{" + p.key + "}"
		re := regexp.MustCompile(
			`<w:r\b[^>]*>(?:<w:rPr>[\s\S]*?</w:rPr>)?<w:t[^>]*>` +
				regexp.QuoteMeta(marker) +
				`</w:t></w:r>`,
		)
		docXML = re.ReplaceAllString(docXML, `<w:r><w:drawing>`+drawingXML+`</w:drawing></w:r>`)
	}

	files[relsFile] = []byte(relsXML)
	files["word/document.xml"] = []byte(docXML)

	return repackDocxZip(zr, files)
}
