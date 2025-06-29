package inotify

import (
	"fmt"
	"path/filepath"
	"syscall"
	"unsafe"
)

const _IN_CLOEXEC = 0x80000

type InotifyEvent uint32

const (
	IN_ACCESS        InotifyEvent = 0x00000001
	IN_MODIFY        InotifyEvent = 0x00000002
	IN_ATTRIB        InotifyEvent = 0x00000004
	IN_CLOSE_WRITE   InotifyEvent = 0x00000008
	IN_CLOSE_NOWRITE InotifyEvent = 0x00000010
	IN_OPEN          InotifyEvent = 0x00000020
	IN_MOVED_FROM    InotifyEvent = 0x00000040
	IN_MOVED_TO      InotifyEvent = 0x00000080
	IN_CREATE        InotifyEvent = 0x00000100
	IN_DELETE        InotifyEvent = 0x00000200
	IN_DELETE_SELF   InotifyEvent = 0x00000400
	IN_MOVE_SELF     InotifyEvent = 0x00000800
)

type inotifyEvent struct {
	Wd     int32
	Mask   uint32
	Cookie uint32
	Len    uint32
}

type Monitor struct {
	fd  int
	wds map[int]string
}

func New() (*Monitor, error) {
	fd, _, errno := syscall.Syscall(syscall.SYS_INOTIFY_INIT1, _IN_CLOEXEC, 0, 0)
	if errno != 0 {
		return nil, errno
	}
	return &Monitor{
		fd:  int(fd),
		wds: map[int]string{},
	}, nil
}

func (m *Monitor) Close() error {
	if m.fd > 0 {
		return syscall.Close(m.fd)
	}
	return nil
}

func (m *Monitor) AddWatch(f string, ev InotifyEvent) error {
	fptr, err := syscall.BytePtrFromString(f)
	if err != nil {
		return err
	}

	wd, _, errno := syscall.Syscall(syscall.SYS_INOTIFY_ADD_WATCH, uintptr(m.fd), uintptr(unsafe.Pointer(fptr)), uintptr(ev))
	if errno != 0 {
		return errno
	}
	m.wds[int(wd)] = f
	return nil
}

func (m *Monitor) Listen(handler func(f string, ev InotifyEvent) error) error {
	buf := make([]byte, 4096)

	for {
		n, err := syscall.Read(m.fd, buf)
		if err != nil {
			return err
		}

		pbuf := buf[:n]
		offset := 0

		for offset < len(pbuf) {
			event := (*inotifyEvent)(unsafe.Pointer(&pbuf[offset]))
			offset += 16 + int(event.Len)

			f, ok := m.wds[int(event.Wd)]
			if !ok {
				return fmt.Errorf("inotify: file to retrieve watched file/directory from wd: %d", event.Wd)
			}

			if event.Len > 0 {
				fnameB := pbuf[offset : offset+int(event.Len)]
				f = filepath.Join(f, string(fnameB[:clen(fnameB)]))
			}

			if err := handler(f, InotifyEvent(event.Mask)); err != nil {
				return err
			}
		}
	}
}

func clen(b []byte) int {
	for i := range b {
		if b[i] == 0 {
			return i
		}
	}
	return len(b)
}
