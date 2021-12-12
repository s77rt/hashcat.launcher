package hashcatlauncher

import "path/filepath"

func (a *App) ScanHashes() (err error) {
	a.Hashes, err = filepath.Glob(filepath.Join(HashesDir, "*"))
	return
}

func (a *App) ScanDictionaries() (err error) {
	a.Dictionaries, err = filepath.Glob(filepath.Join(DictionariesDir, "*"))
	return
}

func (a *App) ScanRules() (err error) {
	a.Rules, err = filepath.Glob(filepath.Join(RulesDir, "*"))
	return
}

func (a *App) ScanMasks() (err error) {
	a.Masks, err = filepath.Glob(filepath.Join(MasksDir, "*"))
	return
}
