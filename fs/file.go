package fs

import (
	"../fm"
	"os"
	"path/filepath"
)

func Stat(path string) (fm.File, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return os.Stat(absPath)
}
