package object

type HashMap map[string]interface{}
type StringMap map[string]string




func MergeHashMap(subMap *HashMap, toMap *HashMap) *HashMap {
	for k, v := range *subMap {
		(*toMap)[k] = v
	}
	return toMap
}


func MergeStringMap(subMap *StringMap, toMap *StringMap) *StringMap {
	for k, v := range *subMap {
		(*toMap)[k] = v
	}
	return toMap
}