package hashcatlauncher

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/zserge/lorca"
)

func (a *App) NewUI() {
	tmpDir, err := ioutil.TempDir("", Name)
	if err != nil {
		log.Fatal(err)
	}

	a.UI, err = lorca.New(
		fmt.Sprintf("data:text/html,<html><head><title>%s</title></head><body>Loading...</body></html>", Name),
		tmpDir,
		1080,
		720,
		[]string{fmt.Sprintf("--class=%s", Name)}...,
	)
	if err != nil {
		log.Fatal(err)
	}

	a.UI.SetBounds(lorca.Bounds{
		WindowState: lorca.WindowStateMaximized,
	})
}

func (a *App) LoadUI() {
	if err := a.UI.Load(fmt.Sprintf("http://%s/frontend/hashcat.launcher/build", a.Server.Addr())); err != nil {
		log.Fatal(err)
	}
}

func (a *App) RestrictUI() {
	a.UI.Eval(`window.addEventListener("contextmenu", function(e) { e.preventDefault(); })`)
	a.UI.Eval(`
		document.onkeydown = function (event) {
			if (event.ctrlKey) { // Ctrl is pressed
				// Allown only A, X, C, V
				if (event.keyCode != 65 && event.keyCode != 88 && event.keyCode != 67 && event.keyCode != 86)
					return false;
			}
			if (event.keyCode >= 112 && event.keyCode <= 123) // F1 to F12
				return false;
		}
	`)
}
