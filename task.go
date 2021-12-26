package hashcatlauncher

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/s77rt/hashcat.launcher/pkg/ansi"
	"github.com/s77rt/hashcat.launcher/pkg/subprocess"
)

type Task struct {
	ID        string                `json:"id"`
	Arguments []string              `json:"arguments"`
	Process   subprocess.Subprocess `json:"process"`
	Priority  int64                 `json:"priority"`
}

func (task *Task) Start() error {
	if task.Process.Status == subprocess.SubprocessStatusRunning {
		return errors.New("Task has been started already")
	}

	go task.Process.Execute()
	return nil
}

func (task *Task) Refresh() error {
	if task.Process.Status != subprocess.SubprocessStatusRunning {
		return errors.New("Task has not been started yet")
	}

	if runtime.GOOS == "windows" {
		task.Process.PostKey(0x53)
	} else {
		io.WriteString(task.Process.StdinStream, "s")
	}
	return nil
}
func (task *Task) Pause() error {
	if task.Process.Status != subprocess.SubprocessStatusRunning {
		return errors.New("Task has not been started yet")
	}

	if runtime.GOOS == "windows" {
		task.Process.PostKey(0x50)
	} else {
		io.WriteString(task.Process.StdinStream, "p")
	}
	return nil
}
func (task *Task) Resume() error {
	if task.Process.Status != subprocess.SubprocessStatusRunning {
		return errors.New("Task has not been started yet")
	}

	if runtime.GOOS == "windows" {
		task.Process.PostKey(0x52)
	} else {
		io.WriteString(task.Process.StdinStream, "r")
	}
	return nil
}
func (task *Task) Checkpoint() error {
	if task.Process.Status != subprocess.SubprocessStatusRunning {
		return errors.New("Task has not been started yet")
	}

	if runtime.GOOS == "windows" {
		task.Process.PostKey(0x43)
	} else {
		io.WriteString(task.Process.StdinStream, "c")
	}
	return nil
}
func (task *Task) Skip() error {
	if task.Process.Status != subprocess.SubprocessStatusRunning {
		return errors.New("Task has not been started yet")
	}

	if runtime.GOOS == "windows" {
		task.Process.PostKey(0x42)
	} else {
		io.WriteString(task.Process.StdinStream, "b")
	}
	return nil
}
func (task *Task) Quit() error {
	if task.Process.Status != subprocess.SubprocessStatusRunning {
		return errors.New("Task has not been started yet")
	}

	if runtime.GOOS == "windows" {
		task.Process.PostKey(0x51)
	} else {
		io.WriteString(task.Process.StdinStream, "q")
	}
	return nil
}

type TaskUpdate struct {
	Task      Task   `json:"task"`
	Message   string `json:"message"`
	Source    string `json:"source"`
	Timestamp int64  `json:"timestamp"`
}

func (a *App) TaskExists(taskID string) bool {
	_, exists := a.Tasks[taskID]
	return exists
}

func (a *App) newTaskID() (taskID string) {
	for {
		taskID = fmt.Sprintf("Task #%d", a.Settings.NextTaskCounter())
		if !a.TaskExists(taskID) {
			break
		}
	}

	return
}

func (a *App) NewTask(args HashcatArgs, priority int64) (err error) {
	task := &Task{
		ID: a.newTaskID(),
	}

	args.Session = &task.ID

	task.Arguments, err = args.Build()
	if err != nil {
		return
	}

	wdir, _ := filepath.Split(a.Hashcat.BinaryFile)
	task.Process = subprocess.Subprocess{
		subprocess.SubprocessStatusNotStarted,
		wdir,
		a.Hashcat.BinaryFile,
		task.Arguments,
		nil,
		nil,
		func(s string) {
			a.TaskUpdateCallback(TaskUpdate{
				Task:      *task,
				Message:   ansi.Strip(s),
				Source:    "stdout",
				Timestamp: time.Now().UnixNano(),
			})
		},
		func(s string) {
			a.TaskUpdateCallback(TaskUpdate{
				Task:      *task,
				Message:   ansi.Strip(s),
				Source:    "stderr",
				Timestamp: time.Now().UnixNano(),
			})
		},
		func() {
			a.TaskPreProcessCallback(TaskUpdate{
				Task:      *task,
				Timestamp: time.Now().UnixNano(),
			})
		},
		func() {
			a.TaskPostProcessCallback(TaskUpdate{
				Task:      *task,
				Timestamp: time.Now().UnixNano(),
			})
			a.StartNextTask()
		},
	}

	task.Priority = priority

	a.Tasks[task.ID] = task

	a.TaskAddCallback(TaskUpdate{
		Task:      *task,
		Timestamp: time.Now().UnixNano(),
	})

	if task.Priority >= 0 {
		a.StartNextTask()
	}

	return
}

