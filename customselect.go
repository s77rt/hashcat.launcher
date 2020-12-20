package hashcatlauncher

import (
	"fmt"
	"io/ioutil"
	filepath_mod "path/filepath"
	"strings"
	"strconv"
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"fyne.io/fyne/dialog"
	"github.com/s77rt/hashcat.launcher/pkg/xfyne/xwidget"
)

// HashType Custom Select:
func customselect_hashtype(hcl_gui *hcl_gui) {
	var modal *widget.PopUp
	data := widget.NewVBox()
	data.Children = []fyne.CanvasObject{widget.NewLabel("Results will appear here...")}
	results_box := 	widget.NewScrollContainer(data)
	search := widget.NewEntry()
	search.SetPlaceHolder("Type to search")
	search.OnChanged = func(keyword string){
		results_box.Offset = fyne.NewPos(0,0)
		go func(){
			customselect_hashtype_options(modal, data, hcl_gui, keyword)
		}()
	}
	c := widget.NewVBox(
		widget.NewHBox(
			fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{485, 40}),
				widget.NewHScrollContainer(search),
			),
			fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{15, 15}),
				widget.NewButton("X", func(){hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){});modal.Hide()}),
				widget.NewButton("?", func(){dialog.ShowInformation("Help", "Type at least (2) two chars to search for hash types.\nIf nothing appears make sure that you have set the hashcat binary file correctly.", hcl_gui.window)}),
			),
		),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{500, 600}),
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
}

func customselect_hashtype_options(modal *widget.PopUp, data *widget.Box, hcl_gui *hcl_gui, keyword string) {
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
				children = append(children, item)
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
	data.Children = children
	data.Refresh()
}

///////////////////////////////////////////////////////////////////////////////////////////////////

