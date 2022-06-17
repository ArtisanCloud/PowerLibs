package request

import (
	"crypto/tls"
	"fmt"
	"github.com/ArtisanCloud/PowerLibs/v2/http/contract"
	"github.com/ArtisanCloud/PowerLibs/v2/http/drivers/gout"
	"github.com/ArtisanCloud/PowerLibs/v2/object"
	"net"
	"net/http"
	"time"
)

type HttpRequest struct {
	httpClient contract.ClientInterface
	BaseURI    string

	Middlewares []interface{}
}

var _defaults *object.HashMap

func NewHttpRequest(config *object.HashMap) (*HttpRequest, error) {

	var httpClient *http.Client

	if (*config)["cert_path"] != nil && (*config)["key_path"] != nil {
		certPath := (*config)["cert_path"].(string)
		keyPath := (*config)["key_path"].(string)
		if certPath != "" && keyPath != "" {
			var err error
			httpClient, err = NewTLSHttpClient(certPath, keyPath)
			if err != nil {
				return nil, err
			}
		}
	}

	return &HttpRequest{
		httpClient: gout.NewClient(config, httpClient),
	}, nil
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
		request.httpClient = gout.NewClient(nil, nil)
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
	returnRaw bool, outHeader interface{}, outBody interface{}) (contract.ResponseInterface, error) {
	// change method string format
	method = object.Lower(method)

	// merge options with default options
	options = object.MergeHashMap(options, _defaults, &object.HashMap{"handler": request.GetMiddlewares()})

	// use request baseUri as final
	if request.BaseURI != "" {
		(*options)["base_uri"] = request.BaseURI
	}

	// use current http client driver to request
	response, err := request.GetHttpClient().Request(method, url, options, returnRaw, outHeader, outBody)
	return response, err
}

func NewTLSHttpClient(certFile string, keyFile string) (httpClient *http.Client, err error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		fmt.Print("can not init cert...")
		return nil, err
	}
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: false,
		ClientAuth:         tls.RequireAndVerifyClientCert,
		MinVersion:         tls.VersionTLS12,
	}

	httpClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig:     tlsConfig,
		},
		Timeout: 60 * time.Second,
	}
	return
}
