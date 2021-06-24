package object

import (
	"encoding/json"
	"fmt"
)

func Map2Xml(obj *HashMap) (strXML string) {

	for k, v := range *obj {
		switch v.(type) {
		case string:
			strXML = strXML + fmt.Sprintf("<%s><![CDATA[%s]]></%s>", k, v, k)
		case int:
			strXML = strXML + fmt.Sprintf("<%s><![CDATA[%d]]></%s>", k, v, k)
		case interface{}:
			b, _ := json.Marshal(v)
			strXML = strXML + fmt.Sprintf("<%s><![CDATA[%s]]></%s>", k, string(b), k)
		}
	}
	return "<xml>" + strXML + "</xml>"
}
