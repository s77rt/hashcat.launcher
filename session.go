package hashcatlauncher

import (
	"os"
	"io"
	"sort"
	"strings"
	"strconv"
	"time"
	"fmt"
	"errors"
	"path/filepath"
	"runtime"
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"github.com/s77rt/hashcat.launcher/pkg/subprocess"
)

type Session struct {
	Id string
	Nickname string
	Arguments []string
	Status SessionStatus
	Process subprocess.Subprocess
	Journal *widget.Label
	Content fyne.CanvasObject
	Notifications_Enabled bool
	Priority int
}

type SessionStatus int
const (
	SessionStatusRunning SessionStatus = iota
	SessionStatusQueued
	SessionStatusPaused
	SessionStatusFailed
	SessionStatusFinished
)

// SessionIdSorter sorts sessions by id (and thus the date since the id includes the date)
type SessionIdSorter []*Session

func (a SessionIdSorter) Len() int           { return len(a) }
func (a SessionIdSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SessionIdSorter) Less(i, j int) bool { return a[i].Id < a[j].Id }


func (session *Session) UpdateJournal(new_msg string) {
	session.Journal.SetText(fmt.Sprintf("%s %s\n%s", time.Now().Format("2006-01-02 15:04:05"), new_msg, session.Journal.Text))
}

func (session *Session) Start() {
	go session.Process.Execute()
	go func() {
		time.Sleep(time.Second)
		if session.Process.Stdin_stream != nil {
			session.Refresh()
		}
	}()
}

func (session *Session) Refresh() {
	if runtime.GOOS == "windows" {
		session.Process.PostKey(0x53)
	} else {
		io.WriteString(session.Process.Stdin_stream, "s")
	}
}
func (session *Session) Pause() {
	if runtime.GOOS == "windows" {
		session.Process.PostKey(0x50)
	} else {
		io.WriteString(session.Process.Stdin_stream, "p")
	}
}
func (session *Session) Resume() {
	if runtime.GOOS == "windows" {
		session.Process.PostKey(0x52)
	} else {
		io.WriteString(session.Process.Stdin_stream, "r")
	}
}
func (session *Session) Checkpoint() {
	if runtime.GOOS == "windows" {
		session.Process.PostKey(0x43)
	} else {
		io.WriteString(session.Process.Stdin_stream, "c")
	}
}
func (session *Session) Skip() {
	if runtime.GOOS == "windows" {
		session.Process.PostKey(0x42)
	} else {
		io.WriteString(session.Process.Stdin_stream, "b")
	}
}
func (session *Session) Quit() {
	if runtime.GOOS == "windows" {
		session.Process.PostKey(0x51)
	} else {
		io.WriteString(session.Process.Stdin_stream, "q")
	}
}

func (session *Session) SetStatus(app fyne.App, hcl_gui *hcl_gui, status SessionStatus) {
	old_status := session.Status
	session.Status = status
	if session.Status != old_status {
		CalculateSessionsStatusStats(hcl_gui)
		if session.Notifications_Enabled == true {
			app.SendNotification(&fyne.Notification{
				Title:   session.Nickname,
				Content: (func() string {
					switch session.Status {
					case SessionStatusRunning:
						return "Task is running..."
					case SessionStatusQueued:
						return "Task is queued."
					case SessionStatusPaused:
						return "Task has been paused."
					case SessionStatusFinished:
						return "Task has been completed."
					case SessionStatusFailed:
						return "Task failure!"
					}
					return "Task unknown status!"
				})(),
			})
		}
	}
	tasks_Refresh(hcl_gui)
}

func RemoveSession(hcl_gui *hcl_gui, session *Session) {
	delete(hcl_gui.sessions, session.Id)
	CalculateSessionsStatusStats(hcl_gui)
	tasks_Refresh(hcl_gui)
}

