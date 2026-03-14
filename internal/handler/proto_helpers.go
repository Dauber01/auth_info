package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	apipb "auth_info/api/gen/api/proto"
	"auth_info/internal/apperr"
	bizdoc "auth_info/internal/biz/document"
	"auth_info/internal/data"
	"auth_info/internal/validation"
)

func writeOperationReply(c *gin.Context, status int, message string) {
	c.JSON(status, &apipb.OperationReply{
		Code:    int32(status),
		Message: message,
	})
}

func writeError(c *gin.Context, err error) {
	if err == nil {
		return
	}
	writeOperationReply(c, apperr.HTTPStatus(err), apperr.Message(err))
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

func structToMap(payload *structpb.Struct) map[string]any {
	if payload == nil {
		return nil
	}
	return payload.AsMap()
}

func structToWordTemplateData(payload *structpb.Struct) bizdoc.WordTemplateData {
	result := bizdoc.WordTemplateData{
		Texts:  make(map[string]bizdoc.RichText),
		Images: make(map[string]bizdoc.ImageValue),
	}
	if payload == nil {
		return result
	}
	m := payload.AsMap()

	if textsRaw, ok := m["texts"].(map[string]any); ok {
		for k, v := range textsRaw {
			vm, ok := v.(map[string]any)
			if !ok {
				if s, ok := v.(string); ok {
					result.Texts[k] = bizdoc.RichText{Runs: []bizdoc.RichRun{{Text: s}}}
				}
				continue
			}

			runsRaw, _ := vm["runs"].([]any)
			rt := bizdoc.RichText{Runs: make([]bizdoc.RichRun, 0, len(runsRaw))}
			for _, r := range runsRaw {
				rm, ok := r.(map[string]any)
				if !ok {
					continue
				}

				run := bizdoc.RichRun{}
				if t, ok := rm["text"].(string); ok {
					run.Text = t
				}
				if b, ok := rm["bold"].(bool); ok {
					run.Bold = b
				}
				if c, ok := rm["color"].(string); ok {
					run.Color = c
				}
				rt.Runs = append(rt.Runs, run)
			}
			result.Texts[k] = rt
		}
	}

	if imagesRaw, ok := m["images"].(map[string]any); ok {
		for k, v := range imagesRaw {
			vm, ok := v.(map[string]any)
			if !ok {
				continue
			}

			iv := bizdoc.ImageValue{}
			if url, ok := vm["image_url"].(string); ok {
				iv.ImageURL = url
			}
			if b, ok := vm["original_size"].(bool); ok {
				iv.OriginalSize = b
			}
			if w, ok := vm["max_width_px"].(float64); ok {
				iv.MaxWidthPx = w
			}
			if h, ok := vm["max_height_px"].(float64); ok {
				iv.MaxHeightPx = h
			}
			result.Images[k] = iv
		}
	}

	return result
}

func dictTypeToProto(model data.DictType) *apipb.DictType {
	return &apipb.DictType{
		Id:          uint64(model.ID),
		Code:        model.Code,
		Name:        model.Name,
		Description: model.Description,
		Sort:        int32(model.Sort),
	}
}

func dictItemToProto(model data.DictItem) *apipb.DictItem {
	return &apipb.DictItem{
		Id:          uint64(model.ID),
		TypeCode:    model.TypeCode,
		ItemKey:     model.ItemKey,
		ItemValue:   model.ItemValue,
		Description: model.Description,
		Sort:        int32(model.Sort),
		Status:      int32(model.Status),
	}
}

func dictTypesToProto(models []data.DictType) []*apipb.DictType {
	items := make([]*apipb.DictType, 0, len(models))
	for _, item := range models {
		items = append(items, dictTypeToProto(item))
	}
	return items
}

func dictItemsToProto(models []data.DictItem) []*apipb.DictItem {
	items := make([]*apipb.DictItem, 0, len(models))
	for _, item := range models {
		items = append(items, dictItemToProto(item))
	}
	return items
}

func validateProtoRules(msg proto.Message) error {
	return validation.ValidateProto(msg)
}

func badRequest(c *gin.Context, err error) {
	writeError(c, apperr.New(apperr.CodeInvalidArgument, err.Error()))
}
