package document

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apipb "auth_info/api/gen/api/proto"
	bizdoc "auth_info/internal/biz/document"
	"auth_info/internal/handler/httpx"
)

// Handler 处理文档生成相关接口
type Handler struct {
	uc *bizdoc.UseCase
}

// NewHandler Wire Provider
func NewHandler(uc *bizdoc.UseCase) *Handler {
	return &Handler{uc: uc}
}

// GeneratePDF 生成 PDF 文档
// @Summary  生成 PDF 文档
// @Tags     Document
// @Accept   json
// @Produce  application/pdf
// @Security BearerAuth
// @Router   /document/generate-pdf [post]
func (h *Handler) GeneratePDF(c *gin.Context) {
	var req apipb.GeneratePDFRequest
	if !httpx.BindAndValidateJSON(c, &req) {
		return
	}

	pdfBytes, err := h.uc.GeneratePDF(c.Request.Context(), req.GetTemplateName(), structToMap(req.GetData()))
	if err != nil {
		httpx.WriteError(c, err)
		return
	}

	c.Header("Content-Disposition", `attachment; filename="document.pdf"`)
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// GenerateWord 生成 Word 文档
// @Summary  生成 Word 文档
// @Tags     Document
// @Accept   json
// @Produce  application/vnd.openxmlformats-officedocument.wordprocessingml.document
// @Security BearerAuth
// @Router   /document/generate-word [post]
func (h *Handler) GenerateWord(c *gin.Context) {
	var req apipb.GenerateWordRequest
	if !httpx.BindAndValidateJSON(c, &req) {
		return
	}

	wordBytes, err := h.uc.GenerateWord(c.Request.Context(), req.GetTemplateName(), structToWordTemplateData(req.GetData()))
	if err != nil {
		httpx.WriteError(c, err)
		return
	}

	c.Header("Content-Disposition", `attachment; filename="document.docx"`)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", wordBytes)
}
