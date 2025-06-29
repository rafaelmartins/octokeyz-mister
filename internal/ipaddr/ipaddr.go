package ipaddr

import (
	"net"
	"time"
)

type Monitor struct {
	itfs []*net.Interface
	ch   chan struct{}
}

func NewMonitor(itfs ...string) (*Monitor, error) {
	rv := &Monitor{
		ch: make(chan struct{}),
	}
	for _, itf := range itfs {
		i, err := net.InterfaceByName(itf)
		if err != nil {
			return nil, err
		}
		rv.itfs = append(rv.itfs, i)
	}
	return rv, nil
}

func (m *Monitor) Close() error {
	if m.ch != nil {
		close(m.ch)
		m.ch = nil
	}
	return nil
}

func (m *Monitor) Run(handler func(itf string, ip net.IP) error) error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		for _, itf := range m.itfs {
			addrs, err := itf.Addrs()
			if err != nil {
				return err
			}

			found := false
			for _, addr := range addrs {
				if netaddr, ok := addr.(*net.IPNet); ok {
					if ip := netaddr.IP.To4(); ip != nil {
						if err := handler(itf.Name, ip); err != nil {
							return err
						}
						found = true
						break
					}
				}
			}
			if !found {
				if err := handler(itf.Name, nil); err != nil {
					return err
				}
			}
		}

		select {
		case <-m.ch:
			return nil
		case <-ticker.C:
			continue
		}
	}
}
