package app

import (
	"log"
	"net"
	"syscall"
	"tinyio/internal"
)

type Connection struct {
	addr    net.Addr
	in, out []byte
}

func accept(eventLoop *eventLoop, event *event) {
	nfd, sa, err := syscall.Accept(event.fd)
	if err != nil {
		if err == syscall.EAGAIN {
			return
		}
		log.Println("ERR:Accept Err " + err.Error())
		return
	}
	if err = syscall.SetNonblock(nfd, true); err != nil {
		return
	}
	if err = eventLoop.newEvent(nfd, &Connection{
		addr: internal.SocketAddrToNetAddr(sa),
	}); err != nil {
		log.Println("ERR:New Event Err " + err.Error())
		return
	}
	return
}

func read(eventLoop *eventLoop, event *event) {
	var (
		out []byte
		n   int
		err error
	)
	in := make([]byte, 0xFFFF)
	n, err = syscall.Read(event.fd, in)
	if n == 0 || err != nil {
		if err == syscall.EAGAIN {
			return
		}
		close(eventLoop, event)
	}
	eventLoop.iter(in, out)
}

func write(eventLoop *eventLoop, event *event) {

}

func close(eventLoop *eventLoop, event *event) {

}
