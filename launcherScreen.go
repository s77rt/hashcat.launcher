package hashcatlauncher

import (
	"os"
	"io/ioutil"
	"bufio"
	"path/filepath"
	"fmt"
	"errors"
	"strings"
	"strconv"
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"fyne.io/fyne/dialog"
	dialog2 "github.com/sqweek/dialog"
	"github.com/s77rt/hashcat.launcher/pkg/xfyne/xwidget"
)

func fake_hash_type_selector_hack(hcl_gui *hcl_gui, hash_type_fakeselector *widget.Box, text string) {
	hash_type_fakeselector.Children = []fyne.CanvasObject{xwidget.NewSelector(text, func(){load_hash_type_selector(hcl_gui, hash_type_fakeselector)})}
	hash_type_fakeselector.Refresh()
}

func load_hash_type_selector(hcl_gui *hcl_gui, hash_type_fakeselector *widget.Box) {
	var modal *widget.PopUp
	data := widget.NewVBox()
	data.Children = []fyne.CanvasObject{widget.NewLabel("Results will appear here...")}
	search := widget.NewEntry()
	search.SetPlaceHolder("Type to search")
	search.OnChanged = func(keyword string){
		if len(keyword) >= 2 {
			go func(){
				load_hash_type_options(modal, data, hcl_gui, hash_type_fakeselector, hcl_gui.hashcat.available_hash_types, keyword)
			}()
		}
	}
	c := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{500, 600}),
		widget.NewScrollContainer(
			widget.NewVBox(
				search,
				data,
			),
		),
	)
    hcl_gui.window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
    	if key.Name == fyne.KeyEscape {
	        modal.Hide()
	        hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){})
    	}
    })
	modal = widget.NewModalPopUp(c, hcl_gui.window.Canvas())
}

func load_hash_type_options(modal *widget.PopUp, data *widget.Box, hcl_gui *hcl_gui, hash_type_fakeselector *widget.Box, items []*xwidget.SelectorOption, keyword string) {
	var children []fyne.CanvasObject
	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Label.Text), strings.ToLower(keyword)) {
			item.OnTapped = func(value string){
				set_hash_type(hcl_gui, hash_type_fakeselector, value)
				modal.Hide()
				hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){})
			}
			children = append(children, item)
		}
	}
	data.Children = children
	data.Refresh()
}

