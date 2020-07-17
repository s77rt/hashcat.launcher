package hashcatlauncher

import (
	"time"
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type hcl_gui_monitor struct {
	hardwares [18]monitor_hardware
}

type monitor_hardware struct {
	temp *widget.Label
	core *widget.Label
	mem *widget.Label
	bus *widget.Label
	fan *widget.ProgressBar
	util *widget.ProgressBar
}

func (monitor *hcl_gui_monitor) Init() {
	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				monitor.Reset()
			}
		}
	}()
	for i, _ := range monitor.hardwares {
		monitor.hardwares[i].temp = widget.NewLabelWithStyle("N/A", fyne.TextAlignLeading, fyne.TextStyle{Bold:true})
		monitor.hardwares[i].core = widget.NewLabelWithStyle("N/A", fyne.TextAlignLeading, fyne.TextStyle{Bold:true})
		monitor.hardwares[i].mem = widget.NewLabelWithStyle("N/A", fyne.TextAlignLeading, fyne.TextStyle{Bold:true})
		monitor.hardwares[i].bus = widget.NewLabelWithStyle("N/A", fyne.TextAlignLeading, fyne.TextStyle{Bold:true})
		monitor.hardwares[i].fan = widget.NewProgressBar()
		monitor.hardwares[i].fan.Min = 0
		monitor.hardwares[i].fan.Max = 100
		monitor.hardwares[i].fan.SetValue(0)
		monitor.hardwares[i].util = widget.NewProgressBar()
		monitor.hardwares[i].util.Min = 0
		monitor.hardwares[i].util.Max = 100
		monitor.hardwares[i].util.SetValue(0)
	}
}

func (monitor *hcl_gui_monitor) Reset() {
	for i, _ := range monitor.hardwares {
		monitor.hardwares[i].temp.SetText("N/A")
		monitor.hardwares[i].fan.SetValue(0)
		monitor.hardwares[i].util.SetValue(0)
		monitor.hardwares[i].core.SetText("N/A")
		monitor.hardwares[i].mem.SetText("N/A")
		monitor.hardwares[i].bus.SetText("N/A")
	}
}
