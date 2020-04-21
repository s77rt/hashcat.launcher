package main

import (
	"os"
	"os/signal"
	"syscall"
	"fyne.io/fyne/app"
	"github.com/s77rt/hashcat.launcher"
)

func main() {
	app := app.NewWithID("com.s77rt.hashcatlauncher.preferences")
	app_gui := hashcatlauncher.NewGUI()
	app_gui.LoadUI(app)
	app_gui.Init(app)
	app.SetIcon(app_gui.Icon)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		<- sigs
		app_gui.Clean()
		os.Exit(0)
	}()

	app.Driver().Run()

	app_gui.Clean()
}
