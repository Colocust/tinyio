package internal

import (
	"syscall"
)

type Epoll struct {
	epfd   int
	events []syscall.EpollEvent
}

func (ep *Epoll) Create() {
	if fd, err := syscall.EpollCreate1(0); err != nil {
		panic(err)
	} else {
		ep.epfd = fd
	}
}

func (ep *Epoll) Add(e *Event) {
	if err := syscall.EpollCtl(ep.epfd, syscall.EPOLL_CTL_ADD, e.fd,
		&syscall.EpollEvent{Fd: int32(e.fd),
			Events: syscall.EPOLLIN,
		},
	); err != nil {
		panic(err)
	}
}

func (ep *Epoll) Delete() {

}

func (ep *Epoll) Poll(iter func(fd int) error) error {
	n, err := syscall.EpollWait(ep.epfd, ep.events, 1000)
	if err != nil && err != syscall.EINTR {
		return err
	}

	for i := 0; i < n; i++ {
		if err = iter(int(ep.events[i].Fd)); err != nil {
			return err
		}
	}

	return nil
}
