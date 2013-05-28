package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const punchCardWidth = 80

var tty *os.File

func init() {
	var err error
	tty, err = os.OpenFile("/dev/tty", syscall.O_WRONLY, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Open /dev/tty failed. Falling back to punch card width (%d).\n", punchCardWidth)
		tty = nil
	}
}

type termSize struct {
	rows uint16
	Cols uint16
	xpx  uint16
	ypx  uint16
}

func getTermWidth() int {
	if tty == nil {
		return punchCardWidth
	} else {
		var tmp termSize
		_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, tty.Fd(), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&tmp)))
		return int(tmp.Cols)
	}
}
