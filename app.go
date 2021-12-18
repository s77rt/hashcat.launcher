package hashcatlauncher

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime"

	"github.com/zserge/lorca"
)

var (
	Version string = "dev"
)

var (
	HashcatDir string

	HashesDir       string
	DictionariesDir string
	RulesDir        string
	MasksDir        string
)

type App struct {
	Server net.Listener
	UI     lorca.UI

	Hashcat      *Hashcat
	Hashes       []string
	Dictionaries []string
	Rules        []string
	Masks        []string

	Tasks                   map[string]*Task
	TaskUpdateCallback      func(TaskUpdate)
	TaskPreProcessCallback  func()
	TaskPostProcessCallback func()
}

func (a *App) Init() {
	exe, err := os.Executable()
	if err != nil {
		panic("unable to get executable location")
	}

	exeDir, _ := filepath.Split(exe)

	HashcatDir = filepath.Join(exeDir, "hashcat")
	err = os.MkdirAll(HashcatDir, 0o755)
	if err != nil {
		panic("unable to create hashcat dir")
	}

	HashesDir = filepath.Join(HashcatDir, "hashes")
	err = os.MkdirAll(HashesDir, 0o755)
	if err != nil {
		panic("unable to create hashes dir")
	}

	DictionariesDir = filepath.Join(HashcatDir, "dictionaries")
	err = os.MkdirAll(DictionariesDir, 0o755)
	if err != nil {
		panic("unable to create dictionaries dir")
	}

	RulesDir = filepath.Join(HashcatDir, "rules")
	err = os.MkdirAll(RulesDir, 0o755)
	if err != nil {
		panic("unable to create rules dir")
	}

	MasksDir = filepath.Join(HashcatDir, "masks")
	err = os.MkdirAll(MasksDir, 0o755)
	if err != nil {
		panic("unable to create masks dir")
	}

	a.Hashcat = &Hashcat{}
	if runtime.GOOS == "windows" {
		a.Hashcat.BinaryFile = filepath.Join(HashcatDir, "hashcat.exe")
	} else {
		a.Hashcat.BinaryFile = filepath.Join(HashcatDir, "hashcat.bin")
	}

	if err := a.Scan(); err != nil {
		log.Println(err)
	}

	a.Tasks = make(map[string]*Task)
	a.TaskUpdateCallback = func(taskUpdate TaskUpdate) {
		a.UI.Eval(`eventBus.dispatch("taskUpdate",` + MarshalJSONS(taskUpdate) + `)`)
	}
	a.TaskPreProcessCallback = func() {}
	a.TaskPostProcessCallback = func() {}
}

func (a *App) Clean() {
}

func (a *App) Scan() (err error) {
	a.Hashcat.GetAlgorithms()

	if err = a.ScanHashes(); err != nil {
		return
	}
	if err = a.ScanDictionaries(); err != nil {
		return
	}
	if err = a.ScanRules(); err != nil {
		return
	}
	if err = a.ScanMasks(); err != nil {
		return
	}

	return
}

func NewApp() *App {
	return &App{}
}
