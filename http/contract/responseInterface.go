package contract

import (
	http2 "net/http"
)

type MessageInterface interface{
	GetBody() *http2.ResponseWriter
	GetHeaders() *http2.ResponseWriter
}

type ResponseContract interface {
	MessageInterface
	GetStatusCode() int
	//WithStatus() object.HashMap
	//GetReasonPhrase() object.HashMap
}

type PromiseInterface interface {

}