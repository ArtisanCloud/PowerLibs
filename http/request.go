package http

import (
	"github.com/ArtisanCloud/go-libs/http/contract"
	"github.com/ArtisanCloud/go-libs/http/drivers/gout"
	"github.com/ArtisanCloud/go-libs/object"
	"github.com/ArtisanCloud/go-libs/str"
)

type HttpRequest struct {
	httpClient contract.ClientInterface

	Middlewares object.HashMap
}

var _defaults *object.HashMap

func NewHttpRequest(config *object.HashMap) *HttpRequest {
	return &HttpRequest{
		httpClient: &gout.Client{
			Config: config,
		},
	}
}

func SetDefaultOptions(defaults *object.HashMap) {
	_defaults = defaults
}

func GetDefaultOptions() *object.HashMap {
	return _defaults
}

func (request *HttpRequest) SetHttpClient(httpClient contract.ClientInterface) *HttpRequest {
	request.httpClient = httpClient
	return request
}

func (request *HttpRequest) GetHttpClient() contract.ClientInterface {

	if request.httpClient == nil {
		//oType:=reflect.TypeOf(request.httpClient)
		//if fmt.Sprintf("%T",oType) != "ClientInterface"{
		//	if request.App.
		//}
		request.httpClient = &gout.Client{}
	}

	return request.httpClient
}

func (request *HttpRequest) GetMiddlewares() object.HashMap {
	return request.Middlewares
}

func (request *HttpRequest) pushMiddleware(middleware interface{}, name string) bool {
	if name != "" {
		request.Middlewares[name] = middleware
		return true
	}
	return false
}

func (request *HttpRequest) PerformRequest(url string, method string, options object.HashMap) *contract.ResponseContract {
	// change method string format
	method = str.Lower(method)

	// merge options with default options
	//options:= merge(_defaults,options,hash)
	response := request.GetHttpClient().Request(method, url, options)
	return response
}
