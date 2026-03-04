package handler

import (
	"net/http"
	"strings"

	"auth_info/internal/biz"

	"github.com/gin-gonic/gin"
)

// DocumentHandler 处理文档生成相关接口
type DocumentHandler struct {
	uc *biz.DocumentUseCase
}

func NewDocumentHandler(uc *biz.DocumentUseCase) *DocumentHandler {
	return &DocumentHandler{uc: uc}
}

type generatePDFRequest struct {
	TemplateName string         `json:"template_name" binding:"required"`
	Data         map[string]any `json:"data"           binding:"required"`
}

// GeneratePDF 生成 PDF 文档
// @Summary  生成 PDF 文档
// @Tags     Document
// @Accept   json
// @Produce  application/pdf
// @Param    body body generatePDFRequest true "模板名称和填充数据"
// @Success  200  {file}   binary "PDF 文件流"
// @Failure  400  {object} map[string]any "参数错误"
// @Failure  404  {object} map[string]any "模板不存在"
// @Failure  500  {object} map[string]any "生成失败"
// @Security BearerAuth
// @Router   /document/generate-pdf [post]
func (h *DocumentHandler) GeneratePDF(c *gin.Context) {
	var req generatePDFRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	pdfBytes, err := h.uc.GeneratePDF(req.TemplateName, req.Data)
	if err != nil {
		status := http.StatusInternalServerError
		if isNotFoundErr(err) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"code":    status,
			"message": err.Error(),
		})
		return
	}

	c.Header("Content-Disposition", `attachment; filename="document.pdf"`)
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// isNotFoundErr 判断是否为模板未找到错误
func isNotFoundErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "template not found")
}
