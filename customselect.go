package hashcatlauncher

import (
	"fmt"
	"io/ioutil"
	filepath_mod "path/filepath"
	"strings"
	"strconv"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	"github.com/s77rt/hashcat.launcher/pkg/xfyne/xwidget"
)

// HashType Custom Select:
func customselect_hashtype(hcl_gui *hcl_gui) {
	var modal *widget.PopUp
	data := container.NewVBox()
	data.Objects = []fyne.CanvasObject{widget.NewLabel("Results will appear here...")}
	results_box := 	container.NewScroll(data)
	search := widget.NewEntry()
	search.SetPlaceHolder("Type to search")
	search.OnChanged = func(keyword string){
		results_box.Offset = fyne.NewPos(0,0)
		go func(){
			customselect_hashtype_options(modal, data, hcl_gui, keyword)
		}()
	}
	c := container.NewVBox(
		container.NewHBox(
			widget.NewLabelWithStyle("Hash Type", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			layout.NewSpacer(),
			widget.NewButtonWithIcon("", theme.HelpIcon(), func() {
				dialog.ShowInformation("Help", "Type at least (2) two chars to search for hash types.\nIf nothing appears make sure that you have set the hashcat bin/exe file correctly.", hcl_gui.window)
			}),
			widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
				hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){})
				modal.Hide()
			}),
		),
		search,
		container.New(layout.NewGridWrapLayout(fyne.Size{500, 600}),
			results_box,
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
}

func customselect_hashtype_options(modal *widget.PopUp, data *fyne.Container, hcl_gui *hcl_gui, keyword string) {
	var children []fyne.CanvasObject
	if len(keyword) >= 2 {
		for _, key := range hcl_gui.hashcat.available_hash_types_sorted_keys {
			option := hcl_gui.hashcat.available_hash_types[key]
			if strings.Contains(strings.ToLower(option), strings.ToLower(keyword)) || strings.Contains(strings.ToLower(key), strings.ToLower(keyword)) {
				item := xwidget.NewSelectorOptionWithStyle(fmt.Sprintf("%s - %s", key, option), key, fyne.TextAlignLeading, fyne.TextStyle{}, func(value string){
					hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){})
					set_hash_type(hcl_gui, value)
					modal.Hide()
				})
				children = append(children, container.NewHScroll(item))
			}
		}
	}
	if len(children) == 0 {
		var error_msg string
		if len(keyword) >= 2 {
			error_msg = "No matching hash types were found"
		} else {
			error_msg = "Keep typing...\n(at least (2) two chars are required)"
		}
		children = append(children, widget.NewLabelWithStyle(error_msg, fyne.TextAlignCenter, fyne.TextStyle{}))
	}
	data.Objects = children
	data.Refresh()
}

///////////////////////////////////////////////////////////////////////////////////////////////////

