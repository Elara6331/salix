package salix

import (
	"errors"

	"go.elara.ws/salix/internal/ast"
)

var (
	ErrIncludeInvalidArgs = errors.New("include expects one string argument")
	ErrNoSuchTemplate     = errors.New("no such template")
)

// forTag represents an #include tag within a Salix template
type includeTag struct{}

func (it includeTag) Run(tc *TagContext, block, args []ast.Node) error {
	if len(args) != 1 {
		return ErrIncludeInvalidArgs
	}

	val, err := tc.GetValue(args[0], nil)
	if err != nil {
		return err
	}

	name, ok := val.(string)
	if !ok {
		return ErrIncludeInvalidArgs
	}

	tmpl, ok := tc.t.ns.GetTemplate(name)
	if !ok {
		return ErrNoSuchTemplate
	}

	return tc.Execute(tmpl.ast, nil)
}
