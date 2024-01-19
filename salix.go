package salix

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"html"
	"io"
	"reflect"
	"strconv"

	"go.elara.ws/salix/ast"
)

// HTML represents unescaped HTML strings
type HTML string

// Template represents a Salix template
type Template struct {
	ns   *Namespace
	name string
	ast  []ast.Node

	escapeHTML *bool
	// WriteOnSuccess indicates whether the output should only be written if generation fully succeeds.
	// This option buffers the output of the template, so it will use more memory. (default: false)
	WriteOnSuccess bool

	tags   map[string]Tag
	vars   map[string]any
	macros map[string][]ast.Node
}

// WithVarMap returns a copy of the template with its variable map set to m.
func (t Template) WithVarMap(m map[string]any) Template {
	if m == nil {
		t.vars = map[string]any{}
	} else {
		t.vars = m
	}
	return t
}

// WithTagMap returns a copy of the template with its tag map set to m.
func (t Template) WithTagMap(m map[string]Tag) Template {
	if m == nil {
		t.tags = map[string]Tag{}
	} else {
		t.tags = m
	}
	return t
}

// WithEscapeHTML returns a copy of the template with HTML escaping enabled or disabled.
// The HTML escaping functionality is NOT context-aware.
// Using the HTML type allows you to get around the escaping if needed.
func (t Template) WithEscapeHTML(b bool) Template {
	t.escapeHTML = &b
	return t
}

// WithWriteOnSuccess enables or disables only writing if generation fully succeeds.
func (t Template) WithWriteOnSuccess(b bool) Template {
	t.WriteOnSuccess = true
	return t
}

// Execute executes a parsed template and writes
// the result to w.
func (t Template) Execute(w io.Writer) error {
	t.macros = map[string][]ast.Node{}
	if t.WriteOnSuccess {
		buf := &bytes.Buffer{}
		err := t.execute(buf, t.ast, nil)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, buf)
		return err
	} else {
		bw := bufio.NewWriterSize(w, 16384)
		defer bw.Flush()
		return t.execute(bw, t.ast, nil)
	}
}

