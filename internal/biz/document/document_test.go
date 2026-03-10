package document

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"os"
	"strings"
	"testing"
)

func TestGeneratePDF_WithImage(t *testing.T) {
	uc := NewUseCase()

	// 1x1 透明 PNG 图片
	pngBytes := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
		0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41,
		0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
		0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00,
		0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE,
		0x42, 0x60, 0x82,
	}
	signature := "data:image/png;base64," + base64.StdEncoding.EncodeToString(pngBytes)

	data := map[string]any{
		"DocumentTitle":   "合同文件",
		"Name":            "李四",
		"Date":            "2026-03-05",
		"Item1":           "咨询费",
		"Amount1":         "5000",
		"Note1":           "不含税",
		"Item2":           "差旅费",
		"Amount2":         "1000",
		"Note2":           "实报实销",
		"SignatureBase64": signature,
	}

	pdfBytes, err := uc.GeneratePDF("example_template", data)
	if err != nil {
		t.Fatalf("GeneratePDF() 错误 = %v", err)
	}
	if len(pdfBytes) == 0 {
		t.Fatal("GeneratePDF() 返回空字节流")
	}
	if string(pdfBytes[:4]) != "%PDF" {
		t.Fatal("GeneratePDF() 返回的不是有效的 PDF 文件")
	}
}

func TestGenerateWord_WithTextAndImage(t *testing.T) {
	// 1x1 透明 PNG 图片
	pngBytes := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
		0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41,
		0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
		0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00,
		0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE,
		0x42, 0x60, 0x82,
	}
	signature := "data:image/png;base64," + base64.StdEncoding.EncodeToString(pngBytes)

	// word_template_test.docx 已常驻于 templates/ 目录，直接使用
	uc := NewUseCase()
	wordBytes, err := uc.GenerateWord("word_template_test", map[string]any{
		"Name":      "张三",
		"Signature": signature,
	})
	if err != nil {
		t.Fatalf("GenerateWord() 错误 = %v", err)
	}
	if len(wordBytes) == 0 {
		t.Fatal("GenerateWord() 返回空字节流")
	}

	// 验证结果是合法的 ZIP（docx 文件头为 PK）
	if wordBytes[0] != 'P' || wordBytes[1] != 'K' {
		t.Fatal("GenerateWord() 返回的不是有效的 docx 文件")
	}

	// 验证文本占位符已替换：document.xml 中不应再含有 {Name}，应含有 张三
	zr, err := zip.NewReader(bytes.NewReader(wordBytes), int64(len(wordBytes)))
	if err != nil {
		t.Fatalf("解析输出 docx zip 失败: %v", err)
	}
	docXML := readZipFile(t, zr, "word/document.xml")
	if strings.Contains(docXML, "{Name}") {
		t.Error("document.xml 中文本占位符 {Name} 未被替换")
	}
	if !strings.Contains(docXML, "张三") {
		t.Error("document.xml 中未找到替换后的文本 张三")
	}

	// 验证图片已注入：media 目录下应有图片文件，document.xml 应含有 w:drawing
	hasMedia := false
	for _, f := range zr.File {
		if strings.HasPrefix(f.Name, "word/media/") {
			hasMedia = true
			break
		}
	}
	if !hasMedia {
		t.Error("输出 docx 中未找到注入的图片文件（word/media/）")
	}
	if !strings.Contains(docXML, "w:drawing") {
		t.Error("document.xml 中未找到图片节点 w:drawing")
	}

	// 所有验证通过，将结果保存到 templates/test.docx 供人工检查
	if err := os.WriteFile("../../../templates/test.docx", wordBytes, 0644); err != nil {
		t.Errorf("保存 test.docx 失败: %v", err)
	}
}

// readZipFile 从 zip.Reader 中读取指定文件内容为字符串
func readZipFile(t *testing.T, zr *zip.Reader, name string) string {
	t.Helper()
	for _, f := range zr.File {
		if f.Name == name {
			rc, err := f.Open()
			if err != nil {
				t.Fatalf("打开 zip 文件 %s 失败: %v", name, err)
			}
			defer rc.Close()
			var buf bytes.Buffer
			buf.ReadFrom(rc)
			return buf.String()
		}
	}
	t.Fatalf("zip 中未找到文件: %s", name)
	return ""
}
