package http

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
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
	textType           = "text/plain"
	appType            = "application/octet-stream"
	contentTypeKey     = "Content-Type"
	acceptEncodingKey  = "Accept-Encoding"
	contentEncodingKey = "Content-Encoding"
	UserAgentKey       = "User-Agent"
	contentLengthKey   = "Content-Length"
)

type Content struct {
	Bytes           []byte
	ContentType     string
	ContentEncoding string
}

func getContentTypeAndLength(body any, encoding string) (*Content, error) {
	if body == nil {
		return nil, nil
	}

	switch body.(type) {
	case string:
		data, contentEncoding, err := encode(encoding, []byte(body.(string)))
		return &Content{
			ContentType:     textType,
			Bytes:           data,
			ContentEncoding: contentEncoding,
		}, err
	case []byte:
		data, contentEncoding, err := encode(encoding, body.([]byte))
		return &Content{
			ContentType:     textType,
			Bytes:           data,
			ContentEncoding: contentEncoding,
		}, err
	case *os.File:
		data, err := io.ReadAll(body.(*os.File))
		if err != nil {
			return nil, err
		}

		data, contentEncoding, err := encode(encoding, data)
		return &Content{
			ContentType:     appType,
			ContentEncoding: contentEncoding,
			Bytes:           data,
		}, err
	default:
		return nil, fmt.Errorf("unsupported content type: %T", body)
	}
}

func encode(encoding string, data []byte) ([]byte, string, error) {
	if encoding == "" {
		return data, "", nil
	}

	if encoding == "gzip" {
		var buf = new(bytes.Buffer)
		gzipWriter := gzip.NewWriter(buf)

		_, err := gzipWriter.Write(data)
		if err != nil {
			return nil, "", err
		}
		err = gzipWriter.Close()
		if err != nil {
			return nil, "", err
		}
		return buf.Bytes(), "gzip", nil
	}

	return data, "", nil
}

func NewResponse(code int, body any, requestHeaders Headers) (*Response, error) {
	headers := make(map[string]string)

	encoding := requestHeaders.Get(acceptEncodingKey)
	content, err := getContentTypeAndLength(body, encoding)
	if err != nil {
		return nil, err
	}

	if content != nil {
		headers[contentTypeKey] = content.ContentType
		headers[contentLengthKey] = strconv.Itoa(len(content.Bytes))
	}

	if content != nil && content.ContentEncoding != "" {
		headers[contentEncodingKey] = content.ContentEncoding
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
	}, nil
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
	case 201:
		reason = "Created"
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
