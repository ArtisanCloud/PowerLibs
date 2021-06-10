package contract

import "github.com/ArtisanCloud/go-libs/object"

type MessageInterface interface{
	GetBody() *object.HashMap
	GetHeaders() *object.HashMap
}

type ResponseContract interface {
	MessageInterface
	GetStatusCode() int
	//WithStatus() object.HashMap
	//GetReasonPhrase() object.HashMap
}

type PromiseInterface interface {

}