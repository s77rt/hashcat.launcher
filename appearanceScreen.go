package hashcatlauncher

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func appearanceScreen(hcl_gui *hcl_gui, hash_type_fakeselector *widget.Box) fyne.CanvasObject {
	return widget.NewVBox(
		widget.NewGroup("Appearance",
			widget.NewForm(
				widget.NewFormItem("Theme:", 
					widget.NewSelect([]string{"Light", "Dark"}, func(theme string) {
						switch theme {
						case "Light":
							hcl_gui.Settings.SetTheme("light")
						case "Dark":
							hcl_gui.Settings.SetTheme("dark")
						}

						go get_available_hash_typess(hcl_gui)
						if hcl_gui.hashcat.args.hash_type >= 0 {
							fake_hash_type_selector_hack(hcl_gui, hash_type_fakeselector, fmt.Sprintf("%d", hcl_gui.hashcat.args.hash_type))
						} else {
							fake_hash_type_selector_hack(hcl_gui, hash_type_fakeselector, "(Select one)")
						}
					}),
				),
				widget.NewFormItem("Scaling:",
					widget.NewSelect([]string{"auto", "50%", "75%", "80%", "90%", "100%", "110%", "125%", "150%", "175%", "200%"}, func(value string) {
						hcl_gui.Settings.SetScale(value)
					}),
				),
			),
		),
	)
}
