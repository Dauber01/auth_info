package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"auth_info/internal/biz"
)

// DictHandler 字典配置 HTTP 处理器
type DictHandler struct {
	uc *biz.DictUseCase
}

// NewDictHandler Wire Provider
func NewDictHandler(uc *biz.DictUseCase) *DictHandler {
	return &DictHandler{uc: uc}
}

// ─── DictType 请求结构体 ────────────────────────────────────────────────────────

type createDictTypeRequest struct {
	Code        string `json:"code"        binding:"required,max=64"`
	Name        string `json:"name"        binding:"required,max=128"`
	Description string `json:"description" binding:"max=256"`
	Sort        int    `json:"sort"`
}

type updateDictTypeRequest struct {
	Name        string `json:"name"        binding:"required,max=128"`
	Description string `json:"description" binding:"max=256"`
	Sort        int    `json:"sort"`
}

// ─── DictItem 请求结构体 ────────────────────────────────────────────────────────

type createDictItemRequest struct {
	TypeCode    string `json:"type_code"   binding:"required,max=64"`
	ItemKey     string `json:"item_key"    binding:"required,max=64"`
	ItemValue   string `json:"item_value"  binding:"required,max=256"`
	Description string `json:"description" binding:"max=256"`
	Sort        int    `json:"sort"`
}

type updateDictItemRequest struct {
	ItemKey     string `json:"item_key"    binding:"required,max=64"`
	ItemValue   string `json:"item_value"  binding:"required,max=256"`
	Description string `json:"description" binding:"max=256"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status"      binding:"min=0,max=1"`
}

// ─── DictType 处理方法 ─────────────────────────────────────────────────────────

// ListDictTypes
// @Summary  获取字典类型列表
// @Tags     Dict
// @Produce  json
// @Success  200 {object} map[string]interface{}
// @Security BearerAuth
// @Router   /dict/types [get]
func (h *DictHandler) ListDictTypes(c *gin.Context) {
	types, err := h.uc.ListDictTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success", "data": types})
}

// CreateDictType
// @Summary  创建字典类型
// @Tags     Dict
// @Accept   json
// @Produce  json
// @Param    body body createDictTypeRequest true "字典类型信息"
// @Success  200 {object} map[string]interface{}
// @Security BearerAuth
// @Router   /dict/types [post]
func (h *DictHandler) CreateDictType(c *gin.Context) {
	var req createDictTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
		return
	}

	if err := h.uc.CreateDictType(req.Code, req.Name, req.Description, req.Sort); err != nil {
		c.JSON(http.StatusConflict, gin.H{"code": http.StatusConflict, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "created successfully"})
}

// UpdateDictType
// @Summary  更新字典类型
// @Tags     Dict
// @Accept   json
// @Produce  json
// @Param    id   path int                  true "字典类型 ID"
// @Param    body body updateDictTypeRequest true "更新信息"
// @Success  200 {object} map[string]interface{}
// @Security BearerAuth
// @Router   /dict/types/{id} [put]
func (h *DictHandler) UpdateDictType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "invalid id"})
		return
	}

	var req updateDictTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
		return
	}

	if err := h.uc.UpdateDictType(uint(id), req.Name, req.Description, req.Sort); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "updated successfully"})
}

// DeleteDictType
// @Summary  删除字典类型
// @Tags     Dict
// @Produce  json
// @Param    id path int true "字典类型 ID"
// @Success  200 {object} map[string]interface{}
// @Security BearerAuth
// @Router   /dict/types/{id} [delete]
func (h *DictHandler) DeleteDictType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "invalid id"})
		return
	}

	if err := h.uc.DeleteDictType(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "deleted successfully"})
}

// ─── DictItem 处理方法 ─────────────────────────────────────────────────────────

// ListDictItems
// @Summary  根据类型编码获取字典数据列表
// @Tags     Dict
// @Produce  json
// @Param    type_code query string true "字典类型编码"
// @Success  200 {object} map[string]interface{}
// @Security BearerAuth
// @Router   /dict/items [get]
func (h *DictHandler) ListDictItems(c *gin.Context) {
	typeCode := c.Query("type_code")
	if typeCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "type_code is required"})
		return
	}

	items, err := h.uc.ListDictItems(typeCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "success", "data": items})
}

// CreateDictItem
// @Summary  创建字典数据
// @Tags     Dict
// @Accept   json
// @Produce  json
// @Param    body body createDictItemRequest true "字典数据信息"
// @Success  200 {object} map[string]interface{}
// @Security BearerAuth
// @Router   /dict/items [post]
func (h *DictHandler) CreateDictItem(c *gin.Context) {
	var req createDictItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
		return
	}

	if err := h.uc.CreateDictItem(req.TypeCode, req.ItemKey, req.ItemValue, req.Description, req.Sort); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "created successfully"})
}

// UpdateDictItem
// @Summary  更新字典数据
// @Tags     Dict
// @Accept   json
// @Produce  json
// @Param    id   path int                  true "字典数据 ID"
// @Param    body body updateDictItemRequest true "更新信息"
// @Success  200 {object} map[string]interface{}
// @Security BearerAuth
// @Router   /dict/items/{id} [put]
func (h *DictHandler) UpdateDictItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "invalid id"})
		return
	}

	var req updateDictItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
		return
	}

	if err := h.uc.UpdateDictItem(uint(id), req.ItemKey, req.ItemValue, req.Description, req.Sort, req.Status); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "updated successfully"})
}

// DeleteDictItem
// @Summary  删除字典数据
// @Tags     Dict
// @Produce  json
// @Param    id path int true "字典数据 ID"
// @Success  200 {object} map[string]interface{}
// @Security BearerAuth
// @Router   /dict/items/{id} [delete]
func (h *DictHandler) DeleteDictItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "invalid id"})
		return
	}

	if err := h.uc.DeleteDictItem(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "deleted successfully"})
}
