package object

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
)

func Str2Xml(in string) string {
	var b bytes.Buffer
	xml.EscapeText(&b, []byte(in))
	return b.String()
}

func Map2Xml(obj *HashMap, isSub bool) (strXML string) {

	for k, v := range *obj {
		switch v.(type) {
		case string:
			strXML = strXML + fmt.Sprintf("<%s><![CDATA[%s]]></%s>", k, v, k)
			break
		case int, int8, int16, int32, int64:
			strXML = strXML + fmt.Sprintf("<%s>%d</%s>", k, v, k)
			break
		case float32:
		case float64:
			strXML = strXML + fmt.Sprintf("<%s>%f</%s>", k, v, k)
			break
		case []*HashMap:
			for _, subV := range v.([]*HashMap) {
				strXML = strXML + fmt.Sprintf("<%s>%s</%s>", k, Map2Xml(subV, true), k)
			}
			break
		case []HashMap:
			for _, subV := range v.([]HashMap) {
				strXML = strXML + fmt.Sprintf("<%s>%s</%s>", k, Map2Xml(&subV, true), k)
			}
			break
		case *HashMap:
			strXML = strXML + fmt.Sprintf("<%s>%s</%s>", k, Map2Xml(v.(*HashMap), true), k)
			break
		case HashMap:
			val := v.(HashMap)
			strXML = strXML + fmt.Sprintf("<%s>%s</%s>", k, Map2Xml(&val, true), k)
			break
		case interface{}:
			b, _ := json.Marshal(v)
			strXML = strXML + fmt.Sprintf("<%s><![CDATA[%s]]></%s>", k, string(b), k)
			break
		}
	}
	if isSub {
		return strXML
	} else {
		return "<xml>" + strXML + "</xml>"
	}
}

func StringMap2Xml(obj *StringMap) (strXML string) {

	for k, v := range *obj {
		strXML = strXML + fmt.Sprintf("<%s><![CDATA[%s]]></%s>", k, v, k)
	}
	return "<xml>" + strXML + "</xml>"
}

func Xml2Map(b []byte) (m HashMap, err error) {

	decoder := xml.NewDecoder(bytes.NewReader(b))
	m = make(HashMap)
	tag := ""
	for {
		token, err := decoder.Token()

		if err != nil {
			break
		}
		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local != "xml" {
				tag = t.Name.Local
			} else {
				tag = ""
			}
			break
		case xml.EndElement:
			break
		case xml.CharData:
			data := strings.TrimSpace(string(t))
			if len(data) != 0 {
				m[tag] = data
			}
			break
		}
	}
	return m, err
}
