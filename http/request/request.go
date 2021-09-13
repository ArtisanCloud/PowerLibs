package request

import (
	"github.com/ArtisanCloud/go-libs/http/contract"
	"github.com/ArtisanCloud/go-libs/http/drivers/gout"
	"github.com/ArtisanCloud/go-libs/object"
)

type HttpRequest struct {
	httpClient contract.ClientInterface
	baseUri    string

	Middlewares []interface{}
}

var _defaults *object.HashMap

func NewHttpRequest(config *object.HashMap) *HttpRequest {
	return &HttpRequest{
		httpClient: gout.NewClient(config),
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
		request.httpClient = gout.NewClient(nil)
	}

	return request.httpClient
}

func (request *HttpRequest) GetMiddlewares() []interface{} {
	return request.Middlewares
}

func (request *HttpRequest) PushMiddleware(middleware interface{}, name string) bool {
	if name != "" {
		request.Middlewares = append(request.Middlewares, middleware)

		return true
	}
	return false
}

func (request *HttpRequest) PerformRequest(url string, method string, options *object.HashMap,
	returnRaw bool, outHeader interface{}, outBody interface{}) (contract.ResponseContract ,error){
	// change method string format
	method = object.Lower(method)

	// merge options with default options
	options = object.MergeHashMap(options, _defaults, &object.HashMap{"handler": request.GetMiddlewares()})

	// use request baseUri as final
	if request.baseUri != "" {
		(*options)["base_uri"] = request.baseUri
	}

	// use current http client driver to request
	response, err := request.GetHttpClient().Request(method, url, options, returnRaw, outHeader, outBody)
	return response, err
}