// Dictionaries Custom Select: (for modes excluding dictionary mode)
func customselect_dictionaries(hcl_gui *hcl_gui, selector *xwidget.Selector) {
	var modal *widget.PopUp
	data := container.NewVBox()
	results_box := 	container.NewScroll(data)
	search := widget.NewEntry()
	search.SetPlaceHolder("Type to search")
	search.OnChanged = func(keyword string){
		results_box.Offset = fyne.NewPos(0,0)
		go func(){
			customselect_dictionaries_options(modal, data, hcl_gui, selector, keyword)
		}()
	}
	c := container.NewVBox(
		container.NewHBox(
			widget.NewLabelWithStyle("Dictionaries", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			layout.NewSpacer(),
			widget.NewButtonWithIcon("", theme.HelpIcon(), func() {
				dialog.ShowInformation("Help", "Select a file, you can use the search box to filter the results", hcl_gui.window)
			}),
			widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
				hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){})
				modal.Hide()
			}),
		),
		search,
		container.New(layout.NewGridWrapLayout(fyne.Size{700, 600}),
			results_box,
		),
		container.NewVBox(
			container.New(layout.NewGridLayout(3),
				widget.NewLabelWithStyle("Left Click: Select", fyne.TextAlignLeading, fyne.TextStyle{}),
				container.New(layout.NewCenterLayout(),
					container.NewHBox(
						widget.NewButton("Add a File", func(){
							go func() {
								file, err := NewFileOpen(hcl_gui)
								if err == nil {
									hcl_gui.data.dictionaries.AddFile(file);
									if file_exists := File_Exists(file); file_exists == true {
										customselect_dictionaries_options(modal, data, hcl_gui, selector, search.Text)
									}
								}
							}()
						}),
						widget.NewButton("Add a Folder", func(){
							go func() {
								dir, err := NewFolderOpen(hcl_gui)
								if err == nil {
									files_added := 0
									found_files, err := ioutil.ReadDir(dir)
									if err == nil {
										for _, file := range found_files {
											if (!file.IsDir()) {
												if file_added := hcl_gui.data.dictionaries.AddFile(filepath_mod.Join(dir, file.Name())); file_added == true {
													files_added++
												}
											}
										}
										if files_added > 0 {
											customselect_dictionaries_options(modal, data, hcl_gui, selector, search.Text)
										}
									}
								}
							}()
						}),
					),
				),
				widget.NewLabelWithStyle("Right Click: Nothing", fyne.TextAlignTrailing, fyne.TextStyle{}),
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
	customselect_dictionaries_options(modal, data, hcl_gui, selector, "")
	modal.Show()
}

func customselect_dictionaries_options(modal *widget.PopUp, data *fyne.Container, hcl_gui *hcl_gui, selector *xwidget.Selector, keyword string) {
	var children []fyne.CanvasObject
	headings := container.NewVBox(
		container.New(layout.NewGridLayout(2),
			widget.NewLabelWithStyle("Filename", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle(fmt.Sprintf("%-37s", "Size"), fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		),
	)
	children = append(children, headings)
	for _, file := range hcl_gui.data.dictionaries.Files {
		filename := file.Name
		filepath := file.Path
		filesize := file.SizeHR
		offset := 73 - len(filename)
		if offset < 27 {
			filename = filename[:len(filename)+offset-27]+"..."
			offset = 73 - len(filename)
		}
		offset_str := strconv.Itoa(offset)
		option := fmt.Sprintf("%s %"+offset_str+"s", filename, filesize)
		if strings.Contains(strings.ToLower(filepath), strings.ToLower(keyword)) {
			item := xwidget.NewSelectorOptionWithStyle(option, filepath, fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}, func(value string){
				hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){})
				selector.SetSelected(value)
				modal.Hide()
			})
			children = append(children, item)
		}
	}
	if len(children) == 1 {
		var error_msg string
		if len(keyword) > 0 {
			error_msg = "No matching dictionaries files were found"
		} else {
			error_msg = "No dictionaries files were added"
		}
		children = append(children, widget.NewLabelWithStyle(error_msg, fyne.TextAlignCenter, fyne.TextStyle{}))
	}
	data.Objects = children
	data.Refresh()
}