func (t *Template) execute(w io.Writer, nodes []ast.Node, local map[string]any) error {
	if local == nil {
		local = map[string]any{}
	}

	for i := 0; i < len(nodes); i++ {
		switch node := nodes[i].(type) {
		case ast.Text:
			_, err := w.Write(node.Data)
			if err != nil {
				return ast.PosError(node, "%w", err)
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
			return ast.PosError(node, "end tag without a matching start tag: %s", node.Name.Value)
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

func (t *Template) getEscapeHTML() bool {
	if t.escapeHTML != nil {
		return *t.escapeHTML
	} else if t.ns.escapeHTML != nil {
		return *t.ns.getEscapeHTML()
	} else {
		return false
	}
}

func (t *Template) toString(v any) string {
	if h, ok := v.(HTML); ok {
		return string(h)
	} else if t.getEscapeHTML() {
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
		return val, nil
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

// valueToString converts an AST node to a textual representation
// for the user to see, such as in error messages. This does not
// directly correlate to Salix source code.
func valueToString(node ast.Node) string {
	if node == nil {
		return "<nil>"
	}

	switch node := node.(type) {
	case ast.Ident:
		return node.Value
	case ast.String:
		return strconv.Quote(node.Value)
	case ast.Integer:
		return strconv.FormatInt(node.Value, 10)
	case ast.Float:
		return strconv.FormatFloat(node.Value, 'g', -1, 64)
	case ast.Bool:
		return strconv.FormatBool(node.Value)
	case ast.Assignment:
		return node.Name.Value + " = " + valueToString(node.Value)
	case ast.Index:
		return valueToString(node.Value) + "[" + valueToString(node.Index) + "]"
	case ast.Ternary:
		return valueToString(node.Condition) + " ? " + valueToString(node.IfTrue) + " : " + valueToString(node.Else)
	case ast.FieldAccess:
		return valueToString(node.Value) + "." + node.Name.Value
	case ast.Value:
		if node.Not {
			return "!" + valueToString(node.Node)
		}
		return valueToString(node.Node)
	case ast.FuncCall:
		if len(node.Params) > 1 {
			return node.Name.Value + "(" + valueToString(node.Params[0]) + ", ...)"
		} else if len(node.Params) == 1 {
			return node.Name.Value + "(" + valueToString(node.Params[0]) + ")"
		} else {
			return node.Name.Value + "()"
		}
	case ast.MethodCall:
		if len(node.Params) > 1 {
			return valueToString(node.Value) + "." + node.Name.Value + "(" + valueToString(node.Params[0]) + ", ...)"
		} else if len(node.Params) == 1 {
			return valueToString(node.Value) + "." + node.Name.Value + "(" + valueToString(node.Params[0]) + ")"
		} else {
			return valueToString(node.Value) + "." + node.Name.Value + "()"
		}
	case ast.Expr:
		if len(node.Rest) == 0 {
			return valueToString(node.First)
		}
		return valueToString(node.First) + " " + node.Rest[0].Operator.Value + " " + valueToString(node.Rest[0])
	case ast.Tag:
		if len(node.Params) > 1 {
			return "#" + node.Name.Value + "(" + valueToString(node.Params[0]) + ", ...)"
		} else if len(node.Params) == 1 {
			return "#" + node.Name.Value + "(" + valueToString(node.Params[0]) + ")"
		} else {
			return "#" + node.Name.Value + "()"
		}
	case ast.EndTag:
		return "#!" + node.Name.Value
	case ast.ExprTag:
		return "#(" + valueToString(node.Value) + ")"
	default:
		return "..."
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
			return nil, ast.PosError(node, "%s: the ! operator can only be used on boolean values", valueToString(node))
		}
		return !rval.Bool(), nil
	}

	return v, err
}

// getVar tries to get a variable from the local map. If it's not found,
// it'll try the global variable map. If it doesn't exist in either map,
// it will return an error.
func (t *Template) getVar(id ast.Ident, local map[string]any) (any, error) {
	if local != nil {
		v, ok := local[id.Value]
		if ok {
			return v, nil
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

	return reflect.Value{}, ast.PosError(id, "no such variable: %s", id.Value)
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
		return 0, ast.PosError(node, "no such tag: %s", node.Name.Value)
	}

	var block []ast.Node
	if node.HasBody {
		block = t.getBlock(nodes, i+1, node.Position.Line, node.Name.Value)
		i += len(block) + 1
	}

	tc := &TagContext{node, w, t, local}

	err = tag.Run(tc, block, node.Params)
	if err != nil {
		return 0, errors.Join(ast.PosError(node, "%s ->", valueToString(node)), err)
	}

	return i, nil
}

// execFuncCall executes a function call
func (t *Template) execFuncCall(fc ast.FuncCall, local map[string]any) (any, error) {
	fn, err := t.getVar(fc.Name, local)
	if err != nil {
		return nil, ast.PosError(fc, "no such function: %s", fc.Name.Value)
	}
	return t.execFunc(reflect.ValueOf(fn), fc, fc.Params, local)
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
	case reflect.Slice, reflect.Array, reflect.String:
		intType := reflect.TypeOf(0)
		if rindex.CanConvert(intType) {
			rindex = rindex.Convert(intType)
		} else {
			return nil, ast.PosError(i, "%s: invalid index type: %T", valueToString(i), index)
		}

		intIndex := rindex.Interface().(int)
		if intIndex < rval.Len() {
			return rval.Index(intIndex).Interface(), nil
		} else {
			return nil, ast.PosError(i, "%s: index out of range: %d", valueToString(i), intIndex)
		}
	case reflect.Map:
		if rindex.CanConvert(rval.Type().Key()) {
			rindex = rindex.Convert(rval.Type().Key())
		} else {
			return nil, ast.PosError(i, "%s: invalid map index type: %T (expected %s)", valueToString(i), index, rval.Type().Key())
		}
		if out := rval.MapIndex(rindex); out.IsValid() {
			return out.Interface(), nil
		} else {
			return nil, ast.PosError(i, "%s: map index not found: %q", valueToString(i), index)
		}
	default:
		return nil, ast.PosError(i, "%s: cannot index type: %T", valueToString(i), val)
	}
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
		return nil, ast.PosError(fa, "%s: no such field: %s", valueToString(fa), fa.Name.Value)
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
	// First, check for a method with the given name
	mtd := rval.MethodByName(mc.Name.Value)
	if mtd.IsValid() {
		return t.execFunc(mtd, mc, mc.Params, local)
	}
	// If the method doesn't exist, also check for a field storing a function.
	field := rval.FieldByName(mc.Name.Value)
	if field.IsValid() && field.Kind() == reflect.Func {
		return t.execFunc(field, mc, mc.Params, local)
	}
	// If neither of those exist, return an error
	return nil, ast.PosError(mc, "no such method: %s", mc.Name.Value)
}

// execFunc executes a function call
func (t *Template) execFunc(fn reflect.Value, node ast.Node, args []ast.Node, local map[string]any) (any, error) {
	fnType := fn.Type()
	lastIndex := fnType.NumIn() - 1
	isVariadic := fnType.IsVariadic()
		
	if !isVariadic && fnType.NumIn() != len(args) {
		return nil, ast.PosError(node, "%s: invalid parameter amount: %d (expected %d)", valueToString(node), len(args), fnType.NumIn())
	}

	if err := validateFunc(fnType, node); err != nil {
		return nil, err
	}

	params := make([]reflect.Value, 0, fnType.NumIn())
	for i, arg := range args {
		if _, ok := arg.(ast.Assignment); ok {
			return nil, ast.PosError(arg, "%s: an assignment cannot be used as a function argument", valueToString(node))
		}
		paramVal, err := t.getValue(arg, local)
		if err != nil {
			return nil, err
		}
		params = append(params, reflect.ValueOf(paramVal))
		if isVariadic && i >= lastIndex {
			if params[i].CanConvert(fnType.In(lastIndex).Elem()) {
				params[i] = params[i].Convert(fnType.In(lastIndex).Elem())
			} else {
				return nil, ast.PosError(node, "%s: invalid parameter type: %T (expected %s)", valueToString(node), paramVal, fnType.In(i))
			}
		} else {
			if params[i].CanConvert(fnType.In(i)) {
				params[i] = params[i].Convert(fnType.In(i))
			} else {
				return nil, ast.PosError(node, "%s: invalid parameter type: %T (expected %s)", valueToString(node), paramVal, fnType.In(i))
			}
		}
	}

	ret := fn.Call(params)
	if len(ret) == 1 {
		retv := ret[0].Interface()
		if err, ok := retv.(error); ok {
			return nil, ast.PosError(node, "%s: %w", valueToString(node), err)
		}
		return ret[0].Interface(), nil
	} else {
		if ret[1].IsNil() {
			return ret[0].Interface(), nil
		}
		return ret[0].Interface(), ast.PosError(node, "%s: %w", valueToString(node), ret[1].Interface().(error))
	}
}

func (t *Template) evalTernary(tr ast.Ternary, local map[string]any) (any, error) {
	condVal, err := t.getValue(tr.Condition, local)
	if err != nil {
		return nil, err
	}

	cond, ok := condVal.(bool)
	if !ok {
		return nil, ast.PosError(tr.Condition, "%s: ternary condition must be a boolean value", valueToString(tr.Condition))
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
	return val, nil
}

func (t *Template) handleAssignment(a ast.Assignment, local map[string]any) error {
	val, err := t.getValue(a.Value, local)
	if err != nil {
		return err
	}
	local[a.Name.Value] = val
	return nil
}

func validateFunc(t reflect.Type, node ast.Node) error {
	numOut := t.NumOut()
	if numOut > 2 {
		return ast.PosError(node, "template functions cannot have more than two return values")
	} else if numOut == 0 {
		return ast.PosError(node, "template functions must have at least one return value")
	}
	if numOut == 2 {
		errType := reflect.TypeOf((*error)(nil)).Elem()
		if !t.Out(1).Implements(errType) {
			return ast.PosError(node, "the second return value of a template function must be an error")
		}
	}

	return nil
}
