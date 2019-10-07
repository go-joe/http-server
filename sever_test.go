package joehttp

import (
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-joe/joe/joetest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestServer_HTTPHandler_GET(t *testing.T) {
	brain := joetest.NewBrain(t)
	conf := config{
		listenAddr: "127.0.0.1:0",
		logger:     zaptest.NewLogger(t),
	}
	s := newServer(conf, brain)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	s.HTTPHandler(resp, req)

	brain.Finish()
	events := brain.RecordedEvents()
	assert.NotEmpty(t, events)

	expectedEvt := RequestEvent{
		Header:     req.Header,
		RemoteAddr: "192.0.2.1:1234",
		Method:     "GET",
		URL:        req.URL,
		Body:       []byte{},
	}
	assert.Equal(t, expectedEvt, events[0])
}

func TestServer_TrustedHeader(t *testing.T) {
	brain := joetest.NewBrain(t)
	conf := config{
		listenAddr:    "127.0.0.1:0",
		logger:        zaptest.NewLogger(t),
		trustedHeader: "x-real-ip",
	}
	s := newServer(conf, brain)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("x-real-ip", "300.300.300.300")
	s.HTTPHandler(resp, req)

	brain.Finish()
	events := brain.RecordedEvents()
	assert.NotEmpty(t, events)

	expectedEvt := RequestEvent{
		Header:     req.Header,
		RemoteAddr: "300.300.300.300:1234",
		Method:     "GET",
		URL:        req.URL,
		Body:       []byte{},
	}
	assert.Equal(t, expectedEvt, events[0])
}

func TestServer_HTTPHandler_POST(t *testing.T) {
	brain := joetest.NewBrain(t)
	conf := config{
		listenAddr: "127.0.0.1:0",
		logger:     zaptest.NewLogger(t),
	}
	s := newServer(conf, brain)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/test", strings.NewReader("hello world"))
	s.HTTPHandler(resp, req)

	brain.Finish()
	events := brain.RecordedEvents()
	assert.NotEmpty(t, events)

	expectedEvt := RequestEvent{
		Header:     req.Header,
		RemoteAddr: "192.0.2.1:1234",
		Method:     "POST",
		URL:        req.URL,
		Body:       []byte("hello world"),
	}
	assert.Equal(t, expectedEvt, events[0])
}

func TestServer(t *testing.T) {
	addr := testAddr(t)
	bot := joetest.NewBot(t, Server(addr))
	bot.Start()

	resp, err := http.Post("http://"+addr+"/test", "text/plain", strings.NewReader("foobar"))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bot.Stop()
}

func testAddr(t *testing.T) string {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}

	s := l.Addr().String()
	err = l.Close()
	if err != nil {
		t.Fatal(err)
	}

	return s
}
