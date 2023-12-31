package salix

import (
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"go.elara.ws/salix/ast"
	"go.elara.ws/salix/parser"
)

// NamedReader is a reader with a name
type NamedReader interface {
	io.Reader
	Name() string
}

// Parse parses a salix template from a NamedReader, which is an io.Reader
// with a Name method that returns a string. os.File implements NamedReader.
func (n *Namespace) Parse(r NamedReader) (Template, error) {
	return n.ParseWithName(r.Name(), r)
}

// ParseWithFilename parses a salix template from r, using the given name.
func (n *Namespace) ParseWithName(name string, r io.Reader) (Template, error) {
	astVal, err := parser.ParseReader(name, r, parser.GlobalStore("name", name))
	if err != nil {
		return Template{}, err
	}

	t := Template{
		ns:             n,
		name:           name,
		ast:            astVal.([]ast.Node),
		tags:           map[string]Tag{},
		vars:           map[string]any{},
		WriteOnSuccess: n.WriteOnSuccess,
	}

	if n.WhitespaceMutations {
		performWhitespaceMutations(t.ast)
	}

	n.mu.Lock()
	defer n.mu.Unlock()
	n.tmpls[name] = t
	return t, nil
}

// ParseFile parses the file at path as a salix template. It uses the path as the name.
func (t *Namespace) ParseFile(path string) (Template, error) {
	fl, err := os.Open(path)
	if err != nil {
		return Template{}, err
	}
	defer fl.Close()
	return t.Parse(fl)
}

// ParseGlob parses all the files that were matched by the given glob
// nd adds them to the namespace.
func (t *Namespace) ParseGlob(glob string) error {
	matches, err := filepath.Glob(glob)
	if err != nil {
		return err
	}

	for _, match := range matches {
		_, err := t.ParseFile(match)
		if err != nil {
			return err
		}
	}

	return nil
}

// ParseFile parses a file at the given path in a filesystem. It uses the path as the name.
func (t *Namespace) ParseFS(fsys fs.FS, path string) (Template, error) {
	fl, err := fsys.Open(path)
	if err != nil {
		return Template{}, err
	}
	defer fl.Close()
	return t.ParseWithName(path, fl)
}

// ParseGlob parses all the files in the filesystem that were matched by the given glob
// and adds them to the namespace.
func (t *Namespace) ParseFSGlob(fsys fs.FS, glob string) error {
	matches, err := fs.Glob(fsys, glob)
	if err != nil {
		return err
	}

	for _, match := range matches {
		_, err := t.ParseFS(fsys, match)
		if err != nil {
			return err
		}
	}

	return nil
}

// ParseString parses a string using the given filename.
func (t *Namespace) ParseString(filename, tmpl string) (Template, error) {
	return t.ParseWithName(filename, strings.NewReader(tmpl))
}

// ParseString parses bytes using the given filename.
func (t *Namespace) ParseBytes(filename string, tmpl []byte) (Template, error) {
	return t.ParseWithName(filename, bytes.NewReader(tmpl))
}

// performWhitespaceMutations mutates nodes in the AST to remove
// whitespace where it isn't needed.
func performWhitespaceMutations(nodes []ast.Node) {
	// lastTag keeps track of which line the
	// last tag was found on, so that if an end
	// tag was found on the same line, we know it's inline
	// and we don't need to handle its whitespace.
	lastTag := 0

	for i := 0; i < len(nodes); i++ {
		switch node := nodes[i].(type) {
		case ast.Tag:
			// If the node has no body, it's an inline tag,
			// so we don't need to handle any whitespace around it.
			if !node.HasBody {
				continue
			}
			handleWhitespace(nodes, i)
			lastTag = node.Position.Line
		case ast.EndTag:
			// If the end tag isn't on the same line as the
			// last start tag, it's not inline, so we can
			// get rid of its whitespace.
			if lastTag != node.Position.Line {
				handleWhitespace(nodes, i)
			}
		}
	}
}

// handleWhitespace mutates nodes above and below tags and end tags
// to remove the unneeded whitespace around them.
func handleWhitespace(nodes []ast.Node, i int) {
	lastIndex := len(nodes) - 1

	var prevNode ast.Text
	var nextNode ast.Text

	if i != 0 {
		if node, ok := nodes[i-1].(ast.Text); ok {
			prevNode = node
		}
	}

	if i != lastIndex {
		if node, ok := nodes[i+1].(ast.Text); ok {
			nextNode = node
		}
	}

	if prevNode.Data != nil && bytes.Contains(nextNode.Data, []byte{'\n'}) {
		prevNode.Data = trimWhitespaceSuffix(prevNode.Data)
		nodes[i-1] = prevNode
	}

	if nextNode.Data != nil {
		nextNode.Data = bytes.TrimPrefix(nextNode.Data, []byte{'\n'})
		nodes[i+1] = nextNode
	}
}

// trimWhitespaceSuffix removes everything up to and including the first newline
// it finds, starting from the end of the slice. If a non-whitespace character is
// encountered before a newline, the data is returned unmodified.
func trimWhitespaceSuffix(data []byte) []byte {
	// Start from the end of the slice
	for i := len(data) - 1; i >= 0; i-- {
		if data[i] == '\n' {
			// If a newline is found, return the slice without the newline and anything after
			return data[:i+1]
		} else if data[i] != ' ' && data[i] != '\t' && data[i] != '\r' {
			// If a non-whitespace character is found, return the original slice
			return data
		}
	}
	// If no newline or non-whitespace character is found, return the original slice
	return data
}
