package main

import (
	"github.com/Colocust/tinyio/app"
)

func main() {
	app.Boot("127.0.0.1:8877", func(in []byte) (out []byte) {
		out = in
		return
	})
}
