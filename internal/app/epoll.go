package app

import (
	"fmt"
	"log"
	"syscall"
)

type Epoll struct {
	epfd   int
	events []syscall.EpollEvent
}

func (ep *Epoll) create() {
	if fd, err := syscall.EpollCreate1(0); err != nil {
		panic(err)
	} else {
		ep.epfd = fd
	}
}

func (ep *Epoll) add(fd int) (err error) {
	if err = syscall.EpollCtl(ep.epfd, syscall.EPOLL_CTL_ADD, fd,
		&syscall.EpollEvent{Fd: int32(fd),
			Events: syscall.EPOLLIN,
		},
	); err != nil {
		return
	}
	return
}

func (ep *Epoll) mod(fd int, flag uint32) (err error) {
	if err = syscall.EpollCtl(ep.epfd, syscall.EPOLL_CTL_MOD, fd,
		&syscall.EpollEvent{Fd: int32(fd),
			Events: flag,
		},
	); err != nil {
		return
	}
	return
}

func (ep *Epoll) delete(fd int) (err error) {
	if err = syscall.EpollCtl(ep.epfd, syscall.EPOLL_CTL_DEL, fd,
		&syscall.EpollEvent{Fd: int32(fd),
			Events: syscall.EPOLLIN | syscall.EPOLLOUT,
		},
	); err != nil {
		return
	}
	return
}

func (ep *Epoll) poll(iter func(fd int)) {
	events := make([]syscall.EpollEvent, 64)
	n, err := syscall.EpollWait(ep.epfd, events, 5000)

	if err != nil && err != syscall.EINTR {
		log.Println(fmt.Println("ERR:Poll Err " + err.Error()))
		return
	}

	for i := 0; i < n; i++ {
		iter(int(events[i].Fd))
	}
	return
}
