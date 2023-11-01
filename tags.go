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

import (
	"bytes"
	"io"

	"go.elara.ws/salix/ast"
)

// Tag represents a tag in a Salix template
type Tag interface {
	Run(tc *TagContext, block, args []ast.Node) error
}

var globalTags = map[string]Tag{
	"if":      ifTag{},
	"for":     forTag{},
	"include": includeTag{},
	"macro":   macroTag{},
}

// TagContext is passed to Tag implementations to allow them to control the interpreter
type TagContext struct {
	w     io.Writer
	t     *Template
	local map[string]any
}

// Execute runs the interpreter on the given AST nodes, with the given local variables.
func (tc *TagContext) Execute(nodes []ast.Node, local map[string]any) error {
	return tc.t.execute(tc.w, nodes, mergeMap(tc.local, local))
}

// ExecuteToMemory runs the interpreter on the given AST nodes, with the given local variables, and
// returns the resulting bytes rather than writing them out.
func (tc *TagContext) ExecuteToMemory(nodes []ast.Node, local map[string]any) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := tc.t.execute(buf, nodes, mergeMap(tc.local, local))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GetValue evaluates the given AST node using the given local variables.
func (tc *TagContext) GetValue(node ast.Node, local map[string]any) (any, error) {
	return tc.t.getValue(node, mergeMap(tc.local, local))
}

// Write writes b to the underlying writer. It implements
// the io.Writer interface.
func (tc *TagContext) Write(b []byte) (int, error) {
	return tc.w.Write(b)
}

func mergeMap(a, b map[string]any) map[string]any {
	out := map[string]any{}
	if a != nil {
		for k, v := range a {
			out[k] = v
		}
	}
	if b != nil {
		for k, v := range b {
			out[k] = v
		}
	}
	return out
}
