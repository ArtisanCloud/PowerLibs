package object

import (
	"encoding/json"
)

type AnyMap map[interface{}]interface{}

// ------------------------------- Conversion ---------------------------------------
func StructToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj) // Convert to a json string

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &newMap) // Convert to a map
	return
}

func StructToJson(obj interface{}) (strJson string, err error) {
	data, err := json.Marshal(obj) // Convert to a json string

	if err != nil {
		return "", err
	}

	return string(data), nil
}




