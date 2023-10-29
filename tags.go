package salix

import (
	"io"

	"go.elara.ws/salix/internal/ast"
)

// Tag represents a tag in a Salix template
type Tag interface {
	Run(tc *TagContext, block, args []ast.Node) error
}

var globalTags = map[string]Tag{
	"if":      ifTag{},
	"for":     forTag{},
	"include": includeTag{},
}

// TagContext is passed to Tag implementations to allow them to control the interpreter
type TagContext struct {
	w     io.Writer
	t     *Template
	local map[string]any
}

// Execute runs the interpreter on the given AST nodes, with the given local variables.
func (tc *TagContext) Execute(nodes []ast.Node, local map[string]any) error {
	return tc.t.execute(tc.w, nodes, mergeMap(tc.local, local))
}

// GetValue evaluates the given AST node using the given local variables.
func (tc *TagContext) GetValue(node ast.Node, local map[string]any) (any, error) {
	return tc.t.getValue(node, mergeMap(tc.local, local))
}

func mergeMap(a, b map[string]any) map[string]any {
	out := map[string]any{}
	if a != nil {
		for k, v := range a {
			out[k] = v
		}
	}
	if b != nil {
		for k, v := range b {
			out[k] = v
		}
	}
	return out
}
