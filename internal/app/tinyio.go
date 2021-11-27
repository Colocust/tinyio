package app

import (
	"log"
	"net"
	"os"
	"syscall"
)

type (
	eventLoop struct {
		events map[int]*event
		poll   *Epoll
		iter   func(in []byte) (out []byte)
		lnFd   int
	}

	event struct {
		fd   int
		data *Connection
	}
)

func newEventLoop(lnFd int, iter func(in []byte) (out []byte)) *eventLoop {
	el := &eventLoop{
		lnFd:   lnFd,
		events: make(map[int]*event),
		iter:   iter,
		poll: &Epoll{
			events: make([]syscall.EpollEvent, 64),
		},
	}
	el.poll.create()
	return el
}

func (eventLoop *eventLoop) newEvent(fd int, data *Connection) error {
	e := &event{
		fd:   fd,
		data: data,
	}

	if err := eventLoop.poll.add(fd); err != nil {
		return err
	}
	eventLoop.events[fd] = e
	return nil
}

func (eventLoop *eventLoop) process() {
	for {
		eventLoop.poll.poll(func(fd int) {
			e := eventLoop.events[fd]
			switch {
			case e.fd == eventLoop.lnFd:
				accept(eventLoop, e)
				break
			case len(e.data.out) > 0:
				write(eventLoop, e)
				break
			default:
				read(eventLoop, e)
				break
			}
		})
	}
}

func Boot(addr string, iter func(in []byte) (out []byte)) {
	if err := boot(addr, iter); err != nil {
		panic(err)
	}
}

func boot(addr string, iter func(in []byte) (out []byte)) (err error) {
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

	el := newEventLoop(fd, iter)
	if err = el.newEvent(fd, nil); err != nil {
		return
	}

	log.Println("INFO：tiny io boot")
	el.process()
	return
}
