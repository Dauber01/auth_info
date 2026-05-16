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
// @Description 根据模板名称和数据生成 PDF 文件。
// @Tags     Document
// @Accept   json
// @Produce  application/pdf
// @Param    request body apipb.GeneratePDFRequest true "生成 PDF 参数"
// @Success  200 {file} file "PDF 文件"
// @Failure  400 {object} apipb.OperationReply "请求参数错误"
// @Failure  401 {object} apipb.OperationReply "未认证"
// @Failure  403 {object} apipb.OperationReply "无访问权限"
// @Failure  404 {object} apipb.OperationReply "模板不存在"
// @Failure  500 {object} apipb.OperationReply "服务器内部错误"
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
// @Description 根据模板名称和数据生成 Word 文档文件。
// @Tags     Document
// @Accept   json
// @Produce  application/vnd.openxmlformats-officedocument.wordprocessingml.document
// @Param    request body apipb.GenerateWordRequest true "生成 Word 参数"
// @Success  200 {file} file "Word 文档文件"
// @Failure  400 {object} apipb.OperationReply "请求参数错误"
// @Failure  401 {object} apipb.OperationReply "未认证"
// @Failure  403 {object} apipb.OperationReply "无访问权限"
// @Failure  404 {object} apipb.OperationReply "模板不存在"
// @Failure  500 {object} apipb.OperationReply "服务器内部错误"
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
