package biz

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

func TestDocumentUseCase_GeneratePDF(t *testing.T) {
	uc := NewDocumentUseCase()

	tests := []struct {
		name         string
		templateName string
		data         map[string]any
		wantErr      bool
		errContains  string
	}{
		{
			name:         "成功生成PDF-基础文本",
			templateName: "example_template",
			data: map[string]any{
				"DocumentTitle": "测试文档",
				"Name":          "张三",
				"Date":          "2026-03-04",
				"Item1":         "服务费",
				"Amount1":       "10000",
				"Note1":         "含税",
				"Item2":         "管理费",
				"Amount2":       "2000",
				"Note2":         "季度",
				"SignatureBase64": "",
			},
			wantErr: false,
		},
		{
			name:         "成功生成PDF-包含图片",
			templateName: "example_template",
			data: map[string]any{
				"DocumentTitle": "合同文件",
				"Name":          "李四",
				"Date":          "2026-03-05",
				"Item1":         "咨询费",
				"Amount1":       "5000",
				"Note1":         "不含税",
				"Item2":         "差旅费",
				"Amount2":       "1000",
				"Note2":         "实报实销",
				"SignatureBase64": generateTestImageBase64(),
			},
			wantErr: false,
		},
		{
			name:         "模板不存在",
			templateName: "non_existent_template",
			data: map[string]any{
				"DocumentTitle": "测试",
			},
			wantErr:     true,
			errContains: "template not found",
		},
		{
			name:         "空数据",
			templateName: "example_template",
			data:         map[string]any{},
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pdfBytes, err := uc.GeneratePDF(tt.templateName, tt.data)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GeneratePDF() 期望错误但没有返回错误")
					return
				}
				if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("GeneratePDF() 错误信息 = %v, 期望包含 %v", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("GeneratePDF() 错误 = %v, 期望成功", err)
				return
			}

			if len(pdfBytes) == 0 {
				t.Errorf("GeneratePDF() 返回空字节流")
				return
			}

			// 验证 PDF 文件头
			if len(pdfBytes) < 4 || string(pdfBytes[:4]) != "%PDF" {
				t.Errorf("GeneratePDF() 返回的不是有效的 PDF 文件")
			}

			// 可选：保存到临时文件用于手动验证
			if os.Getenv("SAVE_TEST_PDF") == "1" {
				tmpDir := filepath.Join("testdata", "output")
				os.MkdirAll(tmpDir, 0755)
				tmpFile := filepath.Join(tmpDir, tt.name+".pdf")
				os.WriteFile(tmpFile, pdfBytes, 0644)
				t.Logf("测试 PDF 已保存到: %s", tmpFile)
			}
		})
	}
}

func TestDocumentUseCase_GeneratePDF_InvalidTemplate(t *testing.T) {
	uc := NewDocumentUseCase()

	// 创建临时目录和无效模板
	tmpDir := t.TempDir()
	uc.templateDir = tmpDir

	invalidJSON := `{"title": "{{.Title}", "sections": [invalid json]}`
	os.WriteFile(filepath.Join(tmpDir, "invalid.json"), []byte(invalidJSON), 0644)

	_, err := uc.GeneratePDF("invalid", map[string]any{"Title": "测试"})
	if err == nil {
		t.Errorf("GeneratePDF() 期望解析错误但成功返回")
	}
}

func TestDocumentUseCase_GeneratePDF_ImageDecodeError(t *testing.T) {
	uc := NewDocumentUseCase()

	// 创建包含无效 base64 图片的模板
	tmpDir := t.TempDir()
	uc.templateDir = tmpDir

	tmpl := `{
		"title": "测试",
		"sections": [
			{
				"type": "image",
				"data": "invalid-base64!!!",
				"width": 60,
				"height": 30
			}
		]
	}`
	os.WriteFile(filepath.Join(tmpDir, "bad_image.json"), []byte(tmpl), 0644)

	_, err := uc.GeneratePDF("bad_image", map[string]any{})
	if err == nil {
		t.Errorf("GeneratePDF() 期望 base64 解码错误但成功返回")
	}
}

// 辅助函数：生成测试���的 1x1 PNG 图片 base64
func generateTestImageBase64() string {
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
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(pngBytes)
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
