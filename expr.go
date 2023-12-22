package salix

import (
	"reflect"
	"strings"

	"go.elara.ws/salix/ast"
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
				return nil, ast.PosError(op, "mismatched types in expression (%s and %s)", a.Type(), b.Type())
			}
		case reflect.Map:
			if a.CanConvert(b.Type().Key()) {
				a = a.Convert(b.Type().Key())
			} else {
				return nil, ast.PosError(op, "mismatched types in expression (%s and %s)", a.Type(), b.Type())
			}
		case reflect.String:
			if a.Kind() != reflect.String {
				return nil, ast.PosError(op, "mismatched types in expression (%s and %s)", a.Type(), b.Type())
			}
		default:
			return nil, ast.PosError(op, "the in operator can only be used on strings, arrays, and slices (got %s and %s)", a.Type(), b.Type())
		}
	} else if b.CanConvert(a.Type()) {
		b = b.Convert(a.Type())
	} else {
		return nil, ast.PosError(op, "mismatched types in expression (%s and %s)", a.Type(), b.Type())
	}

	switch op.Value {
	case "==":
		return a.Equal(b), nil
	case "!=":
		return !a.Equal(b), nil
	case "&&":
		if a.Kind() != reflect.Bool || b.Kind() != reflect.Bool {
			return nil, ast.PosError(op, "logical operations may only be performed on boolean values")
		}
		return a.Bool() && b.Bool(), nil
	case "||":
		if a.Kind() != reflect.Bool || b.Kind() != reflect.Bool {
			return nil, ast.PosError(op, "logical operations may only be performed on boolean values")
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
		case reflect.String:
			return a.String() + b.String(), nil
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
			return nil, ast.PosError(op, "modulus operation cannot be performed on floats")
		}
	case "in":
		if a.Kind() == reflect.String && b.Kind() == reflect.String {
			return strings.Contains(b.String(), a.String()), nil
		} else if b.Kind() == reflect.Map {
			return b.MapIndex(a).IsValid(), nil
		} else if b.Kind() == reflect.Slice || b.Kind() == reflect.Array {
			for i := 0; i < b.Len(); i++ {
				if a.Equal(b.Index(i)) {
					return true, nil
				}
			}
			return false, nil
		}
	}
	return false, ast.PosError(op, "unknown operator: %q", op.Value)
}
