package main

import (
	"fmt"
	"strconv"
	"strings"
)

const CRLF = "\r\n"

type (
	Version struct {
		Mirror, Major int
	}

	StatusCode int
	Headers    map[string]string

	Status struct {
		Version Version
		Code    StatusCode
	}

	Response struct {
	}

	HTTP struct {
		Status   Status
		Headers  Headers
		Response Response
	}
)

func (h HTTP) String() string {
	return fmt.Sprint(h.Status, CRLF, h.Headers, CRLF)
}

func (s Status) String() string {
	return fmt.Sprint(s.Version, s.Code)
}

func (h Headers) String() string {
	var builder strings.Builder
	for k, v := range h {
		builder.WriteString(fmt.Sprintf("%s: %s\n", k, v))
	}

	return builder.String()
}

func (v Version) String() string {
	return fmt.Sprintf("HTTP/%d.%d", v.Major, v.Major)
}

func (s StatusCode) String() string {
	reason := "OK"

	switch s {
	case 200:
		reason = "OK"
	default:
		reason = ""
	}
	return fmt.Sprint(strconv.Itoa(int(s)), " ", reason)
}
