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
	"errors"
	"reflect"
	"strings"

	"go.elara.ws/salix/ast"
)

var (
	ErrModulusFloat     = errors.New("modulo operation cannot be performed on floats")
	ErrTypeMismatch     = errors.New("mismatched types")
	ErrLogicalNonBool   = errors.New("logical operations may only be performed on boolean values")
	ErrInOpInvalidTypes = errors.New("the in operator can only be used on strings, arrays, and slices")
)

func (t *Template) evalExpr(expr ast.Expr, local map[string]any) (any, error) {
	val, err := t.getValue(expr.First, local)
	if err != nil {
		return nil, err
	}
	a := reflect.ValueOf(val)

	for _, exprB := range expr.Rest {
		val, err := t.getValue(exprB.First, local)
		if err != nil {
			return nil, err
		}
		b := reflect.ValueOf(val)

		result, err := t.performOp(a, b, exprB.Operator)
		if err != nil {
			return nil, err
		}

		a = reflect.ValueOf(result)
	}

	return a.Interface(), nil
}

func (t *Template) performOp(a, b reflect.Value, op ast.Operator) (any, error) {
	if op.Value == "in" {
		switch b.Kind() {
		case reflect.Slice, reflect.Array:
			if a.CanConvert(b.Type().Elem()) {
				a = a.Convert(b.Type().Elem())
			} else {
				return nil, t.posError(op, "%w (%s and %s)", ErrTypeMismatch, a.Type(), b.Type())
			}
		case reflect.String:
			if a.Kind() != reflect.String {
				return nil, t.posError(op, "%w (%s and %s)", ErrTypeMismatch, a.Type(), b.Type())
			}
		default:
			return nil, t.posError(op, "%w (got %s and %s)", ErrInOpInvalidTypes, a.Type(), b.Type())
		}
	} else if b.CanConvert(a.Type()) {
		b = b.Convert(a.Type())
	} else {
		return nil, t.posError(op, "%w (%s and %s)", ErrTypeMismatch, a.Type(), b.Type())
	}

	switch op.Value {
	case "==":
		return a.Equal(b), nil
	case "&&":
		if a.Kind() != reflect.Bool || b.Kind() != reflect.Bool {
			return nil, t.posError(op, "%w", ErrLogicalNonBool)
		}
		return a.Bool() && b.Bool(), nil
	case "||":
		if a.Kind() != reflect.Bool || b.Kind() != reflect.Bool {
			return nil, t.posError(op, "%w", ErrLogicalNonBool)
		}
		return a.Bool() || b.Bool(), nil
	case ">=":
		switch a.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return a.Int() >= b.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return a.Uint() >= b.Uint(), nil
		case reflect.Float64, reflect.Float32:
			return a.Float() >= b.Float(), nil
		}
	case "<=":
		switch a.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return a.Int() <= b.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return a.Uint() <= b.Uint(), nil
		case reflect.Float64, reflect.Float32:
			return a.Float() <= b.Float(), nil
		}
	case ">":
		switch a.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return a.Int() > b.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return a.Uint() > b.Uint(), nil
		case reflect.Float64, reflect.Float32:
			return a.Float() > b.Float(), nil
		}
	case "<":
		switch a.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return a.Int() < b.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return a.Uint() < b.Uint(), nil
		case reflect.Float64, reflect.Float32:
			return a.Float() < b.Float(), nil
		}
	case "+":
		switch a.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return a.Int() + b.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return a.Uint() + b.Uint(), nil
		case reflect.Float64, reflect.Float32:
			return a.Float() + b.Float(), nil
		}
	case "-":
		switch a.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return a.Int() - b.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return a.Uint() - b.Uint(), nil
		case reflect.Float64, reflect.Float32:
			return a.Float() - b.Float(), nil
		}
	case "*":
		switch a.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return a.Int() * b.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return a.Uint() * b.Uint(), nil
		case reflect.Float64, reflect.Float32:
			return a.Float() * b.Float(), nil
		}
	case "/":
		switch a.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return a.Int() / b.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return a.Uint() / b.Uint(), nil
		case reflect.Float64, reflect.Float32:
			return a.Float() / b.Float(), nil
		}
	case "%":
		switch a.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return a.Int() % b.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return a.Uint() % b.Uint(), nil
		case reflect.Float64, reflect.Float32:
			return nil, t.posError(op, "%w", ErrModulusFloat)
		}
	case "in":
		if a.Kind() == reflect.String && b.Kind() == reflect.String {
			return strings.Contains(b.String(), a.String()), nil
		} else {
			for i := 0; i < b.Len(); i++ {
				if a.Equal(b.Index(i)) {
					return true, nil
				}
			}
			return false, nil
		}
	}
	return false, nil
}
