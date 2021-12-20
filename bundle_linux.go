package hashcatlauncher

import (
	"embed"
	"io"
	"os"
	"path/filepath"
)

//go:embed resources
var resources embed.FS

func (a *App) bundleIcon() error {
	iconDir := filepath.Join(os.Getenv("HOME"), ".local/share/icons")
	err := os.MkdirAll(iconDir, 0o755)
	if err != nil {
		return err
	}

	icon, err := os.Create(filepath.Join(iconDir, "hashcat.launcher.png"))
	if err != nil {
		return err
	}
	defer icon.Close()

	bundledIcon, err := resources.Open("resources/Icon.png")
	if err != nil {
		return err
	}

	_, err = io.Copy(icon, bundledIcon)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) bundleDesktopEntry() error {
	desktopEntryDir := filepath.Join(os.Getenv("HOME"), ".local/share/applications")
	err := os.MkdirAll(desktopEntryDir, 0o755)
	if err != nil {
		return err
	}

	desktopEntry, err := os.Create(filepath.Join(desktopEntryDir, "hashcat.launcher.desktop"))
	if err != nil {
		return err
	}
	defer desktopEntry.Close()

	bundledDesktopEntry, err := resources.Open("resources/hashcat.launcher.desktop")
	if err != nil {
		return err
	}

	_, err = io.Copy(desktopEntry, bundledDesktopEntry)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) Bundle() error {
	if err := a.bundleIcon(); err != nil {
		return err
	}
	if err := a.bundleDesktopEntry(); err != nil {
		return err
	}

	return nil
}
