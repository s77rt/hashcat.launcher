package subprocess

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
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
	c.Dir = p.WDir

	stdin, _ := c.StdinPipe()
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
	return 0, fmt.Errorf("unsupported os")
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
