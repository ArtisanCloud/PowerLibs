package contract

import (
	"net/http"
	"time"
)

type ClientConfig struct {
	Timeout time.Duration
	Cert    CertConfig
	// 如果需要定制化tls, 设置该属性, 否则请使用Cert
	// TlsConfig *tls.Config
}

type CertConfig struct {
	CertFile string
	KeyFile  string
}

func (c *ClientConfig) Default() {
	if c.Timeout == 0 {
		c.Timeout = time.Second * 30
	}
}

type ClientInterface interface {
	// SetConfig 覆盖原来的 Client 配置
	SetConfig(config *ClientConfig)

	// GetConfig 返回当前客户端配置(副本)
	GetConfig() ClientConfig

	// DoRequest 默认请求接口
	DoRequest(request *http.Request) (response *http.Response, err error)
}
