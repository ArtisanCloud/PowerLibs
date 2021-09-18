package contract

import "net/url"

type RequestInterface interface {
	getMethod() string
	getUri() *url.URL
}


