package salix

import "go.elara.ws/salix/ast"

// macroTag represents an #macro tag within a Salix template
type macroTag struct{}

func (mt macroTag) Run(tc *TagContext, block, args []ast.Node) error {
	if len(args) < 1 {
		return tc.PosError(tc.Tag, "expected at least one argument, got %d", len(args))
	}

	nameVal, err := tc.GetValue(args[0], nil)
	if err != nil {
		return err
	}

	name, ok := nameVal.(string)
	if !ok {
		return tc.PosError(args[0], "invalid first argument type: %T (expected string)", nameVal)
	}

	ignoreMissing := false
	if name[0] == '?' {
		name = name[1:]
		ignoreMissing = true
	}

	if len(block) == 0 {
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
				// If the argument isn't an assigment, return invalid args
				return tc.PosError(arg, "%s: invalid argument type: %T (expected ast.Assignment)", tc.NodeToString(arg), arg)
			}
		}

		macro, ok := tc.t.macros[name]
		if !ok {
			if ignoreMissing {
				return nil
			}
			return tc.PosError(tc.Tag, "no such macro: %q", name)
		}
		return tc.Execute(macro, local)
	} else {
		tc.t.macros[name] = block
	}

	return nil
}
