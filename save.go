package hashcatlauncher

import (
	"os"
	"path/filepath"
)

func (a *App) SaveHash(hash []byte, filename string) (path string, err error) {
	path = NewFilePath(filepath.Join(a.HashesDir, filename))

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.Write(hash)
	if err != nil {
		return
	}

	return
}
