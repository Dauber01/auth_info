package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apipb "auth_info/api/gen/api/proto"
	bizdoc "auth_info/internal/biz/document"
)

// DocumentHandler 处理文档生成相关接口
type DocumentHandler struct {
	uc *bizdoc.UseCase
}

func NewDocumentHandler(uc *bizdoc.UseCase) *DocumentHandler {
	return &DocumentHandler{uc: uc}
}

// GeneratePDF 生成 PDF 文档
// @Summary  生成 PDF 文档
// @Tags     Document
// @Accept   json
// @Produce  application/pdf
// @Security BearerAuth
// @Router   /document/generate-pdf [post]
func (h *DocumentHandler) GeneratePDF(c *gin.Context) {
	var req apipb.GeneratePDFRequest
	if !bindAndValidateJSON(c, &req) {
		return
	}

	pdfBytes, err := h.uc.GeneratePDF(c.Request.Context(), req.GetTemplateName(), structToMap(req.GetData()))
	if err != nil {
		writeError(c, err)
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
func (h *DocumentHandler) GenerateWord(c *gin.Context) {
	var req apipb.GenerateWordRequest
	if !bindAndValidateJSON(c, &req) {
		return
	}

	wordBytes, err := h.uc.GenerateWord(c.Request.Context(), req.GetTemplateName(), structToWordTemplateData(req.GetData()))
	if err != nil {
		writeError(c, err)
		return
	}

	c.Header("Content-Disposition", `attachment; filename="document.docx"`)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", wordBytes)
}
