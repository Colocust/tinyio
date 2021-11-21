package internal

import (
	"net"
	"syscall"
)

const (
	NONE = iota
	READABLE
	WRITABLE
)

type (
	eventLoop struct {
		events map[int]*Event
		stop   bool
		poll   Poll
	}

	Event struct {
		fd        int
		mask      int
		readProc  Proc
		writeProc Proc
		data      *Connection
	}

	Proc func(eventLoop *eventLoop, event *Event) error

	Poll interface {
		create()
		add(fd int) error
		mod(fd int, flag uint32) error
		delete(fd int) error
		poll(iter func(fd int) error) error
	}
)

func newEventLoop() *eventLoop {
	eventLoop := &eventLoop{
		events: make(map[int]*Event),
		stop:   false,
	}
	eventLoop.system()
	eventLoop.poll.create()
	return eventLoop
}

func (eventLoop *eventLoop) NewEvent(fd int, proc Proc, data *Connection) (*Event, error) {
	e := &Event{
		fd:        fd,
		mask:      READABLE,
		readProc:  proc,
		writeProc: Write,
		data:      data,
	}

	if err := eventLoop.poll.add(fd); err != nil {
		return nil, err
	}
	eventLoop.events[fd] = e
	return e, nil
}

func (eventLoop *eventLoop) system() {
	eventLoop.poll = &Epoll{
		events: make([]syscall.EpollEvent, 64),
	}
}

func (eventLoop *eventLoop) process() {
	for !eventLoop.stop {
		if err := eventLoop.poll.poll(func(fd int) error {
			e := eventLoop.events[fd]

			if e.mask&READABLE != 0 {
				if err := e.readProc(eventLoop, e); err != nil {
					return err
				}
			}

			if e.mask&WRITABLE != 0 {
				if err := e.writeProc(eventLoop, e); err != nil {
					return err
				}
			}

			return nil
		}); err != nil {
			eventLoop.stop = true
		}
	}
}

func Serve(addr string) error {
	var (
		ln  net.Listener
		err error
	)
	if ln, err = net.Listen("tcp", addr); err != nil {
		return err
	}

	netLn, _ := ln.(*net.TCPListener)
	f, _ := netLn.File()
	fd := int(f.Fd())

	if err = syscall.SetNonblock(fd, true); err != nil {
		return err
	}

	eventLoop := newEventLoop()

	if _, err = eventLoop.NewEvent(fd, Accept, nil); err != nil {
		return err
	}

	eventLoop.process()
	return nil
}
