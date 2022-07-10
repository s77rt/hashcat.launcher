package hashcatlauncher

import (
	"errors"
	"fmt"

	"github.com/ncruces/zenity"
	"github.com/zserge/lorca"
)

func (a *App) NewUI() error {
	ui, err := lorca.New(
		"data:text/html,<html><head><title>hashcat.launcher</title></head><body>Loading...</body></html>",
		"",
		1080,
		720,
		[]string{"--class=hashcat.launcher"}...,
	)
	if err != nil {
		return err
	}

	ui.SetBounds(lorca.Bounds{
		WindowState: lorca.WindowStateMaximized,
	})

	a.UI = ui

	return nil
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

	a.UI.Bind("GOcreateTask", func(args HashcatArgs, priority int64) error {
		return a.NewTask(args, priority)
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

	a.UI.Bind("GOpriorityTask", func(taskID string, priority int64) error {
		if task, ok := a.Tasks[taskID]; ok {
			task.Priority = priority
			return nil
		}

		return errors.New("task not found")
	})

	a.UI.Bind("GOstartNextTask", func() {
		a.StartNextTask()
	})

	a.UI.Bind("GOrestoreTasks", func() error {
		return a.RestoreTasks()
	})

	a.UI.Bind("GOdeleteTask", func(taskID string) error {
		return a.DeleteTask(taskID)
	})

	a.UI.Bind("GOhashcatDevices", func() (string, error) {
		return a.Hashcat.Devices()
	})

	a.UI.Bind("GOhashcatBenchmark", func(hashMode HashcatHashMode) (string, error) {
		return a.Hashcat.Benchmark(hashMode)
	})

	a.UI.Bind("GOexportConfig", func(config interface{}) (string, error) {
		return a.ExportConfig(config)
	})

	a.UI.Bind("GOsaveHash", func(hash []byte, filename string) (string, error) {
		return a.SaveHash(hash, filename)
	})

	a.UI.Bind("GOsettingsCurrentTaskCounter", func() int {
		return a.Settings.CurrentTaskCounter()
	})

	a.UI.Bind("GOsettingsNextTaskCounter", func() int {
		return a.Settings.NextTaskCounter()
	})

	a.UI.Bind("GOsettingsResetTaskCounter", func() int {
		return a.Settings.ResetTaskCounter()
	})

	a.UI.Bind("GOsettingsCurrentLanguage", func() string {
		return a.Settings.CurrentLanguage()
	})

	a.UI.Bind("GOsettingsChangeLanguage", func(lang string) error {
		return a.Settings.ChangeLanguage(lang)
	})

	a.UI.Bind("GOsaveDialog", func() (string, error) {
		return zenity.SelectFileSave(zenity.ConfirmOverwrite())
	})
}

func (a *App) LoadUI() error {
	return a.UI.Load(fmt.Sprintf("http://%s/frontend/hashcat.launcher/build", a.Server.Addr()))
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
