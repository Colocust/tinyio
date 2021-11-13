package internal

import (
	"fmt"
	"net"
	"syscall"
)

type Connection struct {
	addr    net.Addr
	in, out []byte
}

func Accept(eventLoop *EventLoop, fd int, data interface{}, mask int) error {
	nfd, sa, err := syscall.Accept(fd)
	if err != nil {
		if err == syscall.EAGAIN {
			return nil
		}
		return err
	}

	fmt.Println("s")

	eventLoop.NewEvent(nfd, READABLE, Read, &Connection{
		addr: SocketAddrToNetAddr(sa),
	})
	return nil
}

func Read(eventLoop *EventLoop, fd int, data interface{}, mask int) error {
	return nil
}

func Write(eventLoop *EventLoop, fd int, data interface{}, mask int) error {
	return nil
}

func Close(eventLoop *EventLoop, fd int, data interface{}, mask int) error {
	return nil
}