func launcherScreen(hcl_gui *hcl_gui, hash_type_fakeselector *widget.Box) fyne.CanvasObject {
	// Basic Static Configs...
	hashcat_img := canvas.NewImageFromResource(hcl_gui.Icon)
	hashcat_img.SetMinSize(fyne.Size{100, 100})
	
	hcl_gui.hc_hash_file = widget.NewSelect([]string{}, func(s string) {
		_, file := filepath.Split(s)
		hcl_gui.hc_hash_file.Selected = file
		set_hash_file(hcl_gui, s)
	})

	outfile := widget.NewCheck("Output:", func(bool){})
	outfile.SetChecked(true)
	hcl_gui.hc_outfile = widget.NewSelect([]string{}, func(s string) {
		outfile.SetChecked(true)
		_, file := filepath.Split(s)
		hcl_gui.hc_outfile.Selected = file
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
	outfile_format := widget.NewSelect([]string{"hash[:salt]:plain", "hash[:salt]:plain:hex_plain"}, func(s string) {set_outfile_format(hcl_gui, s)})
	outfile_format.SetSelected("hash[:salt]:plain")

	hcl_gui.hc_attack_mode = widget.NewSelect([]string{}, func(s string) {set_attack_mode(hcl_gui, s)})

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

	optimized_kernel := widget.NewCheck("Enable optimized kernel", func(check bool){set_optimized_kernel(hcl_gui, check)})
	slower_candidate := widget.NewCheck("Enable slower candidate generators", func(check bool){set_slower_candidate(hcl_gui, check)})
	force := widget.NewCheck("Ignore warnings (Not Recommended)", func(check bool){set_force(hcl_gui, check)})
	optimized_kernel.SetChecked(true)

	// Mode Configs start from here...

	// Dictionary Mode
	dictonaries := []string{}
	dictonaries_stats := widget.NewLabel("Loaded 0 files")
	dictonaries_list := widget.NewMultiLineEntry()
	dictonaries_list.SetPlaceHolder("Click [+] to add files... -or- Paste files pathes")
	dictonaries_list.OnChanged = func(s string){
		dictonaries = []string{}
		valid_files := 0
		files_list := strings.Split(strings.Replace(s, "\r\n", "\n", -1), "\n")
		for _, file := range files_list {
			if _, err := os.Stat(file); err == nil {
				dictonaries = append(dictonaries, file)
				valid_files++
			}
		}
		dictonaries_stats.SetText(fmt.Sprintf("Loaded %d files", valid_files))
	}
	// Dictionaries Rules
	dictonaries_rule1 := ""
	dictonaries_rule2 := ""
	dictonaries_rule3 := ""
	dictonaries_rule4 := ""
	// Rule 1
	dictonaries_rule1_select := widget.NewSelect([]string{}, func(string){})
	dictonaries_rule1_check := widget.NewCheck("Rule 1:", func(bool){})
	dictonaries_rule1_select.OnChanged = func(s string) {
		dictonaries_rule1_check.SetChecked(true)
		_, file := filepath.Split(s)
		dictonaries_rule1_select.Selected = file
		dictonaries_rule1 = s
	}
	dictonaries_rule1_check.OnChanged = func(check bool) {
		if check {
			dictonaries_rule1_select.Selected = "(Select one)"
		} else {
			dictonaries_rule1_select.Selected = "None"
		}
		dictonaries_rule1_select.Refresh()
		dictonaries_rule1 = ""
	}
	dictonaries_rule1_button := widget.NewButton("...", func(){
		file, err := dialog2.File().Title("Select Rule").Filter("Rules Files", "txt", "rule").Load()
			if err == nil {
				dictonaries_rule1_check.SetChecked(true)
				dictonaries_rule1_select.Options = append([]string{file}, dictonaries_rule1_select.Options[:min(len(dictonaries_rule1_select.Options), 4)]...)
				dictonaries_rule1_select.SetSelected(file)
			}
	})
	dictonaries_rule1_select.Selected = "None"
	// Rule 2
	dictonaries_rule2_select := widget.NewSelect([]string{}, func(string){})
	dictonaries_rule2_check := widget.NewCheck("Rule 2:", func(bool){})
	dictonaries_rule2_select.OnChanged = func(s string) {
		dictonaries_rule2_check.SetChecked(true)
		_, file := filepath.Split(s)
		dictonaries_rule2_select.Selected = file
		dictonaries_rule2 = s
	}
	dictonaries_rule2_check.OnChanged = func(check bool) {
		if check {
			dictonaries_rule2_select.Selected = "(Select one)"
		} else {
			dictonaries_rule2_select.Selected = "None"
		}
		dictonaries_rule2_select.Refresh()
		dictonaries_rule2 = ""
	}
	dictonaries_rule2_button := widget.NewButton("...", func(){
		file, err := dialog2.File().Title("Select Rule").Filter("Rules Files", "txt", "rule").Load()
			if err == nil {
				dictonaries_rule2_check.SetChecked(true)
				dictonaries_rule2_select.Options = append([]string{file}, dictonaries_rule2_select.Options[:min(len(dictonaries_rule2_select.Options), 4)]...)
				dictonaries_rule2_select.SetSelected(file)
			}
	})
	dictonaries_rule2_select.Selected = "None"
	// Rule 3
	dictonaries_rule3_select := widget.NewSelect([]string{}, func(string){})
	dictonaries_rule3_check := widget.NewCheck("Rule 3:", func(bool){})
	dictonaries_rule3_select.OnChanged = func(s string) {
		dictonaries_rule3_check.SetChecked(true)
		_, file := filepath.Split(s)
		dictonaries_rule3_select.Selected = file
		dictonaries_rule3 = s
	}
	dictonaries_rule3_check.OnChanged = func(check bool) {
		if check {
			dictonaries_rule3_select.Selected = "(Select one)"
		} else {
			dictonaries_rule3_select.Selected = "None"
		}
		dictonaries_rule3_select.Refresh()
		dictonaries_rule3 = ""
	}
	dictonaries_rule3_button := widget.NewButton("...", func(){
		file, err := dialog2.File().Title("Select Rule").Filter("Rules Files", "txt", "rule").Load()
			if err == nil {
				dictonaries_rule3_check.SetChecked(true)
				dictonaries_rule3_select.Options = append([]string{file}, dictonaries_rule3_select.Options[:min(len(dictonaries_rule3_select.Options), 4)]...)
				dictonaries_rule3_select.SetSelected(file)
			}
	})
	dictonaries_rule3_select.Selected = "None"
	// Rule 4
	dictonaries_rule4_select := widget.NewSelect([]string{}, func(string){})
	dictonaries_rule4_check := widget.NewCheck("Rule 4:", func(bool){})
	dictonaries_rule4_select.OnChanged = func(s string) {
		dictonaries_rule4_check.SetChecked(true)
		_, file := filepath.Split(s)
		dictonaries_rule4_select.Selected = file
		dictonaries_rule4 = s
	}
	dictonaries_rule4_check.OnChanged = func(check bool) {
		if check {
			dictonaries_rule4_select.Selected = "(Select one)"
		} else {
			dictonaries_rule4_select.Selected = "None"
		}
		dictonaries_rule4_select.Refresh()
		dictonaries_rule4 = ""
	}
	dictonaries_rule4_button := widget.NewButton("...", func(){
		file, err := dialog2.File().Title("Select Rule").Filter("Rules Files", "txt", "rule").Load()
			if err == nil {
				dictonaries_rule4_check.SetChecked(true)
				dictonaries_rule4_select.Options = append([]string{file}, dictonaries_rule4_select.Options[:min(len(dictonaries_rule4_select.Options), 4)]...)
				dictonaries_rule4_select.SetSelected(file)
			}
	})
	dictonaries_rule4_select.Selected = "None"

	hcl_gui.hc_dictionary_attack_conf = widget.NewVBox(
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{490, 200}),
			widget.NewVBox(
				widget.NewGroup("Dictonaries",
					fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{490, 103}),
						widget.NewScrollContainer(
							dictonaries_list,
						),
					),
					widget.NewHBox(
						spacer(7, 0),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{490, 35}),
							widget.NewVBox(
								widget.NewHBox(
									widget.NewButton("+", func(){
										file, err := dialog2.File().Title("Add Dictionary").Filter("Text Files", "txt", "dict").Load()
										if err == nil {
											dictonaries_list.SetText(dictonaries_list.Text+file+"\n")
										}
									}),
									spacer(10,0),
									widget.NewButton("Load Config", func(){
										file, err := dialog2.File().Title("Load Config").Filter("Config Files", "*").Load()
										if err == nil {
											data, err := ioutil.ReadFile(file)
											if err != nil {
												fmt.Fprintf(os.Stderr, "can't load config: %s\n", err)
												dialog.ShowError(err, hcl_gui.window)
											} else {
												dictonaries_list.SetText(string(data))
											}
										}
									}),
									widget.NewButton("Save Config", func(){
										file, err := dialog2.File().Title("Save Config").Filter("Config Files", "*").Save()
										if err == nil {
											f, err := os.Create(file)
											if err != nil {
												fmt.Fprintf(os.Stderr, "can't save config: %s\n", err)
												dialog.ShowError(err, hcl_gui.window)
											} else {
												defer f.Close()
												w := bufio.NewWriter(f)
												_, err := w.WriteString(dictonaries_list.Text)
												if err != nil {
													fmt.Fprintf(os.Stderr, "can't save config: %s\n", err)
													dialog.ShowError(err, hcl_gui.window)
												} else {
													w.Flush()
												}
											}
										}
									}),
									spacer(135,0),
									widget.NewButton("Clear All", func(){dictonaries_list.SetText("")}),
								),
							),
						),
					),
					fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{490, 35}),
						widget.NewScrollContainer(
							widget.NewHBox(
								spacer(5, 0),
								dictonaries_stats,
							),
						),
					),
				),
			),
		),
		spacer(0, 5),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{490, 140}),
			widget.NewVBox(
				widget.NewGroup("Rules",
					widget.NewHBox(
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{90, 35}),
							widget.NewVBox(
								dictonaries_rule1_check,
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{355, 35}),
							widget.NewVBox(
								dictonaries_rule1_select,
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{30, 35}),
							widget.NewVBox(
								dictonaries_rule1_button,
							),
						),
					),
					widget.NewHBox(
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{90, 35}),
							widget.NewVBox(
								dictonaries_rule2_check,
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{355, 35}),
							widget.NewVBox(
								dictonaries_rule2_select,
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{30, 35}),
							widget.NewVBox(
								dictonaries_rule2_button,
							),
						),
					),
					widget.NewHBox(
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{90, 35}),
							widget.NewVBox(
								dictonaries_rule3_check,
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{355, 35}),
							widget.NewVBox(
								dictonaries_rule3_select,
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{30, 35}),
							widget.NewVBox(
								dictonaries_rule3_button,
							),
						),
					),
					widget.NewHBox(
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{90, 35}),
							widget.NewVBox(
								dictonaries_rule4_check,
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{355, 35}),
							widget.NewVBox(
								dictonaries_rule4_select,
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{30, 35}),
							widget.NewVBox(
								dictonaries_rule4_button,
							),
						),
					),
				),
			),
		),
	)
	hcl_gui.hc_dictionary_attack_conf.Hide()

	// Combinator Mode
	// Left
	combinator_left_wordlist := ""
	combinator_left_wordlist_select := widget.NewSelect([]string{}, func(string){})
	combinator_left_wordlist_select.OnChanged = func(s string){
		_, file := filepath.Split(s)
		combinator_left_wordlist_select.Selected = file
		combinator_left_wordlist = s
	}
	combinator_left_wordlist_button := widget.NewButton("...", func(){
		file, err := dialog2.File().Title("Select Left Wordlist").Filter("Text Files", "txt", "dict").Load()
		if err == nil {
			combinator_left_wordlist_select.Options = append([]string{file}, combinator_left_wordlist_select.Options[:min(len(combinator_left_wordlist_select.Options), 4)]...)
			combinator_left_wordlist_select.SetSelected(file)
		}
	})
	combinator_left_rule := ""
	combinator_left_rule_entry := widget.NewEntry()
	combinator_left_rule_entry.SetText("c")
	combinator_left_rule_entry.Disable()
	combinator_left_rule_entry.OnChanged = func(s string) {
		s = s[:min(len(s), 25)]
		combinator_left_rule_entry.SetText(s)
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
	combinator_right_wordlist_select := widget.NewSelect([]string{}, func(string){})
	combinator_right_wordlist_select.OnChanged = func(s string){
		_, file := filepath.Split(s)
		combinator_right_wordlist_select.Selected = file
		combinator_right_wordlist = s
	}
	combinator_right_wordlist_button := widget.NewButton("...", func(){
		file, err := dialog2.File().Title("Select Right Wordlist").Filter("Text Files", "txt", "dict").Load()
		if err == nil {
			combinator_right_wordlist_select.Options = append([]string{file}, combinator_right_wordlist_select.Options[:min(len(combinator_right_wordlist_select.Options), 4)]...)
			combinator_right_wordlist_select.SetSelected(file)
		}
	})
	combinator_right_rule := ""
	combinator_right_rule_entry := widget.NewEntry()
	combinator_right_rule_entry.SetText("$!")
	combinator_right_rule_entry.Disable()
	combinator_right_rule_entry.OnChanged = func(s string) {
		s = s[:min(len(s), 25)]
		combinator_right_rule_entry.SetText(s)
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

	hcl_gui.hc_combinator_attack_conf = widget.NewVBox(
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{490, 200}),
			widget.NewVBox(
				widget.NewGroup("Wordlists",
					widget.NewHBox(
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{130, 35}),
							widget.NewVBox(
								widget.NewHBox(
									spacer(10, 0),
									widget.NewLabelWithStyle("Left Wordlist:", fyne.TextAlignLeading, fyne.TextStyle{}),
								),
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{300, 35}),
							widget.NewVBox(
								combinator_left_wordlist_select,
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{30, 35}),
							widget.NewVBox(
								combinator_left_wordlist_button,
							),
						),
					),
					widget.NewHBox(
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{130, 35}),
							widget.NewVBox(
								widget.NewHBox(
									spacer(10, 0),
									widget.NewLabelWithStyle("Right Wordlist:", fyne.TextAlignLeading, fyne.TextStyle{}),
								),
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{300, 35}),
							widget.NewVBox(
								combinator_right_wordlist_select,
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{30, 35}),
							widget.NewVBox(
								combinator_right_wordlist_button,
							),
						),
					),
				),
				spacer(0, 3),
				widget.NewGroup("Rules",
					widget.NewHBox(
						spacer(10, 0),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{130, 35}),
							widget.NewVBox(
								combinator_left_rule_check,
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{290, 35}),
							widget.NewVBox(
								combinator_left_rule_entry,
							),
						),
					),
					spacer(0, 5),
					widget.NewHBox(
						spacer(10, 0),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{130, 35}),
							widget.NewVBox(
								combinator_right_rule_check,
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{290, 35}),
							widget.NewVBox(
								combinator_right_rule_entry,
							),
						),
					),
				),
			),
		),
	)
	hcl_gui.hc_combinator_attack_conf.Hide()

	// Mask Mode
	mask := ""
	mask_entry := widget.NewEntry()
	mask_entry.SetPlaceHolder("?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a")
	mask_entry.OnChanged = func(s string) {
		s = s[:min(len(s), 35)]
		mask_entry.SetText(s)
		mask = s
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
		s = s[:min(len(s), 25)]
		mask_customcharset1_entry.SetText(s)
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
		s = s[:min(len(s), 25)]
		mask_customcharset2_entry.SetText(s)
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
		s = s[:min(len(s), 25)]
		mask_customcharset3_entry.SetText(s)
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
		s = s[:min(len(s), 25)]
		mask_customcharset4_entry.SetText(s)
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

	hcl_gui.hc_mask_attack_conf = widget.NewVBox(
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{490, 200}),
			widget.NewVBox(
				widget.NewGroup("Mask",
					widget.NewHBox(
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{100, 35}),
							widget.NewVBox(
								widget.NewHBox(
									spacer(10, 0),
									widget.NewLabelWithStyle("Mask:", fyne.TextAlignLeading, fyne.TextStyle{}),
								),
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{360, 35}),
							widget.NewVBox(
								mask_entry,
							),
						),
					),
					spacer(0, 3),
					widget.NewHBox(
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{300, 35}),
							widget.NewVBox(
								widget.NewHBox(
									spacer(15, 0),
									mask_increment_check,
								),
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{30, 35}),
							widget.NewVBox(
								mask_increment_min_entry,
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{50, 35}),
							widget.NewVBox(
								widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{}),
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{30, 35}),
							widget.NewVBox(
								mask_increment_max_entry,
							),
						),
					),
				),
				widget.NewGroup("Custom charsets",
					widget.NewHBox(
						spacer(10, 0),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{180, 35}),
							widget.NewVBox(
								mask_customcharset1_check,
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{265, 35}),
							widget.NewVBox(
								mask_customcharset1_entry,
							),
						),
					),
					spacer(0, 5),
					widget.NewHBox(
						spacer(10, 0),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{180, 35}),
							widget.NewVBox(
								mask_customcharset2_check,
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{265, 35}),
							widget.NewVBox(
								mask_customcharset2_entry,
							),
						),
					),
					spacer(0, 5),
					widget.NewHBox(
						spacer(10, 0),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{180, 35}),
							widget.NewVBox(
								mask_customcharset3_check,
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{265, 35}),
							widget.NewVBox(
								mask_customcharset3_entry,
							),
						),
					),
					spacer(0, 5),
					widget.NewHBox(
						spacer(10, 0),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{180, 35}),
							widget.NewVBox(
								mask_customcharset4_check,
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{265, 35}),
							widget.NewVBox(
								mask_customcharset4_entry,
							),
						),
					),
					spacer(0, 5),
				),
			),
		),
	)
	hcl_gui.hc_mask_attack_conf.Hide()

	// Hybrid1 Mode
	// Left
	hybrid1_left_wordlist := ""
	hybrid1_left_wordlist_select := widget.NewSelect([]string{}, func(string){})
	hybrid1_left_wordlist_select.OnChanged = func(s string){
		_, file := filepath.Split(s)
		hybrid1_left_wordlist_select.Selected = file
		hybrid1_left_wordlist = s
	}
	hybrid1_left_wordlist_button := widget.NewButton("...", func(){
		file, err := dialog2.File().Title("Select Wordlist").Filter("Text Files", "txt", "dict").Load()
		if err == nil {
			hybrid1_left_wordlist_select.Options = append([]string{file}, hybrid1_left_wordlist_select.Options[:min(len(hybrid1_left_wordlist_select.Options), 4)]...)
			hybrid1_left_wordlist_select.SetSelected(file)
		}
	})
	hybrid1_left_rule := ""
	hybrid1_left_rule_entry := widget.NewEntry()
	hybrid1_left_rule_entry.SetText("^e^h^t")
	hybrid1_left_rule_entry.Disable()
	hybrid1_left_rule_entry.OnChanged = func(s string) {
		s = s[:min(len(s), 25)]
		hybrid1_left_rule_entry.SetText(s)
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
	hybrid1_right_mask_entry := widget.NewEntry()
	hybrid1_right_mask_entry.SetPlaceHolder("?d?d?d?d")
	hybrid1_right_mask_entry.OnChanged = func(s string) {
		s = s[:min(len(s), 35)]
		hybrid1_right_mask_entry.SetText(s)
		hybrid1_right_mask = s
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

	hcl_gui.hc_hybrid1_attack_conf = widget.NewVBox(
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{490, 200}),
			widget.NewVBox(
				widget.NewGroup("Left",
					widget.NewHBox(
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{130, 35}),
							widget.NewVBox(
								widget.NewHBox(
									spacer(10, 0),
									widget.NewLabelWithStyle("Wordlist:", fyne.TextAlignLeading, fyne.TextStyle{}),
								),
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{300, 35}),
							widget.NewVBox(
								hybrid1_left_wordlist_select,
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{30, 35}),
							widget.NewVBox(
								hybrid1_left_wordlist_button,
							),
						),
					),
					widget.NewHBox(
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{130, 35}),
							widget.NewVBox(
								widget.NewHBox(
									spacer(10, 0),
									hybrid1_left_rule_check,
								),
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{300, 35}),
							widget.NewVBox(
								hybrid1_left_rule_entry,
							),
						),
					),
				),
				spacer(0, 3),
				widget.NewGroup("Right",
					widget.NewHBox(
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{100, 35}),
							widget.NewVBox(
								widget.NewHBox(
									spacer(10, 0),
									widget.NewLabelWithStyle("Mask:", fyne.TextAlignLeading, fyne.TextStyle{}),
								),
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{360, 35}),
							widget.NewVBox(
								hybrid1_right_mask_entry,
							),
						),
					),
					spacer(0, 3),
					widget.NewHBox(
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{300, 35}),
							widget.NewVBox(
								widget.NewHBox(
									spacer(15, 0),
									hybrid1_right_mask_increment_check,
								),
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{30, 35}),
							widget.NewVBox(
								hybrid1_right_mask_increment_min_entry,
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{50, 35}),
							widget.NewVBox(
								widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{}),
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{30, 35}),
							widget.NewVBox(
								hybrid1_right_mask_increment_max_entry,
							),
						),
					),
				),
			),
		),
	)
	hcl_gui.hc_hybrid1_attack_conf.Hide()

	// Hybrid2 Mode
	// Left
	hybrid2_left_mask := ""
	hybrid2_left_mask_entry := widget.NewEntry()
	hybrid2_left_mask_entry.SetPlaceHolder("?u?l?l?l?l?l?l?l")
	hybrid2_left_mask_entry.OnChanged = func(s string) {
		s = s[:min(len(s), 35)]
		hybrid2_left_mask_entry.SetText(s)
		hybrid2_left_mask = s
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
	hybrid2_right_wordlist_select := widget.NewSelect([]string{}, func(string){})
	hybrid2_right_wordlist_select.OnChanged = func(s string){
		_, file := filepath.Split(s)
		hybrid2_right_wordlist_select.Selected = file
		hybrid2_right_wordlist = s
	}
	hybrid2_right_wordlist_button := widget.NewButton("...", func(){
		file, err := dialog2.File().Title("Select Wordlist").Filter("Text Files", "txt", "dict").Load()
		if err == nil {
			hybrid2_right_wordlist_select.Options = append([]string{file}, hybrid2_right_wordlist_select.Options[:min(len(hybrid2_right_wordlist_select.Options), 4)]...)
			hybrid2_right_wordlist_select.SetSelected(file)
		}
	})
	hybrid2_right_rule := ""
	hybrid2_right_rule_entry := widget.NewEntry()
	hybrid2_right_rule_entry.SetText("$1 $2 $3 $!")
	hybrid2_right_rule_entry.Disable()
	hybrid2_right_rule_entry.OnChanged = func(s string) {
		s = s[:min(len(s), 25)]
		hybrid2_right_rule_entry.SetText(s)
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

	hcl_gui.hc_hybrid2_attack_conf = widget.NewVBox(
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{490, 200}),
			widget.NewVBox(
				widget.NewGroup("Left",
					widget.NewHBox(
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{100, 35}),
							widget.NewVBox(
								widget.NewHBox(
									spacer(10, 0),
									widget.NewLabelWithStyle("Mask:", fyne.TextAlignLeading, fyne.TextStyle{}),
								),
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{360, 35}),
							widget.NewVBox(
								hybrid2_left_mask_entry,
							),
						),
					),
					spacer(0, 3),
					widget.NewHBox(
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{300, 35}),
							widget.NewVBox(
								widget.NewHBox(
									spacer(15, 0),
									hybrid2_left_mask_increment_check,
								),
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{30, 35}),
							widget.NewVBox(
								hybrid2_left_mask_increment_min_entry,
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{50, 35}),
							widget.NewVBox(
								widget.NewLabelWithStyle("-", fyne.TextAlignCenter, fyne.TextStyle{}),
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{30, 35}),
							widget.NewVBox(
								hybrid2_left_mask_increment_max_entry,
							),
						),
					),
				),
				widget.NewGroup("Right",
					widget.NewHBox(
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{130, 35}),
							widget.NewVBox(
								widget.NewHBox(
									spacer(10, 0),
									widget.NewLabelWithStyle("Wordlist:", fyne.TextAlignLeading, fyne.TextStyle{}),
								),
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{300, 35}),
							widget.NewVBox(
								hybrid2_right_wordlist_select,
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{30, 35}),
							widget.NewVBox(
								hybrid2_right_wordlist_button,
							),
						),
					),
					widget.NewHBox(
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{130, 35}),
							widget.NewVBox(
								widget.NewHBox(
									spacer(10, 0),
									hybrid2_right_rule_check,
								),
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{300, 35}),
							widget.NewVBox(
								hybrid2_right_rule_entry,
							),
						),
					),
				),
			),
		),
	)
	hcl_gui.hc_hybrid2_attack_conf.Hide()
	
	// Run hashcat
	run_hashcat_btn := widget.NewButton("Launch hashcat !", func(){
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
		attack_payload := func() []string {
			attack_payload := []string{}
			// Mode Related Check
			switch hcl_gui.hashcat.args.attack_mode {
			// Dictionary Mode
			case hashcat_attack_mode_Dictionary:
				// Dictionaries
				if len(dictonaries) > 0 {
					attack_payload = append(attack_payload, dictonaries...)
				} else {
					err := errors.New("You must add at least one dictionary")
					fmt.Fprintf(os.Stderr, "%s\n", err)
					dialog.ShowError(err, hcl_gui.window)
					return []string{}
				}
				// Rules
				if len(dictonaries_rule1) > 0 {
					attack_payload = append(attack_payload, []string{"-r", dictonaries_rule1}...)
				}
				if len(dictonaries_rule2) > 0 {
					attack_payload = append(attack_payload, []string{"-r", dictonaries_rule2}...)
				}
				if len(dictonaries_rule3) > 0 {
					attack_payload = append(attack_payload, []string{"-r", dictonaries_rule3}...)
				}
				if len(dictonaries_rule4) > 0 {
					attack_payload = append(attack_payload, []string{"-r", dictonaries_rule4}...)
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
					attack_payload = append(attack_payload, mask)
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
					attack_payload = append(attack_payload, hybrid1_right_mask)
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
				// Left
				if len(hybrid2_left_mask) > 0 {
					attack_payload = append(attack_payload, hybrid2_left_mask)
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
		}()
		if len(attack_payload) > 0 {
			newSession(hcl_gui, attack_payload)
		}
	})

	return widget.NewVBox(
		widget.NewLabelWithStyle("Welcome to hashcat.launcher v"+Version, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		spacer(0,5),
		widget.NewHBox(
			spacer(5,0),
			fyne.NewContainerWithLayout(layout.NewCenterLayout(),
				hashcat_img,
			),
			widget.NewVBox(
				widget.NewHBox(
					fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{80, 35}),
						widget.NewVBox(
							widget.NewLabelWithStyle("Hash File:", fyne.TextAlignTrailing, fyne.TextStyle{}),
						),
					),
					fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{440, 35}),
						widget.NewVBox(
							hcl_gui.hc_hash_file,
						),
					),
					fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{100, 35}),
						widget.NewVBox(
							widget.NewButton("Clipboard", func(){
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
						),
					),
					fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{30, 35}),
						widget.NewVBox(
							widget.NewButton("...", func(){
								file, err := dialog2.File().Title("Select Hash File").Filter("Hash Files", "hash", "txt", "dat", "hccapx").Load()
								if err == nil {
									hcl_gui.hc_hash_file.Options = append([]string{file}, hcl_gui.hc_hash_file.Options[:min(len(hcl_gui.hc_hash_file.Options), 4)]...)
									hcl_gui.hc_hash_file.SetSelected(file)
								}
							}),
						),
					),
				),
				widget.NewHBox(
					fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{80, 35}),
						widget.NewVBox(
							widget.NewLabelWithStyle("Separator:", fyne.TextAlignTrailing, fyne.TextStyle{}),
						),
					),
					fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{50, 35}),
						widget.NewVBox(
							hcl_gui.hc_separator,
						),
					),
					fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{200, 35}),
						widget.NewVBox(
							widget.NewCheck("Remove found hashes", func(check bool){set_remove_found_hashes(hcl_gui, check)}),
						),
					),
					fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{150, 35}),
						widget.NewVBox(
							widget.NewCheck("Disable Pot File", func(check bool){set_disable_potfile(hcl_gui, check)}),
						),
					),
					fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{170, 35}),
						widget.NewVBox(
							widget.NewCheck("Ignore Usernames", func(check bool){set_ignore_usernames(hcl_gui, check)}),
						),
					),
				),
				widget.NewHBox(
					fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{80, 35}),
						widget.NewVBox(
							widget.NewLabelWithStyle("Mode:", fyne.TextAlignTrailing, fyne.TextStyle{}),
						),
					),
					fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{200, 35}),
						widget.NewVBox(
							hcl_gui.hc_attack_mode,
						),
					),
					fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{100, 35}),
						widget.NewVBox(
							widget.NewLabelWithStyle("Hash Type:", fyne.TextAlignTrailing, fyne.TextStyle{}),
						),
					),
					fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{270, 35}),
						hash_type_fakeselector,
					),
				),
			),
		),
		widget.NewHBox(
			hcl_gui.hc_dictionary_attack_conf,
			hcl_gui.hc_combinator_attack_conf,
			hcl_gui.hc_mask_attack_conf,
			hcl_gui.hc_hybrid1_attack_conf,
			hcl_gui.hc_hybrid2_attack_conf,
			fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{280, 390}),
				widget.NewVBox(
					widget.NewGroup("Monitor",
						widget.NewVBox(
							hcl_gui.hc_disable_monitor,
							fyne.NewContainerWithLayout(layout.NewGridLayout(2),
								widget.NewLabelWithStyle("Temp Abort (°C):", fyne.TextAlignLeading, fyne.TextStyle{}),
								hcl_gui.hc_temp_abort,
							),
						),
					),
					spacer(0, 15),
					widget.NewGroup("Reserved",
						reserved(0, 58),
					),
					widget.NewGroup("Devices",
						widget.NewButton("Info", func(){
							var modal *widget.PopUp
							close_btn := widget.NewButton("Close", func(){modal.Hide()})
							close_btn.Disable()
							info_box := widget.NewMultiLineEntry()
							info := "Obtaining info..."
							info_box.SetText(info)
							info_box.OnChanged = func(string) {
								info_box.SetText(info)
							}
							c := widget.NewVBox(
								close_btn,
								fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{600, 500}),
									widget.NewScrollContainer(
										info_box,
									),
								),
							)
							modal = widget.NewModalPopUp(c, hcl_gui.window.Canvas())
							go func() {
								info = get_devices_info(hcl_gui)
								info_box.SetText(info)
								close_btn.Enable()
							}()
						}),
						spacer(0, 8),
						fyne.NewContainerWithLayout(layout.NewGridLayout(2),
							widget.NewVBox(
								widget.NewLabelWithStyle("Devices:", fyne.TextAlignLeading, fyne.TextStyle{}),
								hcl_gui.hc_devices_types,
							),
							widget.NewVBox(
								widget.NewLabelWithStyle("Workload Profile:", fyne.TextAlignLeading, fyne.TextStyle{}),
								hcl_gui.hc_wordload_profiles,
							),
						),
						spacer(0, 2),
						widget.NewButton("Benchmark", func(){
							if hcl_gui.hashcat.args.hash_type == -1 {
								err := errors.New("You must select a hash type")
								fmt.Fprintf(os.Stderr, "%s\n", err)
								dialog.ShowError(err, hcl_gui.window)
								return
							}
							var modal *widget.PopUp
							close_btn := widget.NewButton("Close", func(){modal.Hide()})
							close_btn.Disable()
							benchmark_box := widget.NewMultiLineEntry()
							benchmark := "Benchmarking..."
							benchmark_box.SetText(benchmark)
							benchmark_box.OnChanged = func(string) {
								benchmark_box.SetText(benchmark)
							}
							c := widget.NewVBox(
								close_btn,
								fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{600, 500}),
									widget.NewScrollContainer(
										benchmark_box,
									),
								),
							)
							modal = widget.NewModalPopUp(c, hcl_gui.window.Canvas())
							go func() {
								benchmark = get_benchmark(hcl_gui)
								benchmark_box.SetText(benchmark)
								close_btn.Enable()
							}()
						}),
					),
				),
			),
		),
		widget.NewGroup("Run",
			widget.NewHBox(
				spacer(20, 0),
				fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{90, 35}),
					widget.NewVBox(
						outfile,
					),
				),
				fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{620, 35}),
					widget.NewVBox(
						hcl_gui.hc_outfile,
					),
				),
				fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{30, 35}),
					widget.NewVBox(
						widget.NewButton("...", func(){
							file, err := dialog2.File().Title("Select OutFile").Filter("Out Files", "txt").Save()
							if err == nil {
								outfile.SetChecked(true)
								hcl_gui.hc_outfile.Options = append([]string{file}, hcl_gui.hc_outfile.Options[:min(len(hcl_gui.hc_outfile.Options), 4)]...)
								hcl_gui.hc_outfile.SetSelected(file)
							}
						}),
					),
				),
			),
			widget.NewHBox(
				spacer(40, 0),
				fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{350, 95}),
					widget.NewVBox(
						optimized_kernel,
						slower_candidate,
						force,
					),
				),
				widget.NewVBox(
					widget.NewHBox(
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{70, 35}),
							widget.NewVBox(
								widget.NewLabelWithStyle("Format:", fyne.TextAlignLeading, fyne.TextStyle{}),
							),
						),
						fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{265, 35}),
							widget.NewVBox(
								outfile_format,
							),
						),
					),
					fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{375, 55}),
						run_hashcat_btn,
					),
					spacer(0,0),
				),
			),
		),
	)
}
