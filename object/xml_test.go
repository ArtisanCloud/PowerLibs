package object

import (
	"github.com/ArtisanCloud/PowerLibs/v3/fmt"
	"testing"
)

func Test_StringMap2Xml(t *testing.T) {

	obj := StringMap{
		"ActName":     "红包测试",
		"ClientIp":    "127.0.0.1",
		"MchBillno":   "1634802770-845085000",
		"MchID":       "1613495874",
		"NonceStr":    "sy5ViKYvvT",
		"ReOpenid":    "okQfd5UNMMbDEyJo2ydwsGQ9R4NI",
		"Remark":      "猜越多得越多，快来抢！",
		"RiskInfo":    "",
		"SceneID":     "",
		"SendName":    "技术部-王秦文",
		"Sign":        "581E208372C7F3BACCD81D0458BEAF33AC293B9E88E9EA806C8AC901D55831CE",
		"Text":        "",
		"TotalAmount": "1000",
		"TotalNum":    "1",
		"Wishing":     "技术部测试红包",
		"Wxappid":     "wx94dcb1e3674e84ad",
		"XMLName":     "",
		"cert":        "/private/var/www/html/GO/连续剧待有代币的版本/keys/1613495874/apiclient_cert.pem",
		"ssl_key":     "/private/var/www/html/GO/连续剧待有代币的版本/keys/1613495874/apiclient_key.pem",
	}

	xmlObj := StringMap2Xml(&obj)
	println(xmlObj)
	fmt.Dump(xmlObj)

}
