package salix

import (
	"bytes"
	"io"
	"os"
	"strings"

	"go.elara.ws/salix/internal/ast"
	"go.elara.ws/salix/internal/parser"
)

// Parse parses a salix template from r. If the reader has a Name method
// that returns a string, that will be used as the filename.
func (t *Template) Parse(r io.Reader) (*Template, error) {
	fname := "<input>"
	if r, ok := r.(interface{ Name() string }); ok {
		fname = r.Name()
	}
	return t.ParseWithFilename(fname, r)
}

// ParseWithFilename parses a salix template from r, using the given filename.
func (t *Template) ParseWithFilename(filename string, r io.Reader) (*Template, error) {
	astVal, err := parser.ParseReader(filename, r)
	if err != nil {
		return nil, err
	}
	t.file = filename
	t.ast = astVal.([]ast.Node)
	return t, nil
}

// ParseFile parses the file at path as a salix template
func (t *Template) ParseFile(path string) (*Template, error) {
	fl, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fl.Close()
	return t.Parse(fl)
}

// ParseString parses a string using the given filename.
func (t *Template) ParseString(filename, tmpl string) (*Template, error) {
	return t.ParseWithFilename(filename, strings.NewReader(tmpl))
}

// ParseString parses bytes using the given filename.
func (t *Template) ParseBytes(filename string, tmpl []byte) (*Template, error) {
	return t.ParseWithFilename(filename, bytes.NewReader(tmpl))
}
