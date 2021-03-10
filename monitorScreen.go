package hashcatlauncher

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func monitorScreen(hcl_gui *hcl_gui) fyne.CanvasObject {
	return container.NewVBox(
		container.NewGridWithColumns(4,
			monitorHardwares(hcl_gui)...,
		),
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewLabel("Note: stats will reset every 300s"),
		),
	)
}

func monitorHardware(hcl_gui *hcl_gui, index int) fyne.CanvasObject {
	return	widget.NewCard(fmt.Sprintf("Device #%d", index+1), "monitoring stats",
				container.New(layout.NewFormLayout(),
					widget.NewLabelWithStyle("Temp:", fyne.TextAlignLeading, fyne.TextStyle{}),
					container.NewHScroll(hcl_gui.monitor.hardwares[index].temp),
					widget.NewLabelWithStyle("Core:", fyne.TextAlignLeading, fyne.TextStyle{}),
					container.NewHScroll(hcl_gui.monitor.hardwares[index].core),
					widget.NewLabelWithStyle("Mem:", fyne.TextAlignLeading, fyne.TextStyle{}),
					container.NewHScroll(hcl_gui.monitor.hardwares[index].mem),
					widget.NewLabelWithStyle("Bus:", fyne.TextAlignLeading, fyne.TextStyle{}),
					container.NewHScroll(hcl_gui.monitor.hardwares[index].bus),
					widget.NewLabelWithStyle("Fan (%):", fyne.TextAlignLeading, fyne.TextStyle{}),
					hcl_gui.monitor.hardwares[index].fan,
					widget.NewLabelWithStyle("Util (%):", fyne.TextAlignLeading, fyne.TextStyle{}),
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
