package dict

import (
	apipb "auth_info/api/gen/api/proto"
	bizdict "auth_info/internal/biz/dict"
)

func dictTypeToProto(d bizdict.DictType) *apipb.DictType {
	return &apipb.DictType{
		Id:          uint64(d.ID),
		Code:        d.Code,
		Name:        d.Name,
		Description: d.Description,
		Sort:        int32(d.Sort),
	}
}

func dictItemToProto(d bizdict.DictItem) *apipb.DictItem {
	return &apipb.DictItem{
		Id:          uint64(d.ID),
		TypeCode:    d.TypeCode,
		ItemKey:     d.ItemKey,
		ItemValue:   d.ItemValue,
		Description: d.Description,
		Sort:        int32(d.Sort),
		Status:      int32(d.Status),
	}
}

func dictTypesToProto(items []bizdict.DictType) []*apipb.DictType {
	out := make([]*apipb.DictType, 0, len(items))
	for _, item := range items {
		out = append(out, dictTypeToProto(item))
	}
	return out
}

func dictItemsToProto(items []bizdict.DictItem) []*apipb.DictItem {
	out := make([]*apipb.DictItem, 0, len(items))
	for _, item := range items {
		out = append(out, dictItemToProto(item))
	}
	return out
}
