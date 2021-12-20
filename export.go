package hashcatlauncher

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func (a *App) ExportConfig(config interface{}) error {
	filename := filepath.Join(a.ExportedDir, "config.json")
	_, err := os.Stat(filename)

	count := 0
	for !errors.Is(err, os.ErrNotExist) {
		count++
		filename = filepath.Join(a.ExportedDir, fmt.Sprintf("config (%d).json", count))
		_, err = os.Stat(filename)
	}

	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	if err := encoder.Encode(config); err != nil {
		return err
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, os.ModePerm)
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
