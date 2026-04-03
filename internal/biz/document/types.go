package document

// RichRun 富文本中的一段文字，可独立设置样式
type RichRun struct {
	Text  string `json:"text"`
	Bold  bool   `json:"bold,omitempty"`
	Color string `json:"color,omitempty"`
}

// RichText 富文本值，支持多段样式和换行（\n 自动转为换行符）
type RichText struct {
	Runs []RichRun `json:"runs"`
}

// ImageValue 图片值，支持原始尺寸和最大尺寸限制
type ImageValue struct {
	ImageURL     string  `json:"image_url"`
	OriginalSize bool    `json:"original_size,omitempty"`
	MaxWidthPx   float64 `json:"max_width_px,omitempty"`
	MaxHeightPx  float64 `json:"max_height_px,omitempty"`
}

// WordTemplateData Word 文档生成的数据结构，明确区分文本和图片
type WordTemplateData struct {
	Texts  map[string]RichText   `json:"texts"`
	Images map[string]ImageValue `json:"images"`
}
