package xsettings

import (
	"encoding/json"
	"os"
	"log"
	"strings"
	"strconv"
	"path/filepath"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
)

// Settings gives access to user interfaces to control Fyne settings
type Settings struct {
	fyneSettings app.SettingsSchema
}

func (s *Settings) save() error {
	return s.saveToFile(s.fyneSettings.StoragePath())
}

func (s *Settings) saveToFile(path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil { // this is not an exists error according to docs
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
		file, err = os.Open(path)
		if err != nil {
			return err
		}
	}
	encode := json.NewEncoder(file)

	return encode.Encode(&s.fyneSettings)
}

func (s *Settings) load() {
	err := s.loadFromFile(s.fyneSettings.StoragePath())
	if err != nil {
		fyne.LogError("Settings load error:", err)
	}
}

func (s *Settings) loadFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(path), 0700)
			return nil
		}
		return err
	}
	decode := json.NewDecoder(file)

	return decode.Decode(&s.fyneSettings)
}

func (s *Settings) chooseTheme(name string) {
	s.fyneSettings.ThemeName = name
}

func (s *Settings) chooseScale(value string) {
	if value == "" || strings.EqualFold(value, "auto") {
		s.fyneSettings.Scale = fyne.SettingsScaleAuto
		return
	}

	value = value[:len(value)-1] // Remove "%" (100% => 100)
	scale, err := strconv.ParseFloat(value, 32)
	if err != nil {
		log.Println("Cannot set scale to:", value)
	}
	s.fyneSettings.Scale = float32(scale/100)
}

func (s *Settings) Theme() string {
	return s.fyneSettings.ThemeName
}

func (s *Settings) Scale() float32 {
	return s.fyneSettings.Scale
}

func (s *Settings) SetTheme(name string) {
	s.chooseTheme(name)
	s.save()
}

func (s *Settings) SetScale(value string) {
	s.chooseScale(value)
	s.save()
}

// NewSettings returns a new settings instance with the current configuration loaded
func NewSettings() *Settings {
	s := &Settings{}
	s.load()

	return s
}
