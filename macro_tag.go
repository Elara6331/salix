package salix

import (
	"errors"

	"go.elara.ws/salix/ast"
)

var (
	ErrMacroInvalidArgs = errors.New("macro expects one string argument followed by variable assignments")
	ErrNoSuchMacro      = errors.New("no such template")
)

// macroTag represents an #macro tag within a Salix template
type macroTag struct{}

func (mt macroTag) Run(tc *TagContext, block, args []ast.Node) error {
	if len(args) < 1 {
		return ErrMacroInvalidArgs
	}

	nameVal, err := tc.GetValue(args[0], nil)
	if err != nil {
		return err
	}

	name, ok := nameVal.(string)
	if !ok {
		return ErrMacroInvalidArgs
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
				return ErrMacroInvalidArgs
			}
		}

		macro, ok := tc.t.macros[name]
		if !ok {
			return ErrNoSuchMacro
		}
		return tc.Execute(macro, local)
	} else {
		tc.t.macros[name] = block
	}

	return nil
}
