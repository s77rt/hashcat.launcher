package hashcatlauncher

import (
	"os"
	"io"
	"strings"
	"strconv"
	"time"
	"fmt"
	"errors"
	"path/filepath"
	"runtime"
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"fyne.io/fyne/dialog"
	"github.com/s77rt/hashcat.launcher/pkg/subprocess"
)

type Journal struct {
	msgs [3]*widget.Label
}

func (journal *Journal) UpdateJournal(new_msg string) {
	journal.msgs[2].SetText(journal.msgs[1].Text)
	journal.msgs[1].SetText(journal.msgs[0].Text)
	journal.msgs[0].SetText(fmt.Sprintf("%s %s", time.Now().Format("2006-01-02 15:04:05"), new_msg))
}

type Session struct {
	Id string
	Process subprocess.Subprocess
	Journal Journal
	Tab *widget.TabItem
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
func (session *Session) Quit() {
	if runtime.GOOS == "windows" {
		session.Process.PostKey(0x51)
	} else {
		io.WriteString(session.Process.Stdin_stream, "q")
	}
}

func (session *Session) SetTabTextStatus(hcl_gui *hcl_gui, text string) {
	session.Tab.Text = "["+text+"]"
	hcl_gui.tabs.Refresh()
}

func (hcl_gui *hcl_gui) RemoveSession(session *Session) {
    for i, v := range hcl_gui.sessions {
        if v == session {
            hcl_gui.sessions = append(hcl_gui.sessions[:i], hcl_gui.sessions[i+1:]...)
        }
    }	
}

func newSession(hcl_gui *hcl_gui, attack_payload []string) {
	hcl_gui.count_all_sessions++
	session := &Session{}
	session.Id = fmt.Sprintf("hcl_%d_%s", hcl_gui.count_all_sessions, RandomString(7))
	session.Journal = Journal{[3]*widget.Label{
		widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace:true}),
		widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace:true}),
		widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace:true}),
	}}

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

	start := widget.NewButton("Start", func(){
		if hcl_gui.count_active_sessions < hcl_gui.max_active_sessions {
			session.Start()
		} else {
			dialog.ShowError(errors.New("Max Active Sessions Reached!"), hcl_gui.window)
		}
	})
	refresh := widget.NewButton("Refresh", func(){
		session.Refresh()
	})
	refresh.Disable()
	var pause, resume *widget.Button
	pause = widget.NewButton("Pause", func(){
		session.Pause()
		go func(){
			time.Sleep(100*time.Millisecond)
			session.Refresh()
		}()
		session.Journal.UpdateJournal("Paused")
	})
	pause.Disable()
	resume = widget.NewButton("Resume", func(){
		session.Resume()
		go func(){
			time.Sleep(100*time.Millisecond)
			session.Refresh()
		}()
		session.Journal.UpdateJournal("Resumed")
	})
	resume.Disable()
	stop := widget.NewButton("Stop", func(){
		session.Quit()
		session.Journal.UpdateJournal("Graceful Stop")
	})
	stop.Disable()
	terminate := widget.NewButton("Terminate", func(){
		session.Quit()
		session.Journal.UpdateJournal("Forceful Stop")
		session.Process.Kill()
		terminated = true
	})
	terminate.Disable()
	terminate_n_close := widget.NewButton("Terminate & Close", func(){
		if started {
			session.Quit()
		}
		session.Journal.UpdateJournal("Forceful Stop")
		session.Process.Kill()
		terminated = true
		var tab_index int
		if tab_index = hcl_gui.tabs.CurrentTabIndex()-1; tab_index < 4 {
			tab_index = 3
		} 
		hcl_gui.tabs.SelectTabIndex(tab_index)
		hcl_gui.tabs.Remove(session.Tab)
		hcl_gui.RemoveSession(session)
	})

	args := func() []string {
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
				session.Journal.UpdateJournal(info_line[0])
				return
			}
			status_line := re_status.FindStringSubmatch(s)
			if len(status_line) == 3 {
				switch status_line[1] {
				case "Status":
					if status.Text == "Initializing" {
						session.Journal.UpdateJournal("Running...")
						session.SetTabTextStatus(hcl_gui, "Running")
						pause.Enable()
						stop.Enable()
					} else if status_line[2] == "Paused" {
						session.SetTabTextStatus(hcl_gui, "Paused")
						pause.Disable()
						resume.Enable()
					} else {
						if status_line[2] != "Quit" {
							session.SetTabTextStatus(hcl_gui, "Running")
							pause.Enable()
							resume.Disable()
						} else {
							session.SetTabTextStatus(hcl_gui, "Stopped")
							pause.Disable()
							stop.Disable()
							resume.Disable()
						}
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
							session.Journal.UpdateJournal("Error: can't parse progress percentage")
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
							session.Journal.UpdateJournal("Error: can't parse recovered percentage")
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
				session.SetTabTextStatus(hcl_gui, "Error")
				session.Journal.UpdateJournal("Error: "+re_ansi.ReplaceAllString(s, ""))
			}
		},
		func() {
			hcl_gui.count_active_sessions++
			started = true
			start.Disable()
			refresh.Enable()
			pause.Enable()
			resume.Disable()
			stop.Enable()
			terminate.Enable()
			session.Journal.UpdateJournal("Started.")
			status.SetText("Initializing")
			session.SetTabTextStatus(hcl_gui, "Initializing")
			session.Journal.UpdateJournal("Initializing...")
		},
		func() {
			hcl_gui.count_active_sessions--
			refresh.Disable()
			pause.Disable()
			resume.Disable()
			stop.Disable()
			terminate.Disable()
			if terminated {
				status.SetText("Terminated")
				session.SetTabTextStatus(hcl_gui, "Terminated")
			} else if status.Text == "Initializing" || status.Text == "Running" {
				status.SetText("Exited")
				session.SetTabTextStatus(hcl_gui, "Exited")
			} else if session.Tab.Text != "[Stopped]" {
				session.SetTabTextStatus(hcl_gui, "Finished")
			}
			session.Journal.UpdateJournal("Ended.")
			go AutoStart(hcl_gui)
		},
	}
	session.Tab = widget.NewTabItem("[Queued]", widget.NewVBox(
		widget.NewGroup("Control",
			widget.NewVBox(
				fyne.NewContainerWithLayout(layout.NewGridLayout(5),
					start,
					refresh,
					stop,
					terminate,
					terminate_n_close,
				),
				fyne.NewContainerWithLayout(layout.NewGridLayout(5),
					pause,
					resume,
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
			widget.NewScrollContainer(session.Journal.msgs[0]),
			widget.NewScrollContainer(session.Journal.msgs[1]),
			widget.NewScrollContainer(session.Journal.msgs[2]),
		),
	))
	hcl_gui.tabs.Append(session.Tab)
	hcl_gui.tabs.SelectTab(session.Tab)
	hcl_gui.sessions = append(hcl_gui.sessions, session)
	go AutoStart(hcl_gui)
}

func AutoStart(hcl_gui *hcl_gui) {
	if hcl_gui.autostart_sessions {
		for i := 0; hcl_gui.count_active_sessions < hcl_gui.max_active_sessions && i < len(hcl_gui.sessions); i++ {
			if session := hcl_gui.sessions[i]; session.Process.Status == subprocess.SubprocessStatusNotRunning {
				session.Start()
				time.Sleep(2*time.Second)
			}
		}
	}
}