package uinput

import (
	"syscall"
	"time"
)

func getNowTimeval() syscall.Timeval {
	now := time.Now()
	return syscall.Timeval{
		Sec:  int32(now.Unix()),
		Usec: int32(now.UnixNano() / 1000 % 1000000),
	}
}
