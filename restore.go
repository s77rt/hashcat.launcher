package hashcatlauncher

import (
	"encoding/binary"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2/dialog"
)

func GetRestoreFiles(hcl_gui *hcl_gui) []string {
	dir := filepath.Dir(hcl_gui.hashcat.binary_file)
	files, err := filepath.Glob(filepath.Join(dir, "*.restore"))
	if err != nil {
		dialog.ShowError(err, hcl_gui.window)
		return []string{}
	} else {
		return files
	}
}

type RestoreFile struct {
	Version int32
	Cwd     string // 256 chars

	Dicts_pos uint32
	Masks_pos uint32

	Words_cur uint64

	Argc uint32
	Argv string

	Session_name string
	Time         int64
	Task_id      int

	Path string
}

func ReadRestoreFile(hcl_gui *hcl_gui, path string) *RestoreFile {
	restore_file := &RestoreFile{}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		dialog.ShowError(err, hcl_gui.window)
	} else {
		restore_file.Version = int32(binary.LittleEndian.Uint32(file[0:4]))
		restore_file.Cwd = strings.ToValidUTF8(string(file[4:260]), "")
		restore_file.Dicts_pos = binary.LittleEndian.Uint32(file[260:264])
		restore_file.Masks_pos = binary.LittleEndian.Uint32(file[264:268])
		restore_file.Words_cur = binary.LittleEndian.Uint64(file[272:280])
		restore_file.Argc = binary.LittleEndian.Uint32(file[280:284])
		restore_file.Argv = strings.ToValidUTF8(string(file[288:]), "")

		filename := filepath.Base(path)
		restore_file.Session_name = strings.TrimSuffix(filename, filepath.Ext(filename))
		restore_file_info := re_restore_file_info.FindStringSubmatch(filename)
		if len(restore_file_info) == 3 {
			restore_file.Time, _ = strconv.ParseInt(restore_file_info[1], 10, 64)
			restore_file.Task_id, _ = strconv.Atoi(restore_file_info[2])
		} else {
			restore_file.Time = -1
			restore_file.Task_id = -1
		}
		restore_file.Path = path
	}
	return restore_file
}

func (restore_file *RestoreFile) GetArguments() []string {
	return strings.Split(strings.Replace(restore_file.Argv, "\r\n", "\n", -1), "\n")[1:]
}

func (restore_file *RestoreFile) Delete() error {
	return os.Remove(restore_file.Path)
}
