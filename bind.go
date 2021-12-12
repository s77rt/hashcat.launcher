package hashcatlauncher

import (
	"errors"

	"github.com/sqweek/dialog"
)

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

	a.UI.Bind("GOsaveDialog", func() (string, error) {
		return dialog.File().Save()
	})
}
