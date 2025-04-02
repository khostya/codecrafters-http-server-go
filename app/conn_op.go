package main

import (
	"fmt"
	"github.com/codecrafters-io/http-server-starter-go/app/http"
	"net"
	"os"
	"strings"
	"sync"
)

var bufferPool sync.Pool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 1024)
	},
}

func write(conn net.Conn, data *http.Response) {
	_, err := conn.Write([]byte(data.String()))
	if err != nil {
		fmt.Println("Failed to write response: ", err.Error())
		os.Exit(1)
	}
}
func read(conn net.Conn) (*http.Request, error) {
	buffer := bufferPool.Get().([]byte)
	defer bufferPool.Put(buffer)

	res := strings.Builder{}
	_, err := conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	res.Write(buffer)
	return http.NewRequest(res.String())
}
