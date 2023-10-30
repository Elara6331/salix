package salix

import (
	"errors"

	"go.elara.ws/salix/internal/ast"
)

var (
	ErrIfExpectsOneArg = errors.New("if tag expects one argument")
	ErrIfExpectsBool   = errors.New("if tag expects a bool value")
	ErrIfTwoElse       = errors.New("if tags can only have one else tag inside")
)

// ifTag represents a #if tag within a Salix template
type ifTag struct{}

func (it ifTag) Run(tc *TagContext, block, args []ast.Node) error {
	if len(args) != 1 {
		return ErrIfExpectsOneArg
	}

	inner, err := it.findInner(block)
	if err != nil {
		return err
	}

	val, err := tc.GetValue(args[0], nil)
	if err != nil {
		return err
	}

	cond, ok := val.(bool)
	if !ok {
		return ErrIfExpectsBool
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
					return ErrIfExpectsBool
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
func (it ifTag) findInner(block []ast.Node) (innerTags, error) {
	var out innerTags
	for i, node := range block {
		if tag, ok := node.(ast.Tag); ok {
			switch tag.Name.Value {
			case "elif":
				if out.endRoot == 0 {
					out.endRoot = i
				}
				if len(tag.Params) > 1 {
					return innerTags{}, ErrIfExpectsOneArg
				}
				out.elifTags = append(out.elifTags, elif{
					index: i,
					value: tag.Params[0],
				})
			case "else":
				if out.elseIndex != 0 {
					return innerTags{}, ErrIfTwoElse
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
