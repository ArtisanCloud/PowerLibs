package gout

import (
	"fmt"
	"github.com/ArtisanCloud/go-libs/http/contract"
	"github.com/ArtisanCloud/go-libs/object"
	"github.com/guonaihong/gout/dataflow"
)

const OPTION_SYNCHRONOUS = "synchronous"

type Client struct {
	df *dataflow.DataFlow
	Config *object.HashMap
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
	if options["headers"] != "" {
		headers = options["headers"].(object.StringMap)
	}
	if options["body"] != "" {
		body = options["body"].(object.HashMap)
	}
	//if options["version"] != "" {
	//	version = options["version"].(string)
	//}

	// Merge the URI into the base URI
	//client.buildUri(uri , client.config)

	err := client.df.
		//SetURL().
		SetHeader(headers).
		SetBody(body).
		Do()
	if err != nil {
		fmt.Printf("do request error:",err.Error())
	}

	return nil

}

func (client *Client) RequestAsync(method string, uri string, options object.HashMap) {
	options[OPTION_SYNCHRONOUS] = false

	go client.Request(method, uri, options)

}

func (client *Client) prepareDefaults(options object.HashMap) object.HashMap {
	return options
}
