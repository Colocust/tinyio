package app

import (
	"log"
	"net"
	"os"
	"syscall"
)

type (
	event struct {
		fd      int
		in, out []byte
	}

	eventLoop struct {
		events map[int]*event
		iter   iter
		lnFd   int
		stop   bool
		poll   *Epoll
	}

	iter func(in []byte) (out []byte)
)

func newEventLoop(lnFd int, iter iter) *eventLoop {
	el := &eventLoop{
		events: make(map[int]*event),
		lnFd:   lnFd,
		iter:   iter,
		stop:   false,
		poll: &Epoll{
			events: make([]syscall.EpollEvent, 64),
		},
	}

	el.poll.create()
	return el
}

func (el *eventLoop) addEvent(fd int) (err error) {
	if err = syscall.SetNonblock(fd, true); err != nil {
		return
	}
	el.events[fd] = &event{
		fd: fd,
	}
	el.poll.add(fd)
	return
}

func (e *event) accept(el *eventLoop) {
	nfd, _, err := syscall.Accept(el.lnFd)
	if err != nil {
		return
	}
	el.addEvent(nfd)
}

func (e *event) read(el *eventLoop) {
	in := make([]byte, 0xFFFF)
	var (
		err error
		n   int
		out []byte
	)
	if n, err = syscall.Read(e.fd, in); n == 0 || err != nil {
		if err == syscall.EAGAIN {
			return
		}
		e.close(el)
	}

	out = el.iter(in)
	if len(out) > 0 {
		e.out = append(e.out, out...)
	}
	if len(e.out) > 0 {
		el.poll.modReadWrite(e.fd)
	}
}

func (e *event) write(el *eventLoop) {
	var (
		n   int
		err error
	)
	if n, err = syscall.Write(e.fd, e.out); err != nil {
		if err == syscall.EAGAIN {
			return
		}
		e.close(el)
	}
	if n == len(e.out) {
		e.out = e.out[:0]
	} else {
		e.out = e.out[n:]
	}

	if len(e.out) == 0 {
		el.poll.modRead(e.fd)
	}
}

func (e *event) close(el *eventLoop) {
	delete(el.events, e.fd)
	syscall.Close(e.fd)
}

func (el *eventLoop) boot() {
	for !el.stop {
		el.poll.poll(func(fd int) {
			e := el.events[fd]
			switch {
			case e.fd == el.lnFd:
				e.accept(el)
			case len(e.out) > 0:
				e.write(el)
			default:
				e.read(el)
			}
		})
	}
}

func Boot(address string, iter iter) (err error) {
	var (
		ln net.Listener
		f  *os.File
	)

	if ln, err = net.Listen("tcp", address); err != nil {
		return
	}
	defer ln.Close()

	tcpLn := ln.(*net.TCPListener)
	if f, err = tcpLn.File(); err != nil {
		return
	}

	lnFd := int(f.Fd())

	el := newEventLoop(lnFd, iter)
	if err = el.addEvent(lnFd); err != nil {
		return
	}
	log.Println("tiny io boot success")
	el.boot()
	return
}
