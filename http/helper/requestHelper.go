package helper

import (
	"github.com/ArtisanCloud/PowerLibs/v2/http/contract"
	"github.com/ArtisanCloud/PowerLibs/v2/http/dataflow"
	"github.com/ArtisanCloud/PowerLibs/v2/http/drivers/http"
)

type RequestHelper struct {
	client           contract.ClientInterface
	middlewareHandle contract.RequestMiddleware
	config           *Config
}

type Config struct {
	*contract.ClientConfig
	BaseUrl string
}

func NewRequestHelper(conf *Config) (*RequestHelper, error) {
	client, err := http.NewHttpClient(conf.ClientConfig)
	if err != nil {
		return nil, err
	}
	return &RequestHelper{
		client: client,
		middlewareHandle: func(handle contract.RequestHandle) contract.RequestHandle {
			return handle
		},
		config: conf,
	}, nil
}

func (r *RequestHelper) SetClient(client contract.ClientInterface) {
	//TODO implement me
	panic("implement me")
}

func (r *RequestHelper) GetClient() contract.ClientInterface {
	//TODO implement me
	panic("implement me")
}

func (r *RequestHelper) WithMiddleware(middlewares ...contract.RequestMiddleware) {
	if len(middlewares) == 0 {
		return
	}
	var buildHandle func(md contract.RequestMiddleware, appendMd contract.RequestMiddleware) contract.RequestMiddleware
	buildHandle = func(md contract.RequestMiddleware, appendMd contract.RequestMiddleware) contract.RequestMiddleware {
		return func(handle contract.RequestHandle) contract.RequestHandle {
			return md(appendMd(handle))
		}
	}
	for _, middleware := range middlewares {
		r.middlewareHandle = buildHandle(r.middlewareHandle, middleware)
	}
}

func (r *RequestHelper) Df() contract.RequestDataflowInterface {
	return dataflow.NewDataflow(r.client, r.middlewareHandle, &dataflow.Option{
		BaseUrl: r.config.BaseUrl,
	})
}
