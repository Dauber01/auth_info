package document

import (
	"archive/zip"
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"regexp"
	"sort"
	"strings"
)

// detectImageExtension 根据图片字节流的魔数判断扩展名
func detectImageExtension(imgBytes []byte) string {
	if len(imgBytes) < 4 {
		return "png"
	}
	if imgBytes[0] == 0x89 && imgBytes[1] == 0x50 && imgBytes[2] == 0x4E && imgBytes[3] == 0x47 {
		return "png"
	}
	if imgBytes[0] == 0xFF && imgBytes[1] == 0xD8 && imgBytes[2] == 0xFF {
		return "jpg"
	}
	if imgBytes[0] == 0x47 && imgBytes[1] == 0x49 && imgBytes[2] == 0x46 {
		return "gif"
	}
	return "png"
}

// imageEMU 解码图片获取实际宽高，根据 ImageValue 配置计算 EMU 尺寸
func imageEMU(imgBytes []byte, imgVal ImageValue) (int64, int64) {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(imgBytes))
	if err != nil || cfg.Width == 0 {
		return 1800000, 900000
	}

	var targetW, targetH float64

	if imgVal.OriginalSize {
		targetW = float64(cfg.Width) * 914400 / 96
		targetH = float64(cfg.Height) * 914400 / 96
	} else {
		targetW = 1800000
		targetH = targetW * float64(cfg.Height) / float64(cfg.Width)
	}

	if imgVal.MaxWidthPx > 0 {
		maxW := imgVal.MaxWidthPx * 914400 / 96
		if targetW > maxW {
			scale := maxW / targetW
			targetW = maxW
			targetH *= scale
		}
	}
	if imgVal.MaxHeightPx > 0 {
		maxH := imgVal.MaxHeightPx * 914400 / 96
		if targetH > maxH {
			scale := maxH / targetH
			targetH = maxH
			targetW *= scale
		}
	}

	return int64(targetW), int64(targetH)
}

// buildDrawingXML 构建内联图片的 <wp:inline> XML 片段
func buildDrawingXML(rId string, idx int, cx, cy int64) string {
	name := fmt.Sprintf("Image%d", idx)
	return fmt.Sprintf(
		`<wp:inline xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing" distT="0" distB="0" distL="0" distR="0">`+
			`<wp:extent cx="%d" cy="%d"/>`+
			`<wp:effectExtent l="0" t="0" r="0" b="0"/>`+
			`<wp:docPr id="%d" name="%s"/>`+
			`<wp:cNvGraphicFramePr><a:graphicFrameLocks xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" noChangeAspect="1"/></wp:cNvGraphicFramePr>`+
			`<a:graphic xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">`+
			`<a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/picture">`+
			`<pic:pic xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture">`+
			`<pic:nvPicPr><pic:cNvPr id="%d" name="%s"/><pic:cNvPicPr/></pic:nvPicPr>`+
			`<pic:blipFill>`+
			`<a:blip xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" r:embed="%s"/>`+
			`<a:stretch><a:fillRect/></a:stretch>`+
			`</pic:blipFill>`+
			`<pic:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="%d" cy="%d"/></a:xfrm>`+
			`<a:prstGeom prst="rect"><a:avLst/></a:prstGeom></pic:spPr>`+
			`</pic:pic></a:graphicData></a:graphic></wp:inline>`,
		cx, cy, idx, name, idx, name, rId, cx, cy,
	)
}

// repackDocxZip 按原始 ZIP 条目顺序重新打包，新增文件追加在末尾
func repackDocxZip(original *zip.Reader, files map[string][]byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	written := make(map[string]bool, len(original.File))
	for _, f := range original.File {
		written[f.Name] = true
		w, err := zw.Create(f.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to create zip entry %s: %w", f.Name, err)
		}
		if _, err := w.Write(files[f.Name]); err != nil {
			return nil, fmt.Errorf("failed to write zip entry %s: %w", f.Name, err)
		}
	}

	for name, content := range files {
		if written[name] {
			continue
		}
		w, err := zw.Create(name)
		if err != nil {
			return nil, fmt.Errorf("failed to create new zip entry %s: %w", name, err)
		}
		if _, err := w.Write(content); err != nil {
			return nil, fmt.Errorf("failed to write new zip entry %s: %w", name, err)
		}
	}

	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("failed to close zip writer: %w", err)
	}
	return buf.Bytes(), nil
}

// indexedImagePairs 将 map[string]ImageValue 转为有稳定顺序的序列
func indexedImagePairs(m map[string]ImageValue) []struct {
	idx int
	key string
	val ImageValue
} {
	result := make([]struct {
		idx int
		key string
		val ImageValue
	}, 0, len(m))

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, k := range keys {
		result = append(result, struct {
			idx int
			key string
			val ImageValue
		}{idx: i, key: k, val: m[k]})
	}
	return result
}

// injectRichTextToDocx 将富文本占位符 {key} 替换为带格式的 w:r 片段
func injectRichTextToDocx(docxBytes []byte, richData map[string]RichText) ([]byte, error) {
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

	docXML := string(files["word/document.xml"])

	for k, rt := range richData {
		marker := "{" + k + "}"
		re := regexp.MustCompile(
			`<w:r\b[^>]*>(?:<w:rPr>[\s\S]*?</w:rPr>)?<w:t[^>]*>` +
				regexp.QuoteMeta(marker) +
				`</w:t></w:r>`,
		)
		replacement := buildRichRunsXML(rt.Runs)
		docXML = re.ReplaceAllString(docXML, replacement)
	}

	files["word/document.xml"] = []byte(docXML)
	return repackDocxZip(zr, files)
}

// buildRichRunsXML 将 []RichRun 转为 Word XML 片段，\n 转为 <w:br/>
func buildRichRunsXML(runs []RichRun) string {
	var sb strings.Builder
	for _, run := range runs {
		parts := strings.Split(run.Text, "\n")
		for i, part := range parts {
			if i > 0 {
				sb.WriteString(`<w:r><w:br/></w:r>`)
			}
			if part == "" {
				continue
			}
			sb.WriteString(`<w:r>`)
			if run.Bold || run.Color != "" {
				sb.WriteString(`<w:rPr>`)
				if run.Bold {
					sb.WriteString(`<w:b/>`)
				}
				if run.Color != "" {
					sb.WriteString(`<w:color w:val="` + run.Color + `"/>`)
				}
				sb.WriteString(`</w:rPr>`)
			}
			sb.WriteString(`<w:t xml:space="preserve">`)
			sb.WriteString(xmlEscape(part))
			sb.WriteString(`</w:t></w:r>`)
		}
	}
	return sb.String()
}

// xmlEscape 转义 XML 特殊字符
func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	return s
}
