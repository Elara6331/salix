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
	"errors"

	"go.elara.ws/salix/ast"
)

var (
	ErrIncludeInvalidArgs = errors.New("include expects one string argument")
	ErrNoSuchTemplate     = errors.New("no such template")
)

// includeTag represents an #include tag within a Salix template
type includeTag struct{}

func (it includeTag) Run(tc *TagContext, block, args []ast.Node) error {
	if len(args) < 1 {
		return ErrIncludeInvalidArgs
	}

	val, err := tc.GetValue(args[0], nil)
	if err != nil {
		return err
	}

	name, ok := val.(string)
	if !ok {
		return ErrIncludeInvalidArgs
	}

	tmpl, ok := tc.t.ns.GetTemplate(name)
	if !ok {
		return ErrNoSuchTemplate
	}

	local := map[string]any{}

	// Use the variable assignments after the first argument
	// to set the local variables of the execution
	for _, arg := range args[1:] {
		if a, ok := arg.(ast.Assignment); ok {
			val, err := tc.GetValue(a.Value, local)
			if err != nil {
				return err
			}
			local[a.Name.Value] = val
		} else {
			// If the argument isn't an assigment, return invalid args
			return ErrIncludeInvalidArgs
		}
	}

	return tc.Execute(tmpl.ast, local)
}
