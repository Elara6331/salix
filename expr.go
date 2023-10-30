package salix

import (
	"reflect"
	"strings"

	"go.elara.ws/salix/internal/ast"
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
		a = reflect.ValueOf(t.performOp(a, b, exprB.Operator))
	}

	return a.Interface(), nil
}

func (t *Template) performOp(a, b reflect.Value, op ast.Op) any {
	if op.Op() == "in" {
		switch b.Kind() {
		case reflect.Slice, reflect.Array:
			if a.CanConvert(b.Type().Elem()) {
				a = a.Convert(b.Type().Elem())
			} else {
				panic("todo: invalid in operation")
			}
		case reflect.String:
			if a.Kind() != reflect.String {
				panic("todo: invalid in operation")
			}
		}
	} else if b.CanConvert(a.Type()) {
		b = b.Convert(a.Type())
	} else {
		panic("todo: invalid operation")
	}

	switch op.Op() {
	case "==":
		return a.Equal(b)
	case "&&":
		if a.Kind() != reflect.Bool || b.Kind() != reflect.Bool {
			panic("todo: invalid logical")
		}
		return a.Bool() && b.Bool()
	case "||":
		if a.Kind() != reflect.Bool || b.Kind() != reflect.Bool {
			panic("todo: invalid logical")
		}
		return a.Bool() || b.Bool()
	case ">=":
		switch a.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return a.Int() >= b.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return a.Uint() >= b.Uint()
		case reflect.Float64, reflect.Float32:
			return a.Float() >= b.Float()
		}
	case "<=":
		switch a.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return a.Int() <= b.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return a.Uint() <= b.Uint()
		case reflect.Float64, reflect.Float32:
			return a.Float() <= b.Float()
		}
	case ">":
		switch a.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return a.Int() > b.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return a.Uint() > b.Uint()
		case reflect.Float64, reflect.Float32:
			return a.Float() > b.Float()
		}
	case "<":
		switch a.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return a.Int() < b.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return a.Uint() < b.Uint()
		case reflect.Float64, reflect.Float32:
			return a.Float() < b.Float()
		}
	case "+":
		switch a.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return a.Int() + b.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return a.Uint() + b.Uint()
		case reflect.Float64, reflect.Float32:
			return a.Float() + b.Float()
		}
	case "-":
		switch a.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return a.Int() - b.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return a.Uint() - b.Uint()
		case reflect.Float64, reflect.Float32:
			return a.Float() - b.Float()
		}
	case "*":
		switch a.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return a.Int() * b.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return a.Uint() * b.Uint()
		case reflect.Float64, reflect.Float32:
			return a.Float() * b.Float()
		}
	case "/":
		switch a.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return a.Int() / b.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return a.Uint() / b.Uint()
		case reflect.Float64, reflect.Float32:
			return a.Float() / b.Float()
		}
	case "%":
		switch a.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return a.Int() % b.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return a.Uint() % b.Uint()
		case reflect.Float64, reflect.Float32:
			return a.Float() % b.Float()
	case "in":
		if a.Kind() == reflect.String && b.Kind() == reflect.String {
			return strings.Contains(b.String(), a.String())
		} else {
			for i := 0; i < b.Len(); i++ {
				if a.Equal(b.Index(i)) {
					return true
				}
			}
			return false
		}
	}
	return false
}
