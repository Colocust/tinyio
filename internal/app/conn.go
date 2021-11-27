package app

import (
	"log"
	"net"
	"strconv"
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
	in := make([]byte, 0xFFFF)
	n, err := syscall.Read(event.fd, in)
	if n == 0 || err != nil {
		if err == syscall.EAGAIN {
			return
		}
		close(eventLoop, event)
	}
	event.data.in = in[:n]
	out := eventLoop.iter(event.data.in)

	if len(out) > 0 {
		event.data.out = append(event.data.out, out...)
	}
	if len(event.data.out) > 0 {
		if err := eventLoop.poll.modReadWrite(event.fd); err != nil {
			log.Println("ERR:Mod Read Write Err，Fd" + strconv.Itoa(event.fd))
		}
	}
}

func write(eventLoop *eventLoop, event *event) {
	n, err := syscall.Write(event.fd, event.data.out)
	if err != nil {
		if err == syscall.EAGAIN {
			return
		}
		close(eventLoop, event)
		return
	}

	if n == len(event.data.out) {
		event.data.out = event.data.out[:0]
	} else {
		event.data.out = event.data.out[n:]
	}

	if len(event.data.out) == 0 {
		eventLoop.poll.modRead(event.fd)
	}
}

func close(eventLoop *eventLoop, event *event) {
	delete(eventLoop.events, event.fd)
	syscall.Close(event.fd)
}
