package app

import (
	"net"
	"syscall"
	"tinyio/internal"
)

type Connection struct {
	addr    net.Addr
	in, out []byte
}

func accept(eventLoop *eventLoop, event *event) error {
	nfd, sa, err := syscall.Accept(event.fd)
	if err != nil {
		if err == syscall.EAGAIN {
			return nil
		}
		return err
	}

	conn := &Connection{
		addr: internal.SocketAddrToNetAddr(sa),
	}

	if err = syscall.SetNonblock(nfd, true); err != nil {
		return err
	}

	if err = eventLoop.newEvent(nfd, read, conn); err != nil {
		return err
	}
	return nil
}

func read(eventLoop *eventLoop, event *event) error {
	var (
		out []byte
		n   int
		err error
	)
	in := make([]byte, 0xFFFF)
	n, err = syscall.Read(event.fd, in)
	if n == 0 || err != nil {
		if err == syscall.EAGAIN {
			return nil
		}
		return close(eventLoop, event)
	}

	if err = eventLoop.iter(event.data.in, out); err != nil {
		return err
	}
	return nil
}

func write(eventLoop *eventLoop, event *event) error {
	return nil
}

func close(eventLoop *eventLoop, event *event) error {
	return nil
}