func (a *App) RestoreTask(restoreFile RestoreFile) (err error) {
	_, filename := filepath.Split(restoreFile.File.Name())

	task := &Task{
		ID: strings.TrimSuffix(filename, RestoreFileExt),
	}

	if a.TaskExists(task.ID) {
		err = errors.New(fmt.Sprintf("Task already exists (ID: %s)", task.ID))
		return
	}

	task.Arguments = strings.Split(
		strings.TrimSuffix(
			strings.ReplaceAll(
				string(restoreFile.Data.Argv),
				"\r\n",
				"\n",
			),
			"\n",
		),
		"\n",
	)[1:]

	wdir, _ := filepath.Split(a.Hashcat.BinaryFile)
	task.Process = subprocess.Subprocess{
		subprocess.SubprocessStatusNotStarted,
		wdir,
		a.Hashcat.BinaryFile,
		[]string{fmt.Sprintf("--session=%s", task.ID), "--restore"},
		nil,
		nil,
		func(s string) {
			a.TaskUpdateCallback(TaskUpdate{
				Task:      *task,
				Message:   ansi.Strip(s),
				Source:    "stdout",
				Timestamp: time.Now().UnixNano(),
			})
		},
		func(s string) {
			a.TaskUpdateCallback(TaskUpdate{
				Task:      *task,
				Message:   ansi.Strip(s),
				Source:    "stderr",
				Timestamp: time.Now().UnixNano(),
			})
		},
		func() {
			a.TaskPreProcessCallback(TaskUpdate{
				Task:      *task,
				Timestamp: time.Now().UnixNano(),
			})
		},
		func() {
			a.TaskPostProcessCallback(TaskUpdate{
				Task:      *task,
				Timestamp: time.Now().UnixNano(),
			})
			a.StartNextTask()
		},
	}

	task.Priority = -1

	a.Tasks[task.ID] = task

	a.TaskAddCallback(TaskUpdate{
		Task:      *task,
		Timestamp: time.Now().UnixNano(),
	})

	if task.Priority >= 0 {
		a.StartNextTask()
	}

	return
}

func (a *App) RestoreTasks() (err error) {
	files, err := filepath.Glob(filepath.Join(a.HashcatDir, "*.restore"))
	if err != nil {
		return
	}

	for _, file := range files {
		var f *os.File
		f, err = os.Open(file)
		if err != nil {
			return
		}
		rf := RestoreFile{File: f}
		err = rf.Unpack()
		if err != nil {
			return
		}
		err = a.RestoreTask(rf)
		if err != nil {
			return
		}
	}

	return
}

// StartNextTask starts the next queued task with the highest priority
// given that the task has not been started yet and there is no other running tasks
// tasks that have a priorty less than zero are not eligible
func (a *App) StartNextTask() {
	runningTasks := 0

	tasks := []*Task{}
	for _, task := range a.Tasks {
		if task.Process.Status == subprocess.SubprocessStatusRunning {
			runningTasks++
		}
		tasks = append(tasks, task)
	}

	if runningTasks > 0 {
		return
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Priority > tasks[j].Priority
	})

	for _, task := range tasks {
		if task.Priority < 0 {
			// since tasks are ordered by priority
			// we can stop here as tasks with priority < 0 are not eligible
			// for auto-start
			break
		}
		if task.Process.Status == subprocess.SubprocessStatusNotStarted {
			task.Start()
			break
		}
	}
}

func (a *App) DeleteTask(taskID string) error {
	task, ok := a.Tasks[taskID]
	if !ok {
		return errors.New("Task does not exist")
	}

	if task.Process.Status == subprocess.SubprocessStatusRunning {
		return errors.New("Task is running")
	}

	os.Remove(filepath.Join(a.HashcatDir, fmt.Sprintf("%s.restore", task.ID)))

	delete(a.Tasks, taskID)

	a.TaskDeleteCallback(taskID)

	return nil
}
