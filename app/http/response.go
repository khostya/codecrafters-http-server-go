package http

import (
	"fmt"
	"strconv"
	"strings"
)

const CRLF = "\r\n"

type (
	Version struct {
		Minor, Major int
	}

	StatusCode int
	Headers    map[string]string

	Status struct {
		Version Version
		Code    StatusCode
	}

	Response struct {
		Status   Status
		Headers  Headers
		Response []byte
	}
)

const (
	textType         = "text/plain"
	contentTypeKey   = "Content-Type"
	UserAgentKey     = "User-Agent"
	contentLengthKey = "Content-Length"
)

type Content struct {
	Bytes       []byte
	ContentType string
}

func getContentTypeAndLength(body any) *Content {
	if body == nil {
		return nil
	}

	switch body.(type) {
	case string:
		return &Content{
			ContentType: textType,
			Bytes:       []byte(body.(string)),
		}
	case []byte:
		return &Content{
			ContentType: textType,
			Bytes:       body.([]byte),
		}
	default:
		return nil
	}
}

func NewResponse(code int, body any) *Response {
	headers := make(map[string]string)

	content := getContentTypeAndLength(body)
	if content != nil {
		headers[contentTypeKey] = content.ContentType
		headers[contentLengthKey] = strconv.Itoa(len(content.Bytes))
	}

	var bodyBytes []byte
	if content != nil {
		bodyBytes = content.Bytes
	}

	return &Response{
		Status: Status{
			Version: Version{1, 1},
			Code:    StatusCode(code),
		},
		Headers:  headers,
		Response: bodyBytes,
	}
}

func (h Response) String() string {
	return fmt.Sprint(h.Status, CRLF, h.Headers, CRLF, string(h.Response))
}

func (s Status) String() string {
	return fmt.Sprint(s.Version, s.Code)
}

func (h Headers) String() string {
	var builder strings.Builder
	for k, v := range h {
		builder.WriteString(fmt.Sprintf("%s: %s%s", k, v, CRLF))
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
	case 404:
		reason = "Not Found"
	default:
		reason = ""
	}
	return fmt.Sprint(strconv.Itoa(int(s)), " ", reason)
}

func (h Headers) Get(key string) string {
	_, ok := h[key]
	if !ok {
		return ""
	}
	return h[key]
}
