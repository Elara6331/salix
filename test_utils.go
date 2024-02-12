package salix

import (
	"strings"
	"testing"

	"go.elara.ws/salix/ast"
)

func testTmpl(t *testing.T) Template {
	t.Helper()

	return Template{
		ns: &Namespace{
			tmpls: map[string]Template{},
			tags:  map[string]Tag{},
			vars:  map[string]any{},
		},
		name:   t.Name(),
		tags:   map[string]Tag{},
		vars:   map[string]any{},
		macros: map[string][]ast.Node{},
	}
}

func testPos(t *testing.T) ast.Position {
	t.Helper()

	return ast.Position{
		Name: t.Name(),
		Line: -1,
		Col:  -1,
	}
}

func execStr(t *testing.T, tmplStr string, vars map[string]any) string {
	t.Helper()
	tmpl, err := New().ParseString("test", tmplStr)
	if err != nil {
		t.Fatal(err)
	}
	sb := &strings.Builder{}
	err = tmpl.WithVarMap(vars).Execute(sb)
	if err != nil {
		t.Fatal(err)
	}
	return sb.String()
}
