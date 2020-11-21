package hashcatlauncher

import (
	"sync"
	"errors"
	"strings"
	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
)

func NewFileOpen(hcl_gui *hcl_gui) (string, error) {
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

func NewFileSave(hcl_gui *hcl_gui) (string, error) {
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

func NewFolderOpen(hcl_gui *hcl_gui) (string, error) {
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
