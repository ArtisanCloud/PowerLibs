package contract

import (
	"github.com/ArtisanCloud/PowerLibs/v2/object"
)

type ClientInterface interface {
	Send(request RequestInterface, options *object.HashMap) ResponseInterface
	SendAsync(request RequestInterface, options *object.HashMap) PromiseInterface

	Request(method string, uri string, options *object.HashMap, returnRaw bool, outHeader interface{}, outBody interface{}) (ResponseInterface, error)
	RequestAsync(method string, uri string, options *object.HashMap, returnRaw bool, outHeader interface{}, outBody interface{})

	SetClientConfig(config *object.HashMap) ClientInterface
	GetClientConfig() *object.HashMap
}