// Rules Custom Select:
func customselect_rules(hcl_gui *hcl_gui, selector *xwidget.Selector) {
	var modal *widget.PopUp
	data := container.NewVBox()
	results_box := 	container.NewScroll(data)
	search := widget.NewEntry()
	search.SetPlaceHolder("Type to search")
	search.OnChanged = func(keyword string){
		results_box.Offset = fyne.NewPos(0,0)
		go func(){
			customselect_rules_options(modal, data, hcl_gui, selector, keyword)
		}()
	}
	c := container.NewVBox(
		container.NewHBox(
			widget.NewLabelWithStyle("Rules", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			layout.NewSpacer(),
			widget.NewButtonWithIcon("", theme.HelpIcon(), func() {
				dialog.ShowInformation("Help", "Select a file, you can use the search box to filter the results", hcl_gui.window)
			}),
			widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
				hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){})
				modal.Hide()
			}),
		),
		search,
		container.New(layout.NewGridWrapLayout(fyne.Size{700, 600}),
			results_box,
		),
		container.NewVBox(
			container.New(layout.NewGridLayout(3),
				widget.NewLabelWithStyle("Left Click: Select", fyne.TextAlignLeading, fyne.TextStyle{}),
				container.New(layout.NewCenterLayout(),
					container.NewHBox(
						widget.NewButton("Add a File", func(){
							go func() {
								file, err := NewFileOpen(hcl_gui)
								if err == nil {
									if file_added := hcl_gui.data.rules.AddFile(file); file_added == true {
										customselect_rules_options(modal, data, hcl_gui, selector, search.Text)
									}
								}
							}()
						}),
						widget.NewButton("Add a Folder", func(){
							go func() {
								dir, err := NewFolderOpen(hcl_gui)
								if err == nil {
									files_added := 0
									found_files, err := ioutil.ReadDir(dir)
									if err == nil {
										for _, file := range found_files {
											if (!file.IsDir()) {
												if file_added := hcl_gui.data.rules.AddFile(filepath_mod.Join(dir, file.Name())); file_added == true {
													files_added++
												}
											}
										}
										if files_added > 0 {
											customselect_rules_options(modal, data, hcl_gui, selector, search.Text)
										}
									}
								}
							}()
						}),
					),
				),
				widget.NewLabelWithStyle("Right Click: Nothing", fyne.TextAlignTrailing, fyne.TextStyle{}),
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
	customselect_rules_options(modal, data, hcl_gui, selector, "")
	modal.Show()
}

func customselect_rules_options(modal *widget.PopUp, data *fyne.Container, hcl_gui *hcl_gui, selector *xwidget.Selector, keyword string) {
	var children []fyne.CanvasObject
	headings := container.NewVBox(
		container.New(layout.NewGridLayout(2),
			widget.NewLabelWithStyle("Filename", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle(fmt.Sprintf("%-37s", "Size"), fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		),
	)
	children = append(children, headings)
	for _, file := range hcl_gui.data.rules.Files {
		filename := file.Name
		filepath := file.Path
		filesize := file.SizeHR
		offset := 73 - len(filename)
		if offset < 27 {
			filename = filename[:len(filename)+offset-27]+"..."
			offset = 73 - len(filename)
		}
		offset_str := strconv.Itoa(offset)
		option := fmt.Sprintf("%s %"+offset_str+"s", filename, filesize)
		if strings.Contains(strings.ToLower(filepath), strings.ToLower(keyword)) {
			item := xwidget.NewSelectorOptionWithStyle(option, filepath, fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}, func(value string){
				hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){})
				selector.SetSelected(value)
				modal.Hide()
			})
			children = append(children, item)
		}
	}
	if len(children) == 1 {
		var error_msg string
		if len(keyword) > 0 {
			error_msg = "No matching rules files were found"
		} else {
			error_msg = "No rules files were added"
		}
		children = append(children, widget.NewLabelWithStyle(error_msg, fyne.TextAlignCenter, fyne.TextStyle{}))
	}
	data.Objects = children
	data.Refresh()
}

///////////////////////////////////////////////////////////////////////////////////////////////////

