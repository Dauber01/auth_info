package document

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/go-pdf/fpdf"
	godocx "github.com/lukasjarosch/go-docx"
)

// isImageValue 判断 data 值是否为图片（base64 格式）
func isImageValue(v any) bool {
	s, ok := v.(string)
	if !ok {
		return false
	}
	return strings.HasPrefix(s, "data:image/")
}

// UseCase 处理 PDF 文档生成，无需数据库依赖
type UseCase struct {
	templateDir string
	fontPath    string
}

func NewUseCase() *UseCase {
	return &UseCase{
		templateDir: "D:\\GoProject\\auth_info\\templates",
		fontPath:    "D:\\GoProject\\auth_info\\assets\\fonts\\NotoSansSC-Regular.ttf",
	}
}

// --- 模板结构体 ---

type documentTemplate struct {
	Title    string            `json:"title"`
	Sections []templateSection `json:"sections"`
}

type templateSection struct {
	Type     string     `json:"type"`                // paragraph | table | image
	Content  string     `json:"content,omitempty"`   // paragraph 文本
	FontSize float64    `json:"font_size,omitempty"` // 字号，默认 12
	Bold     bool       `json:"bold,omitempty"`      // 是否粗体
	Headers  []string   `json:"headers,omitempty"`   // table 表头
	Rows     [][]string `json:"rows,omitempty"`      // table 行数据
	Data     string     `json:"data,omitempty"`      // image base64 数据
	Width    float64    `json:"width,omitempty"`     // image 宽度 (mm)
	Height   float64    `json:"height,omitempty"`    // image 高度 (mm)
}

// GeneratePDF 根据模板名称和数据生成 PDF，返回字节流
func (uc *UseCase) GeneratePDF(templateName string, data map[string]any) ([]byte, error) {
	// 1. 读取模板文件
	tmplPath := fmt.Sprintf("%s/%s.json", uc.templateDir, templateName)
	tmplRaw, err := os.ReadFile(tmplPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("template not found: %s", templateName)
		}
		return nil, fmt.Errorf("failed to read template: %w", err)
	}

	// 2. 用 text/template 渲染 JSON（填充占位符）
	rendered, err := renderTemplate(string(tmplRaw), data)
	if err != nil {
		return nil, fmt.Errorf("failed to render template: %w", err)
	}

	// 3. 解析渲染后的 JSON
	var doc documentTemplate
	if err := json.Unmarshal([]byte(rendered), &doc); err != nil {
		return nil, fmt.Errorf("failed to parse template JSON after rendering: %w", err)
	}

	// 4. 生成 PDF
	return uc.buildPDF(&doc)
}

// renderTemplate 使用 text/template 将数据填充到模板字符串中
func renderTemplate(tmplStr string, data map[string]any) (string, error) {
	tmpl, err := template.New("doc").Option("missingkey=zero").Parse(tmplStr)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// buildPDF 使用 fpdf 将解析后的模板绘制成 PDF
func (uc *UseCase) buildPDF(doc *documentTemplate) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")

	// 加载中文字体（如果存在）
	fontFamily := "Helvetica"
	if fontBytes, err := os.ReadFile(uc.fontPath); err == nil {
		pdf.AddUTF8FontFromBytes("NotoSansSC", "", fontBytes)
		pdf.AddUTF8FontFromBytes("NotoSansSC", "B", fontBytes)
		fontFamily = "NotoSansSC"
	}

	pdf.AddPage()

	// 页面边距
	leftMargin, _, rightMargin, _ := pdf.GetMargins()
	pageWidth, _ := pdf.GetPageSize()
	contentWidth := pageWidth - leftMargin - rightMargin

	// 绘制标题
	if doc.Title != "" {
		pdf.SetFont(fontFamily, "B", 18)
		pdf.CellFormat(contentWidth, 12, doc.Title, "", 1, "C", false, 0, "")
		pdf.Ln(4)
	}

	// 逐 section 绘制
	for _, sec := range doc.Sections {
		switch sec.Type {
		case "paragraph":
			if err := uc.drawParagraph(pdf, &sec, fontFamily, contentWidth); err != nil {
				return nil, err
			}
		case "table":
			if err := uc.drawTable(pdf, &sec, fontFamily, contentWidth); err != nil {
				return nil, err
			}
		case "image":
			if err := uc.drawImage(pdf, &sec); err != nil {
				return nil, err
			}
		}
	}

	// 输出到 bytes.Buffer
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to output PDF: %w", err)
	}
	return buf.Bytes(), nil
}

