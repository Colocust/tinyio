package main

import (
	"net"
	"tinyio/internal"
)

func main() {
	eventLoop := internal.NewEventLoop()

	var (
		ln  net.Listener
		err error
	)
	if ln, err = net.Listen("tcp", "127.0.0.1:8877"); err != nil {
		panic(err)
	}

	netLn, _ := ln.(*net.TCPListener)
	f, _ := netLn.File()
	fd := int(f.Fd())

	_, _ = eventLoop.NewEvent(fd, internal.Accept, nil)

	eventLoop.Process()
}
