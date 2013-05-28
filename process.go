package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strconv"
	"sync"
)

func startProcess(cmd string, ch chan<- []byte, callback func(error)) {
	command := exec.Command("/bin/bash", "-c", cmd)
	stdout, stdoutErr := command.StdoutPipe()
	stderr, stderrErr := command.StderrPipe()
	if stdoutErr != nil || stderrErr != nil {
		fmt.Fprintf(os.Stderr, "Error getting stdout/stderr of the command: %v, %v\n", stdoutErr, stderrErr)
		return
	}
	command.Start()
	var logging io.WriteCloser
	if "" != logDir {
		var err error
		logging, err = os.Create(path.Join(logDir, strconv.Itoa(command.Process.Pid)))
		fmt.Fprintf(logging, "Command: %s\n\n", cmd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Open file error: %v\n", err)
			logging = nil
		}
	}
	lineCp := func(reader io.ReadCloser, ch chan<- []byte, wg *sync.WaitGroup) {
		var rd *bufio.Reader
		if logging != nil {
			rd = bufio.NewReader(io.TeeReader(reader, logging))
		} else {
			rd = bufio.NewReader(reader)
		}
		var err error
		var buf []byte
		for err == nil {
			buf, err = rd.ReadBytes('\n')
			if err == nil {
				ch <- buf
			}
		}
		wg.Done()
	}
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go lineCp(stdout, ch, wg)
	go lineCp(stderr, ch, wg)
	go func() {
		wg.Wait()
		close(ch)
		callback(command.Wait())
	}()
}
