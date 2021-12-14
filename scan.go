package hashcatlauncher

import (
	"os"
	"path/filepath"
)

func fileWalk(dir string) ([]string, error) {
	files := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}
		if f.Mode()&os.ModeSymlink != 0 {
			realPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				return err
			}
			realPathFiles, err := fileWalk(realPath)
			if err != nil {
				return err
			}
			files = append(files, realPathFiles...)
		} else {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func (a *App) ScanHashes() (err error) {
	a.Hashes, err = fileWalk(HashesDir)
	return
}

func (a *App) ScanDictionaries() (err error) {
	a.Dictionaries, err = fileWalk(DictionariesDir)
	return
}

func (a *App) ScanRules() (err error) {
	a.Rules, err = fileWalk(RulesDir)
	return
}

func (a *App) ScanMasks() (err error) {
	a.Masks, err = fileWalk(MasksDir)
	return
}
