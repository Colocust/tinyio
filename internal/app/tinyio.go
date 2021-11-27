package app

import (
	"net"
	"os"
	"syscall"
)

const (
	NONE = iota
	READABLE
	WRITABLE
)

type (
	eventLoop struct {
		events map[int]*event
		stop   bool
		poll   *Epoll
		iter   func(in, out []byte) error
	}

	event struct {
		fd   int
		mask int
		proc Proc
		data *Connection
	}

	Proc func(eventLoop *eventLoop, event *event) error
)

func newEventLoop(iter func(in, out []byte) error) *eventLoop {
	el := &eventLoop{
		events: make(map[int]*event),
		stop:   false,
		iter:   iter,
		poll: &Epoll{
			events: make([]syscall.EpollEvent, 64),
		},
	}
	el.poll.create()
	return el
}

func (eventLoop *eventLoop) newEvent(fd int, proc Proc, data *Connection) error {
	e := &event{
		fd:   fd,
		mask: READABLE,
		proc: proc,
		data: data,
	}

	if err := eventLoop.poll.add(fd); err != nil {
		return err
	}
	eventLoop.events[fd] = e
	return nil
}

func (eventLoop *eventLoop) process() {
	for !eventLoop.stop {
		if err := eventLoop.poll.poll(func(fd int) error {
			e := eventLoop.events[fd]

			if e.mask&READABLE != 0 {
				if err := e.proc(eventLoop, e); err != nil {
					return err
				}
			}

			if e.mask&WRITABLE != 0 {
				if err := write(eventLoop, e); err != nil {
					return err
				}
			}

			return nil
		}); err != nil {
			eventLoop.stop = true
		}
	}
}

func Boot(addr string, iter func(in, out []byte) error) {
	if err := boot(addr, iter); err != nil {
		panic(err)
	}
}

func boot(addr string, iter func(in, out []byte) error) (err error) {
	var (
		ln net.Listener
		fd int
		f  *os.File
	)
	if ln, err = net.Listen("tcp", addr); err != nil {
		return
	}
	defer ln.Close()

	netLn := ln.(*net.TCPListener)
	if f, err = netLn.File(); err != nil {
		return
	}
	fd = int(f.Fd())

	if err = syscall.SetNonblock(fd, true); err != nil {
		return
	}

	el := newEventLoop(iter)
	if err = el.newEvent(fd, accept, nil); err != nil {
		return
	}

	el.process()
	return
}
