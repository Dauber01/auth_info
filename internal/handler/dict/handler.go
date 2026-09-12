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

// ListDictTypes 获取字典类型列表
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

// CreateDictType 创建字典类型
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

// UpdateDictType 更新字典类型
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

// DeleteDictType 删除字典类型
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

// ListDictItems 获取字典数据列表
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

// CreateDictItem 创建字典数据
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

// UpdateDictItem 更新字典数据
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

// DeleteDictItem 删除字典数据
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
