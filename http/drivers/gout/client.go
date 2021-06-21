package gout

import (
	"fmt"
	"github.com/ArtisanCloud/go-libs/http/contract"
	"github.com/ArtisanCloud/go-libs/object"
	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/dataflow"
	dataflow2 "github.com/guonaihong/gout/interface"
	"net/url"
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

func (client *Client) Send(request contract.RequestInterface, options *object.HashMap) contract.ResponseContract {
	return nil
}
func (client *Client) SendAsync(request contract.RequestInterface, options *object.HashMap) contract.PromiseInterface {
	return nil
}
func (client *Client) Request(method string, uri string, options *object.HashMap, outResponse interface{}) contract.ResponseContract {

	(*options)[OPTION_SYNCHRONOUS] = true
	options = client.prepareDefaults(options)

	var (
		headers gout.H = gout.H{}
		//body    gout.H = gout.H{}
		//version string = "1.1"
	)
	if (*options)["headers"] != nil {
		headers = (*options)["headers"].(gout.H)
	}
	//if (*options)["body"] != nil {
	//	body = (*options)["body"].(gout.H)
	//}
	// tbd
	//if options["version"] != "" {
	//	version = options["version"].(string)
	//}

	// Merge the URI into the base URI
	parsedURL, _ := url.Parse(uri)
	parsedURL = client.buildUri(parsedURL, options)
	strURL := parsedURL.String()

	// init a dataflow
	df := client.QueryMethod(method, strURL)

	// load middlewares stack
	if (*options)["handler"] != nil {
		middlewares := (*options)["handler"].([]interface{})
		client.useMiddleware(df, middlewares)
	}

	// append query
	queries := (*options)["query"]
	if queries == nil {
		queries = ""
	}

	// debug mode
	debug := false
	//fmt2.Dump(*client.Config)
	if (*client.Config)["debug"] != nil && (*client.Config)["debug"].(bool) == true {
		debug = true
	}

	df = client.applyOptions(df, options)

	err := df.
		Debug(debug).
		SetQuery(queries).
		SetHeader(&headers).
		BindJSON(outResponse).
		Do()

	if err != nil {
		// tbd throw excption
		fmt.Printf("do request error:", err.Error())
	}

	return nil

}

func (client *Client) RequestAsync(method string, uri string, options *object.HashMap, outResponse interface{}) {
	(*options)[OPTION_SYNCHRONOUS] = false

	go client.Request(method, uri, options, outResponse)

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
	return options
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

	return r
}

func (client *Client) buildUri(uri *url.URL, config *object.HashMap) *url.URL {
	var baseUri *url.URL
	if (*config)["base_uri"] != nil {
		strBaseUri := (*config)["base_uri"].(string)
		if strBaseUri != "" {
			baseUri, _ = url.Parse(strBaseUri)
		}
	} else {
		strBaseUri := (*client.Config)["http"].(map[string]string)["base_uri"]
		baseUri, _ = url.Parse(strBaseUri)
	}

	uri = baseUri.ResolveReference(uri)

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

func (client *Client) QueryMethod(method string, url string) (df *dataflow.DataFlow) {

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
		requestMiddleware := middleware.(dataflow2.RequestMiddler)
		df.RequestUse(requestMiddleware)
	}
}
