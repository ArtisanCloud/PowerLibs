package gout

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ArtisanCloud/PowerLibs/v2/http/contract"
	"github.com/ArtisanCloud/PowerLibs/v2/http/response"
	"github.com/ArtisanCloud/PowerLibs/v2/object"
	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/dataflow"
	"github.com/guonaihong/gout/filter"
	dataflow2 "github.com/guonaihong/gout/interface"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
)

const OPTION_SYNCHRONOUS = "synchronous"

type Client struct {
	Config     *object.HashMap
	HttpClient *http.Client
}

type TLSConfig struct {
	CertFile string
	KeyFile  string
}

func NewClient(config *object.HashMap, httpClient *http.Client) *Client {
	client := &Client{
		Config:     config,
		HttpClient: httpClient,
	}

	// set default config
	client.configureDefaults(config)

	return client
}

func (client *Client) Send(request contract.RequestInterface, options *object.HashMap) contract.ResponseInterface {
	return nil
}
func (client *Client) SendAsync(request contract.RequestInterface, options *object.HashMap) contract.PromiseInterface {
	return nil
}

func (client *Client) PrepareRequest(method string, uri string, options *object.HashMap, outBody interface{}) (
	df *dataflow.DataFlow,
	queries *object.StringMap, headers *object.HashMap, body interface{},
	version string, debug bool, err error) {

	(*options)[OPTION_SYNCHRONOUS] = true
	options = client.prepareDefaults(options)

	version = "1.1"

	if (*options)["headers"] != nil {
		headers = (*options)["headers"].(*object.HashMap)
	} else {
		headers = &object.HashMap{}
	}

	if (*options)["json"] != nil {
		(*options)["body"], err = object.JsonEncode((*options)["json"])
		if err != nil {
			return nil, nil, nil, nil, "", false, err
		}
	}
	if (*options)["body"] != nil {
		body = (*options)["body"]
	} else {
		body = nil
	}

	if (*options)["version"] != nil {
		version = (*options)["version"].(string)
	}

	// Merge the URI into the base URI
	parsedURL, _ := url.Parse(uri)
	parsedURL = client.buildUri(parsedURL, options)
	strURL := parsedURL.String()

	// init a dataflow
	df = client.QueryWithMethod(method, strURL)

	// load middlewares stack
	if (*options)["handler"] != nil {
		middlewares := (*options)["handler"].([]interface{})
		err = client.useMiddleware(df, middlewares, outBody)
		if err != nil {
			return df, queries, headers, body, version, debug, err
		}
	}

	// apply handle option
	//df, err = ApplyHandleOptions(df, options)
	//if err != nil {
	//	return df, queries, headers, body, version, debug, err
	//}

	// append query
	queries = &object.StringMap{}
	if !object.IsObjectNil((*options)["query"]) {
		queries = (*options)["query"].(*object.StringMap)
	}

	// debug mode
	debug = false
	if (*client.Config)["http_debug"] != nil && (*client.Config)["http_debug"].(bool) == true {
		debug = true
		fmt.Println("http debug mode open \n")
	}
	if (*client.Config)["debug"] != nil && (*client.Config)["debug"].(bool) == true {
		(*queries)["debug"] = "1"
		fmt.Println("wx debug mode open")
	}

	return df, queries, headers, body, version, debug, err
}

func (client *Client) Request(method string, uri string, options *object.HashMap, returnRaw bool, outHeader interface{}, outBody interface{}) (contract.ResponseInterface, error) {

	df, queries, headers, body, _, debug, err := client.PrepareRequest(method, uri, options, outBody)

	if err != nil {
		return nil, err
	}

	df = client.applyOptions(df, options, headers)

	returnCode := http.StatusOK
	df = df.
		Debug(debug).
		SetQuery(queries).
		SetHeader(headers).
		Code(&returnCode)
	//SetProxy("http://127.0.0.1:1088").

	if body != nil {
		df = df.SetBody(body)
	}

	// bind out header
	if outHeader != nil {
		df = df.BindHeader(outHeader)
	}

	// bind out body
	if outBody != nil {
		if returnRaw {
			df = df.BindBody(outBody)
		} else {
			df = df.BindJSON(outBody)
		}
	}

	err = df.Do()
	if err != nil {
		fmt.Printf("http request error:%s \n", err.Error())
		return nil, err
	}

	rs := client.GetHttpResponseFrom(returnCode, outHeader, outBody, returnRaw)
	return rs, err

}

