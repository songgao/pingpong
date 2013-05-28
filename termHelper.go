package main

import (
	"os"
	"syscall"
	"unsafe"
)

var tty *os.File

func init() {
	var err error
	tty, err = os.OpenFile("/dev/tty", syscall.O_WRONLY, 0)
	if err != nil {
		panic("Open /dev/tty failed")
	}
}

type termSize struct {
	rows uint16
	Cols uint16
	xpx  uint16
	ypx  uint16
}

func getTermWidth() int {
	var tmp termSize
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, tty.Fd(), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&tmp)))
	return int(tmp.Cols)
}
