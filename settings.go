package hashcatlauncher

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type Settings struct {
	mu *sync.Mutex

	TaskCounter int `json:"taskCounter"`
}

func (settings *Settings) CurrentTaskCounter() int {
	settings.mu.Lock()
	defer settings.mu.Unlock()

	return settings.TaskCounter
}

func (settings *Settings) NextTaskCounter() int {
	settings.mu.Lock()
	defer settings.mu.Unlock()

	settings.TaskCounter++
	return settings.TaskCounter
}

func (settings *Settings) ResetTaskCounter() int {
	settings.mu.Lock()
	defer settings.mu.Unlock()

	settings.TaskCounter = 0
	return settings.TaskCounter
}

var DefaultSettings = &Settings{}

func (a *App) LoadSettings() error {
	a.Settings = DefaultSettings
	a.Settings.mu = new(sync.Mutex)

	settingsFile := filepath.Join(a.Dir, "settings.json")
	if _, err := os.Stat(settingsFile); err == nil {
		raw, err := ioutil.ReadFile(settingsFile)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(raw, &a.Settings); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) SaveSettings() error {
	settingsFile := filepath.Join(a.Dir, "settings.json")

	b, err := json.Marshal(a.Settings)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(settingsFile, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(b)
	if err != nil {
		return err
	}

	return nil
}
