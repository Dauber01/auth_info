package biz

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"regexp"
	"strings"
)

// injectImagesToDocx 向已完成文本替换的 docx 字节流中注入图片。
// 找到 document.xml 里的临时标记，将其所在 <w:r> 节点替换为 <w:drawing> 节点，
// 同时向 ZIP 写入图片文件并追加 _rels 关系引用。
func injectImagesToDocx(docxBytes []byte, imageData map[string]string, markerPrefix string) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(docxBytes), int64(len(docxBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to read docx zip: %w", err)
	}

	// 把所有 ZIP 条目读入内存，方便修改后整体重打包
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

	for _, p := range indexedPairs(imageData) {
		imgBytes, imgExt, err := parseDataURI(p.val)
		if err != nil {
			return nil, fmt.Errorf("image key %q: %w", p.key, err)
		}

		mediaName := fmt.Sprintf("image%d.%s", p.idx+1, imgExt)
		files["word/media/"+mediaName] = imgBytes

		rId := fmt.Sprintf("rIdImg%d", p.idx+1)
		relEntry := fmt.Sprintf(
			`<Relationship Id="%s" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="media/%s"/>`,
			rId, mediaName,
		)
		relsXML = strings.Replace(relsXML, "</Relationships>", relEntry+"</Relationships>", 1)

		cx, cy := imageEMU(imgBytes, 1800000) // 默认宽 5cm
		drawingXML := buildDrawingXML(rId, p.idx+1, cx, cy)

		marker := markerPrefix + p.key + "__"
		// 匹配包含 marker 的整个 <w:r> 节点（含可选的 <w:rPr>）
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

// parseDataURI 解析 "data:image/png;base64,..." 格式，返回图片字节和扩展名
func parseDataURI(dataURI string) ([]byte, string, error) {
	semicolon := strings.Index(dataURI, ";")
	comma := strings.Index(dataURI, ",")
	if semicolon == -1 || comma == -1 || comma < semicolon {
		return nil, "", fmt.Errorf("invalid data URI format")
	}
	mimeType := dataURI[5:semicolon] // 截掉 "data:"
	ext := "png"
	if strings.Contains(mimeType, "jpeg") || strings.Contains(mimeType, "jpg") {
		ext = "jpg"
	}
	imgBytes, err := base64.StdEncoding.DecodeString(dataURI[comma+1:])
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode base64: %w", err)
	}
	return imgBytes, ext, nil
}

// imageEMU 解码图片获取实际宽高，按给定目标宽度（EMU）等比缩放返回 cx/cy。
// 1 cm = 360000 EMU；解码失败时返回 defaultCx × defaultCx/2 作为兜底。
func imageEMU(imgBytes []byte, defaultCx int64) (int64, int64) {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(imgBytes))
	if err != nil || cfg.Width == 0 {
		return defaultCx, defaultCx / 2
	}
	cx := defaultCx
	cy := cx * int64(cfg.Height) / int64(cfg.Width)
	return cx, cy
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

// repackDocxZip 按原始 ZIP 条目顺序重新打包，新增文件（如 media 图片）追加在末尾
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

	// 写入新增文件（原 ZIP 中不存在的，如新图片）
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

// indexedPairs 将 map 转为有稳定顺序的 (index, key, value) 序列，
// 保证每次运行图片编号一致（map 迭代顺序不确定）
func indexedPairs(m map[string]string) []struct {
	idx int
	key string
	val string
} {
	pairs := make([]struct {
		idx int
		key string
		val string
	}, 0, len(m))
	i := 0
	for k, v := range m {
		pairs = append(pairs, struct {
			idx int
			key string
			val string
		}{i, k, v})
		i++
	}
	return pairs
}
