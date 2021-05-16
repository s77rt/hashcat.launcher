package hashcatlauncher

import (
	"errors"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	dialog2 "github.com/OpenDiablo2/dialog"
)

const (
	Dialog_OS     string = "os"
	Dialog_Native string = "native"
)

///////////////////////////////////////////////////////////////////

func NewFileOpen(hcl_gui *hcl_gui) (string, error) {
	var file string
	var err error

	if hcl_gui.dialog_handler == Dialog_OS {
		file, err = newFileOpen_OS()
	} else {
		file, err = newFileOpen_Native(hcl_gui)
	}

	return file, err
}

func newFileOpen_OS() (string, error) {
	file, err := dialog2.File().Load()
	return file, err
}

func newFileOpen_Native(hcl_gui *hcl_gui) (string, error) {
	var file string
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, reader_err error) {
		defer wg.Done()
		if reader_err == nil && reader == nil {
			err = errors.New("No selection")
			return
		}
		if reader_err != nil {
			err = reader_err
			dialog.ShowError(reader_err, hcl_gui.window)
			return
		}
		file = strings.Replace(reader.URI().String(), strings.Join([]string{reader.URI().Scheme(), "://"}, ""), "", 1)
	}, hcl_gui.window)
	fd.Show()
	fd.Resize(fyne.NewSize(800, 600))
	wg.Wait()
	return file, err
}

///////////////////////////////////////////////////////////////////

func NewFileSave(hcl_gui *hcl_gui) (string, error) {
	var file string
	var err error

	if hcl_gui.dialog_handler == Dialog_OS {
		file, err = newFileSave_OS()
	} else {
		file, err = newFileSave_Native(hcl_gui)
	}

	return file, err
}

func newFileSave_OS() (string, error) {
	file, err := dialog2.File().Save()
	return file, err
}

func newFileSave_Native(hcl_gui *hcl_gui) (string, error) {
	var file string
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	fd := dialog.NewFileSave(func(reader fyne.URIWriteCloser, reader_err error) {
		defer wg.Done()
		if reader_err == nil && reader == nil {
			err = errors.New("No selection")
			return
		}
		if reader_err != nil {
			err = reader_err
			dialog.ShowError(reader_err, hcl_gui.window)
			return
		}
		file = strings.Replace(reader.URI().String(), strings.Join([]string{reader.URI().Scheme(), "://"}, ""), "", 1)
	}, hcl_gui.window)
	fd.Show()
	fd.Resize(fyne.NewSize(800, 600))
	wg.Wait()
	return file, err
}

///////////////////////////////////////////////////////////////////

func NewFolderOpen(hcl_gui *hcl_gui) (string, error) {
	var folder string
	var err error

	if hcl_gui.dialog_handler == Dialog_OS {
		folder, err = newFolderOpen_OS()
	} else {
		folder, err = newFolderOpen_Native(hcl_gui)
	}

	return folder, err
}

func newFolderOpen_OS() (string, error) {
	folder, err := dialog2.Directory().Browse()
	return folder, err
}

func newFolderOpen_Native(hcl_gui *hcl_gui) (string, error) {
	var folder string
	var err error
	var wg sync.WaitGroup
	wg.Add(1)
	fd := dialog.NewFolderOpen(func(reader fyne.ListableURI, reader_err error) {
		defer wg.Done()
		if reader_err == nil && reader == nil {
			err = errors.New("No selection")
			return
		}
		if reader_err != nil {
			err = reader_err
			dialog.ShowError(reader_err, hcl_gui.window)
			return
		}
		folder = strings.Replace(reader.String(), strings.Join([]string{reader.Scheme(), "://"}, ""), "", 1)
	}, hcl_gui.window)
	fd.Show()
	fd.Resize(fyne.NewSize(800, 600))
	wg.Wait()
	return folder, err
}
