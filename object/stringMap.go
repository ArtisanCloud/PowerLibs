package object

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
)

type StringMap map[string]string


func (m *StringMap) MarshalJSON() ([]byte, error) {

	strMap, err := json.Marshal(m)

	return strMap, err
}


func (m *StringMap) UnmarshalJSON(b []byte) error {

	return json.Unmarshal(b, m)

}





// ------------------------------- Merge --------------------------------------------
func MergeStringMap(toStringMap *StringMap, subStringMaps ...*StringMap) *StringMap {
	if toStringMap == nil {
		toStringMap = &StringMap{}
	}
	for _, subStringMap := range subStringMaps {
		if subStringMap != nil {
			for k, v := range *subStringMap {
				(*subStringMap)[k] = v
			}
		}
	}
	return toStringMap
}

func ConvertStringMapToString(m *StringMap, separate string) string {
	var b bytes.Buffer
	for key, value := range *m {
		fmt.Fprintf(&b, "%s=%s%s", key, value, separate)
	}
	//fmt.Fprint(&b, "/0")
	return b.String()
}


// ------------------------------- Conversion ---------------------------------------

func StructToStringMap(obj interface{}) (newMap *StringMap, err error) {
	data, err := json.Marshal(obj) // Convert to a json string

	if err != nil {
		return
	}

	newMap = &StringMap{}
	err = json.Unmarshal(data, newMap) // Convert to a string map
	return
}



func GetJoinedWithKSort(params *StringMap) string {

	var strJoined string

	// ksort
	var keys []string
	for k := range *params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// join
	for _, k := range keys {
		strJoined += k + "=" + (*params)[k] + "&"
	}

	strJoined = strJoined[0 : len(strJoined)-1]

	return strJoined
}

// ------------------------------- Search --------------------------------------------


func GetStringMapKV(maps StringMap) (keys []string, values []string) {
	mapLen := len(maps)
	keys = make([]string, 0, mapLen)
	values = make([]string, 0, mapLen)

	for k, v := range maps {
		keys = append(keys, k)
		values = append(values, v)
	}

	return keys, values
}

// ------------------------------- Filter --------------------------------------------
func FilterEmptyStringMap(mapData *StringMap) (filteredMap *StringMap) {
	filteredMap = &StringMap{}
	for k, v := range *mapData {
		if v != "" {
			(*filteredMap)[k] = v
		}
	}
	return filteredMap
}
