package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	apipb "auth_info/api/gen/api/proto"
	bizdict "auth_info/internal/biz/dict"
)

// DictHandler handles dictionary HTTP APIs.
type DictHandler struct {
	uc *bizdict.UseCase
}

// NewDictHandler Wire Provider
func NewDictHandler(uc *bizdict.UseCase) *DictHandler {
	return &DictHandler{uc: uc}
}

// ListDictTypes
// @Summary  List dictionary types
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
// @Summary  Create dictionary type
// @Tags     Dict
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Router   /dict/types [post]
func (h *DictHandler) CreateDictType(c *gin.Context) {
	var req apipb.CreateDictTypeRequest
	if !bindAndValidateJSON(c, &req) {
		return
	}

	if err := h.uc.CreateDictType(c.Request.Context(), req.GetCode(), req.GetName(), req.GetDescription(), int(req.GetSort())); err != nil {
		writeError(c, err)
		return
	}

	writeOperationReply(c, http.StatusOK, "created successfully")
}

// UpdateDictType
// @Summary  Update dictionary type
// @Tags     Dict
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Router   /dict/types/{id} [put]
func (h *DictHandler) UpdateDictType(c *gin.Context) {
	var req apipb.UpdateDictTypeRequest
	if !bindPathIDAndValidateJSON(c, "id", &req, func(id uint64) {
		req.Id = id
	}) {
		return
	}

	if err := h.uc.UpdateDictType(c.Request.Context(), uint(req.GetId()), req.GetName(), req.GetDescription(), int(req.GetSort())); err != nil {
		writeError(c, err)
		return
	}

	writeOperationReply(c, http.StatusOK, "updated successfully")
}

// DeleteDictType
// @Summary  Delete dictionary type
// @Tags     Dict
// @Produce  json
// @Security BearerAuth
// @Router   /dict/types/{id} [delete]
func (h *DictHandler) DeleteDictType(c *gin.Context) {
	req := apipb.DeleteDictTypeRequest{}
	if !validatePathIDRequest(c, "id", &req, func(id uint64) {
		req.Id = id
	}) {
		return
	}

	if err := h.uc.DeleteDictType(c.Request.Context(), uint(req.GetId())); err != nil {
		writeError(c, err)
		return
	}

	writeOperationReply(c, http.StatusOK, "deleted successfully")
}

// ListDictItems
// @Summary  List dictionary items by type code
// @Tags     Dict
// @Produce  json
// @Security BearerAuth
// @Router   /dict/items [get]
func (h *DictHandler) ListDictItems(c *gin.Context) {
	req := apipb.ListDictItemsRequest{TypeCode: strings.TrimSpace(c.Query("type_code"))}
	if !validateProtoMessage(c, &req) {
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
// @Summary  Create dictionary item
// @Tags     Dict
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Router   /dict/items [post]
func (h *DictHandler) CreateDictItem(c *gin.Context) {
	var req apipb.CreateDictItemRequest
	if !bindAndValidateJSON(c, &req) {
		return
	}

	if err := h.uc.CreateDictItem(c.Request.Context(), req.GetTypeCode(), req.GetItemKey(), req.GetItemValue(), req.GetDescription(), int(req.GetSort())); err != nil {
		writeError(c, err)
		return
	}

	writeOperationReply(c, http.StatusOK, "created successfully")
}

// UpdateDictItem
// @Summary  Update dictionary item
// @Tags     Dict
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Router   /dict/items/{id} [put]
func (h *DictHandler) UpdateDictItem(c *gin.Context) {
	var req apipb.UpdateDictItemRequest
	if !bindPathIDAndValidateJSON(c, "id", &req, func(id uint64) {
		req.Id = id
	}) {
		return
	}

	if err := h.uc.UpdateDictItem(c.Request.Context(), uint(req.GetId()), req.GetItemKey(), req.GetItemValue(), req.GetDescription(), int(req.GetSort()), int(req.GetStatus())); err != nil {
		writeError(c, err)
		return
	}

	writeOperationReply(c, http.StatusOK, "updated successfully")
}

// DeleteDictItem
// @Summary  Delete dictionary item
// @Tags     Dict
// @Produce  json
// @Security BearerAuth
// @Router   /dict/items/{id} [delete]
func (h *DictHandler) DeleteDictItem(c *gin.Context) {
	req := apipb.DeleteDictItemRequest{}
	if !validatePathIDRequest(c, "id", &req, func(id uint64) {
		req.Id = id
	}) {
		return
	}

	if err := h.uc.DeleteDictItem(c.Request.Context(), uint(req.GetId())); err != nil {
		writeError(c, err)
		return
	}

	writeOperationReply(c, http.StatusOK, "deleted successfully")
}
