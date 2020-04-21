package hashcatlauncher

import (
	"os"
	"runtime"
	"path/filepath"
	"fyne.io/fyne"
)

func GetPreference_hashcat_binary_file(app fyne.App) string {
	var fallback string
	pwd, _ := os.Getwd()
	if runtime.GOOS == "windows" {
		fallback = filepath.Join(pwd, "hashcat.exe")
	} else {
		fallback = filepath.Join(pwd, "hashcat.bin")
	}
	return app.Preferences().StringWithFallback("hashcat_binary_file", fallback)
}

func SetPreference_hashcat_binary_file(app fyne.App, hcl_gui *hcl_gui, value string) {
	hcl_gui.hashcat.binary_file = value
	app.Preferences().SetString("hashcat_binary_file", value)
	go get_available_hash_typess(hcl_gui)
}

func GetPreference_hashcat_status_timer(app fyne.App) int {
	return app.Preferences().IntWithFallback("hashcat_status_timer", 90)
}

func SetPreference_hashcat_status_timer(app fyne.App, hcl_gui *hcl_gui, value int) {
	hcl_gui.hashcat.args.status_timer = value
	app.Preferences().SetInt("hashcat_status_timer", value)
}

func GetPreference_hashcat_extra_args(app fyne.App) string {
	return app.Preferences().StringWithFallback("hashcat_extra_args", "--logfile-disable --restore-disable")
}

func SetPreference_hashcat_extra_args(app fyne.App, hcl_gui *hcl_gui, value string) {
	app.Preferences().SetString("hashcat_extra_args", value)
}

func GetPreference_max_active_sessions(app fyne.App) int {
	return app.Preferences().IntWithFallback("max_active_sessions", 1)
}

func SetPreference_max_active_sessions(app fyne.App, hcl_gui *hcl_gui, value int) {
	hcl_gui.max_active_sessions = value
	app.Preferences().SetInt("max_active_sessions", value)
}

func GetPreference_autostart_sessions(app fyne.App) bool {
	return app.Preferences().BoolWithFallback("autostart_sessions", true)
}

func SetPreference_autostart_sessions(app fyne.App, hcl_gui *hcl_gui, value bool) {
	hcl_gui.autostart_sessions = value
	app.Preferences().SetBool("autostart_sessions", value)
}
