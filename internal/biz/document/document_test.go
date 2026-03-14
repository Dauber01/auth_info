package document

import (
	"archive/zip"
	"bytes"
	"context"
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

	pdfBytes, err := uc.GeneratePDF(context.Background(), "example_template", data)
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
	wordBytes, err := uc.GenerateWord(context.Background(), "word_template_test", WordTemplateData{
		Texts: map[string]RichText{
			"Name": {Runs: []RichRun{{Text: "张三"}}},
		},
		Images: map[string]ImageValue{
			"Signature": {ImageURL: signature}, // 兼容 base64 格式
		},
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

func TestGenerateWord_RichTextAndImageOptions(t *testing.T) {
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

	uc := NewUseCase()
	// 使用模板里实际存在的 key: {Name} 和 {Signature}
	data := WordTemplateData{
		Texts: map[string]RichText{
			"Name": {
				Runs: []RichRun{
					{Text: "重要通知", Bold: true, Color: "FF0000"},
					{Text: "\n"},
					{Text: "请注意以下事项", Bold: false, Color: "0000FF"},
				},
			},
		},
		Images: map[string]ImageValue{
			"Signature": {
				ImageURL:     signature, // 兼容 base64 格式
				OriginalSize: true,
				MaxWidthPx:   200,
				MaxHeightPx:  200,
			},
		},
	}

	wordBytes, err := uc.GenerateWord(context.Background(), "word_template_test", data)
	if err != nil {
		t.Fatalf("GenerateWord() 错误 = %v", err)
	}
	if len(wordBytes) == 0 {
		t.Fatal("GenerateWord() 返回空字节流")
	}

	// 验证 ZIP 格式
	if wordBytes[0] != 'P' || wordBytes[1] != 'K' {
		t.Fatal("GenerateWord() 返回的不是有效�� docx 文件")
	}

	zr, err := zip.NewReader(bytes.NewReader(wordBytes), int64(len(wordBytes)))
	if err != nil {
		t.Fatalf("解析输出 docx zip 失败: %v", err)
	}

	docXML := readZipFile(t, zr, "word/document.xml")

	// 验证富文本：应包含加粗标签 <w:b/> 和颜色标签 <w:color w:val="FF0000"/>
	if !strings.Contains(docXML, "<w:b/>") {
		t.Error("document.xml 中未找到加粗标签 <w:b/>")
	}
	if !strings.Contains(docXML, `<w:color w:val="FF0000"/>`) {
		t.Error("document.xml 中未找到红色标签")
	}
	if !strings.Contains(docXML, `<w:color w:val="0000FF"/>`) {
		t.Error("document.xml 中未找到蓝色标签")
	}
	// 验证换行：应包含 <w:br/>
	if !strings.Contains(docXML, "<w:br/>") {
		t.Error("document.xml 中未找到换行标签 <w:br/>")
	}
	// 验证占位符已被替换
	if strings.Contains(docXML, "{Name}") {
		t.Error("document.xml 中占位符 {Name} 未被替换")
	}

	// 验证图片注入
	mediaCount := 0
	for _, f := range zr.File {
		if strings.HasPrefix(f.Name, "word/media/") {
			mediaCount++
		}
	}
	if mediaCount < 1 {
		t.Errorf("期望至少 1 个图片文件，实际找到 %d 个", mediaCount)
	}
	if !strings.Contains(docXML, "w:drawing") {
		t.Error("document.xml 中未找到图片节点 w:drawing")
	}

	// 保存结果供人工检查
	if err := os.WriteFile("../../../templates/test_rich.docx", wordBytes, 0644); err != nil {
		t.Errorf("保存 test_rich.docx 失败: %v", err)
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
