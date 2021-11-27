package main

import (
	"fmt"
	"tinyio/internal/app"
)

func main() {
	app.Boot("127.0.0.1:8877", func(in, out []byte) error {
		fmt.Println(in)
		return nil
	})
}
