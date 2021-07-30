package object

import (
	"encoding/json"
	"reflect"
)

type HashMap map[string]interface{}


// ------------------------------- Merge --------------------------------------------
func MergeHashMap(toMap *HashMap, subMaps ...*HashMap) *HashMap {
	if toMap == nil {
		toMap = &HashMap{}
	}
	for _, subMap := range subMaps {
		if subMap != nil {
			for k, v := range *subMap {
				(*toMap)[k] = v
			}
		}
	}
	return toMap
}



// ------------------------------- Conversion ---------------------------------------
func StructToHashMap(obj interface{}) (newMap *HashMap, err error) {
	data, err := json.Marshal(obj) // Convert to a json string

	if err != nil {
		return
	}

	newMap = &HashMap{}
	err = json.Unmarshal(data, newMap) // Convert to a map
	return
}



// ------------------------------- Search --------------------------------------------
func InHash(val interface{}, hash *HashMap) (exists bool, key string) {
	exists = false
	key = ""

	switch reflect.TypeOf(hash).Kind() {
	case reflect.Map:
		//s := reflect.ValueOf(hash)

		for k, v := range *hash {
			if reflect.DeepEqual(val, v) == true {
				key = k
				return
			}
		}
	}

	return
}

func GetHashMapKV(maps StringMap) (keys []string, values []interface{}) {
	mapLen := len(maps)
	keys = make([]string, 0, mapLen)
	values = make([]interface{}, 0, mapLen)

	for k, v := range maps {
		keys = append(keys, k)
		values = append(values, v)
	}

	return keys, values
}

// ------------------------------- Filter --------------------------------------------
func FilterEmptyHashMap(mapData *HashMap) (filteredMap *HashMap) {
	filteredMap = &HashMap{}
	for k, v := range *mapData {
		if v != nil {
			(*filteredMap)[k] = v
			switch v.(type) {
			case HashMap:
				o := v.(HashMap)
				v = FilterEmptyHashMap(&o)
				break
			case *HashMap:
				v = FilterEmptyHashMap(v.(*HashMap))
				break
			case string:
				if v.(string) == "" {
					delete(*filteredMap, k)
				}
				break
			}
		}
	}
	return filteredMap
}
