package salix

import (
	"errors"
	"strings"
	"testing"

	"go.elara.ws/salix/ast"
)

func TestFuncCall(t *testing.T) {
	const expected = 1 + 2i
	fn := func() complex128 { return expected }

	// test()
	ast := ast.FuncCall{
		Name:     ast.Ident{Value: "test", Position: testPos(t)},
		Position: testPos(t),
	}

	tmpl := testTmpl(t)
	val, err := tmpl.execFuncCall(ast, map[string]any{"test": fn})
	if err != nil {
		t.Fatalf("execFuncCall error: %s", err)
	}

	if val.(complex128) != expected {
		t.Errorf("Expected %v, got %v", expected, val)
	}
}

func TestFuncCallVariadic(t *testing.T) {
	const expected = "Hello, World"
	concat := func(params ...string) string { return strings.Join(params, ", ") }

	// concat("Hello", "World")
	ast := ast.FuncCall{
		Name: ast.Ident{Value: "concat", Position: testPos(t)},
		Params: []ast.Node{
			ast.String{Value: "Hello", Position: testPos(t)},
			ast.String{Value: "World", Position: testPos(t)},
		},
		Position: testPos(t),
	}

	tmpl := testTmpl(t)
	val, err := tmpl.execFuncCall(ast, map[string]any{"concat": concat})
	if err != nil {
		t.Fatalf("execFuncCall error: %s", err)
	}

	if val.(string) != expected {
		t.Errorf("Expected %q, got %q", expected, val)
	}
}

func TestFuncCallError(t *testing.T) {
	var expectedErr = errors.New("expected error")
	fn := func() error { return expectedErr }

	// test()
	ast := ast.FuncCall{
		Name:     ast.Ident{Value: "test", Position: testPos(t)},
		Position: testPos(t),
	}

	tmpl := testTmpl(t)
	_, err := tmpl.execFuncCall(ast, map[string]any{"test": fn})
	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected %q, got %q", expectedErr, err)
	}
}

func TestFuncCallMultiReturn(t *testing.T) {
	const expected = "test"
	fn := func() (string, error) { return expected, nil }

	// test()
	ast := ast.FuncCall{
		Name:     ast.Ident{Value: "test", Position: testPos(t)},
		Position: testPos(t),
	}

	tmpl := testTmpl(t)
	val, err := tmpl.execFuncCall(ast, map[string]any{"test": fn})
	if err != nil {
		t.Fatalf("execFuncCall error: %s", err)
	}

	if val.(string) != expected {
		t.Errorf("Expected %q, got %q", expected, val)
	}
}

func TestFuncCallMultiReturnError(t *testing.T) {
	var expectedErr = errors.New("expected error")
	fn := func() (string, error) { return "", expectedErr }

	// test()
	ast := ast.FuncCall{
		Name:     ast.Ident{Value: "test", Position: testPos(t)},
		Position: testPos(t),
	}

	tmpl := testTmpl(t)
	_, err := tmpl.execFuncCall(ast, map[string]any{"test": fn})
	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected %q, got %q", expectedErr, err)
	}
}

func TestFuncCallNil(t *testing.T) {
	// test()
	ast := ast.FuncCall{
		Name:     ast.Ident{Value: "test", Position: testPos(t)},
		Position: testPos(t),
	}

	tmpl := testTmpl(t)
	_, err := tmpl.execFuncCall(ast, map[string]any{"test": nil})
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
