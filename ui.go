package hashcatlauncher

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/sqweek/dialog"
	"github.com/zserge/lorca"
)

func (a *App) NewUI() {
	tmpDir, err := ioutil.TempDir("", "hashcat.launcher")
	if err != nil {
		log.Fatal(err)
	}

	a.UI, err = lorca.New(
		"data:text/html,<html><head><title>hashcat.launcher</title></head><body>Loading...</body></html>",
		tmpDir,
		1080,
		720,
		[]string{"--class=hashcat.launcher"}...,
	)
	if err != nil {
		log.Fatal(err)
	}

	a.UI.SetBounds(lorca.Bounds{
		WindowState: lorca.WindowStateMaximized,
	})
}

func (a *App) BindUI() {
	a.UI.Bind("GOgetVersion", func() string {
		return Version
	})

	a.UI.Bind("GOscan", func() error {
		return a.Scan()
	})

	a.UI.Bind("GOgetHashes", func() []string {
		return a.Hashes
	})

	a.UI.Bind("GOgetAlgorithms", func() map[int64]string {
		return a.Hashcat.Algorithms
	})

	a.UI.Bind("GOgetDictionaries", func() []string {
		return a.Dictionaries
	})

	a.UI.Bind("GOgetRules", func() []string {
		return a.Rules
	})

	a.UI.Bind("GOgetMasks", func() []string {
		return a.Masks
	})

	a.UI.Bind("GOcreateTask", func(args HashcatArgs) error {
		return a.NewTask(args)
	})

	a.UI.Bind("GOstartTask", func(taskID string) error {
		if task, ok := a.Tasks[taskID]; ok {
			return task.Start()
		}

		return errors.New("task not found")
	})

	a.UI.Bind("GOrefreshTask", func(taskID string) error {
		if task, ok := a.Tasks[taskID]; ok {
			return task.Refresh()
		}

		return errors.New("task not found")
	})

	a.UI.Bind("GOpauseTask", func(taskID string) error {
		if task, ok := a.Tasks[taskID]; ok {
			return task.Pause()
		}

		return errors.New("task not found")
	})

	a.UI.Bind("GOresumeTask", func(taskID string) error {
		if task, ok := a.Tasks[taskID]; ok {
			return task.Resume()
		}

		return errors.New("task not found")
	})

	a.UI.Bind("GOcheckpointTask", func(taskID string) error {
		if task, ok := a.Tasks[taskID]; ok {
			return task.Checkpoint()
		}

		return errors.New("task not found")
	})

	a.UI.Bind("GOskipTask", func(taskID string) error {
		if task, ok := a.Tasks[taskID]; ok {
			return task.Skip()
		}

		return errors.New("task not found")
	})

	a.UI.Bind("GOquitTask", func(taskID string) error {
		if task, ok := a.Tasks[taskID]; ok {
			return task.Quit()
		}

		return errors.New("task not found")
	})

	a.UI.Bind("GOrestoreTasks", func() error {
		return a.RestoreTasks()
	})

	a.UI.Bind("GOsaveDialog", func() (string, error) {
		return dialog.File().Save()
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
