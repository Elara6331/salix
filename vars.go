/*
 * Salix - Go templating engine
 * Copyright (C) 2023 Elara Musayelyan
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
