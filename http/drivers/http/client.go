package http

import (
	"crypto/tls"
	"github.com/ArtisanCloud/PowerLibs/v3/fmt"
	"github.com/ArtisanCloud/PowerLibs/v3/http/contract"
	"github.com/pkg/errors"
	"net/http"
)

// Client 是 net/http 的封装
type Client struct {
	conf       *contract.ClientConfig
	coreClient *http.Client
}

func NewHttpClient(config *contract.ClientConfig) (*Client, error) {
	if config == nil {
		config = &contract.ClientConfig{}
		config.Default()
	}
	coreClient := http.Client{
		Timeout: config.Timeout,
	}
	if config.Cert.CertFile != "" && config.Cert.KeyFile != "" {
		certPair, err := tls.LoadX509KeyPair(config.Cert.CertFile, config.Cert.KeyFile)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load certificate")
		}
		coreClient.Transport = &http.Transport{TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{certPair},
		}}
	}
	return &Client{
		conf:       config,
		coreClient: &coreClient,
	}, nil
}

// SetConfig 配置客户端
func (c *Client) SetConfig(config *contract.ClientConfig) {
	if config != nil {
		c.conf = config
	}
	// todo set coreClient
	if config.Cert.CertFile != "" && config.Cert.KeyFile != "" {
		coreClient := http.Client{
			Timeout: config.Timeout,
		}
		certPair, err := tls.LoadX509KeyPair(config.Cert.CertFile, config.Cert.KeyFile)
		if err != nil {
			err = errors.Wrap(err, "failed to load certificate")
			fmt.Dump(err)
			return
		}
		coreClient.Transport = &http.Transport{TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{certPair},
		}}
		c.coreClient = &coreClient
	}

}

// GetConfig 返回配置副本
func (c *Client) GetConfig() contract.ClientConfig {
	return *c.conf
}

func (c *Client) DoRequest(request *http.Request) (response *http.Response, err error) {
	return c.coreClient.Do(request)
}
