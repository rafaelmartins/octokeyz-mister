package process

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

type Process int

func New(pidfile string) (Process, error) {
	if _, err := os.Stat(pidfile); err == nil {
		pidB, err := os.ReadFile(pidfile)
		if err != nil {
			return 0, err
		}

		pid, err := strconv.Atoi(strings.TrimSpace(string(pidB)))
		if err != nil {
			return 0, err
		}

		proc, err := os.FindProcess(pid)
		if err != nil {
			return 0, err
		}

		if err := proc.Signal(syscall.Signal(0)); err == nil {
			return Process(pid), nil
		}
	}

	if err := os.WriteFile(pidfile, fmt.Appendf(nil, "%d", os.Getpid()), 0777); err != nil {
		return 0, err
	}
	return -1, nil
}

func (p Process) Kill() error {
	if p == -1 || p == 0 {
		return nil
	}

	proc, err := os.FindProcess(int(p))
	if err != nil {
		return err
	}

	return proc.Signal(syscall.SIGTERM)
}
