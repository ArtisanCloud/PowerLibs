package object

import (
	"bytes"
	"fmt"
)

type HashMap map[string]interface{}
type StringMap map[string]string

func MergeHashMap(toMap *HashMap, subMaps ...*HashMap) *HashMap {
	if toMap == nil{
		toMap = &HashMap{}
	}
	for _, subMap := range subMaps {
		if subMap!=nil{
			for k, v := range *subMap {
				(*toMap)[k] = v
			}
		}
	}
	return toMap
}

func MergeStringMap(toMap *HashMap, subMaps ...*HashMap) *HashMap {
	if toMap == nil{
		toMap = &HashMap{}
	}
	for _, subMap := range subMaps {
		if subMap!=nil{
			for k, v := range *subMap {
				(*toMap)[k] = v
			}
		}
	}
	return toMap
}

func ConvertStringMapToString(m *StringMap) string {
	b := new(bytes.Buffer)
	for key, value := range *m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}