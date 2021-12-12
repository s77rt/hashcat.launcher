package main

import (
	"os"
	"os/signal"
	"syscall"

	hashcatlauncher "github.com/s77rt/hashcat.launcher"
)

func main() {
	app := hashcatlauncher.NewApp()
	app.Init()

	app.NewServer()
	defer app.Server.Close()

	app.NewUI()
	defer app.UI.Close()

	app.BindUI()

	app.LoadUI()

	app.RestrictUI()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	select {
	case <-sigs:
	case <-app.UI.Done():
	}
	app.Clean()
	os.Exit(0)
}
