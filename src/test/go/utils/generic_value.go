package utils

import (
	"fmt"
	"strconv"
)

type GenericValue struct {
	values map[string]any
}

func NewGenericValue() *GenericValue {
	return &GenericValue{
		values: make(map[string]any),
	}
}

func GenericValueOf(values map[string]any) *GenericValue {
	gv := NewGenericValue()
	gv.PutAll(values)
	return gv
}

func (gv *GenericValue) Size() int {
	return len(gv.values)
}

func (gv *GenericValue) IsEmpty() bool {
	return len(gv.values) == 0
}

func (gv *GenericValue) ContainsKey(key string) bool {
	_, ok := gv.values[key]
	return ok
}

func (gv *GenericValue) Get(key string) any {
	return gv.values[key]
}

func (gv *GenericValue) GetChild(key string) *GenericValue {
	v := gv.Get(key)

	switch t := v.(type) {
	case *GenericValue:
		return t
	case map[string]any:
		return GenericValueOf(t)
	default:
		return NewGenericValue()
	}
}

func (gv *GenericValue) GetList(key string) []*GenericValue {
	v := gv.Get(key)
	if v == nil {
		return nil
	}

	list, ok := v.([]any)
	if !ok {
		return nil
	}

	result := make([]*GenericValue, 0, len(list))

	for _, item := range list {
		switch t := item.(type) {
		case *GenericValue:
			result = append(result, t)

		case GenericValue:
			result = append(result, &t)

		case map[string]any:
			result = append(result, GenericValueOf(t))
		}
	}

	return result
}

func (gv *GenericValue) GetString(key string) string {
	v := gv.Get(key)
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

func (gv *GenericValue) GetInt(key string) *int {
	v := gv.Get(key)

	switch t := v.(type) {
	case int:
		return &t
	case int32:
		i := int(t)
		return &i
	case int64:
		i := int(t)
		return &i
	case string:
		if i, err := strconv.Atoi(t); err == nil {
			return &i
		}
	}

	return nil
}

func (gv *GenericValue) GetInt64(key string) *int64 {
	v := gv.Get(key)

	switch t := v.(type) {
	case int64:
		return &t
	case int:
		i := int64(t)
		return &i
	case string:
		if i, err := strconv.ParseInt(t, 10, 64); err == nil {
			return &i
		}
	}

	return nil
}

func (gv *GenericValue) GetBool(key string) *bool {
	v := gv.Get(key)

	switch t := v.(type) {
	case bool:
		return &t
	case string:
		if b, err := strconv.ParseBool(t); err == nil {
			return &b
		}
	}

	return nil
}

func (gv *GenericValue) Put(key string, value any) any {
	old := gv.values[key]
	gv.values[key] = value
	return old
}

func (gv *GenericValue) Remove(key string) any {
	old, exists := gv.values[key]
	if exists {
		delete(gv.values, key)
	}
	return old
}

func (gv *GenericValue) PutAll(m map[string]any) {
	for k, v := range m {
		gv.Put(k, v)
	}
}

func (gv *GenericValue) Clear() {
	gv.values = make(map[string]any)
}

func (gv *GenericValue) Keys() []string {
	keys := make([]string, 0, len(gv.values))
	for k := range gv.values {
		keys = append(keys, k)
	}
	return keys
}

func (gv *GenericValue) Values() []any {
	values := make([]any, 0, len(gv.values))
	for _, v := range gv.values {
		values = append(values, v)
	}
	return values
}

func (gv *GenericValue) Equals(other *GenericValue) bool {
	if other == nil {
		return false
	}

	if len(gv.values) != len(other.values) {
		return false
	}

	for k, v := range gv.values {
		if other.values[k] != v {
			return false
		}
	}

	return true
}

func (gv *GenericValue) String() string {
	return fmt.Sprintf("%v", gv.values)
}

func (gv *GenericValue) Update(source *GenericValue, keyMap map[string]string) *GenericValue {
	if len(keyMap) == 0 {
		for k, v := range source.values {
			gv.Put(k, v)
		}
		return gv
	}

	for sourceKey, targetKey := range keyMap {
		gv.Put(targetKey, source.Get(sourceKey))
	}

	return gv
}

func FindByKeyValue(list []*GenericValue, key string, value any) *GenericValue {
	for _, item := range list {
		if item.Get(key) == value {
			return item
		}
	}
	return nil
}
