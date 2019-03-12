package joehttp

import (
	"net/http"
	"net/url"
)

type RequestEvent struct {
	Header     http.Header
	Method     string
	URL        *url.URL
	RemoteAddr string
	Body       []byte
}
