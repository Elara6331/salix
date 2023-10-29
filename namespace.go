package salix

import (
	"reflect"
	"sync"
)

// Namespace represents a collection of templates that can include each other
type Namespace struct {
	mu    sync.Mutex
	tmpls map[string]*Template
	vars  map[string]reflect.Value
	tags  map[string]Tag
}

// New returns a new template namespace
func New() *Namespace {
	return &Namespace{
		tmpls: map[string]*Template{},
		vars:  map[string]reflect.Value{},
		tags:  map[string]Tag{},
	}
}

// WithVarMap sets the namespace's variable map to m
func (n *Namespace) WithVarMap(m map[string]any) *Namespace {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.vars = map[string]reflect.Value{}
	if m != nil {
		for k, v := range m {
			n.vars[k] = reflect.ValueOf(v)
		}
	}

	return n
}

// WithTagMap sets the namespace's tag map to m
func (n *Namespace) WithTagMap(m map[string]Tag) *Namespace {
	n.mu.Lock()
	defer n.mu.Unlock()

	if m != nil {
		n.tags = m
	} else {
		n.tags = map[string]Tag{}
	}

	return n
}

// GetTemplate tries to get a template from the namespace's template map.
// If it finds the template, it returns the template and true. If it
// doesn't find it, it returns nil and false.
func (n *Namespace) GetTemplate(name string) (*Template, bool) {
	n.mu.Lock()
	defer n.mu.Unlock()
	t, ok := n.tmpls[name]
	return t, ok
}

// getVar tries to get a variable from the namespace's variable map
func (n *Namespace) getVar(name string) (reflect.Value, bool) {
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
