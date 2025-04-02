package main

import (
	"fmt"
	"net"
	"os"
	"slices"
)

var dir string

func init() {
	args := os.Args[1:]
	idx := slices.Index(args, "--directory")
	if idx != -1 && idx+1 < len(args) {
		dir = args[idx+1]
	}
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	listen(l)
}

func listen(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go acceptConnection(conn)
	}
}
