package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"testing"
)

var TestBiliClient *BiliClient

var jsonStr string = `
{
    "code": 0,
    "message": "0",
    "ttl": 1,
    "data": {
        "group": "live",
        "business_id": 0,
        "refresh_row_factor": 0.125,
        "refresh_rate": 100,
        "max_delay": 5000,
        "token": "loyJrzwerIdMul24uxHbmZl-c7t72qXiawC6M7ZHOCxuFkW8ZkEGgA5fPcOyRUSP4ok2wiYxiUp3rReZjBvEgdkhtTdvtJ6DUMAVYz7AuURv6PXdXMd5VqzfuLh40WVTXdfCa4SnNyIxKekP4ODJe8Jf",
        "host_list": [
            {
                "host": "dsa-cn-live-comet-01.chat.bilibili.com",
                "port": 2243,
                "wss_port": 2245,
                "ws_port": 2244
            },
            {
                "host": "hw-sh-live-comet-04.chat.bilibili.com",
                "port": 2243,
                "wss_port": 2245,
                "ws_port": 2244
            },
            {
                "host": "tx-gz-live-comet-02.chat.bilibili.com",
                "port": 2243,
                "wss_port": 2245,
                "ws_port": 2244
            },
            {
                "host": "broadcastlv.chat.bilibili.com",
                "port": 2243,
                "wss_port": 2245,
                "ws_port": 2244
            }
        ]
    }
}
`

func TestMain(m *testing.M) {
	TestBiliClient = NewBiliClient(123)
	TestBiliClient.HTTPClient = client

	os.Exit(m.Run())
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestHTTPClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

var client = NewTestHTTPClient(func(req *http.Request) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(jsonStr)),
		Header:     make(http.Header),
	}
})
