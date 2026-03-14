package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	apipb "auth_info/api/gen/api/proto"
	bizdict "auth_info/internal/biz/dict"
)

// DictHandler 字典配置 HTTP 处理器
type DictHandler struct {
	uc *bizdict.UseCase
}

// NewDictHandler Wire Provider
func NewDictHandler(uc *bizdict.UseCase) *DictHandler {
	return &DictHandler{uc: uc}
}

// ListDictTypes
// @Summary  获取字典类型列表
// @Tags     Dict
// @Produce  json
// @Security BearerAuth
// @Router   /dict/types [get]
func (h *DictHandler) ListDictTypes(c *gin.Context) {
	types, err := h.uc.ListDictTypes(c.Request.Context())
	if err != nil {
		writeError(c, err)
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
// @Tags     Dict
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Router   /dict/types [post]
func (h *DictHandler) CreateDictType(c *gin.Context) {
	var req apipb.CreateDictTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err)
		return
	}
	if err := validateProtoRules(&req); err != nil {
		writeError(c, err)
		return
	}

	if err := h.uc.CreateDictType(c.Request.Context(), req.GetCode(), req.GetName(), req.GetDescription(), int(req.GetSort())); err != nil {
		writeError(c, err)
		return
	}

	writeOperationReply(c, http.StatusOK, "created successfully")
}

// UpdateDictType
// @Summary  更新字典类型
// @Tags     Dict
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Router   /dict/types/{id} [put]
func (h *DictHandler) UpdateDictType(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		badRequest(c, err)
		return
	}

	var req apipb.UpdateDictTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err)
		return
	}
	req.Id = id
	if err := validateProtoRules(&req); err != nil {
		writeError(c, err)
		return
	}

	if err := h.uc.UpdateDictType(c.Request.Context(), uint(req.GetId()), req.GetName(), req.GetDescription(), int(req.GetSort())); err != nil {
		writeError(c, err)
		return
	}

	writeOperationReply(c, http.StatusOK, "updated successfully")
}

// DeleteDictType
// @Summary  删除字典类型
// @Tags     Dict
// @Produce  json
// @Security BearerAuth
// @Router   /dict/types/{id} [delete]
func (h *DictHandler) DeleteDictType(c *gin.Context) {
	req := apipb.DeleteDictTypeRequest{}
	id, err := parseUintParam(c, "id")
	if err != nil {
		badRequest(c, err)
		return
	}
	req.Id = id
	if err := validateProtoRules(&req); err != nil {
		writeError(c, err)
		return
	}

	if err := h.uc.DeleteDictType(c.Request.Context(), uint(req.GetId())); err != nil {
		writeError(c, err)
		return
	}

	writeOperationReply(c, http.StatusOK, "deleted successfully")
}

// ListDictItems
// @Summary  根据类型编码获取字典数据列表
// @Tags     Dict
// @Produce  json
// @Security BearerAuth
// @Router   /dict/items [get]
func (h *DictHandler) ListDictItems(c *gin.Context) {
	req := apipb.ListDictItemsRequest{TypeCode: strings.TrimSpace(c.Query("type_code"))}
	if err := validateProtoRules(&req); err != nil {
		writeError(c, err)
		return
	}

	items, err := h.uc.ListDictItems(c.Request.Context(), req.GetTypeCode())
	if err != nil {
		writeError(c, err)
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
// @Tags     Dict
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Router   /dict/items [post]
func (h *DictHandler) CreateDictItem(c *gin.Context) {
	var req apipb.CreateDictItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err)
		return
	}
	if err := validateProtoRules(&req); err != nil {
		writeError(c, err)
		return
	}

	if err := h.uc.CreateDictItem(c.Request.Context(), req.GetTypeCode(), req.GetItemKey(), req.GetItemValue(), req.GetDescription(), int(req.GetSort())); err != nil {
		writeError(c, err)
		return
	}

	writeOperationReply(c, http.StatusOK, "created successfully")
}

// UpdateDictItem
// @Summary  更新字典数据
// @Tags     Dict
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Router   /dict/items/{id} [put]
func (h *DictHandler) UpdateDictItem(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		badRequest(c, err)
		return
	}

	var req apipb.UpdateDictItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err)
		return
	}
	req.Id = id
	if err := validateProtoRules(&req); err != nil {
		writeError(c, err)
		return
	}

	if err := h.uc.UpdateDictItem(c.Request.Context(), uint(req.GetId()), req.GetItemKey(), req.GetItemValue(), req.GetDescription(), int(req.GetSort()), int(req.GetStatus())); err != nil {
		writeError(c, err)
		return
	}

	writeOperationReply(c, http.StatusOK, "updated successfully")
}

// DeleteDictItem
// @Summary  删除字典数据
// @Tags     Dict
// @Produce  json
// @Security BearerAuth
// @Router   /dict/items/{id} [delete]
func (h *DictHandler) DeleteDictItem(c *gin.Context) {
	req := apipb.DeleteDictItemRequest{}
	id, err := parseUintParam(c, "id")
	if err != nil {
		badRequest(c, err)
		return
	}
	req.Id = id
	if err := validateProtoRules(&req); err != nil {
		writeError(c, err)
		return
	}

	if err := h.uc.DeleteDictItem(c.Request.Context(), uint(req.GetId())); err != nil {
		writeError(c, err)
		return
	}

	writeOperationReply(c, http.StatusOK, "deleted successfully")
}
