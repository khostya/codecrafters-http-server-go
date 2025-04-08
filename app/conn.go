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
			resp, err = http.NewResponse(200, data, req.Headers)
		case req.RequestLine.Target == "/user-agent":
			resp, err = http.NewResponse(200, req.Headers.Get(http.UserAgentKey), req.Headers)
		case req.RequestLine.Target == "/":
			resp, _ = http.NewResponse(200, nil, req.Headers)
		case strings.HasPrefix(req.RequestLine.Target, "/files") && req.RequestLine.Method == "GET":
			resp, err = filesGET(req)
		case strings.HasPrefix(req.RequestLine.Target, "/files") && req.RequestLine.Method == "POST":
			resp, err = filesPOST(req)
		default:
			resp, err = http.NewResponse(404, nil, req.Headers)
		}

		if err != nil {
			conn.Write([]byte(err.Error()))
		} else {
			write(conn, resp)
		}

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

func filesGET(req *http.Request) (*http.Response, error) {
	data := trimLeftUrl(req.RequestLine.Target, "/files")
	filePath := dir + data
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return http.NewResponse(404, nil, req.Headers)
		}
		return http.NewResponse(500, nil, req.Headers)
	}
	defer file.Close()
	return http.NewResponse(200, file, req.Headers)
}

func filesPOST(req *http.Request) (*http.Response, error) {
	data := trimLeftUrl(req.RequestLine.Target, "/files")
	filePath := dir + data
	file, err := os.Create(filePath)
	if err != nil {
		return http.NewResponse(500, nil, req.Headers)
	}
	_, err = file.WriteString(req.Body)
	if err != nil {
		return http.NewResponse(500, nil, req.Headers)
	}

	defer file.Close()
	return http.NewResponse(201, file, req.Headers)
}
