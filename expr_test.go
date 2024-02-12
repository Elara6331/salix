package salix

import (
	"testing"
)

func TestAdd(t *testing.T) {
	res := execStr(t, `#(3 + 1)`, nil)
	if res != "4" {
		t.Errorf("Expected 4, got %s", res)
	}
}

func TestSub(t *testing.T) {
	res := execStr(t, `#(3 - 1)`, nil)
	if res != "2" {
		t.Errorf("Expected 2, got %s", res)
	}
}

func TestMul(t *testing.T) {
	res := execStr(t, `#(3 * 2)`, nil)
	if res != "6" {
		t.Errorf("Expected 6, got %s", res)
	}
}

func TestDiv(t *testing.T) {
	res := execStr(t, `#(8 / 4)`, nil)
	if res != "2" {
		t.Errorf("Expected 2, got %s", res)
	}
}

func TestMod(t *testing.T) {
	res := execStr(t, `#(4 % 4)`, nil)
	if res != "0" {
		t.Errorf("Expected 0, got %s", res)
	}
}

func TestEq(t *testing.T) {
	res := execStr(t, `#("x" == "y")`, nil)
	if res != "false" {
		t.Errorf("Expected false, got %s", res)
	}
}

func TestGEq(t *testing.T) {
	res := execStr(t, `#(2 >= 2)`, nil)
	if res != "true" {
		t.Errorf("Expected true, got %s", res)
	}
}

func TestGt(t *testing.T) {
	res := execStr(t, `#(len("hi") > 2)`, nil)
	if res != "false" {
		t.Errorf("Expected false, got %s", res)
	}
}

func TestLEq(t *testing.T) {
	res := execStr(t, `#(4 <= 4)`, nil)
	if res != "true" {
		t.Errorf("Expected true, got %s", res)
	}
}

func TestLt(t *testing.T) {
	res := execStr(t, `#(4 < 4)`, nil)
	if res != "false" {
		t.Errorf("Expected false, got %s", res)
	}
}

func TestAnd(t *testing.T) {
	res := execStr(t, `#(true && false)`, nil)
	if res != "false" {
		t.Errorf("Expected false, got %s", res)
	}
}

func TestOr(t *testing.T) {
	res := execStr(t, `#(true || false)`, nil)
	if res != "true" {
		t.Errorf("Expected true, got %s", res)
	}
}

func TestInString(t *testing.T) {
	res := execStr(t, `#("h" in "hello")`, nil)
	if res != "true" {
		t.Errorf("Expected true, got %s", res)
	}
}

func TestInSlice(t *testing.T) {
	res := execStr(t, `#(5 in slice) #(6 in slice)`, map[string]any{"slice": []int{1, 2, 3, 4, 5}})
	if res != "true false" {
		t.Errorf("Expected %q, got %q", "true false", res)
	}
}

func TestInMap(t *testing.T) {
	res := execStr(t, `#(3.4234 in map)`, map[string]any{"map": map[float32]uint{3.4234: 0}})
	if res != "true" {
		t.Errorf("Expected %q, got %q", "true", res)
	}
}

func TestParenExpr(t *testing.T) {
	res := execStr(t, `#(5 - 4.0 - 3 - 2) #(5 - (4.0 - 3) - 2)`, nil)
	if res != "-4 2" {
		t.Errorf("Expected %q, got %q", "4 -2", res)
	}
}

func TestCoalescing(t *testing.T) {
	res := execStr(t, `#(hello | "nothing") #(x | "nothing")`, map[string]any{"hello": "world"})
	if res != "world nothing" {
		t.Errorf("Expected %q, got %q", "world nothing", res)
	}
}

func TestTernary(t *testing.T) {
	res := execStr(t, `#(2.0 == 2.0 ? "equal" : "non-equal") #(2.0 == 2.5 ? "equal" : "non-equal")`, nil)
	if res != "equal non-equal" {
		t.Errorf("Expected %q, got %q", "equal non-equal", res)
	}
}

func TestNot(t *testing.T) {
	res := execStr(t, `#(!true)`, nil)
	if res != "false" {
		t.Errorf("Expected %q, got %q", "false", res)
	}
}

func TestAssignment(t *testing.T) {
	res := execStr(t, `#(x = 4)#(x)`, nil)
	if res != "4" {
		t.Errorf("Expected %q, got %q", "4", res)
	}
}
