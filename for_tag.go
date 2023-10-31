package salix

import (
	"errors"
	"reflect"

	"go.elara.ws/salix/ast"
)

var ErrForTagInvalidArgs = errors.New("invalid arguments in for tag")

// forTag represents a #for tag within a Salix template
type forTag struct{}

func (ft forTag) Run(tc *TagContext, block, args []ast.Node) error {
	if len(args) == 0 || len(args) > 2 {
		return ErrForTagInvalidArgs
	}

	var expr ast.Expr
	if len(args) == 1 {
		expr2, ok := args[0].(ast.Expr)
		if !ok {
			return ErrForTagInvalidArgs
		}
		expr = expr2
	} else if len(args) == 2 {
		expr2, ok := args[1].(ast.Expr)
		if !ok {
			return ErrForTagInvalidArgs
		}
		expr = expr2
	}

	var vars []string
	var in reflect.Value

	if len(args) == 2 {
		varName, ok := unwrap(args[0]).(ast.Ident)
		if !ok {
			return ErrForTagInvalidArgs
		}
		vars = append(vars, varName.Value)

	}

	varName, ok := unwrap(expr.First).(ast.Ident)
	if !ok {
		return ErrForTagInvalidArgs
	}
	vars = append(vars, varName.Value)

	if len(expr.Rest) != 1 {
		return ErrForTagInvalidArgs
	}
	rest := expr.Rest[0]

	if rest.Operator.Value != "in" {
		return ErrForTagInvalidArgs
	}

	val, err := tc.GetValue(rest, nil)
	if err != nil {
		return err
	}
	in = reflect.ValueOf(val)

	switch in.Kind() {
	case reflect.Slice, reflect.Array:
		local := map[string]any{}
		for i := 0; i < in.Len(); i++ {
			if len(vars) == 1 {
				local[vars[0]] = in.Index(i).Interface()
			} else if len(vars) == 2 {
				local[vars[0]] = i
				local[vars[1]] = in.Index(i).Interface()
			}

			err = tc.Execute(block, local)
			if err != nil {
				return err
			}
		}
	case reflect.Map:
		local := map[string]any{}
		iter := in.MapRange()
		for iter.Next() {
			if len(vars) == 1 {
				local[vars[0]] = iter.Value().Interface()
			} else if len(vars) == 2 {
				local[vars[0]] = iter.Key().Interface()
				local[vars[1]] = iter.Value().Interface()
			}

			err = tc.Execute(block, local)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func unwrap(n ast.Node) ast.Node {
	if v, ok := n.(ast.Value); ok {
		return v.Node
	}
	return n
}
