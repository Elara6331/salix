package salix

import (
	"io"

	"go.elara.ws/salix/internal/ast"
)

type Tag interface {
	Run(tc *TagContext, block, args []ast.Node) error
}

var defaultTags = map[string]Tag{
	"if":  ifTag{},
	"for": forTag{},
}

type TagContext struct {
	w     io.Writer
	t     *Template
	local map[string]any
}

func (tc *TagContext) Execute(nodes []ast.Node, local map[string]any) error {
	return tc.t.execute(tc.w, nodes, mergeMap(tc.local, local))
}

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
