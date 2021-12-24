package hashcatlauncher

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func NewFilePath(path string) string {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return path
	}

	name := strings.TrimSuffix(path, filepath.Ext(path))
	n := 1
	ext := filepath.Ext(path)

	pathN := fmt.Sprintf("%s (%d)%s", name, n, ext)
	_, err = os.Stat(pathN)
	for !errors.Is(err, os.ErrNotExist) {
		n++
		pathN = fmt.Sprintf("%s (%d)%s", name, n, ext)
		_, err = os.Stat(pathN)
	}

	return pathN
}
