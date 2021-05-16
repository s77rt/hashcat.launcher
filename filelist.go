package hashcatlauncher

import (
	"os"
)

type FileList struct {
	Files []FileListItem
}

type FileListItem struct {
	Path   string
	Name   string
	Size   int64
	SizeHR string
}

func (FileList *FileList) AddFile(filepath string) bool {
	f, err := os.Stat(filepath)
	if err != nil {
		return false
	} else if !f.Mode().IsRegular() {
		return false
	}
	for _, file := range FileList.Files {
		if file.Path == filepath {
			return false
		}
	}
	new_file := &FileListItem{
		Path:   filepath,
		Name:   f.Name(),
		Size:   f.Size(),
		SizeHR: ByteCountIEC(f.Size()),
	}
	FileList.Files = append(FileList.Files, *new_file)
	return true
}

func (FileList *FileList) RemoveFile(filepath string) bool {
	index := -1
	for i, file := range FileList.Files {
		if file.Path == filepath {
			index = i
		}
	}
	if index >= 0 {
		FileList.Files = append(FileList.Files[:index], FileList.Files[index+1:]...)
		return true
	} else {
		return false
	}
}

func (FileList *FileList) Clear() bool {
	FileList.Files = []FileListItem{}
	return true
}
