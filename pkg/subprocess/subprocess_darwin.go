package subprocess

import (
	"os"
	"os/exec"
	"io"
	"bufio"
	"fmt"
)

type SubprocessStatus int
const (
	SubprocessStatusNotRunning SubprocessStatus = iota
	SubprocessStatusRunning
	SubprocessStatusFinished
)

type Subprocess struct {
	Status SubprocessStatus
	WDir string
	Program string
	Args []string
	Process *os.Process
	Stdin_stream io.WriteCloser
	Stdout_callback func(string)
	Stderr_callback func(string)
	Preprocess func()
	Postprocess func()
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
	}

	p.Status = SubprocessStatusRunning
	p.Process = c.Process
	p.Stdin_stream = stdin
	p.Preprocess()

	go func() {
		stdout_scanner := bufio.NewScanner(stdout)
		for stdout_scanner.Scan() {
			p.Stdout_callback(stdout_scanner.Text())
		}

		stderr_scanner := bufio.NewScanner(stderr)
		for stderr_scanner.Scan() {
			p.Stderr_callback(stderr_scanner.Text())
		}
	}()

	c.Wait()
	p.Status = SubprocessStatusFinished
	p.Postprocess()
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
