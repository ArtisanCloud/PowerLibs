package object

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
