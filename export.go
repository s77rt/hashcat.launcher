package hashcatlauncher

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
)

func (a *App) ExportConfig(config interface{}) (path string, err error) {
	path = NewFilePath(filepath.Join(a.ExportedDir, "config.json"))

	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	err = encoder.Encode(config)
	if err != nil {
		return
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.Write(buffer.Bytes())
	if err != nil {
		return
	}

	return
}
