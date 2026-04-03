package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"

	"auth_info/internal/apperr"
	"auth_info/internal/validation"
)

func bindJSON(c *gin.Context, req any) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		badRequest(c, err)
		return false
	}
	return true
}

func validateProtoMessage(c *gin.Context, msg proto.Message) bool {
	if err := validateProtoRules(msg); err != nil {
		writeError(c, err)
		return false
	}
	return true
}

func bindAndValidateJSON(c *gin.Context, msg proto.Message) bool {
	if !bindJSON(c, msg) {
		return false
	}
	return validateProtoMessage(c, msg)
}

func bindPathIDAndValidateJSON(c *gin.Context, param string, msg proto.Message, assignID func(uint64)) bool {
	id, err := parseUintParam(c, param)
	if err != nil {
		badRequest(c, err)
		return false
	}
	if !bindJSON(c, msg) {
		return false
	}
	assignID(id)
	return validateProtoMessage(c, msg)
}

func validatePathIDRequest(c *gin.Context, param string, msg proto.Message, assignID func(uint64)) bool {
	id, err := parseUintParam(c, param)
	if err != nil {
		badRequest(c, err)
		return false
	}
	assignID(id)
	return validateProtoMessage(c, msg)
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

func validateProtoRules(msg proto.Message) error {
	return validation.ValidateProto(msg)
}

func badRequest(c *gin.Context, err error) {
	writeError(c, apperr.New(apperr.CodeInvalidArgument, err.Error()))
}
