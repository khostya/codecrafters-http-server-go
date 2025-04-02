package main

import (
	"fmt"
	"github.com/codecrafters-io/http-server-starter-go/app/http"
	"net"
	"os"
	"strings"
)

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
			data := trimLeftUrl(req.RequestLine.Target, "/echo")
			resp = http.NewResponse(200, data)
		case req.RequestLine.Target == "/user-agent":
			resp = http.NewResponse(200, req.Headers.Get(http.UserAgentKey))
		case req.RequestLine.Target == "/":
			resp = http.NewResponse(200, nil)
		case strings.HasPrefix(req.RequestLine.Target, "/files"):
			data := trimLeftUrl(req.RequestLine.Target, "/files")
			filePath := dir + data
			file, err := os.Open(filePath)
			if err != nil {
				if os.IsNotExist(err) {
					resp = http.NewResponse(404, nil)
				} else {
					resp = http.NewResponse(500, nil)
				}
				continue
			}

			resp = http.NewResponse(200, file)
			file.Close()
		default:
			resp = http.NewResponse(404, nil)
		}

		write(conn, resp)
		conn.Close()
		return
	}
}

func trimLeftUrl(url string, cutset string) string {
	data, _ := strings.CutPrefix(url, cutset)
	if len(data) > 0 && data[0] == '/' {
		data = data[1:]
	}
	return data
}
