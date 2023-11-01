/*
 * Salix - Go templating engine
 * Copyright (C) 2023 Elara Musayelyan
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package salix

import "sync"

// Namespace represents a collection of templates that can include each other
type Namespace struct {
	mu    sync.Mutex
	tmpls map[string]Template
	vars  map[string]any
	tags  map[string]Tag

	escapeHTML *bool
}

// New returns a new template namespace
func New() *Namespace {
	return &Namespace{
		tmpls: map[string]Template{},
		vars:  map[string]any{},
		tags:  map[string]Tag{},
	}
}

// WithVarMap sets the namespace's variable map to m
func (n *Namespace) WithVarMap(m map[string]any) *Namespace {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.vars = m
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

// WithVarMap turns HTML escaping on or off for the namespace
func (n *Namespace) WithEscapeHTML(b bool) *Namespace {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.escapeHTML = &b
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