func (client *Client) GetHttpResponseFrom(returnCode int, outHeader interface{}, outBody interface{}, returnRaw bool) *response.HttpResponse {
	rs := response.NewHttpResponse(returnCode)
	//fmt2.Dump("outHeader:", outHeader)
	//fmt2.Dump("outBody:", outBody)

	// copy body
	if returnRaw {
		switch outBody.(type) {
		case *string:
			rs.Body = ioutil.NopCloser(bytes.NewBufferString(*(outBody.(*string))))
		default:
			rs.Body = nil
		}
	} else {
		bodyBuffer, _ := json.Marshal(outBody)
		rs.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBuffer))
	}

	// copy header
	mapHeader, _ := object.StructToStringMap(outHeader)
	for key, header := range *mapHeader {
		rs.Header.Add(key, header)
	}
	return rs
}

func (client *Client) RequestAsync(method string, uri string, options *object.HashMap, returnRaw bool, outHeader interface{}, outBody interface{}) {
	(*options)[OPTION_SYNCHRONOUS] = false

	go client.Request(method, uri, options, returnRaw, outHeader, outBody)

}

func (client *Client) SetClientConfig(config *object.HashMap) contract.ClientInterface {
	client.Config = config
	return client
}

func (client *Client) GetClientConfig() *object.HashMap {
	return client.Config
}

func (client *Client) prepareDefaults(options *object.HashMap) *object.HashMap {
	// tbd
	defaultOptions := client.Config

	// merge headers
	headers := &object.HashMap{
		"Accept":       "*/*",
		"Content-Type": "application/json",
	}
	if (*options)["headers"] != nil {
		switch (*options)["headers"].(type) {
		case *object.HashMap:
			(*options)["headers"] = object.MergeHashMap((*options)["headers"].(*object.HashMap), headers)
			break
		default:
			println("error header ")
			return nil
		}
	}
	result := object.MergeHashMap(options, defaultOptions)

	return result
}

func (client *Client) applyOptions(r *dataflow.DataFlow, options *object.HashMap, headers *object.HashMap) *dataflow.DataFlow {

	if (*options)["form_params"] != nil {
		var bodyData interface{}
		switch (*options)["form_params"].(type) {
		case string:
			(*options)["body"], _ = (*options)["form_params"]
			bodyData = (*options)["body"].(string)

			r.SetBody(bodyData)
			break

		default:
			(*options)["body"], _ = object.StructToMap((*options)["form_params"])
			bodyData = (*options)["body"].(map[string]interface{})

			r.SetJSON(bodyData)
		}
		(*options)["form_params"] = nil
		(*options)["_conditional"] = &object.StringMap{
			"Content-Type": "application/x-www-form-urlencoded",
		}
	}

	if (*options)["multipart"] != nil {
		formData := gout.H{}

		for _, media := range (*options)["multipart"].([]*object.HashMap) {
			name := (*media)["name"].(string)

			// load data from file
			if (*media)["headers"] != nil {
				value := (*media)["value"].(string)
				headers = (*media)["headers"].(*object.HashMap)
				formData[name] = gout.FormType{FileName: (*headers)["filename"].(string), File: gout.FormFile(value)}

			} else
			// load data from memory
			{
				value := (*media)["value"].([]byte)
				formData[name] = gout.FormMem(value)
			}
		}
		r.SetForm(formData)
	}

	return r
}

func (client *Client) buildUri(uri *url.URL, config *object.HashMap) *url.URL {
	var baseUri *url.URL
	var err error

	// use customer custom url
	if uri.Host != "" {
		baseUri = uri
	} else {

		// use config base uri
		if (*config)["base_uri"] != nil {
			strBaseUri := (*config)["base_uri"].(string)
			if strBaseUri != "" {
				baseUri, _ = url.Parse(strBaseUri)
			}
		} else {
			// use app config base uri
			mapHttp := (*client.Config)["http"].(*object.HashMap)
			strBaseUri := (*mapHttp)["base_uri"].(string)
			baseUri, err = url.Parse(strBaseUri)
			if err != nil {
				print("cannot parse base url, pls make sure base_uri has scheme")
				print(err.Error())
			}
		}

		baseUri.Path = path.Join(baseUri.Path, uri.Path)
		baseUri.RawQuery = uri.RawQuery
		//uri = baseUri.ResolveReference(uri)
		uri = baseUri
	}

	// tbd idn_conversion
	// ...

	if uri.Scheme == "" && uri.Host != "" {
		uri.Scheme = "http"
	}

	return uri
}

