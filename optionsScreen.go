package hashcatlauncher

import (
	"strconv"
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	dialog2 "github.com/sqweek/dialog"
)

func optionsScreen(app fyne.App, hcl_gui *hcl_gui) fyne.CanvasObject {
	hcl_gui.hc_binary_file_select = widget.NewSelect([]string{"Browse..."}, func(string){
		file, err := dialog2.File().Title("Select hashcat binary file").Filter("Bin/Exe Files", "exe", "bin").Load()
		if err == nil {
			SetPreference_hashcat_binary_file(app, hcl_gui, file)
			hcl_gui.hc_binary_file_select.Selected = file
			hcl_gui.hc_binary_file_select.Refresh()
		} else {
			hcl_gui.hc_binary_file_select.Selected = GetPreference_hashcat_binary_file(app)
			hcl_gui.hc_binary_file_select.Refresh()
		}
	})

	hcl_gui.hc_status_timer_select = widget.NewSelect([]string{"10s", "30s", "60s", "90s", "120s", "300s", "Disabled"}, func(s string){
		if s == "Disabled" {
			SetPreference_hashcat_status_timer(app, hcl_gui, 0)
		} else {
			v, _ := strconv.Atoi(s[:len(s)-1])
			SetPreference_hashcat_status_timer(app, hcl_gui, v)
		}

	})

	hcl_gui.autostart_sessions_select = widget.NewSelect([]string{"Enable", "Disable"}, func(s string){
		if s == "Enable" {
			SetPreference_autostart_sessions(app, hcl_gui, true)
		} else {
			SetPreference_autostart_sessions(app, hcl_gui, false)
		}
	})

	hcl_gui.max_active_sessions_select = widget.NewSelect([]string{"1", "2", "3"}, func(s string){
		v, _ := strconv.Atoi(s)
		SetPreference_max_active_sessions(app, hcl_gui, v)
	})

	hcl_gui.hc_extra_args = widget.NewEntry()
	hcl_gui.hc_extra_args.SetPlaceHolder("(Advanced Only)")
	hcl_gui.hc_extra_args.OnChanged = func(s string) {
		s = s[:min(len(s), 40)]
		hcl_gui.hc_extra_args.SetText(s)
		SetPreference_hashcat_extra_args(app, hcl_gui, s)
	}

	return widget.NewVBox(
		widget.NewGroup("hashcat options",
			widget.NewForm(
				widget.NewFormItem("Binary:", hcl_gui.hc_binary_file_select),
				widget.NewFormItem("Status Timer:", hcl_gui.hc_status_timer_select),
				widget.NewFormItem("Extra Args:", hcl_gui.hc_extra_args),
			),
		),
		widget.NewGroup("hashcat.launcher options",
			widget.NewForm(
				widget.NewFormItem("Auto Session Start:", hcl_gui.autostart_sessions_select),
				widget.NewFormItem("Max Active Sessions:", hcl_gui.max_active_sessions_select),
			),
		),
		widget.NewGroup("Notes",
			widget.NewLabel("hashcat version must be v5.1.0+1774 (not v5.1.0) or newer"),
		),
	)
}
