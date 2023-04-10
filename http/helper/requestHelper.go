package helper

import (
	"bytes"
	"encoding/xml"
	"github.com/ArtisanCloud/PowerLibs/v3/http/contract"
	"github.com/ArtisanCloud/PowerLibs/v3/http/dataflow"
	"github.com/ArtisanCloud/PowerLibs/v3/http/drivers/http"
	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"io"
	"io/ioutil"
	http2 "net/http"
	"strings"
)

type RequestDownload struct {
	HashType    string `json:"hash_type""`
	HashValue   string `json:"hash_value""`
	DownloadURL string `json:"download_url""`
}

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
	r.client = client
}

func (r *RequestHelper) GetClient() contract.ClientInterface {
	return r.client
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

func (r *RequestHelper) ParseResponseBodyToMap(rs *http2.Response, outBody *object.HashMap) error {

	b, err := io.ReadAll(rs.Body)
	if err != nil {
		return err
	}
	rs.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	contentType := rs.Header.Get("Content-Type")

	//fmt.Dump(123, contentType)

	if strings.Contains(contentType, "application/xml") || strings.Contains(contentType, "text/xml") {
		*outBody, err = object.Xml2Map(b)
		if err != nil {
			return err
		}
	} else {
		// Handle JSON format.
		err = object.JsonDecode(b, outBody)
		if err != nil {
			return err
		}
	}

	return nil

}

func (r *RequestHelper) ParseResponseBodyContent(rs *http2.Response, outBody interface{}) error {

	b, err := io.ReadAll(rs.Body)
	if err != nil {
		return err
	}
	rs.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	contentType := rs.Header.Get("Content-Type")

	//fmt.Dump(456, contentType)
	if strings.Contains(contentType, "application/xml") || strings.Contains(contentType, "text/xml") {
		err = xml.Unmarshal(b, outBody)
		if err != nil {
			return err
		}
	} else {
		// Handle JSON format.
		err = object.JsonDecode(b, outBody)
		if err != nil {
			return err
		}
	}

	return nil

}

func HttpResponseSend(rs *http2.Response, writer http2.ResponseWriter) (err error) {

	// set header code
	writer.WriteHeader(rs.StatusCode)

	// set write body
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return err
	}
	_, err = writer.Write(body)

	return err
}
