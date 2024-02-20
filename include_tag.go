package salix

import (
	"errors"

	"go.elara.ws/salix/ast"
)

var (
	ErrIncludeInvalidArgs = errors.New("include expects one string argument")
	ErrNoSuchTemplate     = errors.New("no such template")
)

// includeTag represents an #include tag within a Salix template
type includeTag struct{}

func (it includeTag) Run(tc *TagContext, block, args []ast.Node) error {
	if len(args) < 1 {
		return tc.PosError(tc.Tag, "expected at least one argument, got %d", len(args))
	}

	val, err := tc.GetValue(args[0], nil)
	if err != nil {
		return err
	}

	name, ok := val.(string)
	if !ok {
		return tc.PosError(args[0], "invalid first argument type: %T (expected string)", val)
	}

	ignoreMissing := false
	if name[0] == '?' {
		name = name[1:]
		ignoreMissing = true
	}

	tmpl, ok := tc.t.ns.GetTemplate(name)
	if !ok {
		if ignoreMissing {
			return nil
		}
		return tc.PosError(args[0], "no such template: %q", name)
	}

	local := map[string]any{}

	// Use the variable assignments after the first argument
	// to set the local variables of the execution
	for _, arg := range args[1:] {
		if a, ok := arg.(ast.Assignment); ok {
			val, err := tc.GetValue(a.Value, local)
			if err != nil {
				return err
			}
			local[a.Name.Value] = val
		} else {
			// If the argument isn't an assigment, return an error
			return tc.PosError(tc.Tag, "invalid argument type: %T (expected ast.Assignment)", val)
		}
	}

	return tc.Execute(tmpl.ast, local)
}
