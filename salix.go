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
	"fmt"
	"html"
	"io"
	"reflect"

	"go.elara.ws/salix/ast"
)

var (
	ErrNoSuchFunc           = errors.New("no such function")
	ErrNoSuchVar            = errors.New("no such variable")
	ErrNoSuchMethod         = errors.New("no such method")
	ErrNoSuchField          = errors.New("no such field")
	ErrNoSuchTag            = errors.New("no such tag")
	ErrNotOperatorNonBool   = errors.New("not operator cannot be used on a non-bool value")
	ErrParamNumMismatch     = errors.New("incorrect parameter amount")
	ErrIncorrectParamType   = errors.New("incorrect parameter type for function")
	ErrEndTagWithoutStart   = errors.New("end tag without a start tag")
	ErrIncorrectIndexType   = errors.New("incorrect index type")
	ErrIndexOutOfRange      = errors.New("index out of range")
	ErrMapIndexNotFound     = errors.New("map index not found")
	ErrMapInvalidIndexType  = errors.New("invalid map index type")
	ErrFuncTooManyReturns   = errors.New("template functions can only have two return values")
	ErrFuncNoReturns        = errors.New("template functions must return at least one value")
	ErrFuncSecondReturnType = errors.New("the second return value of a template function must be an error")
	ErrTernaryCondBool      = errors.New("ternary condition must be a boolean")
)

// HTML represents unescaped HTML strings
type HTML string

// Template represents a Salix template
type Template struct {
	ns   *Namespace
	name string
	ast  []ast.Node

	escapeHTML bool

	tags   map[string]Tag
	vars   map[string]reflect.Value
	macros map[string][]ast.Node
}

// WithVarMap returns a copy of the template with its variable map set to m.
func (t Template) WithVarMap(m map[string]any) Template {
	t.vars = map[string]reflect.Value{}
	if m != nil {
		for k, v := range m {
			t.vars[k] = reflect.ValueOf(v)
		}
	}
	return t
}

// WithTagMap returns a copy of the template with its tag map set to m.
func (t Template) WithTagMap(m map[string]Tag) Template {
	// Make sure the tag map is never nil to avoid panics
	if m == nil {
		m = map[string]Tag{}
	}
	t.tags = m
	return t
}

// WithEscapeHTML returns a copy of the template with HTML escaping enabled or disabled.
// The HTML escaping functionality is NOT context-aware.
// Using the HTML type allows you to get around the escaping if needed.
func (t Template) WithEscapeHTML(b bool) Template {
	t.escapeHTML = true
	return t
}

// Execute executes a parsed template and writes
// the result to w.
func (t Template) Execute(w io.Writer) error {
	t.macros = map[string][]ast.Node{}
	return t.execute(w, t.ast, nil)
}