// Dictionaries Custom Select:
func customselect_dictionaries(hcl_gui *hcl_gui, selector *xwidget.Selector) {
	var modal *widget.PopUp
	data := widget.NewVBox()
	results_box := 	widget.NewScrollContainer(data)
	search := widget.NewEntry()
	search.SetPlaceHolder("Type to search")
	search.OnChanged = func(keyword string){
		results_box.Offset = fyne.NewPos(0,0)
		go func(){
			customselect_dictionaries_options(modal, data, hcl_gui, selector, keyword)
		}()
	}
	c := widget.NewVBox(
		widget.NewHBox(
			fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{685, 40}),
				widget.NewHScrollContainer(search),
			),
			fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{15, 15}),
				widget.NewButton("X", func(){hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){});modal.Hide()}),
				widget.NewButton("?", func(){dialog.ShowInformation("Help", "Select a file, you can use the search box to filter the results.", hcl_gui.window)}),
			),
		),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{700, 600}),
			results_box,
		),
		widget.NewVBox(
			fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				widget.NewLabelWithStyle("Left Click: Select", fyne.TextAlignLeading, fyne.TextStyle{}),
				fyne.NewContainerWithLayout(layout.NewCenterLayout(),
					widget.NewHBox(
						widget.NewButton("Add a File", func(){
							go func() {
								file, err := NewFileOpen(hcl_gui)
								if err == nil {
									if file_added := hcl_gui.data.dictionaries.AddFile(file); file_added == true {
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
}

func customselect_dictionaries_options(modal *widget.PopUp, data *widget.Box, hcl_gui *hcl_gui, selector *xwidget.Selector, keyword string) {
	var children []fyne.CanvasObject
	headings := widget.NewVBox(
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			widget.NewLabelWithStyle("Filename", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle(fmt.Sprintf("%-37s", "Size"), fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		),
	)
	children = append(children, headings)
	for _, file := range hcl_gui.data.dictionaries.Files {
		filename := file.Name
		filepath := file.Path
		filesize := file.SizeHR
		offset := strconv.Itoa(73 - len(filename))
		option := fmt.Sprintf("%s %"+offset+"s", filename, filesize)
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
	data.Children = children
	data.Refresh()
}

// Rules Custom Select:
func customselect_rules(hcl_gui *hcl_gui, selector *xwidget.Selector) {
	var modal *widget.PopUp
	data := widget.NewVBox()
	results_box := 	widget.NewScrollContainer(data)
	search := widget.NewEntry()
	search.SetPlaceHolder("Type to search")
	search.OnChanged = func(keyword string){
		results_box.Offset = fyne.NewPos(0,0)
		go func(){
			customselect_rules_options(modal, data, hcl_gui, selector, keyword)
		}()
	}
	c := widget.NewVBox(
		widget.NewHBox(
			fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{685, 40}),
				widget.NewHScrollContainer(search),
			),
			fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{15, 15}),
				widget.NewButton("X", func(){hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){});modal.Hide()}),
				widget.NewButton("?", func(){dialog.ShowInformation("Help", "Select a file, you can use the search box to filter the results.", hcl_gui.window)}),
			),
		),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{700, 600}),
			results_box,
		),
		widget.NewVBox(
			fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				widget.NewLabelWithStyle("Left Click: Select", fyne.TextAlignLeading, fyne.TextStyle{}),
				fyne.NewContainerWithLayout(layout.NewCenterLayout(),
					widget.NewHBox(
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
}

func customselect_rules_options(modal *widget.PopUp, data *widget.Box, hcl_gui *hcl_gui, selector *xwidget.Selector, keyword string) {
	var children []fyne.CanvasObject
	headings := widget.NewVBox(
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			widget.NewLabelWithStyle("Filename", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle(fmt.Sprintf("%-37s", "Size"), fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		),
	)
	children = append(children, headings)
	for _, file := range hcl_gui.data.rules.Files {
		filename := file.Name
		filepath := file.Path
		filesize := file.SizeHR
		offset := strconv.Itoa(73 - len(filename))
		option := fmt.Sprintf("%s %"+offset+"s", filename, filesize)
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
	data.Children = children
	data.Refresh()
}

///////////////////////////////////////////////////////////////////////////////////////////////////

// Dictionaries DictionaryList Custom Select: (The one for the [~] button)
func customselect_dictionaries_dictionarylist(hcl_gui *hcl_gui, entry *widget.Entry) {
	var modal *widget.PopUp
	data := widget.NewVBox()
	results_box := 	widget.NewScrollContainer(data)
	search := widget.NewEntry()
	search.SetPlaceHolder("Type to search")
	search.OnChanged = func(keyword string){
		results_box.Offset = fyne.NewPos(0,0)
		go func(){
			customselect_dictionaries_dictionarylist_options(modal, data, hcl_gui, entry, keyword)
		}()
	}
	c := widget.NewVBox(
		widget.NewHBox(
			fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{685, 40}),
				widget.NewHScrollContainer(search),
			),
			fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{15, 15}),
				widget.NewButton("X", func(){hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){});modal.Hide()}),
				widget.NewButton("?", func(){dialog.ShowInformation("Help", "Select a file, you can use the search box to filter the results.", hcl_gui.window)}),
			),
		),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{700, 600}),
			results_box,
		),
		widget.NewVBox(
			fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				widget.NewLabelWithStyle("Left Click: Select", fyne.TextAlignLeading, fyne.TextStyle{}),
				fyne.NewContainerWithLayout(layout.NewCenterLayout(),
					widget.NewHBox(
						widget.NewButton("Add a File", func(){
							go func() {
								file, err := NewFileOpen(hcl_gui)
								if err == nil {
									if file_added := hcl_gui.data.dictionaries.AddFile(file); file_added == true {
										customselect_dictionaries_dictionarylist_options(modal, data, hcl_gui, entry, search.Text)
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
											customselect_dictionaries_dictionarylist_options(modal, data, hcl_gui, entry, search.Text)
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
	customselect_dictionaries_dictionarylist_options(modal, data, hcl_gui, entry, "")
}

func customselect_dictionaries_dictionarylist_options(modal *widget.PopUp, data *widget.Box, hcl_gui *hcl_gui, entry *widget.Entry, keyword string) {
	var children []fyne.CanvasObject
	headings := widget.NewVBox(
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			widget.NewLabelWithStyle("Filename", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle(fmt.Sprintf("%-37s", "Size"), fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		),
	)
	children = append(children, headings)
	for _, file := range hcl_gui.data.dictionaries.Files {
		filename := file.Name
		filepath := file.Path
		filesize := file.SizeHR
		offset := strconv.Itoa(73 - len(filename))
		option := fmt.Sprintf("%s %"+offset+"s", filename, filesize)
		if strings.Contains(strings.ToLower(filepath), strings.ToLower(keyword)) {
			item := xwidget.NewSelectorOptionWithStyle(option, filepath, fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}, func(value string){
				hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){})
				entry.SetText(entry.Text+filepath+"\n")
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
	data.Children = children
	data.Refresh()
}

///////////////////////////////////////////////////////////////////////////////////////////////////

// Dictionaries Edit Custom Select: (The one for the open filebase feature)
func customselect_dictionaries_edit(hcl_gui *hcl_gui) {
	var modal *widget.PopUp
	data := widget.NewVBox()
	results_box := 	widget.NewScrollContainer(data)
	search := widget.NewEntry()
	search.SetPlaceHolder("Type to search")
	search.OnChanged = func(keyword string){
		results_box.Offset = fyne.NewPos(0,0)
		go func(){
			customselect_dictionaries_edit_options(modal, data, hcl_gui, keyword)
		}()
	}
	c := widget.NewVBox(
		widget.NewHBox(
			fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{685, 40}),
				widget.NewHScrollContainer(search),
			),
			fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{15, 15}),
				widget.NewButton("X", func(){hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){});modal.Hide()}),
				widget.NewButton("?", func(){dialog.ShowInformation("Help", "View/Remove a file, you can use the search box to filter the results.", hcl_gui.window)}),
			),
		),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{700, 600}),
			results_box,
		),
		widget.NewVBox(
			fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				widget.NewLabelWithStyle("Left Click: View FileInfo", fyne.TextAlignLeading, fyne.TextStyle{}),
				fyne.NewContainerWithLayout(layout.NewCenterLayout(),
					widget.NewHBox(
						widget.NewButton("Add a File", func(){
							go func() {
								file, err := NewFileOpen(hcl_gui)
								if err == nil {
									if file_added := hcl_gui.data.dictionaries.AddFile(file); file_added == true {
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
}

func customselect_dictionaries_edit_options(modal *widget.PopUp, data *widget.Box, hcl_gui *hcl_gui,keyword string) {
	var children []fyne.CanvasObject
	headings := widget.NewVBox(
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			widget.NewLabelWithStyle("Filename", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle(fmt.Sprintf("%-37s", "Size"), fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		),
	)
	children = append(children, headings)
	for _, file := range hcl_gui.data.dictionaries.Files {
		filename := file.Name
		filepath := file.Path
		filesize := file.SizeHR
		offset := strconv.Itoa(73 - len(filename))
		option := fmt.Sprintf("%s %"+offset+"s", filename, filesize)
		if strings.Contains(strings.ToLower(filepath), strings.ToLower(keyword)) {
			item := xwidget.NewSelectorOptionWithStyle(option, filepath, fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}, func(value string){
				dialog.ShowCustom("FileInfo", "OK", 
					fyne.NewContainerWithLayout(layout.NewFormLayout(),
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
	data.Children = children
	data.Refresh()
}

// Dictionaries Edit Custom Select: (The one for the open filebase feature)
func customselect_rules_edit(hcl_gui *hcl_gui) {
	var modal *widget.PopUp
	data := widget.NewVBox()
	results_box := 	widget.NewScrollContainer(data)
	search := widget.NewEntry()
	search.SetPlaceHolder("Type to search")
	search.OnChanged = func(keyword string){
		results_box.Offset = fyne.NewPos(0,0)
		go func(){
			customselect_rules_edit_options(modal, data, hcl_gui, keyword)
		}()
	}
	c := widget.NewVBox(
		widget.NewHBox(
			fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{685, 40}),
				widget.NewHScrollContainer(search),
			),
			fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{15, 15}),
				widget.NewButton("X", func(){hcl_gui.window.Canvas().SetOnTypedKey(func(*fyne.KeyEvent){});modal.Hide()}),
				widget.NewButton("?", func(){dialog.ShowInformation("Help", "View/Remove a file, you can use the search box to filter the results.", hcl_gui.window)}),
			),
		),
		fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{700, 600}),
			results_box,
		),
		widget.NewVBox(
			fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				widget.NewLabelWithStyle("Left Click: View FileInfo", fyne.TextAlignLeading, fyne.TextStyle{}),
				fyne.NewContainerWithLayout(layout.NewCenterLayout(),
					widget.NewHBox(
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
}

func customselect_rules_edit_options(modal *widget.PopUp, data *widget.Box, hcl_gui *hcl_gui,keyword string) {
	var children []fyne.CanvasObject
	headings := widget.NewVBox(
		fyne.NewContainerWithLayout(layout.NewGridLayout(2),
			widget.NewLabelWithStyle("Filename", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle(fmt.Sprintf("%-37s", "Size"), fyne.TextAlignTrailing, fyne.TextStyle{Bold: true}),
		),
	)
	children = append(children, headings)
	for _, file := range hcl_gui.data.rules.Files {
		filename := file.Name
		filepath := file.Path
		filesize := file.SizeHR
		offset := strconv.Itoa(73 - len(filename))
		option := fmt.Sprintf("%s %"+offset+"s", filename, filesize)
		if strings.Contains(strings.ToLower(filepath), strings.ToLower(keyword)) {
			item := xwidget.NewSelectorOptionWithStyle(option, filepath, fyne.TextAlignLeading, fyne.TextStyle{Monospace: true}, func(value string){
				dialog.ShowCustom("FileInfo", "OK", 
					fyne.NewContainerWithLayout(layout.NewFormLayout(),
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
	data.Children = children
	data.Refresh()
}
