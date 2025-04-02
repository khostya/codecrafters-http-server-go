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

func acceptConnection(conn net.Conn) {
	for {
		req, err := read(conn)
		if err != nil {
			conn.Write([]byte(err.Error()))
			continue
		}
		fmt.Println(req)

		var resp *http.Response
		switch {
		case strings.HasPrefix(req.RequestLine.Target, "/echo"):
			data := strings.TrimLeft(req.RequestLine.Target, "/echo")
			if len(data) > 0 && data[0] == '/' {
				data = data[1:]
			}

			resp = http.NewResponse(200, data)
		case req.RequestLine.Target == "/user-agent":
			resp = http.NewResponse(200, req.Headers.Get(http.UserAgentKey))
		case req.RequestLine.Target == "/":
			write(conn, http.NewResponse(200, nil))
		default:
			write(conn, http.NewResponse(404, nil))
		}

		write(conn, resp)
		conn.Close()
		return
	}
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
