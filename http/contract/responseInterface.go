package contract

import (
	"io"
	"net/http"
)

type MessageInterface interface {
	GetBody() io.ReadCloser
	GetBodyData() ([]byte, error)
	GetHeader() http.Header
}

type ResponseInterface interface {
	MessageInterface
	GetStatusCode() int
	//WithStatus() object.HashMap
	//GetReasonPhrase() object.HashMap
}

type PromiseInterface interface {
}
