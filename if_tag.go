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

import "go.elara.ws/salix/ast"

// ifTag represents a #if tag within a Salix template
type ifTag struct{}

func (it ifTag) Run(tc *TagContext, block, args []ast.Node) error {
	if len(args) != 1 {
		return tc.PosError(tc.Tag, "expected one argument, got %d", len(args))
	}

	inner, err := it.findInner(tc, block)
	if err != nil {
		return err
	}

	val, err := tc.GetValue(args[0], nil)
	if err != nil {
		return err
	}

	cond, ok := val.(bool)
	if !ok {
		return tc.PosError(args[0], "expected boolean argument, got %T", val)
	}

	if cond {
		return tc.Execute(block[:inner.endRoot], nil)
	} else {
		if len(inner.elifTags) > 0 {
			for i, elifTag := range inner.elifTags {
				val, err := tc.GetValue(elifTag.value, nil)
				if err != nil {
					return err
				}

				cond, ok := val.(bool)
				if !ok {
					return tc.PosError(elifTag.value, "expected boolean argument, got %T", val)
				}

				nextIndex := len(block)
				if i < len(inner.elifTags)-1 {
					nextIndex = inner.elifTags[i+1].index
				} else if inner.elseIndex != 0 {
					nextIndex = inner.elseIndex
				}

				if cond {
					return tc.Execute(block[elifTag.index+1:nextIndex], nil)
				}
			}
		}

		if inner.elseIndex != 0 {
			return tc.Execute(block[inner.elseIndex+1:], nil)
		}
	}

	return nil
}

type innerTags struct {
	endRoot   int
	elifTags  []elif
	elseIndex int
}

type elif struct {
	index int
	value ast.Node
}

// findInner finds the inner elif and else tags in a block
// passed to the if tag.
func (it ifTag) findInner(tc *TagContext, block []ast.Node) (innerTags, error) {
	var out innerTags
	for i, node := range block {
		if tag, ok := node.(ast.Tag); ok {
			switch tag.Name.Value {
			case "elif":
				if out.endRoot == 0 {
					out.endRoot = i
				}
				if len(tag.Params) > 1 {
					return innerTags{}, tc.PosError(tag.Params[1], "expected one argument, got %d", len(tag.Params))
				}
				out.elifTags = append(out.elifTags, elif{
					index: i,
					value: tag.Params[0],
				})
			case "else":
				if out.elseIndex != 0 {
					return innerTags{}, tc.PosError(tag, "cannot have more than one else tag in an if tag")
				}
				if out.endRoot == 0 {
					out.endRoot = i
				}
				out.elseIndex = i
				break
			}
		}
	}
	if out.endRoot == 0 {
		out.endRoot = len(block)
	}
	return out, nil
}
