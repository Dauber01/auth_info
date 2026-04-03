package handler

import (
	apipb "auth_info/api/gen/api/proto"
	"auth_info/internal/data"
)

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
