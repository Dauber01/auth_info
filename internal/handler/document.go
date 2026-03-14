package handler

import (
	"net/http"
	"strings"

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
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err)
		return
	}
	if err := validateDocumentTemplateRequest(req.GetTemplateName(), req.GetData()); err != nil {
		badRequest(c, err)
		return
	}

	pdfBytes, err := h.uc.GeneratePDF(req.GetTemplateName(), structToMap(req.GetData()))
	if err != nil {
		status := http.StatusInternalServerError
		if isNotFoundErr(err) {
			status = http.StatusNotFound
		}
		writeOperationReply(c, status, err.Error())
		return
	}

	c.Header("Content-Disposition", `attachment; filename="document.pdf"`)
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// isNotFoundErr 判断是否为模板未找到错误
func isNotFoundErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "template not found")
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
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err)
		return
	}
	if err := validateDocumentTemplateRequest(req.GetTemplateName(), req.GetData()); err != nil {
		badRequest(c, err)
		return
	}

	wordBytes, err := h.uc.GenerateWord(req.GetTemplateName(), structToWordTemplateData(req.GetData()))
	if err != nil {
		status := http.StatusInternalServerError
		if isNotFoundErr(err) {
			status = http.StatusNotFound
		}
		writeOperationReply(c, status, err.Error())
		return
	}

	c.Header("Content-Disposition", `attachment; filename="document.docx"`)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", wordBytes)
}
