package salix

import (
	"fmt"
	"testing"
	"time"

	"go.elara.ws/salix/ast"
)

func TestSliceGetIndex(t *testing.T) {
	testSlice := []any{1, "2", 3.0}

	tmpl := testTmpl(t)
	for index, expected := range testSlice {
		t.Run(fmt.Sprint(index), func(t *testing.T) {
			// test[index]
			ast := ast.Index{
				Value:    ast.Ident{Value: "test", Position: testPos(t)},
				Index:    ast.Ident{Value: "index", Position: testPos(t)},
				Position: testPos(t),
			}

			val, err := tmpl.getIndex(ast, map[string]any{"test": testSlice, "index": index})
			if err != nil {
				t.Fatalf("getIndex error: %s", err)
			}

			if val != expected {
				t.Errorf("Expected %v, got %v", expected, val)
			}
		})
	}
}

func TestSliceGetIndexOutOfRange(t *testing.T) {
	testSlice := []any{}
	tmpl := testTmpl(t)

	// test[0.0]
	ast := ast.Index{
		Value:    ast.Ident{Value: "test", Position: testPos(t)},
		Index:    ast.Float{Value: 0, Position: testPos(t)},
		Position: testPos(t),
	}

	_, err := tmpl.getIndex(ast, map[string]any{"test": testSlice})
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestSliceGetIndexInvalidType(t *testing.T) {
	testSlice := []any{}
	tmpl := testTmpl(t)

	// test[0.0]
	ast := ast.Index{
		Value:    ast.Ident{Value: "test", Position: testPos(t)},
		Index:    ast.String{Value: "key", Position: testPos(t)},
		Position: testPos(t),
	}

	_, err := tmpl.getIndex(ast, map[string]any{"test": testSlice})
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestMapGetIndex(t *testing.T) {
	testMap := map[any]any{1: "2", 3.0: 4, "5": 6.0}

	tmpl := testTmpl(t)
	for index, expected := range testMap {
		t.Run(fmt.Sprint(index), func(t *testing.T) {
			// test[index]
			ast := ast.Index{
				Value:    ast.Ident{Value: "test", Position: testPos(t)},
				Index:    ast.Ident{Value: "index", Position: testPos(t)},
				Position: testPos(t),
			}

			val, err := tmpl.getIndex(ast, map[string]any{"test": testMap, "index": index})
			if err != nil {
				t.Fatalf("getIndex error: %s", err)
			}

			if val != expected {
				t.Errorf("Expected %v, got %v", expected, val)
			}
		})
	}
}

func TestMapGetIndexNotFound(t *testing.T) {
	testMap := map[string]any{}
	tmpl := testTmpl(t)

	// test["key"]
	ast := ast.Index{
		Value:    ast.Ident{Value: "test", Position: testPos(t)},
		Index:    ast.String{Value: "key", Position: testPos(t)},
		Position: testPos(t),
	}

	_, err := tmpl.getIndex(ast, map[string]any{"test": testMap})
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestMapGetIndexInvalidType(t *testing.T) {
	testMap := map[int]any{}
	tmpl := testTmpl(t)

	// test["key"]
	ast := ast.Index{
		Value:    ast.Ident{Value: "test", Position: testPos(t)},
		Index:    ast.String{Value: "key", Position: testPos(t)},
		Position: testPos(t),
	}

	_, err := tmpl.getIndex(ast, map[string]any{"test": testMap})
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestGetIndexNil(t *testing.T) {
	tmpl := testTmpl(t)

	// test["key"]
	ast := ast.Index{
		Value:    ast.Ident{Value: "test", Position: testPos(t)},
		Index:    ast.String{Value: "key", Position: testPos(t)},
		Position: testPos(t),
	}

	_, err := tmpl.getIndex(ast, map[string]any{"test": nil})
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestGetIndexNilIndex(t *testing.T) {
	testMap := map[string]any{}
	tmpl := testTmpl(t)

	// test[index]
	ast := ast.Index{
		Value:    ast.Ident{Value: "test", Position: testPos(t)},
		Index:    ast.Ident{Value: "index", Position: testPos(t)},
		Position: testPos(t),
	}

	_, err := tmpl.getIndex(ast, map[string]any{"test": testMap, "index": nil})
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestGetIndexInvalidContainer(t *testing.T) {
	testStruct := struct{}{}
	tmpl := testTmpl(t)

	// test["key"]
	ast := ast.Index{
		Value:    ast.Ident{Value: "test", Position: testPos(t)},
		Index:    ast.String{Value: "key", Position: testPos(t)},
		Position: testPos(t),
	}

	_, err := tmpl.getIndex(ast, map[string]any{"test": testStruct})
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestGetField(t *testing.T) {
	testStruct := &struct {
		X int
		Y string
		Z struct{ A int }
	}{
		X: 1,
		Y: "2",
		Z: struct{ A int }{A: 1},
	}

	testCases := map[string]any{
		"X": testStruct.X,
		"Y": testStruct.Y,
		"Z": testStruct.Z,
	}

	tmpl := testTmpl(t)
	for fieldName, expected := range testCases {
		t.Run(fieldName, func(t *testing.T) {
			// test.Field
			ast := ast.FieldAccess{
				Value:    ast.Ident{Value: "test", Position: testPos(t)},
				Name:     ast.Ident{Value: fieldName, Position: testPos(t)},
				Position: testPos(t),
			}

			val, err := tmpl.getField(ast, map[string]any{"test": testStruct})
			if err != nil {
				t.Fatalf("getField error: %s", err)
			}

			if val != expected {
				t.Errorf("Expected %v, got %v", expected, val)
			}
		})
	}
}

func TestGetFieldNil(t *testing.T) {
	tmpl := testTmpl(t)

	// test.Field
	ast := ast.FieldAccess{
		Value:    ast.Ident{Value: "test", Position: testPos(t)},
		Name:     ast.Ident{Value: "Field", Position: testPos(t)},
		Position: testPos(t),
	}

	_, err := tmpl.getField(ast, map[string]any{"test": nil})
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestGetFieldNoFields(t *testing.T) {
	testStruct := struct{}{}
	tmpl := testTmpl(t)

	// test.Field
	ast := ast.FieldAccess{
		Value:    ast.Ident{Value: "test", Position: testPos(t)},
		Name:     ast.Ident{Value: "Field", Position: testPos(t)},
		Position: testPos(t),
	}

	_, err := tmpl.getField(ast, map[string]any{"test": testStruct})
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestGetFieldNotFound(t *testing.T) {
	testStruct := struct{ X int }{0}
	tmpl := testTmpl(t)

	// test.Field
	ast := ast.FieldAccess{
		Value:    ast.Ident{Value: "test", Position: testPos(t)},
		Name:     ast.Ident{Value: "Field", Position: testPos(t)},
		Position: testPos(t),
	}

	_, err := tmpl.getField(ast, map[string]any{"test": testStruct})
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestMethodCall(t *testing.T) {
	executed := false
	testStruct := &struct{ Method func(bool) bool }{
		Method: func(b bool) bool {
			executed = b
			return executed
		},
	}

	// test.Method(true)
	ast := ast.MethodCall{
		Value: ast.Ident{Value: "test", Position: testPos(t)},
		Name:  ast.Ident{Value: "Method", Position: testPos(t)},
		Params: []ast.Node{
			ast.Bool{Value: true, Position: testPos(t)},
		},
		Position: testPos(t),
	}

	tmpl := testTmpl(t)
	val, err := tmpl.execMethodCall(ast, map[string]any{"test": testStruct})
	if err != nil {
		t.Fatalf("execMethodCall error: %s", err)
	}

	if !executed || !val.(bool) {
		t.Error("Expected method to execute, got false")
	}
}

func TestMethodCallTime(t *testing.T) {
	const expected = "2023-04-26T00:00:00Z"
	testStruct := time.Date(2023, time.April, 26, 0, 0, 0, 0, time.UTC)
	tmpl := testTmpl(t)

	// t.Format("2006-01-02T15:04:05Z07:00)
	ast := ast.MethodCall{
		Value: ast.Ident{Value: "t", Position: testPos(t)},
		Name:  ast.Ident{Value: "Format", Position: testPos(t)},
		Params: []ast.Node{
			ast.String{Value: time.RFC3339, Position: testPos(t)},
		},
		Position: testPos(t),
	}

	val, err := tmpl.execMethodCall(ast, map[string]any{"t": testStruct})
	if err != nil {
		t.Fatalf("execMethodCall error: %s", err)
	}

	if val != expected {
		t.Errorf("Expected %q, got %q", expected, val)
	}
}

func TestMethodCallNotFound(t *testing.T) {
	tmpl := testTmpl(t)

	// test.Method(true)
	ast := ast.MethodCall{
		Value: ast.Ident{Value: "test", Position: testPos(t)},
		Name:  ast.Ident{Value: "Method", Position: testPos(t)},
		Params: []ast.Node{
			ast.Bool{Value: true, Position: testPos(t)},
		},
		Position: testPos(t),
	}

	_, err := tmpl.execMethodCall(ast, map[string]any{"test": struct{}{}})
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestMethodCallNil(t *testing.T) {
	tmpl := testTmpl(t)

	// test.Method(true)
	ast := ast.MethodCall{
		Value: ast.Ident{Value: "test", Position: testPos(t)},
		Name:  ast.Ident{Value: "Method", Position: testPos(t)},
		Params: []ast.Node{
			ast.Bool{Value: true, Position: testPos(t)},
		},
		Position: testPos(t),
	}

	_, err := tmpl.execMethodCall(ast, map[string]any{"test": nil})
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
