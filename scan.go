package hashcatlauncher

import (
	"errors"
	"os"
	"path/filepath"
)

const MaxRecursiveFileWalk = 1000

func fileWalk(dir string, level *int) ([]string, error) {
	if level == nil {
		level = new(int)
		*level = 0
	} else {
		*level++
	}

	if *level > MaxRecursiveFileWalk {
		return nil, errors.New("Too many recursive file walk (cyclic import?)")
	}

	files := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return nil
		}
		if f.IsDir() {
			return nil
		}
		if f.Mode()&os.ModeSymlink != 0 {
			realPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				return err
			}
			realPathFiles, err := fileWalk(realPath, level)
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
	a.Hashes, err = fileWalk(a.HashesDir, nil)
	return
}

func (a *App) ScanDictionaries() (err error) {
	a.Dictionaries, err = fileWalk(a.DictionariesDir, nil)
	return
}

func (a *App) ScanRules() (err error) {
	a.Rules, err = fileWalk(a.RulesDir, nil)
	return
}

func (a *App) ScanMasks() (err error) {
	a.Masks, err = fileWalk(a.MasksDir, nil)
	return
}