// Dictionaries DictionaryList Custom Select: (for dictionary mode)
func customselect_dictionaries_dictionarylist(hcl_gui *hcl_gui, dictionaries *[]string, entry *widget.Entry) {
	var modal *widget.PopUp
	data := container.NewVBox()
	var table *widget.Table
	results_box := 	container.NewScroll(data)
	search := widget.NewEntry()
	search.SetPlaceHolder("Type to search")
	search.OnChanged = func(keyword string){
		results_box.Offset = fyne.NewPos(0,0)
		go func(){
			customselect_dictionaries_dictionarylist_options(modal, data, hcl_gui, table, dictionaries, entry, keyword)
		}()
	}
	c := container.NewVBox(
		container.NewHBox(
			widget.NewLabelWithStyle("Dictionaries Selection", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			layout.NewSpacer(),
			widget.NewButtonWithIcon("", theme.HelpIcon(), func() {
				dialog.ShowInformation("Help", "In the left you can see your dictionaries collection:\nthose are the dictionaries that can be selected.\n(you can also add new dictionaries to the collection)\n\nIn the right there is a listing of your selected dictionaries:\nthose are the ones that will be used in the attack.", hcl_gui.window)
			}),
			widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
				hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){})
				modal.Hide()
			}),
		),
		container.NewHSplit(
			container.NewVBox(
				container.NewPadded(
					search,
				),
				container.New(layout.NewGridWrapLayout(fyne.Size{500, 600}),
					results_box,
				),
				container.NewPadded(
					container.New(layout.NewGridLayout(2),
						widget.NewLabelWithStyle("Left Click to Select", fyne.TextAlignLeading, fyne.TextStyle{}),
						container.NewHBox(
							layout.NewSpacer(),
							widget.NewButtonWithIcon("File", theme.ContentAddIcon(), func(){
								go func() {
									file, err := NewFileOpen(hcl_gui)
									if err == nil {
										hcl_gui.data.dictionaries.AddFile(file);
										if file_exists := File_Exists(file); file_exists == true {
											customselect_dictionaries_dictionarylist_options(modal, data, hcl_gui, table, dictionaries, entry, search.Text)
										}
									}
								}()
							}),
							widget.NewButtonWithIcon("Folder", theme.ContentAddIcon(), func(){
								go func() {
									dir, err := NewFolderOpen(hcl_gui)
									if err == nil {
										files_added := 0
										found_files, err := ioutil.ReadDir(dir)
										if err == nil {
											for _, file := range found_files {
												if (!file.IsDir()) {
													if file_added := hcl_gui.data.dictionaries.AddFile(filepath_mod.Join(dir, file.Name())); file_added == true {
														files_added++
													}
												}
											}
											if files_added > 0 {
												customselect_dictionaries_dictionarylist_options(modal, data, hcl_gui, table, dictionaries, entry, search.Text)
											}
										}
									}
								}()
							}),
						),
					),
				),
			),
			container.NewVBox(
				container.New(layout.NewGridWrapLayout(fyne.Size{0, 45}),
					widget.NewLabelWithStyle("Selected Dictionaries:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				),
				container.New(layout.NewGridWrapLayout(fyne.Size{500, 600}),
					func () fyne.CanvasObject {
						table = widget.NewTable(
							func() (int, int) { return 1, 2 },
							func() fyne.CanvasObject {
								return container.NewHBox(
									widget.NewLabel("Filename label"),
									container.NewGridWithColumns(3,
										widget.NewButtonWithIcon("", theme.MoveUpIcon(), func(){}),
										widget.NewButtonWithIcon("", theme.MoveDownIcon(), func(){}),
										widget.NewButtonWithIcon("", theme.DeleteIcon(), func(){}),
									),
								)
							},
							func(id widget.TableCellID, cell fyne.CanvasObject) {
								label := cell.(*fyne.Container).Objects[0].(*widget.Label)
								buttons := cell.(*fyne.Container).Objects[1].(*fyne.Container)
								label.Show()
								buttons.Hide()
								if id.Row == 0 {
									switch id.Col {
									case 0:
										label.TextStyle.Bold = true
										label.SetText("Dictionary")
									case 1:
										label.TextStyle.Bold = true
										label.SetText("Actions")
									}
								} else {
									files_list := strings.Split(strings.Replace(entry.Text, "\r\n", "\n", -1), "\n")
									if len(files_list) > id.Row - 1 {
										fileindex := id.Row - 1
										filename := files_list[fileindex]
										switch id.Col {
										case 0:
											_, filename_short := filepath_mod.Split(filename)
											offset := 51 - len(filename_short)
											if offset < 19 {
												filename_short = filename_short[:len(filename_short)+offset-19]+"..."
												offset = 51 - len(filename_short)
											}
											label.SetText(filename_short)
										case 1:
											label.Hide()
											buttons.Show()
											// Move Up Button
											buttons.Objects[0].(*widget.Button).OnTapped = func() {
												if fileindex == 0 {
													return
												}
												files_list := append(files_list[:fileindex-1], append([]string{files_list[fileindex], files_list[fileindex-1]}, files_list[fileindex:]...)...)
												entry.SetText(strings.Join(files_list, "\n")+"\n")
												valid_files := len(*dictionaries)
												table.Length = func() (int, int) { return valid_files+1, 2 }
												table.Refresh()
											}
											// Move Down Button
											buttons.Objects[1].(*widget.Button).OnTapped = func() {
												if fileindex == len(files_list) - 1 {
													return
												}
												files_list := append(files_list[:fileindex], append([]string{files_list[fileindex+1], files_list[fileindex]}, files_list[fileindex+1:]...)...)
												entry.SetText(strings.Join(files_list, "\n")+"\n")
												valid_files := len(*dictionaries)
												table.Length = func() (int, int) { return valid_files+1, 2 }
												table.Refresh()
											}
											// Delete Button
											buttons.Objects[2].(*widget.Button).OnTapped = func() {
												files_list := append(files_list[:fileindex], files_list[fileindex+1:]...)
												entry.SetText(strings.Join(files_list, "\n")+"\n")
												valid_files := len(*dictionaries)
												table.Length = func() (int, int) { return valid_files+1, 2 }
												table.Refresh()
											}
										}
									}
								}
							})
						valid_files := len(*dictionaries)
						table.Length = func() (int, int) { return valid_files+1, 2 }
						table.SetColumnWidth(0, 350)
						table.SetColumnWidth(1, 100)
						return table
					}(),
				),
				container.NewPadded(
					container.NewGridWithColumns(2,
						container.NewHBox(
							widget.NewButtonWithIcon("File", theme.ContentAddIcon(), func() {
								go func() {
									file, err := NewFileOpen(hcl_gui)
									if err == nil {
										hcl_gui.data.dictionaries.AddFile(file);
										if file_exists := File_Exists(file); file_exists == true {
											customselect_dictionaries_dictionarylist_options(modal, data, hcl_gui, table, dictionaries, entry, search.Text)
											entry.SetText(entry.Text+file+"\n")
											valid_files := len(*dictionaries)
											table.Length = func() (int, int) { return valid_files+1, 2 }
											table.Refresh()
										}
									}
								}()
							}),	
							widget.NewButtonWithIcon("Folder", theme.ContentAddIcon(), func() {
								go func() {
									dir, err := NewFolderOpen(hcl_gui)
									if err == nil {
										found_files, err := ioutil.ReadDir(dir)
										files := []string{}
										if err == nil {
											if len(found_files) > 0 {
												for _, found_file := range found_files {
													if (!found_file.IsDir()) {
														file := filepath_mod.Join(dir, found_file.Name())
														hcl_gui.data.dictionaries.AddFile(file);
														files = append(files, file)
													}
												}
												customselect_dictionaries_dictionarylist_options(modal, data, hcl_gui, table, dictionaries, entry, search.Text)
												entry.SetText(entry.Text+strings.Join(files, "\n")+"\n")
												valid_files := len(*dictionaries)
												table.Length = func() (int, int) { return valid_files+1, 2 }
												table.Refresh()
											}
										}
									}
								}()
							}),
						),
						widget.NewButtonWithIcon("Clear All", theme.ContentClearIcon(), func() {
							entry.SetText("")
							valid_files := len(*dictionaries)
							table.Length = func() (int, int) { return valid_files+1, 2 }
							table.Refresh()
						}),
					),
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
	customselect_dictionaries_dictionarylist_options(modal, data, hcl_gui, table, dictionaries, entry, "")
	modal.Show()
}

func customselect_dictionaries_dictionarylist_options(modal *widget.PopUp, data *fyne.Container, hcl_gui *hcl_gui, table *widget.Table, dictionaries *[]string, entry *widget.Entry, keyword string) {
	var children []fyne.CanvasObject
	headings := container.NewVBox(
		container.New(layout.NewGridLayout(2),
			widget.NewLabelWithStyle("Filename", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle(fmt.Sprintf("%-37s", "Size"), fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		),
	)
	children = append(children, headings)
	for _, file := range hcl_gui.data.dictionaries.Files {
		filename := file.Name
		filepath := file.Path
		filesize := file.SizeHR
		offset := 51 - len(filename)
		if offset < 19 {
			filename = filename[:len(filename)+offset-19]+"..."
			offset = 51 - len(filename)
		}
		offset_str := strconv.Itoa(offset)
		option := fmt.Sprintf("%s %"+offset_str+"s", filename, filesize)
		if strings.Contains(strings.ToLower(filepath), strings.ToLower(keyword)) {
			item := xwidget.NewSelectorOptionWithStyle(option, filepath, fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}, func(value string){
				hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){})
				entry.SetText(entry.Text+filepath+"\n")
				valid_files := len(*dictionaries)
				table.Length = func() (int, int) { return valid_files+1, 2 }
				table.Refresh()
			})
			children = append(children, item)
		}
	}
	if len(children) == 1 {
		var error_msg string
		if len(keyword) > 0 {
			error_msg = "No matching dictionaries files were found"
		} else {
			error_msg = "No dictionaries files were added"
		}
		children = append(children, widget.NewLabelWithStyle(error_msg, fyne.TextAlignCenter, fyne.TextStyle{}))
	}
	data.Objects = children
	data.Refresh()
}

///////////////////////////////////////////////////////////////////////////////////////////////////

// Dictionaries Edit Custom Select: (The one for the open filebase feature)
func customselect_dictionaries_edit(hcl_gui *hcl_gui) {
	var modal *widget.PopUp
	data := container.NewVBox()
	results_box := 	container.NewScroll(data)
	search := widget.NewEntry()
	search.SetPlaceHolder("Type to search")
	search.OnChanged = func(keyword string){
		results_box.Offset = fyne.NewPos(0,0)
		go func(){
			customselect_dictionaries_edit_options(modal, data, hcl_gui, keyword)
		}()
	}
	c := container.NewVBox(
		container.NewHBox(
			widget.NewLabelWithStyle("Dictionaries", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			layout.NewSpacer(),
			widget.NewButtonWithIcon("", theme.HelpIcon(), func() {
				dialog.ShowInformation("Help", "View/Remove a file, you can use the search box to filter the results", hcl_gui.window)
			}),
			widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
				hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){})
				modal.Hide()
			}),
		),
		search,
		container.New(layout.NewGridWrapLayout(fyne.Size{700, 600}),
			results_box,
		),
		container.NewVBox(
			container.New(layout.NewGridLayout(3),
				widget.NewLabelWithStyle("Left Click: View FileInfo", fyne.TextAlignLeading, fyne.TextStyle{}),
				container.New(layout.NewCenterLayout(),
					container.NewHBox(
						widget.NewButton("Add a File", func(){
							go func() {
								file, err := NewFileOpen(hcl_gui)
								if err == nil {
									hcl_gui.data.dictionaries.AddFile(file);
									if file_exists := File_Exists(file); file_exists == true {
										customselect_dictionaries_edit_options(modal, data, hcl_gui, search.Text)
									}
								}
							}()
						}),
						widget.NewButton("Add a Folder", func(){
							go func() {
								dir, err := NewFolderOpen(hcl_gui)
								if err == nil {
									files_added := 0
									found_files, err := ioutil.ReadDir(dir)
									if err == nil {
										for _, file := range found_files {
											if (!file.IsDir()) {
												if file_added := hcl_gui.data.dictionaries.AddFile(filepath_mod.Join(dir, file.Name())); file_added == true {
													files_added++
												}
											}
										}
										if files_added > 0 {
											customselect_dictionaries_edit_options(modal, data, hcl_gui, search.Text)
										}
									}
								}
							}()
						}),
						widget.NewButton("Clear", func() {
							hcl_gui.data.dictionaries.Clear()
							customselect_dictionaries_edit_options(modal, data, hcl_gui, "")
						}),
					),
				),
				widget.NewLabelWithStyle("Right Click: Remove File", fyne.TextAlignTrailing, fyne.TextStyle{}),
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
	customselect_dictionaries_edit_options(modal, data, hcl_gui, "")
	modal.Show()
}

func customselect_dictionaries_edit_options(modal *widget.PopUp, data *fyne.Container, hcl_gui *hcl_gui,keyword string) {
	var children []fyne.CanvasObject
	headings := container.NewVBox(
		container.New(layout.NewGridLayout(2),
			widget.NewLabelWithStyle("Filename", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle(fmt.Sprintf("%-37s", "Size"), fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		),
	)
	children = append(children, headings)
	for _, file := range hcl_gui.data.dictionaries.Files {
		filename := file.Name
		filepath := file.Path
		filesize := file.SizeHR
		offset := 73 - len(filename)
		if offset < 27 {
			filename = filename[:len(filename)+offset-27]+"..."
			offset = 73 - len(filename)
		}
		offset_str := strconv.Itoa(offset)
		option := fmt.Sprintf("%s %"+offset_str+"s", filename, filesize)
		if strings.Contains(strings.ToLower(filepath), strings.ToLower(keyword)) {
			item := xwidget.NewSelectorOptionWithStyle(option, filepath, fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}, func(value string){
				dialog.ShowCustom("FileInfo", "OK", 
					container.New(layout.NewFormLayout(),
						widget.NewLabelWithStyle("Filename:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
						widget.NewLabelWithStyle(filename, fyne.TextAlignLeading, fyne.TextStyle{}),
						widget.NewLabelWithStyle("Filepath:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
						widget.NewLabelWithStyle(filepath, fyne.TextAlignLeading, fyne.TextStyle{}),
						widget.NewLabelWithStyle("Filesize:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
						widget.NewLabelWithStyle(filesize, fyne.TextAlignLeading, fyne.TextStyle{}),
					),
				hcl_gui.window)
			})
			item.OnTappedSecondary = func(value string) {
				if file_removed := hcl_gui.data.dictionaries.RemoveFile(value); file_removed == true {
					item.Hide()
				}
			}
			children = append(children, item)
		}
	}
	if len(children) == 1 {
		var error_msg string
		if len(keyword) > 0 {
			error_msg = "No matching dictionaries files were found"
		} else {
			error_msg = "No dictionaries files were added"
		}
		children = append(children, widget.NewLabelWithStyle(error_msg, fyne.TextAlignCenter, fyne.TextStyle{}))
	}
	data.Objects = children
	data.Refresh()
}

// Dictionaries Edit Custom Select: (The one for the open filebase feature)
func customselect_rules_edit(hcl_gui *hcl_gui) {
	var modal *widget.PopUp
	data := container.NewVBox()
	results_box := 	container.NewScroll(data)
	search := widget.NewEntry()
	search.SetPlaceHolder("Type to search")
	search.OnChanged = func(keyword string){
		results_box.Offset = fyne.NewPos(0,0)
		go func(){
			customselect_rules_edit_options(modal, data, hcl_gui, keyword)
		}()
	}
	c := container.NewVBox(
		container.NewHBox(
			widget.NewLabelWithStyle("Rules", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			layout.NewSpacer(),
			widget.NewButtonWithIcon("", theme.HelpIcon(), func() {
				dialog.ShowInformation("Help", "View/Remove a file, you can use the search box to filter the results", hcl_gui.window)
			}),
			widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
				hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){})
				modal.Hide()
			}),
		),
		search,
		container.New(layout.NewGridWrapLayout(fyne.Size{700, 600}),
			results_box,
		),
		container.NewVBox(
			container.New(layout.NewGridLayout(3),
				widget.NewLabelWithStyle("Left Click: View FileInfo", fyne.TextAlignLeading, fyne.TextStyle{}),
				container.New(layout.NewCenterLayout(),
					container.NewHBox(
						widget.NewButton("Add a File", func(){
							go func() {
								file, err := NewFileOpen(hcl_gui)
								if err == nil {
									if file_added := hcl_gui.data.rules.AddFile(file); file_added == true {
										customselect_rules_edit_options(modal, data, hcl_gui, search.Text)
									}
								}
							}()
						}),
						widget.NewButton("Add a Folder", func(){
							go func() {
								dir, err := NewFolderOpen(hcl_gui)
								if err == nil {
									files_added := 0
									found_files, err := ioutil.ReadDir(dir)
									if err == nil {
										for _, file := range found_files {
											if (!file.IsDir()) {
												if file_added := hcl_gui.data.rules.AddFile(filepath_mod.Join(dir, file.Name())); file_added == true {
													files_added++
												}
											}
										}
										if files_added > 0 {
											customselect_rules_edit_options(modal, data, hcl_gui, search.Text)
										}
									}
								}
							}()
						}),
					),
				),
				widget.NewLabelWithStyle("Right Click: Remove File", fyne.TextAlignTrailing, fyne.TextStyle{}),
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
	customselect_rules_edit_options(modal, data, hcl_gui, "")
	modal.Show()
}

func customselect_rules_edit_options(modal *widget.PopUp, data *fyne.Container, hcl_gui *hcl_gui,keyword string) {
	var children []fyne.CanvasObject
	headings := container.NewVBox(
		container.New(layout.NewGridLayout(2),
			widget.NewLabelWithStyle("Filename", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle(fmt.Sprintf("%-37s", "Size"), fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		),
	)
	children = append(children, headings)
	for _, file := range hcl_gui.data.rules.Files {
		filename := file.Name
		filepath := file.Path
		filesize := file.SizeHR
		offset := 73 - len(filename)
		if offset < 27 {
			filename = filename[:len(filename)+offset-27]+"..."
			offset = 73 - len(filename)
		}
		offset_str := strconv.Itoa(offset)
		option := fmt.Sprintf("%s %"+offset_str+"s", filename, filesize)
		if strings.Contains(strings.ToLower(filepath), strings.ToLower(keyword)) {
			item := xwidget.NewSelectorOptionWithStyle(option, filepath, fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}, func(value string){
				dialog.ShowCustom("FileInfo", "OK", 
					container.New(layout.NewFormLayout(),
						widget.NewLabelWithStyle("Filename:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
						widget.NewLabelWithStyle(filename, fyne.TextAlignLeading, fyne.TextStyle{}),
						widget.NewLabelWithStyle("Filepath:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
						widget.NewLabelWithStyle(filepath, fyne.TextAlignLeading, fyne.TextStyle{}),
						widget.NewLabelWithStyle("Filesize:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
						widget.NewLabelWithStyle(filesize, fyne.TextAlignLeading, fyne.TextStyle{}),
					),
				hcl_gui.window)
			})
			item.OnTappedSecondary = func(value string) {
				if file_removed := hcl_gui.data.rules.RemoveFile(value); file_removed == true {
					item.Hide()
				}
			}
			children = append(children, item)
		}
	}
	if len(children) == 1 {
		var error_msg string
		if len(keyword) > 0 {
			error_msg = "No matching rules files were found"
		} else {
			error_msg = "No rules files were added"
		}
		children = append(children, widget.NewLabelWithStyle(error_msg, fyne.TextAlignCenter, fyne.TextStyle{}))
	}
	data.Objects = children
	data.Refresh()
}