func (client *Client) configureDefaults(config *object.HashMap) {
	defaults := &object.HashMap{
		//"allow_redirects": RedirectMiddleware::$defaultSettings,
		"http_errors":    true,
		"decode_content": true,
		"verify":         true,
		"cookies":        false,
		"idn_conversion": false,
	}

	object.MergeHashMap(client.Config, defaults, config)
}

func (client *Client) QueryWithMethod(method string, url string) (df *dataflow.DataFlow) {

	goutClient := gout.New(client.HttpClient)

	switch method {
	case "get":
		df = goutClient.GET(url)
		break
	case "post":
		df = goutClient.POST(url)
		break
	case "put":
		df = goutClient.PUT(url)
		break
	default:
		df = goutClient.GET(url)
	}
	return df
}

func (client *Client) useMiddleware(df *dataflow.DataFlow, middlewares []interface{}, outBody interface{}) (err error) {
	for _, middleware := range middlewares {

		md := middleware.(contract.MiddlewareInterface)
		mdName := md.GetName()
		if mdName == "retry" {
			df, err = client.handleRetryMiddleware(df, md, outBody)
		}

		requestMiddleware := middleware.(dataflow2.RequestMiddler)
		df.RequestUse(requestMiddleware)
	}

	return nil
}

func (client *Client) handleRetryMiddleware(df *dataflow.DataFlow, retryMiddleware contract.MiddlewareInterface, outBody interface{}) (*dataflow.DataFlow, error) {

	retries := retryMiddleware.Retries()
	delay := retryMiddleware.Delay()
	df.F().Retry().Attempt(retries).WaitTime(delay).MaxWaitTime(delay).Func(func(c *gout.Context) error {

		resp, err := c.Response()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil
		}
		mapResponse := &object.HashMap{}
		err = json.Unmarshal(data, mapResponse)
		if err != nil {
			return nil
		}
		if (*mapResponse)["errcode"] != nil {
			errCode := (*mapResponse)["errcode"].(int)
			conditions := &object.HashMap{
				"code": errCode,
			}
			if c.Error != nil || retryMiddleware.RetryDecider(conditions) {
				return filter.ErrRetry
			}

		}

		return nil
	})
	//.Do()

	return df, nil

}

func ApplyHandleOptions(df *dataflow.DataFlow, options *object.HashMap) (*dataflow.DataFlow, error) {

	var err error
	var cert, sslKey string
	if (*options)["cert"] != nil {
		cert = (*options)["cert"].(string)

		if _, err = os.Stat(cert); os.IsNotExist(err) {
			err = errors.New("SSL certificate not found:" + cert)
			return df, err
		}
	}
	if (*options)["ssl_key"] != nil {
		sslKey = (*options)["ssl_key"].(string)

		if _, err = os.Stat(cert); os.IsNotExist(err) {
			err = errors.New("SSL certificate not found:" + cert)
			return df, err
		}
	}
	if cert == "" || sslKey == "" {
		return df, nil
	}

	tlsConfig, err := NewTLSConfig(cert, sslKey, "")
	if err != nil {
		return df, err
	}
	tr := &http.Transport{TLSClientConfig: tlsConfig}

	df.Client().Transport = tr

	return df, err
}

func NewTLSConfig(clientCertFile string, clientKeyFile string, caCertFile string) (*tls.Config, error) {

	// Load client cert
	tlsCert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		return nil, err
	}

	tlsConfig := tls.Config{
		Certificates:       []tls.Certificate{tlsCert},
		InsecureSkipVerify: false,
		ClientAuth:         tls.RequireAndVerifyClientCert,
		MinVersion:         tls.VersionTLS12,
	}

	// Load CA cert
	if caCertFile != "" {
		certBytes, err := ioutil.ReadFile(caCertFile)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(certBytes); !ok {
			return nil, errors.New("Unable to load caCert")
		}
		tlsConfig.RootCAs = caCertPool
		tlsConfig.ClientCAs = caCertPool
	}

	return &tlsConfig, nil
}