func (uc *UseCase) drawParagraph(pdf *fpdf.Fpdf, sec *templateSection, fontFamily string, contentWidth float64) error {
	fontSize := sec.FontSize
	if fontSize <= 0 {
		fontSize = 12
	}
	style := ""
	if sec.Bold {
		style = "B"
	}
	pdf.SetFont(fontFamily, style, fontSize)
	lineHeight := fontSize * 0.45
	pdf.MultiCell(contentWidth, lineHeight, sec.Content, "", "L", false)
	pdf.Ln(2)
	return nil
}

func (uc *UseCase) drawTable(pdf *fpdf.Fpdf, sec *templateSection, fontFamily string, contentWidth float64) error {
	if len(sec.Headers) == 0 {
		return nil
	}

	colCount := len(sec.Headers)
	colWidth := contentWidth / float64(colCount)
	rowHeight := 7.0

	// 绘制表头
	pdf.SetFont(fontFamily, "B", 11)
	pdf.SetFillColor(220, 220, 220)
	for _, h := range sec.Headers {
		pdf.CellFormat(colWidth, rowHeight, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// 绘制数据行
	pdf.SetFont(fontFamily, "", 11)
	pdf.SetFillColor(255, 255, 255)
	for i, row := range sec.Rows {
		if i%2 == 0 {
			pdf.SetFillColor(248, 248, 248)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}
		for j, cell := range row {
			if j >= colCount {
				break
			}
			pdf.CellFormat(colWidth, rowHeight, cell, "1", 0, "L", true, 0, "")
		}
		pdf.Ln(-1)
	}
	pdf.Ln(4)
	return nil
}

func (uc *UseCase) drawImage(pdf *fpdf.Fpdf, sec *templateSection) error {
	if sec.Data == "" {
		return nil
	}

	raw := sec.Data
	imgType := "PNG"
	if idx := strings.Index(raw, ";base64,"); idx != -1 {
		mimeStr := raw[5:idx]
		switch {
		case strings.Contains(mimeStr, "jpeg") || strings.Contains(mimeStr, "jpg"):
			imgType = "JPG"
		case strings.Contains(mimeStr, "png"):
			imgType = "PNG"
		}
		raw = raw[idx+8:]
	}

	imgBytes, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return fmt.Errorf("failed to decode image base64: %w", err)
	}

	w := sec.Width
	h := sec.Height
	if w <= 0 {
		w = 60
	}
	if h <= 0 {
		h = 0
	}

	imgName := fmt.Sprintf("img_%p", sec)
	pdf.RegisterImageOptionsReader(imgName, fpdf.ImageOptions{ImageType: imgType}, bytes.NewReader(imgBytes))
	pdf.Image(imgName, pdf.GetX(), pdf.GetY(), w, h, true, imgType, 0, "")
	pdf.Ln(4)
	return nil
}

// GenerateWord 根据 .docx 模板和数据生成 Word 文档，返回字节流。
func (uc *UseCase) GenerateWord(templateName string, data map[string]any) ([]byte, error) {
	tmplPath := fmt.Sprintf("%s/%s.docx", uc.templateDir, templateName)

	docBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("template not found: %s", templateName)
		}
		return nil, fmt.Errorf("failed to read template: %w", err)
	}

	// 分离文本数据和图片数据
	textData := make(godocx.PlaceholderMap)
	imageData := make(map[string]string)
	for k, v := range data {
		if isImageValue(v) {
			imageData[k] = v.(string)
		} else {
			textData[k] = v
		}
	}

	const imgMarkerPrefix = "__IMGPLACEHOLDER_"
	for k := range imageData {
		textData[k] = imgMarkerPrefix + k + "__"
	}

	doc, err := godocx.OpenBytes(docBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to open docx template: %w", err)
	}
	defer doc.Close()

	if err := doc.ReplaceAll(textData); err != nil {
		return nil, fmt.Errorf("failed to replace text placeholders: %w", err)
	}

	var buf bytes.Buffer
	if err := doc.Write(&buf); err != nil {
		return nil, fmt.Errorf("failed to write docx after text replace: %w", err)
	}

	if len(imageData) > 0 {
		result, err := injectImagesToDocx(buf.Bytes(), imageData, imgMarkerPrefix)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	return buf.Bytes(), nil
}
