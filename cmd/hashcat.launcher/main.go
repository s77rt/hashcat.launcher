package main

import (
	"os"
	"os/signal"
	"syscall"
	"fyne.io/fyne/v2/app"
	"github.com/s77rt/hashcat.launcher"
	dialog2 "github.com/OpenDiablo2/dialog"
)

func main() {
	app := app.NewWithID("s77rt.hashcat.launcher")
	app_gui := hashcatlauncher.NewGUI()
	app_gui.Pre(app)
	app_gui.LoadUI(app)
	app_gui.Post(app)
	app.SetIcon(app_gui.Icon)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		<- sigs
		app_gui.Clean()
		os.Exit(0)
	}()

	defer func() {
		app_gui.Clean()
	}()
 
	dialog2.Init()
	app.Driver().Run()
}
