package http

import (
	"fmt"
	"strconv"
	"strings"
)

type (
	RequestLine struct {
		Method  string
		Target  string
		Version Version
	}

	Request struct {
		RequestLine RequestLine
		Headers     Headers
		Body        string
	}
)

func NewRequest(s string) (*Request, error) {
	println(s)

	sp := strings.Split(s, CRLF)
	requestLine, err := parseRequestLine(sp[0])
	if err != nil {
		return nil, err
	}

	sp = sp[1:]

	body := sp[len(sp)-1]
	sp = sp[:len(sp)-2]

	headers := make(Headers)
	for _, s := range sp {
		sp := strings.Split(s, ": ")
		headers[sp[0]] = sp[1]
	}

	var res []rune
	for _, v := range body {
		if v == '\x00' {
			continue
		}
		res = append(res, v)
	}

	return &Request{
		RequestLine: requestLine,
		Headers:     headers,
		Body:        string(res),
	}, nil
}

func parseRequestLine(s string) (RequestLine, error) {
	sp := strings.Split(s, " ")

	method := sp[0]
	target := sp[1]
	httpVersion := sp[2]

	version, err := parseVersion(httpVersion)
	if err != nil {
		return RequestLine{}, err
	}

	return RequestLine{
		Method:  method,
		Target:  target,
		Version: version,
	}, nil
}

func parseHeader(s string) (Headers, error) {
	return Headers{}, nil
}

func parseVersion(httpVersion string) (Version, error) {
	httpVersionSp := strings.Split(httpVersion, "/")
	httpVersionMajorMinor := strings.Split(httpVersionSp[1], ".")

	major, err := strconv.Atoi(httpVersionMajorMinor[0])
	if err != nil {
		return Version{}, err
	}

	minor, err := strconv.Atoi(httpVersionMajorMinor[1])
	if err != nil {
		return Version{}, err
	}

	return Version{
		Major: major,
		Minor: minor,
	}, nil
}

func (h Request) String() string {
	return fmt.Sprint(h.RequestLine, CRLF, h.Headers, CRLF, string(h.Body))
}

func (r RequestLine) String() string {
	return fmt.Sprint(r.Method, " ", r.Target, " ", r.Version.String())
}
