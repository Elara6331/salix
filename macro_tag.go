package salix

import (
	"errors"

	"go.elara.ws/salix/ast"
)

var (
	ErrMacroInvalidArgs = errors.New("macro expects one string argument")
	ErrNoSuchMacro      = errors.New("no such template")
)

// macroTag represents an #macro tag within a Salix template
type macroTag struct{}

func (mt macroTag) Run(tc *TagContext, block, args []ast.Node) error {
	if len(args) != 1 {
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
		tc.t.ns.mu.Lock()
		macro, ok := tc.t.ns.macros[name]
		if !ok {
			return ErrNoSuchMacro
		}
		tc.t.ns.mu.Unlock()
		return tc.Execute(macro, nil)
	} else {
		tc.t.ns.mu.Lock()
		tc.t.ns.macros[name] = block
		tc.t.ns.mu.Unlock()
	}

	return nil
}
