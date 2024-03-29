package salix

import (
	"fmt"
	"io"
	"sync"
)

// Namespace represents a collection of templates that can include each other
type Namespace struct {
	mu    sync.Mutex
	tmpls map[string]Template
	vars  map[string]any
	tags  map[string]Tag

	// WhitespaceMutations enables postprocessing to remove whitespace where it isn't needed
	// to make the resulting document look better. Postprocessing is only done once when the
	// template is parsed, so it will not affect performance. (default: true)
	WhitespaceMutations bool
	// WriteOnSuccess indicates whether the output should only be written if generation fully succeeds.
	// This option buffers the output of the template, so it will use more memory. (default: false)
	WriteOnSuccess bool
	// NilToZero indictes whether nil pointer values should be converted to zero values of their underlying
	// types.
	NilToZero bool
	escapeHTML     *bool
}

// New returns a new template namespace
func New() *Namespace {
	return &Namespace{
		tmpls:               map[string]Template{},
		vars:                map[string]any{},
		tags:                map[string]Tag{},
		WhitespaceMutations: true,
		WriteOnSuccess:      false,
	}
}

// WithVarMap sets the namespace's variable map to m
func (n *Namespace) WithVarMap(m map[string]any) *Namespace {
	n.mu.Lock()
	defer n.mu.Unlock()

	if m == nil {
		n.vars = map[string]any{}
	} else {
		n.vars = m
	}

	return n
}

// WithTagMap sets the namespace's tag map to m
func (n *Namespace) WithTagMap(m map[string]Tag) *Namespace {
	n.mu.Lock()
	defer n.mu.Unlock()

	if m == nil {
		n.tags = map[string]Tag{}
	} else {
		n.tags = m
	}

	return n
}

// WithEscapeHTML turns HTML escaping on or off for the namespace
func (n *Namespace) WithEscapeHTML(b bool) *Namespace {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.escapeHTML = &b
	return n
}

// WithWriteOnSuccess enables or disables only writing if generation fully succeeds.
func (n *Namespace) WithWriteOnSuccess(b bool) *Namespace {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.WriteOnSuccess = true
	return n
}

// WithWhitespaceMutations turns whitespace mutations on or off for the namespace
func (n *Namespace) WithWhitespaceMutations(b bool) *Namespace {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.WhitespaceMutations = true
	return n
}

// WithNilToZero enables or disables conversion of nil values to zero values for the namespace
func (n *Namespace) WithNilToZero(b bool) *Namespace {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.NilToZero = true
	return n
}

// GetTemplate tries to get a template from the namespace's template map.
// If it finds the template, it returns the template and true. If it
// doesn't find it, it returns nil and false.
func (n *Namespace) GetTemplate(name string) (Template, bool) {
	n.mu.Lock()
	defer n.mu.Unlock()
	t, ok := n.tmpls[name]
	return t, ok
}

// MustGetTemplate is the same as GetTemplate but it panics if the template
// doesn't exist in the namespace.
func (n *Namespace) MustGetTemplate(name string) Template {
	tmpl, ok := n.GetTemplate(name)
	if !ok {
		panic(fmt.Errorf("no such template: %q", name))
	}
	return tmpl
}

// ExecuteTemplate gets and executes a template with the given name.
func (n *Namespace) ExecuteTemplate(w io.Writer, name string, vars map[string]any) error {
	tmpl, ok := n.GetTemplate(name)
	if !ok {
		return fmt.Errorf("no such template: %q", name)
	}
	return tmpl.WithVarMap(vars).Execute(w)
}

// getVar tries to get a variable from the namespace's variable map
func (n *Namespace) getVar(name string) (any, bool) {
	n.mu.Lock()
	defer n.mu.Unlock()
	v, ok := n.vars[name]
	return v, ok
}

// getTag tries to get a tag from the namespace's tag map
func (n *Namespace) getTag(name string) (Tag, bool) {
	n.mu.Lock()
	defer n.mu.Unlock()
	t, ok := n.tags[name]
	return t, ok
}

// getEscapeHTML returns the namespace's escapeHTML value
func (n *Namespace) getEscapeHTML() *bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.escapeHTML
}
