package hashcatlauncher

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func optionsScreen(app fyne.App, hcl_gui *hcl_gui) fyne.CanvasObject {
	hcl_gui.hc_binary_file_select = widget.NewSelect([]string{"Browse..."}, func(string) {
		go func() {
			file, err := NewFileOpen(hcl_gui)
			if err == nil {
				SetPreference_hashcat_binary_file(app, hcl_gui, file)
				hcl_gui.hc_binary_file_select.Selected = file
			} else {
				hcl_gui.hc_binary_file_select.Selected = GetPreference_hashcat_binary_file(app)
			}
			hcl_gui.hc_binary_file_select.Refresh()
		}()
	})

	hcl_gui.hc_status_timer_select = widget.NewSelect([]string{"10s", "30s", "60s", "90s", "120s", "300s", "Disabled"}, func(s string) {
		if s == "Disabled" {
			SetPreference_hashcat_status_timer(app, hcl_gui, 0)
		} else {
			v, _ := strconv.Atoi(s[:len(s)-1])
			SetPreference_hashcat_status_timer(app, hcl_gui, v)
		}

	})

	hcl_gui.autostart_sessions_select = widget.NewSelect([]string{"Enable", "Disable"}, func(s string) {
		if s == "Enable" {
			SetPreference_autostart_sessions(app, hcl_gui, true)
		} else {
			SetPreference_autostart_sessions(app, hcl_gui, false)
		}
	})

	hcl_gui.max_active_sessions_select = widget.NewSelect([]string{"1", "2", "3", "4", "5"}, func(s string) {
		v, _ := strconv.Atoi(s)
		SetPreference_max_active_sessions(app, hcl_gui, v)
	})

	hcl_gui.hc_extra_args = widget.NewEntry()
	hcl_gui.hc_extra_args.SetPlaceHolder("(Advanced Only)")
	hcl_gui.hc_extra_args.OnChanged = func(s string) {
		SetPreference_hashcat_extra_args(app, hcl_gui, s)
	}

	primarycolornames := []string{}
	for _, c := range theme.PrimaryColorNames() {
		primarycolornames = append(primarycolornames, c)
	}

	return container.NewVBox(
		widget.NewCard("Settings", "hashcat settings",
			container.NewVBox(
				widget.NewForm(
					widget.NewFormItem("Hashcat:", hcl_gui.hc_binary_file_select),
					widget.NewFormItem("Status Timer:", hcl_gui.hc_status_timer_select),
					widget.NewFormItem("Extra Args:", hcl_gui.hc_extra_args),
				),
				container.NewHBox(
					layout.NewSpacer(),
					widget.NewLabel("Note: hashcat version must be "+hashcat_min_version+" or higher"),
				),
			),
		),
		widget.NewCard("Options", "hashcat.launcher Options",
			widget.NewForm(
				widget.NewFormItem("Auto Task Start:", hcl_gui.autostart_sessions_select),
				widget.NewFormItem("Max Active Tasks:", hcl_gui.max_active_sessions_select),
				widget.NewFormItem("Theme:",
					func() fyne.CanvasObject {
						w := widget.NewSelect([]string{"Light", "Dark"}, func(theme string) {
							switch theme {
							case "Light":
								hcl_gui.Settings.SetTheme("light")
							case "Dark":
								hcl_gui.Settings.SetTheme("dark")
							}
						})
						if hcl_gui.Settings.Theme() == "light" {
							w.SetSelected("Light")
						} else if hcl_gui.Settings.Theme() == "dark" {
							w.SetSelected("Dark")
						}
						return w
					}(),
				),
				widget.NewFormItem("Primary Color:",
					func() fyne.CanvasObject {
						w := widget.NewSelect(primarycolornames, func(primarycolorname string) {
							hcl_gui.Settings.SetPrimaryColor(primarycolorname)
						})
						w.SetSelected(hcl_gui.Settings.PrimaryColor())
						return w
					}(),
				),
				widget.NewFormItem("Scaling:",
					func() fyne.CanvasObject {
						w := widget.NewSelect([]string{"auto", "50%", "70%", "75%", "80%", "85%", "90%", "95%", "100%", "130%", "180%"}, func(value string) {
							hcl_gui.Settings.SetScale(value)
						})
						w.SetSelected(fmt.Sprintf("%.0f%%", hcl_gui.Settings.Scale()*100))
						return w
					}(),
				),
				widget.NewFormItem("Dialog Handler:",
					func() fyne.CanvasObject {
						w := widget.NewSelect([]string{"OS", "Native"}, func(s string) {
							if s == "OS" {
								SetPreference_dialog_handler(app, hcl_gui, Dialog_OS)
							} else {
								SetPreference_dialog_handler(app, hcl_gui, Dialog_Native)
							}
						})
						if GetPreference_dialog_handler(app) == Dialog_OS {
							w.SetSelected("OS")
						} else if GetPreference_dialog_handler(app) == Dialog_Native {
							w.SetSelected("Native")
						}
						return w
					}(),
				),
			),
		),
		widget.NewCard("Data", "clear data and reset stats",
			container.NewGridWithColumns(4,
				widget.NewButton("Reset Task Id Counter", func() {
					dialog.ShowConfirm(
						"Reset Task Id Counter?",
						"Are you sure you want to reset task id counter?",
						func(confirm bool) {
							if confirm {
								SetPreference_next_task_id(app, hcl_gui, 1)
								dialog.ShowInformation("Success", "Task Id Counter has been reset.", hcl_gui.window)
							}
						},
						hcl_gui.window,
					)
				}),
				widget.NewButton("Reset Monitor Stats", func() {
					dialog.ShowConfirm(
						"Reset Monitor Stats?",
						"Are you sure you want to reset monitor stats?",
						func(confirm bool) {
							if confirm {
								hcl_gui.monitor.Reset()
								dialog.ShowInformation("Success", "Monitor Stats has been reset.", hcl_gui.window)
							}
						},
						hcl_gui.window,
					)
				}),
			),
		),
	)
}
