package salix

import (
	"reflect"
	"strings"
)

var defaultFuncs = map[string]any{
	"len":        tmplLen,
	"toUpper":    strings.ToUpper,
	"toLower":    strings.ToLower,
	"hasPrefix":  strings.HasPrefix,
	"trimPrefix": strings.TrimPrefix,
	"hasSuffix":  strings.HasSuffix,
	"trimSuffix": strings.TrimSuffix,
	"trimSpace":  strings.TrimSpace,
	"equalFold":  strings.EqualFold,
	"count":      strings.Count,
	"split":      strings.Split,
	"join":       strings.Join,
}

func tmplLen(v any) int {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Array, reflect.Slice, reflect.String, reflect.Map:
		return val.Len()
	default:
		return -1
	}
}
