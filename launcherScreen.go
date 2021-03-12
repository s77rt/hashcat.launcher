package hashcatlauncher

import (
	"os"
	"io/ioutil"
	"bufio"
	"path/filepath"
	"fmt"
	"time"
	"errors"
	"strings"
	"strconv"
	"net/url"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"github.com/s77rt/hashcat.launcher/pkg/xfyne/xwidget"
)

func launcherScreen(app fyne.App, hcl_gui *hcl_gui) fyne.CanvasObject {
	// Basic Static Configs...
	hashcat_img := canvas.NewImageFromResource(hcl_gui.Icon)
	hashcat_img.SetMinSize(fyne.Size{50, 50})
	
	hcl_gui.hc_hash_file = widget.NewSelect([]string{}, func(s string) {
		_, file := filepath.Split(s)
		hcl_gui.hc_hash_file.Selected = file
		hcl_gui.hc_hash_file.Refresh()
		set_hash_file(hcl_gui, s)
	})

	outfile := widget.NewCheck("Output:", func(bool){})
	outfile.SetChecked(true)
	hcl_gui.hc_outfile = widget.NewSelect([]string{}, func(s string) {
		outfile.SetChecked(true)
		_, file := filepath.Split(s)
		hcl_gui.hc_outfile.Selected = file
		hcl_gui.hc_outfile.Refresh()
		set_outfile(hcl_gui, s)
	})
	outfile.OnChanged = func(check bool) {
		if check {
			hcl_gui.hc_outfile.Selected = "(Select one)"
		} else {
			hcl_gui.hc_outfile.Selected = "None"
		}
		hcl_gui.hc_outfile.Refresh()
		set_outfile(hcl_gui, "")
	}
	outfile_format := widget.NewSelect([]string{"hash[:salt]", "plain", "hash[:salt]:plain", "hash[:salt]:plain:hex_plain", "hash[:salt]:plain:crack_pos", "hash[:salt]:plain:hex_plain:crack_pos"}, func(s string) {set_outfile_format(hcl_gui, s)})
	outfile_format.SetSelected("hash[:salt]:plain")

	hcl_gui.hc_attack_mode = widget.NewSelect([]string{}, func(s string) {set_attack_mode(hcl_gui, s)})

	hcl_gui.hc_hash_type = xwidget.NewSelector("(Select one)", func(){customselect_hashtype(hcl_gui)})

	hcl_gui.hc_separator = widget.NewEntry()
	hcl_gui.hc_separator.SetText("   :   ")
	hcl_gui.hc_separator.CursorColumn = 4
	hcl_gui.hc_separator.OnCursorChanged = func() {
		hcl_gui.hc_separator.CursorColumn = 4
	}
	hcl_gui.hc_separator.OnChanged = func(s string) {
		if len(s) == 8 {
			text := "   "+string(s[4])+"   "
			hcl_gui.hc_separator.SetText(text)
		} else if len(s) != 7 || s == "       " {
			hcl_gui.hc_separator.SetText("   :   ")
		}
		set_separator(hcl_gui)
	}
	set_separator(hcl_gui)

	hcl_gui.hc_disable_monitor = widget.NewCheck("Disable Monitor", func(check bool){set_disable_monitor(hcl_gui, check)})

	hcl_gui.hc_temp_abort = widget.NewSelect([]string{"60", "70", "75", "80", "85", "90", "95", "100"}, func(s string) {set_temp_abort(hcl_gui, s)})

	hcl_gui.hc_devices_types = widget.NewSelect([]string{"GPU", "CPU", "FPGA", "GPU+CPU", "GPU+FPGA", "CPU+FPGA", "All"}, func(s string) {set_devices_types(hcl_gui, s)})

	hcl_gui.hc_wordload_profiles = widget.NewSelect([]string{"Low", "Default", "High", "Nightmare"}, func(s string) {set_workload_profile(hcl_gui, s)})

	remove_found_hashes := widget.NewCheck("Remove found hashes", func(check bool){set_remove_found_hashes(hcl_gui, check)})
	disable_potfile := widget.NewCheck("Disable Pot File", func(check bool){set_disable_potfile(hcl_gui, check)})
	ignore_usernames := widget.NewCheck("Ignore Usernames", func(check bool){set_ignore_usernames(hcl_gui, check)})
	optimized_kernel := widget.NewCheck("Enable optimized kernel", func(check bool){set_optimized_kernel(hcl_gui, check)})
	slower_candidate := widget.NewCheck("Enable slower candidate generators", func(check bool){set_slower_candidate(hcl_gui, check)})
	disable_self_test := widget.NewCheck("Disable self-test (Not Recommended)", func(check bool){set_disable_self_test(hcl_gui, check)})
	force := widget.NewCheck("Ignore warnings (Not Recommended)", func(check bool){set_force(hcl_gui, check)})
	optimized_kernel.SetChecked(true)

	// Notifications
	enable_notifications := false
	enable_notifications_check := widget.NewCheck("Enable Notifications", func(check bool){
		enable_notifications = check
	})
	enable_notifications_check.SetChecked(true)

	// Priority
	priority := 0
	priority_entry := widget.NewEntry()
	priority_entry.OnChanged = func(s string) {
		extract_priority := re_priority.FindStringSubmatch(s)
		if len(extract_priority) == 2 {
			priority, _ = strconv.Atoi(extract_priority[1])
		} else {
			priority = 0
		}
		priority_entry.SetText(fmt.Sprintf("%d", priority))
	}
	priority_entry.SetText(fmt.Sprintf("%d", priority))

	// Mode Configs start from here...

	// Dictionary Mode
	dictionaries := []string{}
	dictionaries_stats := widget.NewLabel("Selected 0 dictionaries")
	dictionaries_entry := widget.NewMultiLineEntry()
	dictionaries_entry.OnChanged = func(s string){
		dictionaries = []string{}
		files_list := strings.Split(strings.Replace(s, "\r\n", "\n", -1), "\n")
		for _, file := range files_list {
			hcl_gui.data.dictionaries.AddFile(file);
			if file_exists := File_Exists(file); file_exists == true {
				if !(StringArrayIncludes(dictionaries, file)) {
					dictionaries = append(dictionaries, file)
				}
			}
		}
		valid_files := len(dictionaries)
		dictionaries_entry.Text = strings.Join(dictionaries, "\n")+"\n"
		dictionaries_stats.SetText(fmt.Sprintf("Selected %d dictionaries", valid_files))
	}
	// Dictionaries Rules
	dictionaries_rule1 := ""
	dictionaries_rule2 := ""
	dictionaries_rule3 := ""
	dictionaries_rule4 := ""
	// Rule 1
	var dictionaries_rule1_select *xwidget.Selector
	dictionaries_rule1_select = xwidget.NewSelector("(Select one)", func(){customselect_rules(hcl_gui, dictionaries_rule1_select)})
	dictionaries_rule1_check := widget.NewCheck("Rule 1:", func(bool){})
	dictionaries_rule1_select.OnChanged = func(s string) {
		dictionaries_rule1_check.SetChecked(true)
		_, file := filepath.Split(s)
		dictionaries_rule1_select.Selected = file
		dictionaries_rule1_select.Refresh()
		dictionaries_rule1 = s
	}
	dictionaries_rule1_check.OnChanged = func(check bool) {
		if check {
			dictionaries_rule1_select.Selected = "(Select one)"
		} else {
			dictionaries_rule1_select.Selected = "None"
		}
		dictionaries_rule1_select.Refresh()
		dictionaries_rule1 = ""
	}
	dictionaries_rule1_select.Selected = "None"
	// Rule 2
	var dictionaries_rule2_select *xwidget.Selector
	dictionaries_rule2_select = xwidget.NewSelector("(Select one)", func(){customselect_rules(hcl_gui, dictionaries_rule2_select)})
	dictionaries_rule2_check := widget.NewCheck("Rule 2:", func(bool){})
	dictionaries_rule2_select.OnChanged = func(s string) {
		dictionaries_rule2_check.SetChecked(true)
		_, file := filepath.Split(s)
		dictionaries_rule2_select.Selected = file
		dictionaries_rule2_select.Refresh()
		dictionaries_rule2 = s
	}
	dictionaries_rule2_check.OnChanged = func(check bool) {
		if check {
			dictionaries_rule2_select.Selected = "(Select one)"
		} else {
			dictionaries_rule2_select.Selected = "None"
		}
		dictionaries_rule2_select.Refresh()
		dictionaries_rule2 = ""
	}
	dictionaries_rule2_select.Selected = "None"
	// Rule 3
	var dictionaries_rule3_select *xwidget.Selector
	dictionaries_rule3_select = xwidget.NewSelector("(Select one)", func(){customselect_rules(hcl_gui, dictionaries_rule3_select)})
	dictionaries_rule3_check := widget.NewCheck("Rule 3:", func(bool){})
	dictionaries_rule3_select.OnChanged = func(s string) {
		dictionaries_rule3_check.SetChecked(true)
		_, file := filepath.Split(s)
		dictionaries_rule3_select.Selected = file
		dictionaries_rule3_select.Refresh()
		dictionaries_rule3 = s
	}
	dictionaries_rule3_check.OnChanged = func(check bool) {
		if check {
			dictionaries_rule3_select.Selected = "(Select one)"
		} else {
			dictionaries_rule3_select.Selected = "None"
		}
		dictionaries_rule3_select.Refresh()
		dictionaries_rule3 = ""
	}
	dictionaries_rule3_select.Selected = "None"
	// Rule 4
	var dictionaries_rule4_select *xwidget.Selector
	dictionaries_rule4_select = xwidget.NewSelector("(Select one)", func(){customselect_rules(hcl_gui, dictionaries_rule4_select)})
	dictionaries_rule4_check := widget.NewCheck("Rule 4:", func(bool){})
	dictionaries_rule4_select.OnChanged = func(s string) {
		dictionaries_rule4_check.SetChecked(true)
		_, file := filepath.Split(s)
		dictionaries_rule4_select.Selected = file
		dictionaries_rule4_select.Refresh()
		dictionaries_rule4 = s
	}
	dictionaries_rule4_check.OnChanged = func(check bool) {
		if check {
			dictionaries_rule4_select.Selected = "(Select one)"
		} else {
			dictionaries_rule4_select.Selected = "None"
		}
		dictionaries_rule4_select.Refresh()
		dictionaries_rule4 = ""
	}
	dictionaries_rule4_select.Selected = "None"

	hcl_gui.hc_dictionary_attack_conf = container.NewVBox(
		container.NewGridWithColumns(2,
			widget.NewCard("", "Dictionaries",
				container.NewGridWithRows(2,
					widget.NewButton("Select Dictionaries", func(){
						customselect_dictionaries_dictionarylist(hcl_gui, &dictionaries, dictionaries_entry)
					}),
					container.NewGridWithRows(2,
						dictionaries_stats,
						container.NewGridWithColumns(4,
							widget.NewButton("Load Config", func(){
								go func() {
									file, err := NewFileOpen(hcl_gui)
									if err == nil {
										data, err := ioutil.ReadFile(file)
										if err != nil {
											fmt.Fprintf(os.Stderr, "can't load config: %s\n", err)
											dialog.ShowError(err, hcl_gui.window)
										} else {
											dictionaries_entry.SetText(string(data))
										}
									}
								}()
							}),
							widget.NewButton("Save Config", func(){
								go func() {
									file, err := NewFileSave(hcl_gui)
									if err == nil {
										f, err := os.Create(file)
										if err != nil {
											fmt.Fprintf(os.Stderr, "can't save config: %s\n", err)
											dialog.ShowError(err, hcl_gui.window)
										} else {
											defer f.Close()
											w := bufio.NewWriter(f)
											_, err := w.WriteString(dictionaries_entry.Text)
											if err != nil {
												fmt.Fprintf(os.Stderr, "can't save config: %s\n", err)
												dialog.ShowError(err, hcl_gui.window)
											} else {
												w.Flush()
											}
										}
									}
								}()
							}),
							layout.NewSpacer(),
							widget.NewButton("Clear All", func(){dictionaries_entry.SetText("")}),
						),
					),
				),
			),
			widget.NewCard("", "Rules",
				container.NewVBox(
					container.New(layout.NewFormLayout(),
						dictionaries_rule1_check,
						dictionaries_rule1_select,
						dictionaries_rule2_check,
						dictionaries_rule2_select,
						dictionaries_rule3_check,
						dictionaries_rule3_select,
						dictionaries_rule4_check,
						dictionaries_rule4_select,
					),
				),
			),
		),
	)
	hcl_gui.hc_dictionary_attack_conf.Hide()

	// Combinator Mode
	// Left
	combinator_left_wordlist := ""
	var combinator_left_wordlist_select *xwidget.Selector
	combinator_left_wordlist_select = xwidget.NewSelector("(Select one)", func(){customselect_dictionaries(hcl_gui, combinator_left_wordlist_select)})
	combinator_left_wordlist_select.OnChanged = func(s string){
		_, file := filepath.Split(s)
		combinator_left_wordlist_select.Selected = file
		combinator_left_wordlist_select.Refresh()
		combinator_left_wordlist = s
	}
	combinator_left_rule := ""
	combinator_left_rule_entry := widget.NewEntry()
	combinator_left_rule_entry.SetText("c")
	combinator_left_rule_entry.Disable()
	combinator_left_rule_entry.OnChanged = func(s string) {
		combinator_left_rule = s
	}
	combinator_left_rule_check := widget.NewCheck("Left Rule:", func(check bool){
		if check {
			combinator_left_rule_entry.Enable()
			combinator_left_rule = combinator_left_rule_entry.Text
		} else {
			combinator_left_rule_entry.Disable()
			combinator_left_rule = ""
		}
	})
	// Right
	combinator_right_wordlist := ""
	var combinator_right_wordlist_select *xwidget.Selector
	combinator_right_wordlist_select = xwidget.NewSelector("(Select one)", func(){customselect_dictionaries(hcl_gui, combinator_right_wordlist_select)})
	combinator_right_wordlist_select.OnChanged = func(s string){
		_, file := filepath.Split(s)
		combinator_right_wordlist_select.Selected = file
		combinator_right_wordlist_select.Refresh()
		combinator_right_wordlist = s
	}
	combinator_right_rule := ""
	combinator_right_rule_entry := widget.NewEntry()
	combinator_right_rule_entry.SetText("$!")
	combinator_right_rule_entry.Disable()
	combinator_right_rule_entry.OnChanged = func(s string) {
		combinator_right_rule = s
	}
	combinator_right_rule_check := widget.NewCheck("Right Rule:", func(check bool){
		if check {
			combinator_right_rule_entry.Enable()
			combinator_right_rule = combinator_right_rule_entry.Text
		} else {
			combinator_right_rule_entry.Disable()
			combinator_right_rule = ""
		}
	})

	hcl_gui.hc_combinator_attack_conf = container.NewVBox(
		container.NewGridWithColumns(2,
			widget.NewCard("", "Wordlists",
				container.NewVBox(
					container.New(layout.NewFormLayout(),
						widget.NewLabelWithStyle("Left Wordlist:", fyne.TextAlignLeading, fyne.TextStyle{}),
						combinator_left_wordlist_select,
						widget.NewLabelWithStyle("Right Wordlist:", fyne.TextAlignLeading, fyne.TextStyle{}),
						combinator_right_wordlist_select,
					),
				),
			),
			widget.NewCard("", "Rules",
				container.NewVBox(
					container.New(layout.NewFormLayout(),
						combinator_left_rule_check,
						combinator_left_rule_entry,
						combinator_right_rule_check,
						combinator_right_rule_entry,
					),
				),
			),
		),
	)
	hcl_gui.hc_combinator_attack_conf.Hide()

	// Mask Mode
	mask := ""
	mask_length := widget.NewLabelWithStyle("0", fyne.TextAlignLeading, fyne.TextStyle{})
	mask_entry := widget.NewEntry()
	mask_entry.SetPlaceHolder("?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a")
	mask_entry.OnChanged = func(s string) {
		if mask_length.Text == "[F]" {
			mask = ""
			mask_entry.SetText("")
			mask_length.SetText("0")
		} else {
			l := get_mask_length(s);
			mask_length.SetText(fmt.Sprintf("%d", l))
			mask = s
		}
	}
	mask_increment_min := ""
	mask_increment_min_entry := widget.NewEntry()
	mask_increment_min_entry.SetText("1")
	mask_increment_min_entry.Disable()
	mask_increment_min_entry.OnChanged = func(s string) {
		s = s[:min(len(s), 2)]
		if _, err := strconv.Atoi(s); err != nil {
			s = ""
		}
		mask_increment_min_entry.SetText(s)
		mask_increment_min = s
	}
	mask_increment_max := ""
	mask_increment_max_entry := widget.NewEntry()
	mask_increment_max_entry.SetText("16")
	mask_increment_max_entry.Disable()
	mask_increment_max_entry.OnChanged = func(s string) {
		s = s[:min(len(s), 2)]
		if _, err := strconv.Atoi(s); err != nil {
			s = ""
		}
		mask_increment_max_entry.SetText(s)
		mask_increment_max = s
	}
	mask_increment_check := widget.NewCheck("Enable mask increment mode:", func(check bool){
		if check {
			mask_increment_min_entry.Enable()
			mask_increment_min = mask_increment_min_entry.Text
			mask_increment_max_entry.Enable()
			mask_increment_max = mask_increment_max_entry.Text
		} else {
			mask_increment_min_entry.Disable()
			mask_increment_min = ""
			mask_increment_max_entry.Disable()
			mask_increment_max = ""
		}
	})
	mask_customcharset1 := ""
	mask_customcharset1_entry := widget.NewEntry()
	mask_customcharset1_entry.SetText("?l?u?d")
	mask_customcharset1_entry.Disable()
	mask_customcharset1_entry.OnChanged = func(s string) {
		mask_customcharset1 = s
	}
	mask_customcharset1_check := widget.NewCheck("Custom charset 1:", func(check bool){
		if check {
			mask_customcharset1_entry.Enable()
			mask_customcharset1 = mask_customcharset1_entry.Text
		} else {
			mask_customcharset1_entry.Disable()
			mask_customcharset1 = ""
		}
	})
	mask_customcharset2 := ""
	mask_customcharset2_entry := widget.NewEntry()
	mask_customcharset2_entry.SetText("?l?d")
	mask_customcharset2_entry.Disable()
	mask_customcharset2_entry.OnChanged = func(s string) {
		mask_customcharset2 = s
	}
	mask_customcharset2_check := widget.NewCheck("Custom charset 2:", func(check bool){
		if check {
			mask_customcharset2_entry.Enable()
			mask_customcharset2 = mask_customcharset2_entry.Text
		} else {
			mask_customcharset2_entry.Disable()
			mask_customcharset2 = ""
		}
	})
	mask_customcharset3 := ""
	mask_customcharset3_entry := widget.NewEntry()
	mask_customcharset3_entry.SetText("?d?s")
	mask_customcharset3_entry.Disable()
	mask_customcharset3_entry.OnChanged = func(s string) {
		mask_customcharset3 = s
	}
	mask_customcharset3_check := widget.NewCheck("Custom charset 3:", func(check bool){
		if check {
			mask_customcharset3_entry.Enable()
			mask_customcharset3 = mask_customcharset3_entry.Text
		} else {
			mask_customcharset3_entry.Disable()
			mask_customcharset3 = ""
		}
	})
	mask_customcharset4 := ""
	mask_customcharset4_entry := widget.NewEntry()
	mask_customcharset4_entry.SetText("ABCDabcd1234")
	mask_customcharset4_entry.Disable()
	mask_customcharset4_entry.OnChanged = func(s string) {
		mask_customcharset4 = s
	}
	mask_customcharset4_check := widget.NewCheck("Custom charset 4:", func(check bool){
		if check {
			mask_customcharset4_entry.Enable()
			mask_customcharset4 = mask_customcharset4_entry.Text
		} else {
			mask_customcharset4_entry.Disable()
			mask_customcharset4 = ""
		}
	})

	hcl_gui.hc_mask_attack_conf = container.NewVBox(
		container.NewGridWithColumns(2,
			widget.NewCard("", "Mask",
				container.NewVBox(
					container.New(layout.NewFormLayout(),
						widget.NewLabelWithStyle("Mask:", fyne.TextAlignLeading, fyne.TextStyle{}),
						mask_entry,
					),
					container.NewGridWithColumns(3,
						container.NewHBox(
							widget.NewLabelWithStyle("Length:", fyne.TextAlignLeading, fyne.TextStyle{}),
							mask_length,
						),
						layout.NewSpacer(),
						widget.NewButtonWithIcon("Load .hcmask file", theme.FolderOpenIcon(), func(){
							go func() {
								file, err := NewFileOpen(hcl_gui)
								if err == nil {
									mask_entry.SetText("[hcmask file]")
									mask_length.SetText("[F]")
									mask = file
								}
							}()
						}),
					),
					container.NewHBox(
						mask_increment_check,
						mask_increment_min_entry,
						widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{}),
						mask_increment_max_entry,
					),
				),
			),
			widget.NewCard("", "Custom charsets",
				container.NewVBox(
					container.New(layout.NewFormLayout(),
						mask_customcharset1_check,
						mask_customcharset1_entry,
						mask_customcharset2_check,
						mask_customcharset2_entry,
						mask_customcharset3_check,
						mask_customcharset3_entry,
						mask_customcharset4_check,
						mask_customcharset4_entry,
					),
				),
			),
		),
	)
	hcl_gui.hc_mask_attack_conf.Hide()

	// Hybrid1 Mode
	// Left
	hybrid1_left_wordlist := ""
	var hybrid1_left_wordlist_select *xwidget.Selector
	hybrid1_left_wordlist_select = xwidget.NewSelector("(Select one)", func(){customselect_dictionaries(hcl_gui, hybrid1_left_wordlist_select)})
	hybrid1_left_wordlist_select.OnChanged = func(s string){
		_, file := filepath.Split(s)
		hybrid1_left_wordlist_select.Selected = file
		hybrid1_left_wordlist_select.Refresh()
		hybrid1_left_wordlist = s
	}
	hybrid1_left_rule := ""
	hybrid1_left_rule_entry := widget.NewEntry()
	hybrid1_left_rule_entry.SetText("^e^h^t")
	hybrid1_left_rule_entry.Disable()
	hybrid1_left_rule_entry.OnChanged = func(s string) {
		hybrid1_left_rule = s
	}
	hybrid1_left_rule_check := widget.NewCheck("Rule:", func(check bool){
		if check {
			hybrid1_left_rule_entry.Enable()
			hybrid1_left_rule = hybrid1_left_rule_entry.Text
		} else {
			hybrid1_left_rule_entry.Disable()
			hybrid1_left_rule = ""
		}
	})
	// Right
	hybrid1_right_mask := ""
	hybrid1_right_mask_length := widget.NewLabelWithStyle("0", fyne.TextAlignLeading, fyne.TextStyle{})
	hybrid1_right_mask_entry := widget.NewEntry()
	hybrid1_right_mask_entry.SetPlaceHolder("?d?d?d?d")
	hybrid1_right_mask_entry.OnChanged = func(s string) {
		if hybrid1_right_mask_length.Text == "[F]" {
			hybrid1_right_mask = ""
			hybrid1_right_mask_entry.SetText("")
			hybrid1_right_mask_length.SetText("0")
		} else {
			l := get_mask_length(s);
			hybrid1_right_mask_length.SetText(fmt.Sprintf("%d", l))
			hybrid1_right_mask = s
		}
	}
	hybrid1_right_mask_increment_min := ""
	hybrid1_right_mask_increment_min_entry := widget.NewEntry()
	hybrid1_right_mask_increment_min_entry.SetText("1")
	hybrid1_right_mask_increment_min_entry.Disable()
	hybrid1_right_mask_increment_min_entry.OnChanged = func(s string) {
		s = s[:min(len(s), 2)]
		if _, err := strconv.Atoi(s); err != nil {
			s = ""
		}
		hybrid1_right_mask_increment_min_entry.SetText(s)
		hybrid1_right_mask_increment_min = s
	}
	hybrid1_right_mask_increment_max := ""
	hybrid1_right_mask_increment_max_entry := widget.NewEntry()
	hybrid1_right_mask_increment_max_entry.SetText("4")
	hybrid1_right_mask_increment_max_entry.Disable()
	hybrid1_right_mask_increment_max_entry.OnChanged = func(s string) {
		s = s[:min(len(s), 2)]
		if _, err := strconv.Atoi(s); err != nil {
			s = ""
		}
		hybrid1_right_mask_increment_max_entry.SetText(s)
		hybrid1_right_mask_increment_max = s
	}
	hybrid1_right_mask_increment_check := widget.NewCheck("Enable mask increment mode:", func(check bool){
		if check {
			hybrid1_right_mask_increment_min_entry.Enable()
			hybrid1_right_mask_increment_min = hybrid1_right_mask_increment_min_entry.Text
			hybrid1_right_mask_increment_max_entry.Enable()
			hybrid1_right_mask_increment_max = hybrid1_right_mask_increment_max_entry.Text
		} else {
			hybrid1_right_mask_increment_min_entry.Disable()
			hybrid1_right_mask_increment_min = ""
			hybrid1_right_mask_increment_max_entry.Disable()
			hybrid1_right_mask_increment_max = ""
		}
	})

	hcl_gui.hc_hybrid1_attack_conf = container.NewVBox(
		container.NewGridWithColumns(2,
			widget.NewCard("", "Left",
				container.NewVBox(
					container.New(layout.NewFormLayout(),
						widget.NewLabelWithStyle("Wordlist:", fyne.TextAlignLeading, fyne.TextStyle{}),
						hybrid1_left_wordlist_select,
						hybrid1_left_rule_check,
						hybrid1_left_rule_entry,
					),
				),
			),
			widget.NewCard("", "Right",
				container.NewVBox(
					container.New(layout.NewFormLayout(),
						widget.NewLabelWithStyle("Mask:", fyne.TextAlignLeading, fyne.TextStyle{}),
						hybrid1_right_mask_entry,
					),
					container.NewGridWithColumns(3,
						container.NewHBox(
							widget.NewLabelWithStyle("Length:", fyne.TextAlignLeading, fyne.TextStyle{}),
							hybrid1_right_mask_length,
						),
						layout.NewSpacer(),
						widget.NewButtonWithIcon("Load .hcmask file", theme.FolderOpenIcon(), func(){
							go func() {
								file, err := NewFileOpen(hcl_gui)
								if err == nil {
									hybrid1_right_mask_entry.SetText("[hcmask file]")
									hybrid1_right_mask_length.SetText("[F]")
									hybrid1_right_mask = file
								}
							}()
						}),
					),
					container.NewHBox(
						hybrid1_right_mask_increment_check,
						hybrid1_right_mask_increment_min_entry,
						widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{}),
						hybrid1_right_mask_increment_max_entry,
					),
				),
			),
		),
	)
	hcl_gui.hc_hybrid1_attack_conf.Hide()

	// Hybrid2 Mode
	// Left
	hybrid2_left_mask := ""
	hybrid2_left_mask_length := widget.NewLabelWithStyle("0", fyne.TextAlignLeading, fyne.TextStyle{})
	hybrid2_left_mask_entry := widget.NewEntry()
	hybrid2_left_mask_entry.SetPlaceHolder("?u?l?l?l?l?l?l?l")
	hybrid2_left_mask_entry.OnChanged = func(s string) {
		if hybrid2_left_mask_length.Text == "[F]" {
			hybrid2_left_mask = ""
			hybrid2_left_mask_entry.SetText("")
			hybrid2_left_mask_length.SetText("0")
		} else {
			l := get_mask_length(s);
			hybrid2_left_mask_length.SetText(fmt.Sprintf("%d", l))
			hybrid2_left_mask = s
		}
	}
	hybrid2_left_mask_increment_min := ""
	hybrid2_left_mask_increment_min_entry := widget.NewEntry()
	hybrid2_left_mask_increment_min_entry.SetText("1")
	hybrid2_left_mask_increment_min_entry.Disable()
	hybrid2_left_mask_increment_min_entry.OnChanged = func(s string) {
		s = s[:min(len(s), 2)]
		if _, err := strconv.Atoi(s); err != nil {
			s = ""
		}
		hybrid2_left_mask_increment_min_entry.SetText(s)
		hybrid2_left_mask_increment_min = s
	}
	hybrid2_left_mask_increment_max := ""
	hybrid2_left_mask_increment_max_entry := widget.NewEntry()
	hybrid2_left_mask_increment_max_entry.SetText("4")
	hybrid2_left_mask_increment_max_entry.Disable()
	hybrid2_left_mask_increment_max_entry.OnChanged = func(s string) {
		s = s[:min(len(s), 2)]
		if _, err := strconv.Atoi(s); err != nil {
			s = ""
		}
		hybrid2_left_mask_increment_max_entry.SetText(s)
		hybrid2_left_mask_increment_max = s
	}
	hybrid2_left_mask_increment_check := widget.NewCheck("Enable mask increment mode:", func(check bool){
		if check {
			hybrid2_left_mask_increment_min_entry.Enable()
			hybrid2_left_mask_increment_min = hybrid2_left_mask_increment_min_entry.Text
			hybrid2_left_mask_increment_max_entry.Enable()
			hybrid2_left_mask_increment_max = hybrid2_left_mask_increment_max_entry.Text
		} else {
			hybrid2_left_mask_increment_min_entry.Disable()
			hybrid2_left_mask_increment_min = ""
			hybrid2_left_mask_increment_max_entry.Disable()
			hybrid2_left_mask_increment_max = ""
		}
	})
	// Right
	hybrid2_right_wordlist := ""
	var hybrid2_right_wordlist_select *xwidget.Selector
	hybrid2_right_wordlist_select = xwidget.NewSelector("(Select one)", func(){customselect_dictionaries(hcl_gui, hybrid2_right_wordlist_select)})
	hybrid2_right_wordlist_select.OnChanged = func(s string){
		_, file := filepath.Split(s)
		hybrid2_right_wordlist_select.Selected = file
		hybrid2_right_wordlist_select.Refresh()
		hybrid2_right_wordlist = s
	}
	hybrid2_right_rule := ""
	hybrid2_right_rule_entry := widget.NewEntry()
	hybrid2_right_rule_entry.SetText("$1 $2 $3 $!")
	hybrid2_right_rule_entry.Disable()
	hybrid2_right_rule_entry.OnChanged = func(s string) {
		hybrid2_right_rule = s
	}
	hybrid2_right_rule_check := widget.NewCheck("Rule:", func(check bool){
		if check {
			hybrid2_right_rule_entry.Enable()
			hybrid2_right_rule = hybrid2_right_rule_entry.Text
		} else {
			hybrid2_right_rule_entry.Disable()
			hybrid2_right_rule = ""
		}
	})

	hcl_gui.hc_hybrid2_attack_conf = container.NewVBox(
		container.NewGridWithColumns(2,
			widget.NewCard("", "Left",
				container.NewVBox(
					container.New(layout.NewFormLayout(),
						widget.NewLabelWithStyle("Mask:", fyne.TextAlignLeading, fyne.TextStyle{}),
						hybrid2_left_mask_entry,
					),
					container.NewGridWithColumns(3,
						container.NewHBox(
							widget.NewLabelWithStyle("Length:", fyne.TextAlignLeading, fyne.TextStyle{}),
							hybrid2_left_mask_length,
						),
						layout.NewSpacer(),
						widget.NewButtonWithIcon("Load .hcmask file", theme.FolderOpenIcon(), func(){
							go func() {
								file, err := NewFileOpen(hcl_gui)
								if err == nil {
									hybrid2_left_mask_entry.SetText("[hcmask file]")
									hybrid2_left_mask_length.SetText("[F]")
									hybrid2_left_mask = file
								}
							}()
						}),
					),
					container.NewHBox(
						hybrid2_left_mask_increment_check,
						hybrid2_left_mask_increment_min_entry,
						widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{}),
						hybrid2_left_mask_increment_max_entry,
					),
				),
			),
			widget.NewCard("", "Right",
				container.NewVBox(
					container.New(layout.NewFormLayout(),
						widget.NewLabelWithStyle("Wordlist:", fyne.TextAlignLeading, fyne.TextStyle{}),
						hybrid2_right_wordlist_select,
						hybrid2_right_rule_check,
						hybrid2_right_rule_entry,
					),
				),
			),
		),
	)
	hcl_gui.hc_hybrid2_attack_conf.Hide()
	
	// Attack Payload Builder
	build_attack_payload := func() []string {
		attack_payload := []string{}
		// Mode Related Check
		switch hcl_gui.hashcat.args.attack_mode {
		// Dictionary Mode
		case hashcat_attack_mode_Dictionary:
			// Dictionaries
			if len(dictionaries) > 0 {
				attack_payload = append(attack_payload, dictionaries...)
			} else {
				err := errors.New("You must add at least one dictionary")
				fmt.Fprintf(os.Stderr, "%s\n", err)
				dialog.ShowError(err, hcl_gui.window)
				return []string{}
			}
			// Rules
			if len(dictionaries_rule1) > 0 {
				attack_payload = append(attack_payload, []string{"-r", dictionaries_rule1}...)
			}
			if len(dictionaries_rule2) > 0 {
				attack_payload = append(attack_payload, []string{"-r", dictionaries_rule2}...)
			}
			if len(dictionaries_rule3) > 0 {
				attack_payload = append(attack_payload, []string{"-r", dictionaries_rule3}...)
			}
			if len(dictionaries_rule4) > 0 {
				attack_payload = append(attack_payload, []string{"-r", dictionaries_rule4}...)
			}
		case hashcat_attack_mode_Combinator:
			// Wordlists
			if len(combinator_left_wordlist) > 0 && len(combinator_right_wordlist) > 0 {
				attack_payload = append(attack_payload, []string{combinator_left_wordlist, combinator_right_wordlist}...)
			} else {
				err := errors.New("You must add both of left and right wordlists")
				fmt.Fprintf(os.Stderr, "%s\n", err)
				dialog.ShowError(err, hcl_gui.window)
				return []string{}
			}
			// Rules
			if len(combinator_left_rule) > 0 {
				attack_payload = append(attack_payload, []string{"-j", combinator_left_rule}...)
			}
			if len(combinator_right_rule) > 0 {
				attack_payload = append(attack_payload, []string{"-k", combinator_right_rule}...)
			}
		case hashcat_attack_mode_Mask:
			// Custom Charsets
			if len(mask_customcharset1) > 0 {
				attack_payload = append(attack_payload, []string{"-1", mask_customcharset1}...)
			}
			if len(mask_customcharset2) > 0 {
				attack_payload = append(attack_payload, []string{"-2", mask_customcharset2}...)
			}
			if len(mask_customcharset3) > 0 {
				attack_payload = append(attack_payload, []string{"-3", mask_customcharset3}...)
			}
			if len(mask_customcharset4) > 0 {
				attack_payload = append(attack_payload, []string{"-4", mask_customcharset4}...)
			}
			// Mask
			if len(mask) > 0 {
				attack_payload = append(attack_payload, strings.Split(mask, " ")...)
			} else {
				err := errors.New("You must specify a mask")
				fmt.Fprintf(os.Stderr, "%s\n", err)
				dialog.ShowError(err, hcl_gui.window)
				return []string{}
			}
			// Increment
			if len(mask_increment_min) > 0 && len(mask_increment_max) > 0 {
				attack_payload = append(attack_payload, []string{"-i", fmt.Sprintf("--increment-min=%s", mask_increment_min), fmt.Sprintf("--increment-max=%s", mask_increment_max)}...)
			}
		case hashcat_attack_mode_Hybrid1:
			// Custom Charsets
			if len(mask_customcharset1) > 0 {
				attack_payload = append(attack_payload, []string{"-1", mask_customcharset1}...)
			}
			if len(mask_customcharset2) > 0 {
				attack_payload = append(attack_payload, []string{"-2", mask_customcharset2}...)
			}
			if len(mask_customcharset3) > 0 {
				attack_payload = append(attack_payload, []string{"-3", mask_customcharset3}...)
			}
			if len(mask_customcharset4) > 0 {
				attack_payload = append(attack_payload, []string{"-4", mask_customcharset4}...)
			}
			// Left
			if len(hybrid1_left_wordlist) > 0 {
				attack_payload = append(attack_payload, hybrid1_left_wordlist)
			} else {
				err := errors.New("You must add a wordlist")
				fmt.Fprintf(os.Stderr, "%s\n", err)
				dialog.ShowError(err, hcl_gui.window)
				return []string{}
			}
			if len(hybrid1_left_rule) > 0 {
				attack_payload = append(attack_payload, []string{"-j", hybrid1_left_rule}...)
			}
			// Right
			if len(hybrid1_right_mask) > 0 {
				attack_payload = append(attack_payload, strings.Split(hybrid1_right_mask, " ")...)
			} else {
				err := errors.New("You must specify a mask")
				fmt.Fprintf(os.Stderr, "%s\n", err)
				dialog.ShowError(err, hcl_gui.window)
				return []string{}
			}
			if len(hybrid1_right_mask_increment_min) > 0 && len(hybrid1_right_mask_increment_max) > 0 {
				attack_payload = append(attack_payload, []string{"-i", fmt.Sprintf("--increment-min=%s", hybrid1_right_mask_increment_min), fmt.Sprintf("--increment-max=%s", hybrid1_right_mask_increment_max)}...)
			}
		case hashcat_attack_mode_Hybrid2:
			// Custom Charsets
			if len(mask_customcharset1) > 0 {
				attack_payload = append(attack_payload, []string{"-1", mask_customcharset1}...)
			}
			if len(mask_customcharset2) > 0 {
				attack_payload = append(attack_payload, []string{"-2", mask_customcharset2}...)
			}
			if len(mask_customcharset3) > 0 {
				attack_payload = append(attack_payload, []string{"-3", mask_customcharset3}...)
			}
			if len(mask_customcharset4) > 0 {
				attack_payload = append(attack_payload, []string{"-4", mask_customcharset4}...)
			}
			// Left
			if len(hybrid2_left_mask) > 0 {
				attack_payload = append(attack_payload, strings.Split(hybrid2_left_mask, " ")...)
			} else {
				err := errors.New("You must specify a mask")
				fmt.Fprintf(os.Stderr, "%s\n", err)
				dialog.ShowError(err, hcl_gui.window)
				return []string{}
			}
			if len(hybrid2_left_mask_increment_min) > 0 && len(hybrid2_left_mask_increment_max) > 0 {
				attack_payload = append(attack_payload, []string{"-i", fmt.Sprintf("--increment-min=%s", hybrid2_left_mask_increment_min), fmt.Sprintf("--increment-max=%s", hybrid2_left_mask_increment_max)}...)
			}
			// Right
			if len(hybrid2_right_wordlist) > 0 {
				attack_payload = append(attack_payload, hybrid2_right_wordlist)
			} else {
				err := errors.New("You must add a wordlist")
				fmt.Fprintf(os.Stderr, "%s\n", err)
				dialog.ShowError(err, hcl_gui.window)
				return []string{}
			}
			if len(hybrid2_right_rule) > 0 {
				attack_payload = append(attack_payload, []string{"-j", hybrid2_right_rule}...)
			}
		}
		return attack_payload
	}

	// Run hashcat
	run_hashcat_btn := widget.NewButtonWithIcon("Create Task", theme.ContentAddIcon(), func(){
		// Basic Configs Check
		if len(hcl_gui.hashcat.args.hash_file) == 0 {
			err := errors.New("You must select a hash file")
			fmt.Fprintf(os.Stderr, "%s\n", err)
			dialog.ShowError(err, hcl_gui.window)
			return
		}
		if hcl_gui.hashcat.args.hash_type == -1 {
			err := errors.New("You must select a hash type")
			fmt.Fprintf(os.Stderr, "%s\n", err)
			dialog.ShowError(err, hcl_gui.window)
			return
		}
		if outfile.Checked && len(hcl_gui.hashcat.args.outfile) == 0 {
			err := errors.New("Output is enabled, but no outfile has been specified")
			fmt.Fprintf(os.Stderr, "%s\n", err)
			dialog.ShowError(err, hcl_gui.window)
			return
		}
		attack_payload := build_attack_payload()
		if len(attack_payload) > 0 {
			NewSession(app, hcl_gui, -1, "", nil, attack_payload, enable_notifications, priority)
		}
	})
	run_hashcat_btn.Importance = widget.HighImportance

	run_hashcat_restore_btn := widget.NewButtonWithIcon("Restore Task", theme.HistoryIcon(), func(){
		var modal *widget.PopUp
		c := container.NewVBox(
			container.NewHBox(
				widget.NewLabelWithStyle("Restore Task", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				layout.NewSpacer(),
				widget.NewButtonWithIcon("", theme.HelpIcon(), func() {
					dialog.ShowInformation("Help", "Select which task to restore", hcl_gui.window)
				}),
				widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
					hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){})
					modal.Hide()
				}),
			),
			container.New(layout.NewGridWrapLayout(fyne.Size{900, 700}),
				(func() fyne.CanvasObject {
					var data []string
					var restore_file *RestoreFile
					var list *widget.List
					var selected_id widget.ListItemID
					data = GetRestoreFiles(hcl_gui)
					session_name_label := widget.NewLabel("N/A")
					time_label := widget.NewLabel("N/A")
					task_id_label := widget.NewLabel("N/A")
					argv_label := widget.NewLabel("N/A")
					restore_btn := widget.NewButtonWithIcon("Restore", theme.HistoryIcon(), func() {
						hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){})
						modal.Hide()
						NewSession(app, hcl_gui, restore_file.Task_id, restore_file.Session_name, restore_file, []string{}, enable_notifications, priority)
					})
					restore_btn.Disable()
					delete_btn := widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
						dialog.ShowConfirm(
							"Delete Restore File?",
							"Are you sure you want to delete the selected restore file?",
							func (confirm bool) {
								if confirm {
									list.Unselect(selected_id)
									restore_file.Delete()
									data = GetRestoreFiles(hcl_gui)
									list.Refresh()
								}
							},
							hcl_gui.window,
						)
					})
					delete_btn.Disable()
					vbox := container.NewVBox(
						container.NewGridWithRows(2,
							container.NewVBox(
								widget.NewLabelWithStyle("Session Name", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
								container.NewHScroll(session_name_label),
								widget.NewLabelWithStyle("Time", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
								container.NewHScroll(time_label),
								widget.NewLabelWithStyle("Task ID", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
								container.NewHScroll(task_id_label),
								widget.NewLabelWithStyle("Arguments", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
							),
							container.NewScroll(argv_label),
						),
						container.NewPadded(
							container.NewVBox(
								widget.NewLabelWithStyle("Actions", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
								restore_btn,
								container.NewHBox(
									layout.NewSpacer(),
									delete_btn,
								),
							),
						),
					)
					list = widget.NewList(
						func() int {
							return len(data)
						},
						func() fyne.CanvasObject {
							return container.New(layout.NewHBoxLayout(), widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
						},
						func(id widget.ListItemID, item fyne.CanvasObject) {
							item.(*fyne.Container).Objects[1].(*widget.Label).SetText(filepath.Base(data[id]))
						},
					)
					list.OnSelected = func(id widget.ListItemID) {
						selected_id = id
						restore_file = ReadRestoreFile(hcl_gui, data[id])
						session_name_label.SetText(restore_file.Session_name)
						if restore_file.Time > 0 {
							time_label.SetText(fmt.Sprintf("%s", time.Unix(0, restore_file.Time)))
						} else {
							time_label.SetText("N/A")
						}
						if restore_file.Task_id > 0 {
							task_id_label.SetText(fmt.Sprintf("%d", restore_file.Task_id))
						} else {
							task_id_label.SetText("N/A")
						}
						argv_label.SetText(strings.Join(restore_file.GetArguments(), "\n"))
						restore_btn.Enable()
						delete_btn.Enable()
					}
					list.OnUnselected = func(id widget.ListItemID) {
						session_name_label.SetText("N/A")
						time_label.SetText("N/A")
						task_id_label.SetText("N/A")
						argv_label.SetText("N/A")
						restore_btn.Disable()
						delete_btn.Disable()
					}
					return container.NewHSplit(container.NewMax(list), container.NewMax(vbox))
				})(),
			),
		)
		hcl_gui.window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
			if key.Name == fyne.KeyEscape {
				hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){})
				modal.Hide()
			}
		})
		modal = widget.NewModalPopUp(c, hcl_gui.window.Canvas())
		modal.Show()
	})

	return container.NewVBox(
		container.NewPadded(
			container.NewGridWithColumns(2,
				container.NewMax(
					banner(),
					container.NewHBox(
						hashcat_img,
					),
				),
				container.NewGridWithColumns(2,
					container.NewGridWithColumns(2,
						widget.NewButtonWithIcon("Help", theme.QuestionIcon(), func() {
							dialog.ShowInformation(
								"Help",
								strings.Join(
									[]string{
										"The usage is pretty simple;",
										"1. Select Hash File => This is the file that contains your hashes",
										"2. Select Hash Type => The hashing algorithm",
										"3. Configure Attack => The attack mode",
										"4. Create Task",
										"That's all!",
										"",
										"Got a question? Ask the community on forum.hashkiller.io",
										"Something is wrong or missing? Report on our github repo (see \"About\" tab)",
									},
									"\n",
								),
								hcl_gui.window,
							)
						}),
						run_hashcat_restore_btn,
					),
					run_hashcat_btn,
				),
			),
		),
		widget.NewCard("Target", "choose target hash and type",
			container.NewGridWithColumns(2,
				container.NewGridWithRows(2,
					container.NewGridWithColumns(4,
						widget.NewLabelWithStyle("Hash File:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
						layout.NewSpacer(),
						widget.NewButtonWithIcon("Clipboard", theme.ContentPasteIcon(), func(){
							data := hcl_gui.window.Clipboard().Content()
							if len(data) > 0 {
								pwd, _ := os.Getwd()
								file := filepath.Join(pwd, "clipboard.txt")
								f, err := os.Create(file)
								if err != nil {
									fmt.Fprintf(os.Stderr, "clipboard: %s\n", err)
									dialog.ShowError(err, hcl_gui.window)
								} else {
									defer f.Close()
									w := bufio.NewWriter(f)
									_, err := w.WriteString(data)
									if err != nil {
										fmt.Fprintf(os.Stderr, "clipboard: %s\n", err)
										dialog.ShowError(err, hcl_gui.window)
									} else {
										w.Flush()
										hcl_gui.hc_hash_file.Options = append([]string{file}, hcl_gui.hc_hash_file.Options[:min(len(hcl_gui.hc_hash_file.Options), 4)]...)
										hcl_gui.hc_hash_file.SetSelected(file)
									}
								}
							} else {
								err := errors.New("clipboard is empty")
								fmt.Fprintf(os.Stderr, "clipboard: %s\n", err)
								dialog.ShowError(err, hcl_gui.window)
							}
						}),
						widget.NewButtonWithIcon("Browse", theme.FolderOpenIcon(), func(){
							go func() {
								file, err := NewFileOpen(hcl_gui)
								if err == nil {
									hcl_gui.hc_hash_file.Options = append([]string{file}, hcl_gui.hc_hash_file.Options[:min(len(hcl_gui.hc_hash_file.Options), 4)]...)
									hcl_gui.hc_hash_file.SetSelected(file)
								}
							}()
						}),
					),
					hcl_gui.hc_hash_file,
				),
				container.NewGridWithRows(2,
					widget.NewLabelWithStyle("Hash Type:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
					hcl_gui.hc_hash_type,
				),
			),
		),
		widget.NewCard("Attack", "configure attack",
			container.NewVBox(
				container.New(layout.NewFormLayout(),
					widget.NewLabelWithStyle("Mode:", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
					hcl_gui.hc_attack_mode,
				),
				hcl_gui.hc_dictionary_attack_conf,
				hcl_gui.hc_combinator_attack_conf,
				hcl_gui.hc_mask_attack_conf,
				hcl_gui.hc_hybrid1_attack_conf,
				hcl_gui.hc_hybrid2_attack_conf,
			),
		),
		container.NewGridWithColumns(2,
			widget.NewCard("Options", "hashcat options",
				container.NewVBox(
					container.NewGridWithRows(4,
						optimized_kernel,
						slower_candidate,
						remove_found_hashes,
						disable_potfile,
						ignore_usernames,
						disable_self_test,
						force,
					),
				),
			),
			container.NewGridWithColumns(2,
				widget.NewCard("Devices", "choose which devices to use",
					container.NewVBox(
						container.NewGridWithColumns(2,
							container.NewVBox(
								widget.NewLabelWithStyle("Devices:", fyne.TextAlignLeading, fyne.TextStyle{}),
								hcl_gui.hc_devices_types,
							),
							container.NewVBox(
								widget.NewLabelWithStyle("Workload Profile:", fyne.TextAlignLeading, fyne.TextStyle{}),
								hcl_gui.hc_wordload_profiles,
							),
						),
						widget.NewButton("Info", func(){
							var modal *widget.PopUp
							var info_box *widget.TextGrid
							var copy_btn *widget.Button
							var close_btn *widget.Button
							info_box = widget.NewTextGrid()
							info := "Obtaining info..."
							info_box.SetText(info)
							copy_btn = widget.NewButton("Copy", func(){
								hcl_gui.window.Clipboard().SetContent(info_box.Text())
								copy_btn.SetText("Copied!")
							})
							copy_btn.Disable()
							close_btn = widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
								hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){})
								modal.Hide()
							})
							close_btn.Disable()
							c := container.NewVBox(
								container.NewHBox(
									widget.NewLabelWithStyle("Info", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
									layout.NewSpacer(),
									close_btn,
								),
								container.NewPadded(
									container.New(layout.NewGridWrapLayout(fyne.Size{600, 500}),
										container.NewScroll(
											info_box,
										),
									),
								),
								container.NewHBox(
									copy_btn,
								),
							)
							hcl_gui.window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
								if key.Name == fyne.KeyEscape {
									hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){})
									modal.Hide()
								}
							})
							modal = widget.NewModalPopUp(c, hcl_gui.window.Canvas())
							modal.Show()
							go func() {
								info = get_devices_info(hcl_gui)
								info_box.SetText(info)
								copy_btn.Enable()
								close_btn.Enable()
							}()
						}),
						widget.NewButton("Benchmark", func(){
							if hcl_gui.hashcat.args.hash_type == -1 {
								err := errors.New("You must select a hash type")
								fmt.Fprintf(os.Stderr, "%s\n", err)
								dialog.ShowError(err, hcl_gui.window)
								return
							}
							var modal *widget.PopUp
							var benchmark_box *widget.TextGrid
							var copy_btn *widget.Button
							var close_btn *widget.Button
							copy_btn = widget.NewButton("Copy", func(){
								hcl_gui.window.Clipboard().SetContent(benchmark_box.Text())
								copy_btn.SetText("Copied!")
							})
							copy_btn.Disable()
							close_btn = widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
								hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){})
								modal.Hide()
							})
							close_btn.Disable()
							benchmark_box = widget.NewTextGrid()
							benchmark := "Benchmarking..."
							benchmark_box.SetText(benchmark)
							c := container.NewVBox(
								container.NewHBox(
									widget.NewLabelWithStyle("Benchmark", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
									layout.NewSpacer(),
									close_btn,
								),
								container.NewPadded(
									container.New(layout.NewGridWrapLayout(fyne.Size{600, 500}),
										container.NewScroll(
											benchmark_box,
										),
									),
								),
								container.NewHBox(
									copy_btn,
								),
							)
							hcl_gui.window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
								if key.Name == fyne.KeyEscape {
									hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){})
									modal.Hide()
								}
							})
							modal = widget.NewModalPopUp(c, hcl_gui.window.Canvas())
							modal.Show()
							go func() {
								benchmark = get_benchmark(hcl_gui)
								benchmark_box.SetText(benchmark)
								copy_btn.Enable()
								close_btn.Enable()
							}()
						}),
					),
				),
				widget.NewCard("Monitor", "monitoring options",
					container.NewVBox(
						hcl_gui.hc_disable_monitor,
						container.New(layout.NewFormLayout(),
							widget.NewLabelWithStyle("Temp Abort (°C):", fyne.TextAlignLeading, fyne.TextStyle{}),
							hcl_gui.hc_temp_abort,
						),
					),
				),
			),
		),
		container.NewGridWithColumns(2,
			widget.NewCard("Output", "output file and format",
				container.NewVBox(
					container.NewGridWithColumns(2,
						container.New(layout.NewFormLayout(),
							outfile,
							hcl_gui.hc_outfile,
						),
						widget.NewButtonWithIcon("Set Output File", theme.DocumentSaveIcon(), func(){
							go func() {
								file, err := NewFileSave(hcl_gui)
								if err == nil {
									outfile.SetChecked(true)
									hcl_gui.hc_outfile.Options = append([]string{file}, hcl_gui.hc_outfile.Options[:min(len(hcl_gui.hc_outfile.Options), 4)]...)
									hcl_gui.hc_outfile.SetSelected(file)
								}
							}()
						}),
					),
					container.New(layout.NewFormLayout(),
						widget.NewLabelWithStyle("Format:", fyne.TextAlignLeading, fyne.TextStyle{}),
						outfile_format,
					),
				),
			),
			widget.NewCard("Features", "task options and features",
				container.NewVBox(
					enable_notifications_check,
					container.New(layout.NewFormLayout(),
						widget.NewLabelWithStyle("Priority:", fyne.TextAlignLeading, fyne.TextStyle{}),
						priority_entry,
					),
				),
			),
		),
		layout.NewSpacer(),
		container.NewHBox(
			layout.NewSpacer(),
			func () fyne.CanvasObject {
				contribute_url, _ := url.Parse("https://github.com/s77rt/hashcat.launcher/issues/new/")
				w := widget.NewHyperlinkWithStyle("Report a bug / Request a feature", contribute_url, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
				return w
			}(),
		),
	)
}