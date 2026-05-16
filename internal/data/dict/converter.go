package dict

import bizdict "auth_info/internal/biz/dict"

func dictTypeToBiz(m *DictType) *bizdict.DictType {
	if m == nil {
		return nil
	}
	return &bizdict.DictType{
		ID:          m.ID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		Sort:        m.Sort,
	}
}

func dictTypeFromBiz(d *bizdict.DictType) *DictType {
	if d == nil {
		return nil
	}
	m := &DictType{
		Code:        d.Code,
		Name:        d.Name,
		Description: d.Description,
		Sort:        d.Sort,
	}
	m.ID = d.ID
	m.CreatedAt = d.CreatedAt
	m.UpdatedAt = d.UpdatedAt
	return m
}

func dictTypesToBiz(models []DictType) []bizdict.DictType {
	items := make([]bizdict.DictType, 0, len(models))
	for i := range models {
		items = append(items, *dictTypeToBiz(&models[i]))
	}
	return items
}

func dictItemToBiz(m *DictItem) *bizdict.DictItem {
	if m == nil {
		return nil
	}
	return &bizdict.DictItem{
		ID:          m.ID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		TypeCode:    m.TypeCode,
		ItemKey:     m.ItemKey,
		ItemValue:   m.ItemValue,
		Description: m.Description,
		Sort:        m.Sort,
		Status:      m.Status,
	}
}

func dictItemFromBiz(d *bizdict.DictItem) *DictItem {
	if d == nil {
		return nil
	}
	m := &DictItem{
		TypeCode:    d.TypeCode,
		ItemKey:     d.ItemKey,
		ItemValue:   d.ItemValue,
		Description: d.Description,
		Sort:        d.Sort,
		Status:      d.Status,
	}
	m.ID = d.ID
	m.CreatedAt = d.CreatedAt
	m.UpdatedAt = d.UpdatedAt
	return m
}

func dictItemsToBiz(models []DictItem) []bizdict.DictItem {
	items := make([]bizdict.DictItem, 0, len(models))
	for i := range models {
		items = append(items, *dictItemToBiz(&models[i]))
	}
	return items
}
