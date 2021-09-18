package gout

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ArtisanCloud/go-libs/http/contract"
	"github.com/ArtisanCloud/go-libs/http/response"
	"github.com/ArtisanCloud/go-libs/object"
	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/dataflow"
	"github.com/guonaihong/gout/filter"
	dataflow2 "github.com/guonaihong/gout/interface"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

const OPTION_SYNCHRONOUS = "synchronous"

type Client struct {
	Config *object.HashMap
}

func NewClient(config *object.HashMap) *Client {
	client := &Client{
		Config: config,
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

func (client *Client) PrepareRequest(method string, uri string, options *object.HashMap) (
	df *dataflow.DataFlow,
	queries *object.StringMap, headers interface{}, body interface{},
	version string, debug bool) {

	(*options)[OPTION_SYNCHRONOUS] = true
	options = client.prepareDefaults(options)

	version = "1.1"

	if (*options)["headers"] != nil {
		headers = (*options)["headers"]
	}
	if (*options)["body"] != nil {
		body = (*options)["body"]
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
		client.useMiddleware(df, middlewares)
	}

	// append query
	queries = &object.StringMap{}
	if (*options)["query"] != nil {
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

	return df, queries, headers, body, version, debug
}

func (client *Client) Request(method string, uri string, options *object.HashMap, returnRaw bool, outHeader interface{}, outBody interface{}) (contract.ResponseInterface, error) {

	df, queries, headers, body, _, debug := client.PrepareRequest(method, uri, options)

	df = client.applyOptions(df, options)

	returnCode := http.StatusOK
	df = df.
		Debug(debug).
		SetQuery(queries).
		SetHeader(headers)
	//Code(&returnCode).
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

	err := df.Do()
	if err != nil {
		fmt.Printf("do request error:%s \n", err.Error())
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
		case string:
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
			(*options)["headers"] = object.MergeHashMap(headers, (*options)["headers"].(*object.HashMap))
			break
		default:
			println("error header ")
			return nil
		}
	}
	result := object.MergeHashMap(defaultOptions, options)

	return result
}

func (client *Client) applyOptions(r *dataflow.DataFlow, options *object.HashMap) *dataflow.DataFlow {

	if (*options)["form_params"] != nil {
		(*options)["body"], _ = object.StructToMap((*options)["form_params"])
		(*options)["form_params"] = nil

		(*options)["_conditional"] = &object.StringMap{
			"Content-Type": "application/x-www-form-urlencoded",
		}

		bodyData := (*options)["body"].(map[string]interface{})
		r.SetJSON(bodyData)

	}

	if (*options)["multipart"] != nil {
		for _, media := range (*options)["multipart"].([]*object.HashMap) {
			name := (*media)["name"].(string)

			// load data from file
			if (*media)["headers"] != nil {
				value := (*media)["value"].(string)
				//headers := (*media)["headers"].(string)
				r.SetForm(gout.H{
					name: gout.FormFile(value),
				}).SetHeader(gout.H{
				})
			} else
			// load data from memory
			{
				value := (*media)["value"].([]byte)
				r.SetForm(gout.H{
					name: gout.FormMem(value),
				})
			}
		}

	}

	return r
}

func (client *Client) buildUri(uri *url.URL, config *object.HashMap) *url.URL {
	var baseUri *url.URL

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
			strBaseUri := (*client.Config)["http"].(object.HashMap)["base_uri"].(string)
			baseUri, _ = url.Parse(strBaseUri)
		}

		baseUri.Path = path.Join(baseUri.Path, uri.Path)
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

	switch method {
	case "get":
		df = gout.GET(url)
		break
	case "post":
		df = gout.POST(url)
		break
	case "put":
		df = gout.PUT(url)
		break
	default:
		df = gout.GET(url)
	}
	return df
}

func (client *Client) useMiddleware(df *dataflow.DataFlow, middlewares []interface{}) {
	for _, middleware := range middlewares {

		md := middleware.(contract.MiddlewareInterface)
		mdName := md.GetName()
		if mdName == "retry" {
			client.handleRetryMiddleware(df, md)
		}

		requestMiddleware := middleware.(dataflow2.RequestMiddler)
		df.RequestUse(requestMiddleware)
	}
}

func (client *Client) handleRetryMiddleware(df *dataflow.DataFlow, retryMiddleware contract.MiddlewareInterface) *dataflow.DataFlow {

	retries := retryMiddleware.Retries()
	delay := retryMiddleware.Delay()
	df.F().Retry().Attempt(retries).WaitTime(delay).MaxWaitTime(delay).Func(func(c *gout.Context) error {
		conditions := &object.HashMap{
			"code": c.Code,
		}
		if c.Error != nil || retryMiddleware.RetryDecider(conditions) {
			return filter.ErrRetry
		}

		return nil
	})

	return df

}
