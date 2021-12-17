package hashcatlauncher

import (
	"io"
	"io/ioutil"
	"os"
)

type Restore struct {
	Version int32
	Cwd     [256]byte

	DictsPos uint32
	MasksPos uint32

	WordsCur uint64

	Argc uint32
	Argv []byte
}

func UnpackRestore(r io.Reader) (restore *Restore, err error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	restore = &Restore{}

	restore.Version = int32(nativeEndian.Uint32(data[0:4]))
	copy(restore.Cwd[:], data[4:260])

	restore.DictsPos = nativeEndian.Uint32(data[260:264])
	restore.MasksPos = nativeEndian.Uint32(data[264:268])

	restore.WordsCur = nativeEndian.Uint64(data[272:280])

	restore.Argc = nativeEndian.Uint32(data[280:284])
	restore.Argv = data[288:]

	return
}

const RestoreFileExt = ".restore"

type RestoreFile struct {
	*os.File
	Data *Restore
}

func (rf *RestoreFile) Unpack() error {
	restore, err := UnpackRestore(rf.File)
	if err != nil {
		return err
	}
	rf.Data = restore

	return nil
}
