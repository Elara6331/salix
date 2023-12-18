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
