package object

import (
	fmt2 "fmt"
	"github.com/ArtisanCloud/PowerLibs/v3/fmt"
	"reflect"
	"testing"
)

func Test_Map2Xml(t *testing.T) {

	obj := HashMap{
		"ActName":   "红包测试",
		"ClientIp":  "127.0.0.1",
		"MchBillno": "1634802770-845085000",
		"MchID":     "1613495874",
		"NonceStr":  "sy5ViKYvvT",
		"ReOpenid":  "okQfd5UNMMbDEyJo2ydwsGQ9R4NI",
		"Remark":    "猜越多得越多，快来抢！",
		"RiskInfo":  "",
		"SceneID":   "",
		"SendName":  "技术部-王秦文",
		"Sign":      "581E208372C7F3BACCD81D0458BEAF33AC293B9E88E9EA806C8AC901D55831CE",
		"Text":      "",
		"Array": &HashMap{
			"key1": "value",
			"keyArray": HashMap{
				"key2": 123,
			},
		},
		"ArrayStr": StringMap{
			"key1": "value",
		},
		"ArrayStrP": &StringMap{
			"key1": "value",
		},
		"TotalAmount": "1000",
		"TotalNum":    "1",
		"Wishing":     "技术部测试红包",
		"Wxappid":     "wx94dcb1e3674e84ad",
		"XMLName":     "",
		"cert":        "/private/var/www/html/GO/连续剧待有代币的版本/keys/1613495874/apiclient_cert.pem",
		"ssl_key":     "/private/var/www/html/GO/连续剧待有代币的版本/keys/1613495874/apiclient_key.pem",
	}

	xmlObj := Map2Xml(&obj, false)
	println(xmlObj)
	fmt.Dump(xmlObj)

}

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

func Test_Xml2HashMap(t *testing.T) {
	xmlData := []byte(`
        <root>
            <person1>
                <name>John</name>
                <age>30</age>
            </person1>
            <person2>
                <name>Alice</name>
                <age>25</age>
            </person2>
        </root>
    `)

	result, err := Xml2Map(xmlData)

	if err != nil {
		t.Errorf("Xml2Map error: %v", err)
	}

	expected := HashMap{
		"root": HashMap{
			"person1": HashMap{
				"name": "John",
				"age":  "30",
			},
			"person2": HashMap{
				"name": "Alice",
				"age":  "25",
			},
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Result mismatch.\nExpected: %#v\nActual: %#v", expected, result)
	}
}

func Test_DumpXML(t *testing.T) {
	//str := `<?xml version="1.0" encoding="UTF-8" ?>
	//<test>中文</test>`
	str := `<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<book>\n  <title>中文来的</title>\n
	<author>J.K. 中文来的</author>\n  <publisher>中文来的</publisher>\n
	<publishedYear>1997</publishedYear>\n</book>"`
	fmt2.Println(str)
	//fmt.Dump(str)
}