func (t *Template) execute(w io.Writer, nodes []ast.Node, local map[string]any) error {
	for i := 0; i < len(nodes); i++ {
		switch node := nodes[i].(type) {
		case ast.Text:
			_, err := w.Write(node.Data)
			if err != nil {
				return t.posError(node, "%w", err)
			}
		case ast.Tag:
			newOffset, err := t.execTag(node, w, nodes, i, local)
			if err != nil {
				return err
			}
			i = newOffset
		case ast.EndTag:
			// We should never see an end tag here because it
			// should be taken care of by execTag, so if we do,
			// return an error because execTag was never called,
			// which means there was no start tag.
			return ErrEndTagWithoutStart
		case ast.ExprTag:
			v, err := t.getValue(node.Value, local)
			if err != nil {
				return err
			}
			if _, ok := v.(ast.Assignment); ok {
				continue
			}
			_, err = io.WriteString(w, t.toString(v))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *Template) toString(v any) string {
	if h, ok := v.(HTML); ok {
		return string(h)
	} else if t.escapeHTML {
		return html.EscapeString(fmt.Sprint(v))
	}
	return fmt.Sprint(v)
}

// getBlock gets all the nodes in the input, up to the end tag with the given name
func (t *Template) getBlock(nodes []ast.Node, offset, startLine int, name string) []ast.Node {
	var out []ast.Node
	tagAmount := 1
	for i := offset; i < len(nodes); i++ {
		switch node := nodes[i].(type) {
		case ast.Tag:
			// If we encounter another tag with the same name,
			// increment tagAmount so that we know that the next
			// end tag isn't the end of this tag.
			if node.Name.Value == name {
				tagAmount++
			}
			out = append(out, node)
		case ast.EndTag:
			if node.Name.Value == name {
				tagAmount--
			}
			// Once tagAmount is zero (all the tags of the same name
			// have been closed with an end tag), we can return
			// the nodes we've accumulated.
			if tagAmount == 0 {
				return out
			} else {
				out = append(out, node)
			}
		default:
			out = append(out, node)
		}
	}
	return out
}

// getValue gets a Go value from an AST node
func (t *Template) getValue(node ast.Node, local map[string]any) (any, error) {
	switch node := node.(type) {
	case ast.Value:
		return t.unwrapASTValue(node, local)
	case ast.Ident:
		val, err := t.getVar(node, local)
		if err != nil {
			return nil, err
		}
		return val.Interface(), nil
	case ast.String:
		return node.Value, nil
	case ast.Float:
		return node.Value, nil
	case ast.Integer:
		return node.Value, nil
	case ast.Bool:
		return node.Value, nil
	case ast.Expr:
		return t.evalExpr(node, local)
	case ast.FuncCall:
		return t.execFuncCall(node, local)
	case ast.Index:
		return t.getIndex(node, local)
	case ast.FieldAccess:
		return t.getField(node, local)
	case ast.MethodCall:
		return t.execMethodCall(node, local)
	case ast.Ternary:
		return t.evalTernary(node, local)
	case ast.VariableOr:
		return t.evalVariableOr(node, local)
	case ast.Assignment:
		return node, t.handleAssignment(node, local)
	default:
		return nil, nil
	}
}

// unwrapASTValue unwraps an ast.Value node into its underlying value
func (t *Template) unwrapASTValue(node ast.Value, local map[string]any) (any, error) {
	v, err := t.getValue(node.Node, local)
	if err != nil {
		return nil, err
	}

	if node.Not {
		rval := reflect.ValueOf(v)
		if rval.Kind() != reflect.Bool {
			return nil, ErrNotOperatorNonBool
		}
		return !rval.Bool(), nil
	}

	return v, err
}

// getVar tries to get a variable from the local map. If it's not found,
// it'll try the global variable map. If it doesn't exist in either map,
// it will return an error.
func (t *Template) getVar(id ast.Ident, local map[string]any) (reflect.Value, error) {
	if local != nil {
		v, ok := local[id.Value]
		if ok {
			return reflect.ValueOf(v), nil
		}
	}

	v, ok := t.vars[id.Value]
	if ok {
		return v, nil
	}

	v, ok = t.ns.getVar(id.Value)
	if ok {
		return v, nil
	}

	v, ok = globalVars[id.Value]
	if ok {
		return v, nil
	}

	return reflect.Value{}, t.posError(id, "%w: %s", ErrNoSuchVar, id.Value)
}

func (t *Template) getTag(name string) (Tag, bool) {
	tag, ok := t.tags[name]
	if ok {
		return tag, true
	}

	tag, ok = t.ns.getTag(name)
	if ok {
		return tag, true
	}

	tag, ok = globalTags[name]
	if ok {
		return tag, true
	}

	return nil, false
}

// execTag executes a tag
func (t *Template) execTag(node ast.Tag, w io.Writer, nodes []ast.Node, i int, local map[string]any) (newOffset int, err error) {
	tag, ok := t.getTag(node.Name.Value)
	if !ok {
		return 0, t.posError(node, "%w: %s", ErrNoSuchTag, node.Name.Value)
	}

	var block []ast.Node
	if node.HasBody {
		block = t.getBlock(nodes, i+1, node.Position.Line, node.Name.Value)
		i += len(block) + 1
	}

	tc := &TagContext{w, t, local}

	err = tag.Run(tc, block, node.Params)
	if err != nil {
		return 0, err
	}

	return i, nil
}

// execFuncCall executes a function call
func (t *Template) execFuncCall(fc ast.FuncCall, local map[string]any) (any, error) {
	fn, err := t.getVar(fc.Name, local)
	if err != nil {
		return nil, t.posError(fc, "%w: %s", ErrNoSuchFunc, fc.Name.Value)
	}
	return t.execFunc(fn, fc, fc.Params, local)
}

// getIndex tries to evaluate an ast.Index node by indexing the underlying value.
func (t *Template) getIndex(i ast.Index, local map[string]any) (any, error) {
	val, err := t.getValue(i.Value, local)
	if err != nil {
		return nil, err
	}

	index, err := t.getValue(i.Index, local)
	if err != nil {
		return nil, err
	}

	rval := reflect.ValueOf(val)
	rindex := reflect.ValueOf(index)
	switch rval.Kind() {
	case reflect.Slice, reflect.Array:
		intType := reflect.TypeOf(0)
		if rindex.CanConvert(intType) {
			rindex = rindex.Convert(intType)
		} else {
			return nil, ErrIncorrectIndexType
		}

		intIndex := rindex.Interface().(int)
		if intIndex < rval.Len() {
			return rval.Index(intIndex).Interface(), nil
		} else {
			return nil, t.posError(i, "%w: %d", ErrIndexOutOfRange, intIndex)
		}
	case reflect.Map:
		if rindex.CanConvert(rval.Type().Key()) {
			rindex = rindex.Convert(rval.Type().Key())
		} else {
			return nil, t.posError(i, "%w: %T (expected %s)", ErrMapInvalidIndexType, index, rval.Type().Key())
		}
		if out := rval.MapIndex(rindex); out.IsValid() {
			return out.Interface(), nil
		} else {
			return nil, t.posError(i, "%w: %q", ErrMapIndexNotFound, index)
		}
	}
	return nil, nil
}

// getField tries to get a struct field from the underlying value
func (t *Template) getField(fa ast.FieldAccess, local map[string]any) (any, error) {
	val, err := t.getValue(fa.Value, local)
	if err != nil {
		return nil, err
	}
	rval := reflect.ValueOf(val)
	for rval.Kind() == reflect.Pointer {
		rval = rval.Elem()
	}
	field := rval.FieldByName(fa.Name.Value)
	if !field.IsValid() {
		return nil, t.posError(fa, "%w: %s", ErrNoSuchField, fa.Name.Value)
	}
	return field.Interface(), nil
}

// execMethodCall executes a method call on the underlying value
func (t *Template) execMethodCall(mc ast.MethodCall, local map[string]any) (any, error) {
	val, err := t.getValue(mc.Value, local)
	if err != nil {
		return nil, err
	}
	rval := reflect.ValueOf(val)
	for rval.Kind() == reflect.Pointer {
		rval = rval.Elem()
	}
	mtd := rval.MethodByName(mc.Name.Value)
	if !mtd.IsValid() {
		return nil, t.posError(mc, "%w: %s", ErrNoSuchMethod, mc.Name.Value)
	}
	return t.execFunc(mtd, mc, mc.Params, local)
}

// execFunc executes a function call
func (t *Template) execFunc(fn reflect.Value, node ast.Node, args []ast.Node, local map[string]any) (any, error) {
	fnType := fn.Type()
	if fnType.NumIn() != len(args) {
		return nil, t.posError(node, "%w: %d (expected %d)", ErrParamNumMismatch, len(args), fnType.NumIn())
	}

	if err := validateFunc(fnType); err != nil {
		return nil, t.posError(node, "%w", err)
	}

	params := make([]reflect.Value, fnType.NumIn())
	for i, arg := range args {
		if _, ok := arg.(ast.Assignment); ok {
			return nil, t.posError(arg, "assignment cannot be used as a function argument")
		}
		paramVal, err := t.getValue(arg, local)
		if err != nil {
			return nil, err
		}
		params[i] = reflect.ValueOf(paramVal)
		if params[i].CanConvert(fnType.In(i)) {
			params[i] = params[i].Convert(fnType.In(i))
		} else {
			return nil, t.posError(node, "%w", ErrIncorrectParamType)
		}
	}

	ret := fn.Call(params)
	if len(ret) == 1 {
		retv := ret[0].Interface()
		if err, ok := retv.(error); ok {
			return nil, err
		}
		return ret[0].Interface(), nil
	} else {
		return ret[0].Interface(), ret[1].Interface().(error)
	}
}

func (t *Template) evalTernary(tr ast.Ternary, local map[string]any) (any, error) {
	condVal, err := t.getValue(tr.Condition, local)
	if err != nil {
		return nil, t.posError(tr.Condition, "%w", ErrTernaryCondBool)
	}

	cond, ok := condVal.(bool)
	if !ok {
		return nil, errors.New("ternary condition must be a boolean")
	}

	if cond {
		return t.getValue(tr.IfTrue, local)
	} else {
		return t.getValue(tr.Else, local)
	}
}

func (t *Template) evalVariableOr(vo ast.VariableOr, local map[string]any) (any, error) {
	val, err := t.getVar(vo.Variable, local)
	if err != nil {
		return t.getValue(vo.Or, local)
	}
	return val.Interface(), nil
}

func (t *Template) handleAssignment(a ast.Assignment, local map[string]any) error {
	val, err := t.getValue(a.Value, local)
	if err != nil {
		return err
	}
	local[a.Name.Value] = val
	return nil
}

func (t *Template) posError(n ast.Node, format string, v ...any) error {
	return ast.PosError(n, t.name, format, v...)
}

func validateFunc(t reflect.Type) error {
	numOut := t.NumOut()
	if numOut > 2 {
		return ErrFuncTooManyReturns
	} else if numOut == 0 {
		return ErrFuncNoReturns
	}

	if numOut == 2 {
		if !t.Out(1).Implements(reflect.TypeOf(error(nil))) {
			return ErrFuncSecondReturnType
		}
	}

	return nil
}
