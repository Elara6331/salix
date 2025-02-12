package salix

import (
	"errors"
	"reflect"

	"go.elara.ws/salix/ast"
)

// forTag represents a #for tag within a Salix template
type forTag struct{}

func (ft forTag) Run(tc *TagContext, block, args []ast.Node) error {
	if len(args) == 0 || len(args) > 3 {
		return tc.PosError(tc.Tag, "invalid argument amount")
	}

	expr, ok := args[len(args)-1].(ast.Expr)
	if !ok {
		return tc.PosError(args[0], "invalid argument type: %T (expected ast.Expr)", args[0])
	}

	var vars []string
	var in reflect.Value

	if len(args) > 1 {
		for _, arg := range args[:len(args)-1] {
			varName, ok := unwrap(arg).(ast.Ident)
			if !ok {
				return tc.PosError(arg, "invalid argument type: %T (expected ast.Ident)", expr.First)
			}
			vars = append(vars, varName.Value)
		}
	}

	varName, ok := unwrap(expr.First).(ast.Ident)
	if !ok {
		return tc.PosError(expr.First, "invalid argument type: %T (expected ast.Ident)", args[0])
	}
	vars = append(vars, varName.Value)

	if len(expr.Rest) != 1 {
		return tc.PosError(expr.First, "invalid expression (expected 1 element, got %d)", len(expr.Rest))
	}
	rest := expr.Rest[0]

	if rest.Operator.Value != "in" {
		return tc.PosError(expr.First, `invalid operator in expression (expected "in", got %q)`, rest.Operator.Value)
	}

	val, err := tc.GetValue(rest, nil)
	if err != nil {
		return err
	}
	in = reflect.ValueOf(val)

	switch in.Kind() {
	case reflect.Int:
		local := map[string]any{}
		for i := range in.Int() {
			local[vars[0]] = i
			err = tc.Execute(block, local)
			if err != nil {
				return err
			}
		}
	case reflect.Slice, reflect.Array:
		local := map[string]any{}
		for i := 0; i < in.Len(); i++ {
			if len(vars) == 1 {
				local[vars[0]] = in.Index(i).Interface()
			} else if len(vars) == 2 {
				local[vars[0]] = i
				local[vars[1]] = in.Index(i).Interface()
			} else {
				return errors.New("slices and arrays can only use two for loop variables")
			}

			err = tc.Execute(block, local)
			if err != nil {
				return err
			}
		}
	case reflect.Map:
		local := map[string]any{}
		iter := in.MapRange()
		i := 0
		for iter.Next() {
			if len(vars) == 1 {
				local[vars[0]] = iter.Value().Interface()
			} else if len(vars) == 2 {
				local[vars[0]] = iter.Key().Interface()
				local[vars[1]] = iter.Value().Interface()
			} else if len(vars) == 3 {
				local[vars[0]] = i
				local[vars[1]] = iter.Key().Interface()
				local[vars[2]] = iter.Value().Interface()
			}

			err = tc.Execute(block, local)
			if err != nil {
				return err
			}

			i++
		}
	case reflect.Func:
		local := map[string]any{}
		i := 0
		if len(vars) == 1 {
			for val := range in.Seq() {
				local[vars[0]] = val.Interface()
				err = tc.Execute(block, local)
				if err != nil {
					return err
				}
			}
		} else if len(vars) == 2 {
			for val1, val2 := range in.Seq2() {
				local[vars[0]] = val1.Interface()
				local[vars[1]] = val2.Interface()
				err = tc.Execute(block, local)
				if err != nil {
					return err
				}
			}
		} else {
			for val1, val2 := range in.Seq2() {
				local[vars[0]] = i
				local[vars[1]] = val1.Interface()
				local[vars[2]] = val2.Interface()
				err = tc.Execute(block, local)
				if err != nil {
					return err
				}
				i++
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
