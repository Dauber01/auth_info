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
	types, err := h.uc.ListDictTypes()
	if err != nil {
		writeOperationReply(c, http.StatusInternalServerError, err.Error())
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
		badRequest(c, err)
		return
	}

	if err := h.uc.CreateDictType(req.GetCode(), req.GetName(), req.GetDescription(), int(req.GetSort())); err != nil {
		writeOperationReply(c, http.StatusConflict, err.Error())
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
		badRequest(c, err)
		return
	}

	if err := h.uc.UpdateDictType(uint(req.GetId()), req.GetName(), req.GetDescription(), int(req.GetSort())); err != nil {
		writeOperationReply(c, http.StatusNotFound, err.Error())
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
		badRequest(c, err)
		return
	}

	if err := h.uc.DeleteDictType(uint(req.GetId())); err != nil {
		writeOperationReply(c, http.StatusNotFound, err.Error())
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
		badRequest(c, err)
		return
	}

	items, err := h.uc.ListDictItems(req.GetTypeCode())
	if err != nil {
		writeOperationReply(c, http.StatusInternalServerError, err.Error())
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
		badRequest(c, err)
		return
	}

	if err := h.uc.CreateDictItem(req.GetTypeCode(), req.GetItemKey(), req.GetItemValue(), req.GetDescription(), int(req.GetSort())); err != nil {
		writeOperationReply(c, http.StatusInternalServerError, err.Error())
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
		badRequest(c, err)
		return
	}

	if err := h.uc.UpdateDictItem(uint(req.GetId()), req.GetItemKey(), req.GetItemValue(), req.GetDescription(), int(req.GetSort()), int(req.GetStatus())); err != nil {
		writeOperationReply(c, http.StatusNotFound, err.Error())
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
		badRequest(c, err)
		return
	}

	if err := h.uc.DeleteDictItem(uint(req.GetId())); err != nil {
		writeOperationReply(c, http.StatusNotFound, err.Error())
		return
	}

	writeOperationReply(c, http.StatusOK, "deleted successfully")
}
