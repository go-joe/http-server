package joehttp

import "net/url"

type RequestEvent struct {
	Method     string
	URL        *url.URL
	RemoteAddr string
	Body       []byte
}