func NewSession(app fyne.App, hcl_gui *hcl_gui, task_id int, session_id string, restore_file *RestoreFile, attack_payload []string, enable_notifications bool, priority int) {
	if task_id < 0 {
		task_id = hcl_gui.next_task_id
		SetPreference_next_task_id(app, hcl_gui, (hcl_gui.next_task_id+1))
	}
	if session_id == "" {
		session_id = fmt.Sprintf("hcl_%d_%d", time.Now().UnixNano(), task_id)
	}
	_, ok := hcl_gui.sessions[session_id]
	if ok == true { // session exists
		dialog.ShowError(errors.New("Task already exists"), hcl_gui.window)
		return
	}
	session := &Session{}
	session.Id = session_id
	session.Nickname = fmt.Sprintf("T#%d", task_id)
	session.Status = SessionStatusQueued
	session.Journal = widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace:true})
	session.Notifications_Enabled = enable_notifications
	session.Priority = priority

	multiple_devices := false
	started := false
	terminated := false

	status := widget.NewLabelWithStyle("N/A", fyne.TextAlignLeading, fyne.TextStyle{Bold:true})
	speed := widget.NewLabelWithStyle("N/A", fyne.TextAlignLeading, fyne.TextStyle{Bold:true})
	hash_type := widget.NewLabelWithStyle("N/A", fyne.TextAlignLeading, fyne.TextStyle{Bold:true})
	hash_target := widget.NewLabelWithStyle("N/A", fyne.TextAlignLeading, fyne.TextStyle{Bold:true})
	time_started := widget.NewLabelWithStyle("N/A", fyne.TextAlignLeading, fyne.TextStyle{Bold:true})
	time_estimated := widget.NewLabelWithStyle("N/A", fyne.TextAlignLeading, fyne.TextStyle{Bold:true})
	progress := widget.NewProgressBar()
	progress.Min = 0
	progress.Max = 100
	progress_text := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{})
	recovered := widget.NewProgressBar()
	recovered.Min = 0
	recovered.Max = 100
	recovered_text := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{})
	guess_queue := widget.NewLabelWithStyle("N/A", fyne.TextAlignLeading, fyne.TextStyle{Bold:true})
	guess_base := widget.NewLabelWithStyle("N/A", fyne.TextAlignLeading, fyne.TextStyle{Bold:true})
	guess_mod := widget.NewLabelWithStyle("N/A", fyne.TextAlignLeading, fyne.TextStyle{Bold:true})
	guess_mask := widget.NewLabelWithStyle("N/A", fyne.TextAlignLeading, fyne.TextStyle{Bold:true})
	guess_charset := widget.NewLabelWithStyle("N/A", fyne.TextAlignLeading, fyne.TextStyle{Bold:true})

	start := widget.NewButtonWithIcon("Start", theme.MediaPlayIcon(), func(){
		if hcl_gui.count_sessions_running < hcl_gui.max_active_sessions {
			session.Start()
		} else {
			dialog.ShowError(errors.New("Max Active Sessions Reached!"), hcl_gui.window)
		}
	})
	info := widget.NewButtonWithIcon("Info", theme.InfoIcon(), func(){
		var modal *widget.PopUp
		c := widget.NewVBox(
			widget.NewHBox(
				fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{560, 30}),
					widget.NewLabel("Info"),
				),
				fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{40, 30}),
					widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
						modal.Hide()
					}),
				),
			),
			fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{600, 300}),
				container.NewMax(
					container.NewScroll(
						widget.NewLabel(strings.Join(session.Arguments, "\n")),
					),
				),
			),
		)
		hcl_gui.window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
			if key.Name == fyne.KeyEscape {
				hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){})
				modal.Hide()
			}
		})
		modal = widget.NewModalPopUp(c, hcl_gui.window.Canvas())
	})
	refresh := widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func(){
		session.Refresh()
	})
	refresh.Disable()
	var pause, resume *widget.Button
	pause = widget.NewButtonWithIcon("Pause", theme.MediaPauseIcon(), func(){
		session.Pause()
		go func(){
			time.Sleep(100*time.Millisecond)
			session.Refresh()
		}()
		session.UpdateJournal("Paused")
	})
	pause.Disable()
	resume = widget.NewButtonWithIcon("Resume", theme.MediaPlayIcon(), func(){
		if hcl_gui.count_sessions_running < hcl_gui.max_active_sessions {
			session.Resume()
			go func(){
				time.Sleep(100*time.Millisecond)
				session.Refresh()
			}()
			session.UpdateJournal("Resumed")
		} else {
			dialog.ShowError(errors.New("Max Active Sessions Reached!"), hcl_gui.window)
		}
	})
	resume.Disable()
	checkpoint := widget.NewButtonWithIcon("Checkpoint", theme.MediaRecordIcon(), func(){
		session.Checkpoint()
	})
	checkpoint.Disable()
	skip := widget.NewButtonWithIcon("Skip", theme.MediaSkipNextIcon(), func(){
		session.Skip()
	})
	skip.Disable()
	stop := widget.NewButtonWithIcon("Stop", theme.CheckButtonIcon(), func(){
		if started {
			session.Quit()
		}
		session.UpdateJournal("Graceful Stop")
	})
	stop.Disable()
	terminate := widget.NewButtonWithIcon("Terminate", theme.CancelIcon(), func(){
		if started {
			session.Quit()
		}
		session.UpdateJournal("Forceful Stop")
		session.Process.Kill()
		terminated = true
	})
	terminate.Disable()
	terminate_n_close := widget.NewButtonWithIcon("Terminate & Close", theme.CancelIcon(), func(){
		if started {
			session.Quit()
		}
		session.UpdateJournal("Forceful Stop")
		session.Process.Kill()
		terminated = true
		hcl_gui.tasks_tree.Unselect(session.Id)
		RemoveSession(hcl_gui, session)
	})

	args := []string{}
	if restore_file != nil {
		args = append(args, fmt.Sprintf("--session=%s", string(session.Id)))
		args = append(args, "--restore")
		session.Arguments = restore_file.GetArguments()
	} else {
		args = func() []string {
			args := []string{}

			args = append(args, fmt.Sprintf("--session=%s", string(session.Id)))

			if hcl_gui.hashcat.args.status_timer > 0 {
				args = append(args, []string{"--status", fmt.Sprintf("--status-timer=%d", hcl_gui.hashcat.args.status_timer)}...)
			}

			if hcl_gui.hashcat.args.optimized_kernel {
				args = append(args, "-O")
			}

			if hcl_gui.hashcat.args.slower_candidate {
				args = append(args, "-S")
			}

			if hcl_gui.hashcat.args.disable_self_test {
				args = append(args, "--self-test-disable")
			}

			if hcl_gui.hashcat.args.force {
				args = append(args, "--force")
			}

			if hcl_gui.hashcat.args.remove_found_hashes {
				args = append(args, "--remove")
			}

			if hcl_gui.hashcat.args.disable_potfile {
				args = append(args, "--potfile-disable")
			}

			if hcl_gui.hashcat.args.ignore_usernames {
				args = append(args, "--username")
			}

			if hcl_gui.hashcat.args.disable_monitor {
				args = append(args, "--hwmon-disable")
			} else {
				args = append(args, fmt.Sprintf("--hwmon-temp-abort=%d", hcl_gui.hashcat.args.temp_abort))
			}

			args = append(args, fmt.Sprintf("-w%d", hcl_gui.hashcat.args.workload_profile))

			args = append(args, fmt.Sprintf("-m%d", hcl_gui.hashcat.args.hash_type))

			args = append(args, fmt.Sprintf("-a%d", hcl_gui.hashcat.args.attack_mode))

			args = append(args, hcl_gui.hashcat.args.hash_file)

			args = append(args, fmt.Sprintf("--separator=%s", string(hcl_gui.hashcat.args.separator)))

			args = append(args, []string{"-D", intSliceToString(hcl_gui.hashcat.args.devices_types,",")}...)

			if len(hcl_gui.hashcat.args.outfile) > 0 {
				args = append(args, []string{"-o", hcl_gui.hashcat.args.outfile}...)
			}

			args = append(args, fmt.Sprintf("--outfile-format=%s", intSliceToString(hcl_gui.hashcat.args.outfile_format,",")))

			if len(hcl_gui.hc_extra_args.Text) > 0 {
				args = append(args, strings.Split(hcl_gui.hc_extra_args.Text, " ")...)
			}
				
			args = append(args, attack_payload...)

			return args
		}()
		session.Arguments = args
	}

	hardware_list := []string{}
	for i, _ := range hcl_gui.monitor.hardwares {
		hardware_list = append(hardware_list, fmt.Sprintf("Hardware.Mon.#%d", i+1))
	}

	wdir, _ := filepath.Split(hcl_gui.hashcat.binary_file)
	session.Process = subprocess.Subprocess{
		subprocess.SubprocessStatusNotRunning,
		wdir,
		hcl_gui.hashcat.binary_file,
		args,
		nil,
		nil,
		func(s string) {
			fmt.Fprintf(os.Stdout, "%s\n", s)
			info_line := re_info.FindStringSubmatch(s)
			if len(info_line) == 1 {
				session.UpdateJournal(info_line[0])
				return
			}
			checkpoint_enabled_line := re_checkpoint_enabled.FindStringSubmatch(s)
			if len(checkpoint_enabled_line) == 1 {
				session.UpdateJournal(checkpoint_enabled_line[0])
				return
			}
			checkpoint_disabled_line := re_checkpoint_disabled.FindStringSubmatch(s)
			if len(checkpoint_disabled_line) == 1 {
				session.UpdateJournal(checkpoint_disabled_line[0])
				return
			}
			status_line := re_status.FindStringSubmatch(s)
			if len(status_line) == 3 {
				switch status_line[1] {
				case "Status":
					if status_line[2] == "Running" && status.Text == "Initializing" {
						session.UpdateJournal("Running...")
						session.SetStatus(app, hcl_gui, SessionStatusRunning)
						pause.Enable()
						checkpoint.Enable()
						skip.Enable()
						stop.Enable()
					} else if status_line[2] == "Paused" {
						session.SetStatus(app, hcl_gui, SessionStatusPaused)
						pause.Disable()
						resume.Enable()
					} else if status_line[2] == "Quit" {
						session.SetStatus(app, hcl_gui, SessionStatusFinished)
						pause.Disable()
						checkpoint.Disable()
						skip.Disable()
						stop.Disable()
						resume.Disable()
					} else {
						if status_line[2] == "Bypass" {
							session.UpdateJournal("Skipped")
						}
						session.SetStatus(app, hcl_gui, SessionStatusRunning)
						pause.Enable()
						resume.Disable()
					}
					status.SetText(status_line[2])
				case "Hash.Name", "Hash.Type":
					hash_type.SetText(status_line[2])
				case "Hash.Target":
					hash_target.SetText(status_line[2])
				case "Guess.Queue.Base":
					guess_queue.SetText("Base: "+status_line[2])
				case "Guess.Queue.Mod":
					guess_queue.SetText(guess_queue.Text+", Mod: "+status_line[2])
				case "Guess.Queue":
					guess_queue.SetText(status_line[2])
				case "Guess.Base":
					guess_base.SetText(status_line[2])
				case "Guess.Mod":
					guess_mod.SetText(status_line[2])
				case "Guess.Mask":
					guess_mask.SetText(status_line[2])
				case "Guess.Charset":
					guess_charset.SetText(status_line[2])
				case "Progress":
					progress_line := re_progress.FindStringSubmatch(status_line[2])
					if len(progress_line) == 3 {
						progress_text.SetText(progress_line[1])
						perc, err := strconv.ParseFloat(progress_line[2], 64)
						if err != nil {
							fmt.Fprintf(os.Stderr, "can't parse progress percentage : %s\n", err)
							session.UpdateJournal("Error: can't parse progress percentage")
						} else {
							progress.SetValue(perc)
						}
					}
				case "Recovered":
					recovered_line := re_recovered.FindStringSubmatch(status_line[2])
					if len(recovered_line) == 3 {
						recovered_text.SetText(recovered_line[1])
						perc, err := strconv.ParseFloat(recovered_line[2], 64)
						if err != nil {
							fmt.Fprintf(os.Stderr, "can't parse recovered percentage : %s\n", err)
							session.UpdateJournal("Error: can't parse recovered percentage")
						} else {
							recovered.SetValue(perc)
						}
					}
				case "Time.Started":
					time_started.SetText(status_line[2])
				case "Time.Estimated":
					time_estimated.SetText(status_line[2])
				case "Speed.#1":
					if (!multiple_devices) {
						speed_line := re_speed.FindStringSubmatch(status_line[2])
						if len(speed_line) == 1 {
							speed.SetText(speed_line[0])
						}
					}
				case "Speed.#*":
					multiple_devices = true
					speed_line := re_speed.FindStringSubmatch(status_line[2])
					if len(speed_line) == 1 {
						speed.SetText(speed_line[0])
					}
				default:
					if StringArrayIncludes(hardware_list, status_line[1]) {
						hwmon := re_hwmon.FindStringSubmatch(s)
						var hwmon_id int64 = -1
						var hwmon_temp string = "N/A"
						var hwmon_fan float64 = 0
						var hwmon_util float64 = 0
						var hwmon_core string = "N/A"
						var hwmon_mem string = "N/A"
						var hwmon_bus string = "N/A"
						for i := 0; i < len(hwmon)-1; i++ {
							switch hwmon[i] {
							case "Temp":
								i++
								hwmon_temp = hwmon[i]
							case "Fan":
								i++
								hwmon_fan_tmp, err := strconv.ParseFloat(hwmon[i], 64)
								if err == nil {
									hwmon_fan = hwmon_fan_tmp
								}
							case "Util":
								i++
								hwmon_util_tmp, err := strconv.ParseFloat(hwmon[i], 64)
								if err == nil {
									hwmon_util = hwmon_util_tmp
								}
							case "Core":
								i++
								hwmon_core = hwmon[i]
							case "Mem":
								i++
								hwmon_mem = hwmon[i]
							case "Bus":
								i++
								hwmon_bus = hwmon[i]
							default:
								hwmon_id_tmp, err := strconv.ParseInt(hwmon[i], 10, 64)
								if err == nil {
									hwmon_id = hwmon_id_tmp - 1
								}
							}
						}
						if (hwmon_id >= 0) {
							hcl_gui.monitor.hardwares[hwmon_id].temp.SetText(hwmon_temp)
							hcl_gui.monitor.hardwares[hwmon_id].fan.SetValue(hwmon_fan)
							hcl_gui.monitor.hardwares[hwmon_id].util.SetValue(hwmon_util)
							hcl_gui.monitor.hardwares[hwmon_id].core.SetText(hwmon_core)
							hcl_gui.monitor.hardwares[hwmon_id].mem.SetText(hwmon_mem)
							hcl_gui.monitor.hardwares[hwmon_id].bus.SetText(hwmon_bus)
						}
					}
				}
			}
		},
		func(s string) {
			fmt.Fprintf(os.Stderr, "%s\n", s)
			if len(s) > 0 {
				status.SetText("An error occurred")
				session.SetStatus(app, hcl_gui, SessionStatusFailed)
				session.UpdateJournal("Error: "+re_ansi.ReplaceAllString(s, ""))
			}
		},
		func() {
			started = true
			start.Disable()
			refresh.Enable()
			pause.Enable()
			resume.Disable()
			checkpoint.Enable()
			skip.Enable()
			stop.Enable()
			terminate.Enable()
			session.UpdateJournal("Started.")
			status.SetText("Initializing")
			session.SetStatus(app, hcl_gui, SessionStatusRunning)
			session.UpdateJournal("Initializing...")
		},
		func() {
			refresh.Disable()
			pause.Disable()
			resume.Disable()
			checkpoint.Disable()
			skip.Disable()
			stop.Disable()
			terminate.Disable()
			if terminated {
				status.SetText("Terminated")
			} else if session.Status == SessionStatusRunning {
				status.SetText("Exited")
			}
			if session.Status != SessionStatusFailed {
				session.SetStatus(app, hcl_gui, SessionStatusFinished)
			}
			session.UpdateJournal("Ended.")
			go AutoStart(hcl_gui)
		},
	}
	session.Content = widget.NewVBox(
		widget.NewGroup("Main",
			container.NewGridWithColumns(2,
				(func() *widget.Check {
					c := widget.NewCheck("Enable Notifications", func(check bool){
						session.Notifications_Enabled = check
					})
					c.SetChecked(session.Notifications_Enabled)
					return c
				})(),
				(func() *widget.Entry {
					e := widget.NewEntry()
					e.OnChanged = func(s string) {
						extract_priority := re_priority.FindStringSubmatch(s)
						if len(extract_priority) == 2 {
							session.Priority, _ = strconv.Atoi(extract_priority[1])
						} else {
							session.Priority = 0
						}
						e.SetText(fmt.Sprintf("Priority: %d", session.Priority))
					}
					e.SetText(fmt.Sprintf("Priority: %d", session.Priority))
					return e
				})(),
			),
		),
		widget.NewGroup("Control",
			widget.NewVBox(
				fyne.NewContainerWithLayout(layout.NewGridLayout(4),
					start,
					refresh,
					pause,
					resume,
				),
				fyne.NewContainerWithLayout(layout.NewGridLayout(3),
					info,
					checkpoint,
					skip,
				),
				fyne.NewContainerWithLayout(layout.NewGridLayout(3),
					stop,
					terminate,
					terminate_n_close,
				),
			),
		),
		widget.NewGroup("Stats",
			fyne.NewContainerWithLayout(layout.NewGridLayout(2),
				widget.NewVBox(
					widget.NewLabelWithStyle("Status:", fyne.TextAlignLeading, fyne.TextStyle{}),
					widget.NewScrollContainer(status),
				),
				widget.NewVBox(
					widget.NewLabelWithStyle("Speed:", fyne.TextAlignLeading, fyne.TextStyle{}),
					widget.NewScrollContainer(speed),
				),
				widget.NewVBox(
					widget.NewLabelWithStyle("Hash Type:", fyne.TextAlignLeading, fyne.TextStyle{}),
					widget.NewScrollContainer(hash_type),
				),
				widget.NewVBox(
					widget.NewLabelWithStyle("Hash Target:", fyne.TextAlignLeading, fyne.TextStyle{}),
					widget.NewScrollContainer(hash_target),
				),
				widget.NewVBox(
					widget.NewLabelWithStyle("Time Started:", fyne.TextAlignLeading, fyne.TextStyle{}),
					widget.NewScrollContainer(time_started),
				),
				widget.NewVBox(
					widget.NewLabelWithStyle("Time Estimated:", fyne.TextAlignLeading, fyne.TextStyle{}),
					widget.NewScrollContainer(time_estimated),
				),
				widget.NewVBox(
					widget.NewScrollContainer(
						widget.NewHBox(
							widget.NewLabelWithStyle("Progress:", fyne.TextAlignLeading, fyne.TextStyle{}),
							progress_text,
						),
					),
					progress,
				),
				widget.NewVBox(
					widget.NewScrollContainer(
						widget.NewHBox(
							widget.NewLabelWithStyle("Recovered:", fyne.TextAlignLeading, fyne.TextStyle{}),
							recovered_text,
						),
					),
					recovered,
				),
			),
		),
		widget.NewGroup("Attack Details",
			widget.NewVBox(
				widget.NewScrollContainer(
					widget.NewHBox(
						widget.NewLabelWithStyle("Guess Queue:", fyne.TextAlignLeading, fyne.TextStyle{}),
						guess_queue,
					),
				),
			),
			fyne.NewContainerWithLayout(layout.NewGridLayout(2),
				widget.NewVBox(
					widget.NewScrollContainer(
						widget.NewHBox(
							widget.NewLabelWithStyle("Guess Base:", fyne.TextAlignLeading, fyne.TextStyle{}),
							guess_base,
						),
					),
				),
				widget.NewVBox(
					widget.NewScrollContainer(
						widget.NewHBox(
							widget.NewLabelWithStyle("Guess Mod:", fyne.TextAlignLeading, fyne.TextStyle{}),
							guess_mod,
						),
					),
				),
			),
			widget.NewVBox(
				widget.NewScrollContainer(
					widget.NewHBox(
						widget.NewLabelWithStyle("Guess Mask:", fyne.TextAlignLeading, fyne.TextStyle{}),
						guess_mask,
					),
				),
			),
			widget.NewVBox(
				widget.NewScrollContainer(
					widget.NewHBox(
						widget.NewLabelWithStyle("Guess Charset:", fyne.TextAlignLeading, fyne.TextStyle{}),
						guess_charset,
					),
				),
			),
		),
		widget.NewGroup("Journal",
			container.NewMax(
				container.NewScroll(
					session.Journal,
				),
			),
		),
	)
	hcl_gui.sessions[session.Id] = session
	CalculateSessionsStatusStats(hcl_gui)
	tasks_Refresh(hcl_gui)
	hcl_gui.tabs_container.SelectTab(hcl_gui.tabs["Tasks"])
	hcl_gui.tasks_tree.Select(session.Id)
	go AutoStart(hcl_gui)
}

