package subprocess

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"unsafe"
)

type SubprocessStatus int

const (
	SubprocessStatusNotStarted SubprocessStatus = iota
	SubprocessStatusRunning
	SubprocessStatusFinished
)

type Subprocess struct {
	Status         SubprocessStatus `json:"status"`
	WDir           string           `json:"-"`
	Program        string           `json:"-"`
	Args           []string         `json:"-"`
	Process        *os.Process      `json:"-"`
	StdinStream    io.WriteCloser   `json:"-"`
	StdoutCallback func(string)     `json:"-"`
	StderrCallback func(string)     `json:"-"`
	PreProcess     func()           `json:"-"`
	PostProcess    func()           `json:"-"`
}

func (p *Subprocess) Execute() {
	c := exec.Command(p.Program, p.Args...)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	c.Dir = p.WDir

	var stdin io.WriteCloser
	c.Stdin = os.Stdin
	stderr, _ := c.StderrPipe()
	stdout, _ := c.StdoutPipe()

	err := c.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't execute subprocess: %s\n", err)
		return
	}

	p.Status = SubprocessStatusRunning
	p.Process = c.Process
	p.StdinStream = stdin
	p.PreProcess()

	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		stdoutScanner := bufio.NewScanner(stdout)
		for stdoutScanner.Scan() {
			p.StdoutCallback(stdoutScanner.Text())
		}
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		stderrScanner := bufio.NewScanner(stderr)
		for stderrScanner.Scan() {
			p.StderrCallback(stderrScanner.Text())
		}
	}(&wg)

	wg.Wait()

	c.Wait()
	p.Status = SubprocessStatusFinished
	p.PostProcess()
}

func (p *Subprocess) PostKey(key uint8) (uintptr, error) {
	var user32 = syscall.NewLazyDLL("user32.dll")

	EnumWindows := func(enumFunc uintptr, lparam uintptr) {
		user32.NewProc("EnumWindows").Call(uintptr(enumFunc), uintptr(lparam))
	}

	var hwnd uintptr
	cb := syscall.NewCallback(func(h uintptr, prm uintptr) uintptr {
		var itrPid uint32
		itrPid = 0x0001
		user32.NewProc("GetWindowThreadProcessId").Call(uintptr(h), uintptr(unsafe.Pointer(&itrPid)))
		if int(itrPid) == p.Process.Pid {
			hwnd = h
			user32.NewProc("PostMessageW").Call(h, 0x0100, uintptr(key), 0)
			//return 0 // stop enumeration (commented to make sure all windows created by our process get's the message)
		}
		return 1 // continue enumeration
	})
	EnumWindows(cb, 0)
	if hwnd == 0 {
		return 0, fmt.Errorf("No window with pid %d found", p.Process.Pid)
	}
	return hwnd, nil
}

func (p *Subprocess) Kill() {
	if p.Process != nil {
		err := p.Process.Kill()
		if err != nil {
			if p.Status == SubprocessStatusRunning {
				fmt.Fprintf(os.Stderr, "can't kill subprocess: %s\n", err)
			}
		}
	}
}
