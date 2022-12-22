package object

import (
	"encoding/json"
	"github.com/ArtisanCloud/PowerLibs/v3/fmt"
	"reflect"
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

// --- ---
func GetModelTags(t reflect.Type, tagName string) (tags []string) {

	// TypeOf returns the reflection Type that represents the dynamic type of variable.
	// If variable is a nil interface value, TypeOf returns nil.
	//t := reflect.TypeOf(model)

	// Get the type and kind of our user variable
	//fmt.Println("Type:", t.Name())
	//fmt.Println("Kind:", t.Kind())

	tags = []string{}
	subTags := []string{}
	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		// Get the field, returns https://golang.org/pkg/reflect/#StructField
		field := t.Field(i)
		kind := field.Type.Kind().String()
		if kind == "struct" || kind == "ptr" {
			subTags = GetModelTags(field.Type.Elem(), tagName)
			tags = append(tags, subTags...)

		} else {
			// Get the field tag value
			tag := field.Tag.Get(tagName)
			tags = append(tags, tag)
		}
		fmt.Dump(tags)

	}
	return tags
}

func IsObjectNil(obj interface{}) bool {
	return obj == nil || (reflect.ValueOf(obj).Kind() == reflect.Ptr && reflect.ValueOf(obj).IsNil())
}
