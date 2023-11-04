package salix

import (
	"strings"
	"testing"
)

func TestIf(t *testing.T) {
	const tmplStr = `#if(weather.Temp > 30):
    <p>It's a hot day!</p>
#elif(weather.Temp < 0):
    <p>It's freezing!</p>
#else:
    <p>The temperature is #(weather.Temp)</p>
#!if`

	type weather struct {
		Temp int
	}

	res := execStr(t, tmplStr, map[string]any{
		"weather": weather{Temp: 40},
	})
	res = strings.TrimSpace(res)
	if res != "<p>It's a hot day!</p>" {
		t.Errorf("Expected %q, got %q", "<p>It's a hot day!</p>", res)
	}

	res = execStr(t, tmplStr, map[string]any{
		"weather": weather{Temp: -1},
	})
	res = strings.TrimSpace(res)
	if res != "<p>It's freezing!</p>" {
		t.Errorf("Expected %q, got %q", "<p>It's freezing!</p>", res)
	}

	res = execStr(t, tmplStr, map[string]any{
		"weather": weather{Temp: 25},
	})
	res = strings.TrimSpace(res)
	if res != "<p>The temperature is 25</p>" {
		t.Errorf("Expected %q, got %q", "<p>The temperature is 25</p>", res)
	}
}

func TestFor(t *testing.T) {
	const tmplStr = `#for(item in items):
#(item)
#!for`

	res := execStr(t, tmplStr, map[string]any{
		"items": []any{1.2, 3, "Hello", 4 + 2i, "-2", false, nil},
	})

	expected := `1.2
3
Hello
(4+2i)
-2
false
<nil>
`

	if res != expected {
		t.Errorf("Expected %q, got %q", expected, res)
	}
}

func TestForTwoArgs(t *testing.T) {
	const tmplStr = `#for(i, item in items):
#(i) #(item)
#!for`

	res := execStr(t, tmplStr, map[string]any{
		"items": []any{1.2, 3, "Hello", 4 + 2i, "-2", false, nil},
	})

	expected := `0 1.2
1 3
2 Hello
3 (4+2i)
4 -2
5 false
6 <nil>
`

	if res != expected {
		t.Errorf("Expected %q, got %q", expected, res)
	}
}
