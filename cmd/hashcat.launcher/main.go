package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	hashcatlauncher "github.com/s77rt/hashcat.launcher"
)

func main() {
	app := hashcatlauncher.NewApp()
	if err := app.Init(); err != nil {
		panic(err)
	}

	if err := app.NewServer(); err != nil {
		panic(err)
	}
	defer app.Server.Close()

	if err := app.NewUI(); err != nil {
		panic(err)
	}
	defer app.UI.Close()

	app.BindUI()

	if err := app.LoadUI(); err != nil {
		panic(err)
	}

	app.RestrictUI()

	if err := app.NewWatcher(); err != nil {
		panic(err)
	}
	defer app.Watcher.Close()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	select {
	case <-sigs:
	case <-app.UI.Done():
	}

	if err := app.Clean(); err != nil {
		log.Println(err)
	}
}
