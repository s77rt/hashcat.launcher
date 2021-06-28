// +build ignore

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func version() {
	f, err := os.Create("version.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	out, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
	if err != nil {
		panic(err)
	}
	version := strings.TrimRight(string(out), "\r\n")

	w := bufio.NewWriter(f)
	_, err = w.WriteString(fmt.Sprintf("// Auto generated. DO NOT EDIT!\n\npackage hashcatlauncher\n\nconst Version = \"%s\"\n", version))
	if err != nil {
		panic(err)
	}
	w.Flush()
}

func main() {
	version()
}
