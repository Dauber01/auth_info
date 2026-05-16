package httpx

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"

	apipb "auth_info/api/gen/api/proto"
	"auth_info/internal/apperr"
	"auth_info/internal/validation"
)

// BindJSON 仅做 JSON 绑定，失败时把错误塞进 gin.Context 并返回 false。
func BindJSON(c *gin.Context, req any) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		BadRequest(c, err)
		return false
	}
	return true
}

// ValidateProto 调用 protovalidate 规则校验。
func ValidateProto(c *gin.Context, msg proto.Message) bool {
	if err := validation.ValidateProto(msg); err != nil {
		WriteError(c, err)
		return false
	}
	return true
}

// BindAndValidateJSON 先绑定 JSON 再做 proto 校验。
func BindAndValidateJSON(c *gin.Context, msg proto.Message) bool {
	if !BindJSON(c, msg) {
		return false
	}
	return ValidateProto(c, msg)
}

// BindPathIDAndValidateJSON 解析路径上的 id 参数，回填到 proto 后再校验。
func BindPathIDAndValidateJSON(c *gin.Context, param string, msg proto.Message, assignID func(uint64)) bool {
	id, err := parseUintParam(c, param)
	if err != nil {
		BadRequest(c, err)
		return false
	}
	if !BindJSON(c, msg) {
		return false
	}
	assignID(id)
	return ValidateProto(c, msg)
}

// ValidatePathIDRequest 仅解析路径 id，常用于 DELETE 等无 Body 请求。
func ValidatePathIDRequest(c *gin.Context, param string, msg proto.Message, assignID func(uint64)) bool {
	id, err := parseUintParam(c, param)
	if err != nil {
		BadRequest(c, err)
		return false
	}
	assignID(id)
	return ValidateProto(c, msg)
}

// WriteOperationReply 统一的成功（操作类）响应。
func WriteOperationReply(c *gin.Context, status int, message string) {
	c.JSON(status, &apipb.OperationReply{
		Code:    int32(status),
		Message: message,
	})
}

// WriteError 把错误塞进 gin.Context，由 middleware/error.go 统一渲染。
func WriteError(c *gin.Context, err error) {
	if err == nil {
		return
	}
	_ = c.Error(err)
	c.Abort()
}

// BadRequest 包装为 CodeInvalidArgument 错误。
func BadRequest(c *gin.Context, err error) {
	WriteError(c, apperr.New(apperr.CodeInvalidArgument, err.Error()))
}

func parseUintParam(c *gin.Context, key string) (uint64, error) {
	value := strings.TrimSpace(c.Param(key))
	if value == "" {
		return 0, fmt.Errorf("%s is required", key)
	}
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s", key)
	}
	return id, nil
}
