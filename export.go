package hashcatlauncher

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
)

func (a *App) ExportConfig(config interface{}) error {
	path := NewFilePath(filepath.Join(a.ExportedDir, "config.json"))

	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	if err := encoder.Encode(config); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}
