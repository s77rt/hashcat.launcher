package hashcatlauncher

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func tasksScreen(hcl_gui *hcl_gui) fyne.CanvasObject {
	// Content
	hcl_gui.tasks_content = container.NewVBox()

	// Tree
	hcl_gui.tasks_tree = widget.NewTreeWithStrings(
		map[string][]string{
			"":         {},
			"Running":  {},
			"Queued":   {},
			"Paused":   {},
			"Failed":   {},
			"Finished": {},
		},
	)

	hcl_gui.tasks_tree.OnSelected = func(uid string) {
		_, ok := hcl_gui.sessions[uid]
		if ok == true {
			hcl_gui.tasks_content.Objects = []fyne.CanvasObject{
				hcl_gui.sessions[uid].Content,
			}
			hcl_gui.tasks_content.Refresh()
		} else {
			hcl_gui.tasks_content.Objects = []fyne.CanvasObject{}
			hcl_gui.tasks_content.Refresh()
		}
	}

	hcl_gui.tasks_tree.OnUnselected = func(uid string) {
		hcl_gui.tasks_content.Objects = []fyne.CanvasObject{}
		hcl_gui.tasks_content.Refresh()
	}

	hcl_gui.tasks_tree.ChildUIDs = func(uid string) (children []string) {
		if uid == "" {
			children = []string{"Running", "Queued", "Paused", "Failed", "Finished"}
		} else {
			var status SessionStatus
			switch uid {
			case "Running":
				status = SessionStatusRunning
			case "Queued":
				status = SessionStatusQueued
			case "Paused":
				status = SessionStatusPaused
			case "Failed":
				status = SessionStatusFailed
			case "Finished":
				status = SessionStatusFinished
			}
			for _, session := range SortTasksByPriorityAndDate(hcl_gui) {
				if session.Status == status {
					children = append(children, session.Id)
				}
			}
		}
		return
	}

	hcl_gui.tasks_tree.UpdateNode = func(uid string, branch bool, node fyne.CanvasObject) {
		_, ok := hcl_gui.sessions[uid]
		if ok == true {
			node.(*widget.Label).SetText(hcl_gui.sessions[uid].Nickname)
		} else {
			text := uid
			switch uid {
			case "Running":
				text = fmt.Sprintf("%s (%d)", uid, hcl_gui.count_sessions_running)
			case "Queued":
				text = fmt.Sprintf("%s (%d)", uid, hcl_gui.count_sessions_queued)
			case "Paused":
				text = fmt.Sprintf("%s (%d)", uid, hcl_gui.count_sessions_paused)
			case "Failed":
				text = fmt.Sprintf("%s (%d)", uid, hcl_gui.count_sessions_failed)
			case "Finished":
				text = fmt.Sprintf("%s (%d)", uid, hcl_gui.count_sessions_finished)
			}
			node.(*widget.Label).SetText(text)
		}
	}

	hcl_gui.tasks_tree.OpenAllBranches()

	return container.NewHSplit(
		hcl_gui.tasks_tree,
		hcl_gui.tasks_content,
	)
}

func tasks_Refresh(hcl_gui *hcl_gui) {
	hcl_gui.tabs["Tasks"].Text = fmt.Sprintf("Tasks (%d)", len(hcl_gui.sessions))
	hcl_gui.tabs_container.Refresh()
	hcl_gui.tasks_tree.Refresh()
	hcl_gui.tasks_content.Refresh()
}
