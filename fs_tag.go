package salix

import (
	"io"
	"io/fs"
	"path"

	"go.elara.ws/salix/ast"
)

// FSTag writes files from an fs.FS to a template
//
// No escaping is done on the files, so make sure to avoid user-generated data.
type FSTag struct {
	// FS is the filesystem that files will be loaded from.
	FS fs.FS

	// PathPrefix is joined to the path string before a file is read.
	PathPrefix string

	// Extension is appended to the end of the path string before a file is read.
	Extension string
}

func (ft FSTag) Run(tc *TagContext, block, args []ast.Node) error {
	if len(args) != 1 {
		return tc.PosError(tc.Tag, "expected one argument, got %d", len(args))
	}

	pathVal, err := tc.GetValue(args[0], nil)
	if err != nil {
		return err
	}

	pathStr, ok := pathVal.(string)
	if !ok {
		return tc.PosError(args[0], "expected string argument, got %T", pathVal)
	}

	fl, err := ft.FS.Open(path.Join(ft.PathPrefix, pathStr) + ft.Extension)
	if err != nil {
		return err
	}
	defer fl.Close()

	_, err = io.Copy(tc, fl)
	return err
}
