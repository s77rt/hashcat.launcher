package main

import (
	"fyne.io/fyne/app"
	"github.com/s77rt/hashcat.launcher"
)

func main() {
	app := app.NewWithID("com.s77rt.hashcatlauncher.preferences")
	app_gui := hashcatlauncher.NewGUI()
	app_gui.LoadUI(app)
	app_gui.Init(app)
	app.SetIcon(app_gui.Icon)
	app.Driver().Run()
}
