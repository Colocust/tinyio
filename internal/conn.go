package internal

import (
	"net"
	"syscall"
)

type Connection struct {
	addr    net.Addr
	in, out []byte
}

func Accept(eventLoop *EventLoop, event *Event) error {
	nfd, sa, err := syscall.Accept(event.fd)
	if err != nil {
		if err == syscall.EAGAIN {
			return nil
		}
		return err
	}

	if err = syscall.SetNonblock(nfd, true); err != nil {
		return err
	}

	if _, err = eventLoop.NewEvent(nfd, Read, &Connection{
		addr: SocketAddrToNetAddr(sa),
	}); err != nil {
		return err
	}
	return nil
}

func Read(eventLoop *EventLoop, event *Event) error {
	return nil
}

func Write(eventLoop *EventLoop, event *Event) error {
	return nil
}

func Close(eventLoop *EventLoop, event *Event) error {
	return nil
}