func CalculateSessionsStatusStats(hcl_gui *hcl_gui) {
	tmp_count_sessions_running := 0
	tmp_count_sessions_queued := 0
	tmp_count_sessions_paused := 0
	tmp_count_sessions_finished := 0
	tmp_count_sessions_failed := 0
	for _, session := range hcl_gui.sessions {
		switch session.Status {
		case SessionStatusRunning:
			tmp_count_sessions_running++
		case SessionStatusQueued:
			tmp_count_sessions_queued++
		case SessionStatusPaused:
			tmp_count_sessions_paused++
		case SessionStatusFinished:
			tmp_count_sessions_finished++
		case SessionStatusFailed:
			tmp_count_sessions_failed++
		}
	}
	hcl_gui.count_sessions_running = tmp_count_sessions_running
	hcl_gui.count_sessions_queued = tmp_count_sessions_queued
	hcl_gui.count_sessions_paused = tmp_count_sessions_paused
	hcl_gui.count_sessions_finished = tmp_count_sessions_finished
	hcl_gui.count_sessions_failed = tmp_count_sessions_failed
}

func SortTasksByPriorityAndDate(hcl_gui *hcl_gui) []*Session {
	sessions := []*Session{}
	for _, session := range hcl_gui.sessions {
		sessions = append(sessions, session)
	}
	sort.Sort(SessionIdSorter(sessions)) // THe date is in the session id
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].Priority > sessions[j].Priority
	})
	return sessions
}

func AutoStart(hcl_gui *hcl_gui) {
	if hcl_gui.autostart_sessions {
		if hcl_gui.count_sessions_running < hcl_gui.max_active_sessions {
			for _, session := range SortTasksByPriorityAndDate(hcl_gui) {
				if session.Process.Status == subprocess.SubprocessStatusNotRunning {
					session.Start()
					time.Sleep(2*time.Second)
					break
				}
			}
		}
	}
}
