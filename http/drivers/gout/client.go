package gout

import (
	"fmt"
	fmt2 "github.com/ArtisanCloud/go-libs/fmt"
	"github.com/ArtisanCloud/go-libs/http/contract"
	"github.com/ArtisanCloud/go-libs/object"
	"github.com/guonaihong/gout/dataflow"
	"net/url"
)

const OPTION_SYNCHRONOUS = "synchronous"

type Client struct {
	df     *dataflow.DataFlow
	Config *object.HashMap
}

func NewClient(config *object.HashMap) *Client {
	client := &Client{
		df:     &dataflow.DataFlow{},
		Config: config,
	}
	return client
}

func (client *Client) Send(request contract.RequestInterface, options object.HashMap) *contract.ResponseContract {
	return nil
}
func (client *Client) SendAsync(request contract.RequestInterface, options object.HashMap) *contract.PromiseInterface {
	return nil
}
func (client *Client) Request(method string, uri string, options object.HashMap) *contract.ResponseContract {

	options[OPTION_SYNCHRONOUS] = true
	options = client.prepareDefaults(options)

	var (
		headers object.StringMap = object.StringMap{}
		body    object.HashMap   = object.HashMap{}
		//version string = "1.1"
	)
	if options["headers"] != nil {
		headers = options["headers"].(object.StringMap)
	}
	if options["body"] != nil {
		body = options["body"].(object.HashMap)
	}
	//if options["version"] != "" {
	//	version = options["version"].(string)
	//}

	// Merge the URI into the base URI
	parsedURL, _ := url.Parse(uri)
	parsedURL = client.buildUri(parsedURL, client.Config)
	strURL := parsedURL.String()
	println(strURL)
	result := ""
	err := client.df.
		SetURL(strURL).
		SetHeader(headers).
		SetBody(body).
		BindBody(&result).
		Do()

	if err != nil {
		fmt.Printf("do request error:", err.Error())
	}
	fmt2.Dump(result)
	return nil

}

func (client *Client) RequestAsync(method string, uri string, options object.HashMap) {
	options[OPTION_SYNCHRONOUS] = false

	go client.Request(method, uri, options)

}

func (client *Client) prepareDefaults(options object.HashMap) object.HashMap {
	return options
}

func (client *Client) buildUri(url *url.URL, config *object.HashMap) *url.URL {
	arrayConfig := *config
	fmt2.Dump(arrayConfig)
	strUri:=arrayConfig["http"].(map[string]string)["base_uri"]
	if strUri!= "" {
		url, _ = url.Parse(strUri)
	}

	// idn_conversion
	// ...

	if url.Scheme == "" && url.Host != "" {
		url.Scheme = "http"
	}

	return url
}
