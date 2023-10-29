package salix

import (
	"reflect"
	"strings"
)

var globalVars = map[string]reflect.Value{
	"len":        reflect.ValueOf(tmplLen),
	"toUpper":    reflect.ValueOf(strings.ToUpper),
	"toLower":    reflect.ValueOf(strings.ToLower),
	"hasPrefix":  reflect.ValueOf(strings.HasPrefix),
	"trimPrefix": reflect.ValueOf(strings.TrimPrefix),
	"hasSuffix":  reflect.ValueOf(strings.HasSuffix),
	"trimSuffix": reflect.ValueOf(strings.TrimSuffix),
	"trimSpace":  reflect.ValueOf(strings.TrimSpace),
	"equalFold":  reflect.ValueOf(strings.EqualFold),
	"count":      reflect.ValueOf(strings.Count),
	"split":      reflect.ValueOf(strings.Split),
	"join":       reflect.ValueOf(strings.Join),
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
