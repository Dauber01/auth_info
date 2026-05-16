package dict

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	apipb "auth_info/api/gen/api/proto"
	bizdict "auth_info/internal/biz/dict"
	"auth_info/internal/handler/httpx"
)

// Handler 字典模块 HTTP 处理器。
type Handler struct {
	uc *bizdict.UseCase
}

// NewHandler Wire Provider
func NewHandler(uc *bizdict.UseCase) *Handler {
	return &Handler{uc: uc}
}

// ListDictTypes
// @Summary  获取字典类型列表
// @Description 获取全部字典类型，按 sort 正序返回。
// @Tags     Dict
// @Produce  json
// @Success  200 {object} apipb.ListDictTypesReply "请求成功"
// @Failure  401 {object} apipb.OperationReply "未认证"
// @Failure  403 {object} apipb.OperationReply "无访问权限"
// @Failure  500 {object} apipb.OperationReply "服务器内部错误"
// @Security BearerAuth
// @Router   /dict/types [get]
func (h *Handler) ListDictTypes(c *gin.Context) {
	types, err := h.uc.ListDictTypes(c.Request.Context())
	if err != nil {
		httpx.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, &apipb.ListDictTypesReply{
		Code:    http.StatusOK,
		Message: "success",
		Data:    dictTypesToProto(types),
	})
}

// CreateDictType
// @Summary  创建字典类型
// @Description 创建一个新的字典类型，code 必须唯一。
// @Tags     Dict
// @Accept   json
// @Produce  json
// @Param    request body apipb.CreateDictTypeRequest true "创建字典类型参数"
// @Success  200 {object} apipb.OperationReply "创建成功"
// @Failure  400 {object} apipb.OperationReply "请求参数错误"
// @Failure  401 {object} apipb.OperationReply "未认证"
// @Failure  403 {object} apipb.OperationReply "无访问权限"
// @Failure  409 {object} apipb.OperationReply "字典类型编码已存在"
// @Failure  500 {object} apipb.OperationReply "服务器内部错误"
// @Security BearerAuth
// @Router   /dict/types [post]
func (h *Handler) CreateDictType(c *gin.Context) {
	var req apipb.CreateDictTypeRequest
	if !httpx.BindAndValidateJSON(c, &req) {
		return
	}

	if err := h.uc.CreateDictType(c.Request.Context(), req.GetCode(), req.GetName(), req.GetDescription(), int(req.GetSort())); err != nil {
		httpx.WriteError(c, err)
		return
	}

	httpx.WriteOperationReply(c, http.StatusOK, "created successfully")
}

// UpdateDictType
// @Summary  更新字典类型
// @Description 根据路径 ID 更新字典类型名称、描述和排序，code 不可修改。
// @Tags     Dict
// @Accept   json
// @Produce  json
// @Param    id path int true "字典类型 ID"
// @Param    request body apipb.UpdateDictTypeRequest true "更新字典类型参数"
// @Success  200 {object} apipb.OperationReply "更新成功"
// @Failure  400 {object} apipb.OperationReply "请求参数错误"
// @Failure  401 {object} apipb.OperationReply "未认证"
// @Failure  403 {object} apipb.OperationReply "无访问权限"
// @Failure  404 {object} apipb.OperationReply "字典类型不存在"
// @Failure  500 {object} apipb.OperationReply "服务器内部错误"
// @Security BearerAuth
// @Router   /dict/types/{id} [put]
func (h *Handler) UpdateDictType(c *gin.Context) {
	var req apipb.UpdateDictTypeRequest
	if !httpx.BindPathIDAndValidateJSON(c, "id", &req, func(id uint64) {
		req.Id = id
	}) {
		return
	}

	if err := h.uc.UpdateDictType(c.Request.Context(), uint(req.GetId()), req.GetName(), req.GetDescription(), int(req.GetSort())); err != nil {
		httpx.WriteError(c, err)
		return
	}

	httpx.WriteOperationReply(c, http.StatusOK, "updated successfully")
}

// DeleteDictType
// @Summary  删除字典类型
// @Description 根据路径 ID 软删除字典类型。
// @Tags     Dict
// @Produce  json
// @Param    id path int true "字典类型 ID"
// @Success  200 {object} apipb.OperationReply "删除成功"
// @Failure  400 {object} apipb.OperationReply "请求参数错误"
// @Failure  401 {object} apipb.OperationReply "未认证"
// @Failure  403 {object} apipb.OperationReply "无访问权限"
// @Failure  404 {object} apipb.OperationReply "字典类型不存在"
// @Failure  500 {object} apipb.OperationReply "服务器内部错误"
// @Security BearerAuth
// @Router   /dict/types/{id} [delete]
func (h *Handler) DeleteDictType(c *gin.Context) {
	req := apipb.DeleteDictTypeRequest{}
	if !httpx.ValidatePathIDRequest(c, "id", &req, func(id uint64) {
		req.Id = id
	}) {
		return
	}

	if err := h.uc.DeleteDictType(c.Request.Context(), uint(req.GetId())); err != nil {
		httpx.WriteError(c, err)
		return
	}

	httpx.WriteOperationReply(c, http.StatusOK, "deleted successfully")
}

