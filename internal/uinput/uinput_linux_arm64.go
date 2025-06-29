package uinput

import (
	"syscall"
	"time"
)

func getNowTimeval() syscall.Timeval {
	// mister fpga is arm 32 bits only
	// build on 64 bits to make it easier to develop on rpi5
	now := time.Now()
	return syscall.Timeval{
		Sec:  now.Unix(),
		Usec: now.UnixNano() / 1000 % 1000000,
	}
}
