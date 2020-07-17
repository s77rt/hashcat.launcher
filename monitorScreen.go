package hashcatlauncher

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"fyne.io/fyne/layout"
)

func monitorScreen(hcl_gui *hcl_gui) fyne.CanvasObject {
	return widget.NewVBox(
		fyne.NewContainerWithLayout(layout.NewGridLayout(6),
			monitorHardwares(hcl_gui)...,
		),
		widget.NewLabelWithStyle("IMPORTANT: The first progress bar is for the Fan and the second is for the Util | Stats will reset every: 60s", fyne.TextAlignCenter, fyne.TextStyle{}),
	)
}

func monitorHardware(hcl_gui *hcl_gui, index int) fyne.CanvasObject {
	return	widget.NewVBox(
				widget.NewGroup(fmt.Sprintf("#%d", index+1),
					widget.NewHBox(
						widget.NewLabelWithStyle("Temp:", fyne.TextAlignLeading, fyne.TextStyle{}),
						hcl_gui.monitor.hardwares[index].temp,
					),
					widget.NewHBox(
						widget.NewLabelWithStyle("Core:", fyne.TextAlignLeading, fyne.TextStyle{}),
						hcl_gui.monitor.hardwares[index].core,
					),
					widget.NewHBox(
						widget.NewLabelWithStyle("Mem:", fyne.TextAlignLeading, fyne.TextStyle{}),
						hcl_gui.monitor.hardwares[index].mem,
					),
					widget.NewHBox(
						widget.NewLabelWithStyle("Bus:", fyne.TextAlignLeading, fyne.TextStyle{}),
						hcl_gui.monitor.hardwares[index].bus,
					),
					hcl_gui.monitor.hardwares[index].fan,
					hcl_gui.monitor.hardwares[index].util,
				),
			)
}

func monitorHardwares(hcl_gui *hcl_gui) []fyne.CanvasObject {
	canvas_objects := []fyne.CanvasObject{}
	for i, _ := range hcl_gui.monitor.hardwares {
		canvas_objects = append(canvas_objects, monitorHardware(hcl_gui, i))
	}
	return canvas_objects
}
