//go:build linux && !arm

package uinput

import (
	"syscall"
	"time"
)

func getNowTimeval() syscall.Timeval {
	now := time.Now()
	return syscall.Timeval{
		Sec:  now.Unix(),
		Usec: now.UnixNano() / 1000 % 1000000,
	}
}