// ListDictItems
// @Summary  获取字典数据列表
// @Description 根据字典类型编码获取字典数据，按 sort 正序返回。
// @Tags     Dict
// @Produce  json
// @Param    type_code query string true "字典类型编码"
// @Success  200 {object} apipb.ListDictItemsReply "请求成功"
// @Failure  400 {object} apipb.OperationReply "请求参数错误"
// @Failure  401 {object} apipb.OperationReply "未认证"
// @Failure  403 {object} apipb.OperationReply "无访问权限"
// @Failure  500 {object} apipb.OperationReply "服务器内部错误"
// @Security BearerAuth
// @Router   /dict/items [get]
func (h *Handler) ListDictItems(c *gin.Context) {
	req := apipb.ListDictItemsRequest{TypeCode: strings.TrimSpace(c.Query("type_code"))}
	if !httpx.ValidateProto(c, &req) {
		return
	}

	items, err := h.uc.ListDictItems(c.Request.Context(), req.GetTypeCode())
	if err != nil {
		httpx.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, &apipb.ListDictItemsReply{
		Code:    http.StatusOK,
		Message: "success",
		Data:    dictItemsToProto(items),
	})
}

// CreateDictItem
// @Summary  创建字典数据
// @Description 创建一个新的字典数据项，默认启用。
// @Tags     Dict
// @Accept   json
// @Produce  json
// @Param    request body apipb.CreateDictItemRequest true "创建字典数据参数"
// @Success  200 {object} apipb.OperationReply "创建成功"
// @Failure  400 {object} apipb.OperationReply "请求参数错误"
// @Failure  401 {object} apipb.OperationReply "未认证"
// @Failure  403 {object} apipb.OperationReply "无访问权限"
// @Failure  500 {object} apipb.OperationReply "服务器内部错误"
// @Security BearerAuth
// @Router   /dict/items [post]
func (h *Handler) CreateDictItem(c *gin.Context) {
	var req apipb.CreateDictItemRequest
	if !httpx.BindAndValidateJSON(c, &req) {
		return
	}

	if err := h.uc.CreateDictItem(c.Request.Context(), req.GetTypeCode(), req.GetItemKey(), req.GetItemValue(), req.GetDescription(), int(req.GetSort())); err != nil {
		httpx.WriteError(c, err)
		return
	}

	httpx.WriteOperationReply(c, http.StatusOK, "created successfully")
}

// UpdateDictItem
// @Summary  更新字典数据
// @Description 根据路径 ID 更新字典数据的键、值、描述、排序和状态。
// @Tags     Dict
// @Accept   json
// @Produce  json
// @Param    id path int true "字典数据 ID"
// @Param    request body apipb.UpdateDictItemRequest true "更新字典数据参数"
// @Success  200 {object} apipb.OperationReply "更新成功"
// @Failure  400 {object} apipb.OperationReply "请求参数错误"
// @Failure  401 {object} apipb.OperationReply "未认证"
// @Failure  403 {object} apipb.OperationReply "无访问权限"
// @Failure  404 {object} apipb.OperationReply "字典数据不存在"
// @Failure  500 {object} apipb.OperationReply "服务器内部错误"
// @Security BearerAuth
// @Router   /dict/items/{id} [put]
func (h *Handler) UpdateDictItem(c *gin.Context) {
	var req apipb.UpdateDictItemRequest
	if !httpx.BindPathIDAndValidateJSON(c, "id", &req, func(id uint64) {
		req.Id = id
	}) {
		return
	}

	if err := h.uc.UpdateDictItem(c.Request.Context(), uint(req.GetId()), req.GetItemKey(), req.GetItemValue(), req.GetDescription(), int(req.GetSort()), int(req.GetStatus())); err != nil {
		httpx.WriteError(c, err)
		return
	}

	httpx.WriteOperationReply(c, http.StatusOK, "updated successfully")
}

// DeleteDictItem
// @Summary  删除字典数据
// @Description 根据路径 ID 软删除字典数据。
// @Tags     Dict
// @Produce  json
// @Param    id path int true "字典数据 ID"
// @Success  200 {object} apipb.OperationReply "删除成功"
// @Failure  400 {object} apipb.OperationReply "请求参数错误"
// @Failure  401 {object} apipb.OperationReply "未认证"
// @Failure  403 {object} apipb.OperationReply "无访问权限"
// @Failure  404 {object} apipb.OperationReply "字典数据不存在"
// @Failure  500 {object} apipb.OperationReply "服务器内部错误"
// @Security BearerAuth
// @Router   /dict/items/{id} [delete]
func (h *Handler) DeleteDictItem(c *gin.Context) {
	req := apipb.DeleteDictItemRequest{}
	if !httpx.ValidatePathIDRequest(c, "id", &req, func(id uint64) {
		req.Id = id
	}) {
		return
	}

	if err := h.uc.DeleteDictItem(c.Request.Context(), uint(req.GetId())); err != nil {
		httpx.WriteError(c, err)
		return
	}

	httpx.WriteOperationReply(c, http.StatusOK, "deleted successfully")
}
