package utils

import (
	"io/fs"
	"os"
)

func HasOnlySubDirectory(parent, name string) (bool, error) {
	var has, hasAnother bool
	err := fs.WalkDir(os.DirFS(parent), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path == "." {
				return nil
			}
			if path == name {
				has = true
				return fs.SkipDir
			}
		}
		hasAnother = true
		return fs.SkipAll
	})

	return has && !hasAnother, err
}
