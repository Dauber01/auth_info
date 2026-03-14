package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/structpb"

	apipb "auth_info/api/gen/api/proto"
	bizdoc "auth_info/internal/biz/document"
	"auth_info/internal/data"
)

func writeOperationReply(c *gin.Context, status int, message string) {
	c.JSON(status, &apipb.OperationReply{
		Code:    int32(status),
		Message: message,
	})
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

// structToWordTemplateData 将 protobuf Struct 转换为 WordTemplateData。
// 期望的 JSON 结构：
//
//	{
//	  "texts":  { "key": { "runs": [{"text":"...","bold":true,"color":"FF0000"}] } },
//	  "images": { "key": { "image_url":"https://example.com/image.png","original_size":true,"max_width_px":400 } }
//	}
func structToWordTemplateData(payload *structpb.Struct) bizdoc.WordTemplateData {
	result := bizdoc.WordTemplateData{
		Texts:  make(map[string]bizdoc.RichText),
		Images: make(map[string]bizdoc.ImageValue),
	}
	if payload == nil {
		return result
	}
	m := payload.AsMap()

	// 解析 texts
	if textsRaw, ok := m["texts"].(map[string]any); ok {
		for k, v := range textsRaw {
			vm, ok := v.(map[string]any)
			if !ok {
				// 纯字符串简写：{"texts":{"key":"hello"}}
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

	// 解析 images
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

func validateRegisterRequest(req *apipb.RegisterRequest) error {
	username := strings.TrimSpace(req.GetUsername())
	password := req.GetPassword()
	if len(username) < 3 || len(username) > 32 {
		return fmt.Errorf("username length must be between 3 and 32")
	}
	if len(password) < 6 || len(password) > 64 {
		return fmt.Errorf("password length must be between 6 and 64")
	}
	return nil
}

func validateLoginRequest(req *apipb.LoginRequest) error {
	if strings.TrimSpace(req.GetUsername()) == "" {
		return fmt.Errorf("username is required")
	}
	if req.GetPassword() == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

func validateProtoRules(msg proto.Message) error {
	if msg == nil {
		return fmt.Errorf("request is required")
	}

	m := msg.ProtoReflect()
	fields := m.Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		options, ok := field.Options().(*descriptorpb.FieldOptions)
		if !ok || options == nil || !proto.HasExtension(options, apipb.E_Rules) {
			continue
		}

		ext := proto.GetExtension(options, apipb.E_Rules)
		rules, ok := ext.(*apipb.FieldRules)
		if !ok || rules == nil {
			continue
		}

		value := m.Get(field)
		if err := validateFieldRules(field, value, rules); err != nil {
			return err
		}
	}

	return nil
}

func validateFieldRules(field protoreflect.FieldDescriptor, value protoreflect.Value, rules *apipb.FieldRules) error {
	fieldName := string(field.Name())

	switch field.Kind() {
	case protoreflect.StringKind:
		str := strings.TrimSpace(value.String())
		if rules.GetRequired() && str == "" {
			return fmt.Errorf("%s is required", fieldName)
		}
		if rules.MaxLen != nil && len(str) > int(rules.GetMaxLen()) {
			return fmt.Errorf("%s must be at most %d characters", fieldName, rules.GetMaxLen())
		}
		if rules.MinLen != nil && len(str) < int(rules.GetMinLen()) {
			return fmt.Errorf("%s must be at least %d characters", fieldName, rules.GetMinLen())
		}
	case protoreflect.Int32Kind, protoreflect.Int64Kind, protoreflect.Sint32Kind, protoreflect.Sint64Kind,
		protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind:
		if err := validateNumericRange(fieldName, value.Int(), rules); err != nil {
			return err
		}
	case protoreflect.Uint32Kind, protoreflect.Uint64Kind, protoreflect.Fixed32Kind, protoreflect.Fixed64Kind:
		if err := validateNumericRange(fieldName, int64(value.Uint()), rules); err != nil {
			return err
		}
	}

	return nil
}

func validateNumericRange(fieldName string, num int64, rules *apipb.FieldRules) error {
	if rules.Gte != nil && num < rules.GetGte() {
		return fmt.Errorf("%s must be greater than or equal to %d", fieldName, rules.GetGte())
	}
	if rules.Lte != nil && num > rules.GetLte() {
		return fmt.Errorf("%s must be less than or equal to %d", fieldName, rules.GetLte())
	}
	return nil
}

func validateDocumentTemplateRequest(templateName string, payload *structpb.Struct) error {
	if strings.TrimSpace(templateName) == "" {
		return fmt.Errorf("template_name is required")
	}
	if payload == nil {
		return fmt.Errorf("data is required")
	}
	return nil
}

func badRequest(c *gin.Context, err error) {
	writeOperationReply(c, http.StatusBadRequest, err.Error())
}
