package internal

import (
	"syscall"
)

const (
	NONE = iota
	READABLE
	WRITABLE
)

type (
	EventLoop struct {
		events map[int]*Event
		stop   bool
		Poll   Poll
	}

	Event struct {
		fd   int
		mask int
		Proc Proc
		Data interface{}
	}

	Proc func(eventLoop *EventLoop, fd int, data interface{}, mask int) error

	Poll interface {
		Create()
		Add(e *Event)
		Delete()
		Poll(iter func(fd int) error) error
	}
)

func NewEventLoop() *EventLoop {
	eventLoop := &EventLoop{
		events: make(map[int]*Event),
		stop:   false,
	}
	eventLoop.system()
	eventLoop.Poll.Create()
	return eventLoop
}

func (eventLoop *EventLoop) NewEvent(fd int, mask int, proc Proc, data interface{}) *Event {
	e := &Event{
		fd:   fd,
		mask: mask,
		Data: data,
		Proc: proc,
	}
	eventLoop.Poll.Add(e)
	eventLoop.events[fd] = e
	return e
}

func (eventLoop *EventLoop) system() {
	eventLoop.Poll = &Epoll{
		events: make([]syscall.EpollEvent, 64),
	}
}

func (eventLoop *EventLoop) Process() {
	for !eventLoop.stop {
		eventLoop.process()
	}
}

func (eventLoop *EventLoop) process() {
	if err := eventLoop.Poll.Poll(func(fd int) error {
		e := eventLoop.events[fd]

		// 判断当前mask

		return e.Proc(eventLoop, e.fd, e.Data, e.mask)
	}); err != nil {
		eventLoop.stop = true
	}
	return
}
