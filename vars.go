package salix

import (
	"fmt"
	"reflect"
	"strings"
)

var globalVars = map[string]any{
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
	"sprintf":    fmt.Sprintf,
}

func tmplLen(v any) (int, error) {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Array, reflect.Slice, reflect.String, reflect.Map:
		return val.Len(), nil
	default:
		return 0, fmt.Errorf("cannot get length of %T", v)
	}
}
