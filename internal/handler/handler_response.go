package handler

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/structpb"

	apipb "auth_info/api/gen/api/proto"
	bizdoc "auth_info/internal/biz/document"
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
	_ = c.Error(err)
	c.Abort()
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
