package http

import "testing"

func TestNewRequest(t *testing.T) {
	type args struct {
		s string
	}

	tests := []struct {
		name string
		args args
		want Request
	}{
		{
			name: "GET with target",
			args: args{s: "GET /index.html HTTP/1.1\r\n\r\n"},
			want: Request{
				RequestLine: RequestLine{
					Method: "GET",
					Target: "/index.html",
					Version: Version{
						Minor: 1,
						Major: 1,
					},
				},
			},
		},
		{
			name: "GET with target",
			args: args{s: "GET /index.html HTTP/1.1\r\nHost: localhost:4221\r\nUser-Agent: curl/7.64.1\r\nAccept: */*\r\n\r\n"},
			want: Request{
				RequestLine: RequestLine{
					Method: "GET",
					Target: "/index.html",
					Version: Version{
						Minor: 1,
						Major: 1,
					},
				},
				Headers: map[string]string{
					"Host":       "localhost:4221",
					"User-Agent": "curl/7.64.1",
					"Accept":     "*/*",
				},
			},
		},
		{
			name: "POST with body",
			args: args{s: "GET /index.html HTTP/1.1\r\nAccept: */*\r\n123\r\n"},
			want: Request{
				RequestLine: RequestLine{
					Method: "GET",
					Target: "/index.html",
					Version: Version{
						Minor: 1,
						Major: 1,
					},
				},
				Headers: map[string]string{
					"Accept": "*/*",
				},
				Body: []byte("123"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req, err := NewRequest(tt.args.s)
			if err != nil {
				t.Fatal(err)
			}

			if req.String() != tt.want.String() {
				t.Fatal(req.String(), "\n", tt.want.String())
			}
		})
	}
}
