package object

import (
	"strings"
	"time"
)

type Collection struct {
	items *HashMap
}

func NewCollection(items *HashMap) *Collection {

	if items == nil {
		items = &HashMap{}
	}

	return &Collection{
		items: items,
	}
}

func (c *Collection) All() *HashMap {
	return c.items
}

func (c *Collection) Only(keys []string) (result *HashMap) {

	result = &HashMap{}

	for key, value := range *c.items {
		value = c.Get(key, nil)
		if value != nil {
			(*result)[key] = value
		}
	}

	return result
}

func (c *Collection) Except(keys []string) HashMap {
	return nil
}

func (c *Collection) Merge(items *HashMap) HashMap {
	return nil
}

func (c *Collection) Has(key string) bool {
	return false
}

func (c *Collection) First() interface{} {
	return nil
}

func (c *Collection) Last() interface{} {
	return nil
}
func (c *Collection) Add(key string, value interface{}) {

}

func (c *Collection) Set(key string, value interface{}) {
	segments := strings.Split(key, ".")
	newItem := c.items

	var segment string
	for len(segments) > 1 {
		segment, segments = segments[0], segments[1:]
		if (*newItem)[segment] == nil {
			(*newItem)[segment] = &HashMap{}
		}
		newItem = (*newItem)[segment].(*HashMap)
	}

	(*newItem)[segments[0]] = value
}

func (c *Collection) GetBoolPointer(key string, defaultValue bool) *bool {
	value := c.GetBool(key, defaultValue)
	return &value
}

func (c *Collection) GetIntPointer(key string, defaultValue int) *int {
	value := c.GetInt(key, defaultValue)
	return &value
}

func (c *Collection) GetInt8Pointer(key string, defaultValue int8) *int8 {
	value := c.GetInt8(key, defaultValue)
	return &value
}

func (c *Collection) GetInt16Pointer(key string, defaultValue int16) *int16 {
	value := c.GetInt16(key, defaultValue)
	return &value
}

func (c *Collection) GetInt32Pointer(key string, defaultValue int32) *int32 {
	value := c.GetInt32(key, defaultValue)
	return &value
}

func (c *Collection) GetInt64Pointer(key string, defaultValue int64) *int64 {
	value := c.GetInt64(key, defaultValue)
	return &value
}

func (c *Collection) GetStringPointer(key string, defaultValue string) *string {
	value := c.GetString(key, defaultValue)
	return &value

}

func (c *Collection) GetFloat64Pointer(key string, defaultValue float64) *float64 {
	value := c.GetFloat64(key, defaultValue)
	return &value
}

func (c *Collection) GetFloat32Pointer(key string, defaultValue float64) *float32 {
	value := c.GetFloat32(key, defaultValue)
	return &value
}

func (c *Collection) GetDateTimePointer(key string, defaultValue time.Time) *time.Time {
	value := c.GetDateTime(key, defaultValue)
	return &value
}

// ----------------------------------------------------------------------------------------

func (c *Collection) GetBool(key string, defaultValue bool) bool {
	return c.Get(key, defaultValue).(bool)
}

func (c *Collection) GetIntArray(key string, defaultValue []int) []int {
	return c.Get(key, defaultValue).([]int)
}
func (c *Collection) GetFloat64Array(key string, defaultValue []float64) []float64 {
	return c.Get(key, defaultValue).([]float64)
}

func (c *Collection) GetInterfaceArray(key string, defaultValue []interface{}) []interface{} {
	return c.Get(key, defaultValue).([]interface{})
}

func (c *Collection) GetStringArray(key string, defaultValue []string) []string {
	return c.Get(key, defaultValue).([]string)
}

func (c *Collection) GetInt(key string, defaultValue int) int {
	return c.Get(key, defaultValue).(int)
}

func (c *Collection) GetInt8(key string, defaultValue int8) int8 {
	return c.Get(key, defaultValue).(int8)
}

func (c *Collection) GetInt16(key string, defaultValue int16) int16 {
	return c.Get(key, defaultValue).(int16)
}

func (c *Collection) GetInt32(key string, defaultValue int32) int32 {
	return c.Get(key, defaultValue).(int32)
}

func (c *Collection) GetInt64(key string, defaultValue int64) int64 {
	return c.Get(key, defaultValue).(int64)
}

func (c *Collection) GetString(key string, defaultValue string) string {

	strResult := c.Get(key, defaultValue).(string)
	if strResult == "" {
		strResult = defaultValue
	}
	return strResult
}

func (c *Collection) GetNullString(key string, defaultValue string) NullString {

	value := c.Get(key, nil)

	switch value.(type) {
	case NullString:
		return value.(NullString)

	case string:
		strValue := value.(string)
		if strValue != "" {
			return NewNullString(strValue, true)
		} else {
			if defaultValue != "" {
				return NewNullString(defaultValue, true)
			}
		}

	}
	return NewNullString("", false)
}

func (c *Collection) GetFloat64(key string, defaultValue float64) float64 {
	return c.Get(key, defaultValue).(float64)
}

func (c *Collection) GetFloat32(key string, defaultValue float64) float32 {
	return c.Get(key, defaultValue).(float32)
}

func (c *Collection) GetDateTime(key string, defaultValue time.Time) time.Time {
	return c.Get(key, defaultValue).(time.Time)
}

// Get an item from an hashMap using "dot" notation.
func (c *Collection) Get(key string, defaultValue interface{}) interface{} {

	var result interface{}

	hashedObject := c.items

	if key == "" {
		return &hashedObject
	}

	if (*hashedObject)[key] != nil {
		return (*hashedObject)[key]
	} else {
		result = defaultValue
	}

	segments := strings.Split(key, ".")
	if len(segments) > 1 {
		for _, segment := range segments {
			if (*hashedObject)[segment] == nil {
				return defaultValue
			} else {
				switch (*hashedObject)[segment].(type) {
				case *HashMap:
					hashedObject = (*hashedObject)[segment].(*HashMap)
				case HashMap:
					*hashedObject = (*hashedObject)[segment].(HashMap)
				case map[string]interface{}:
					*hashedObject = (*hashedObject)[segment].(map[string]interface{})
				default:
					return (*hashedObject)[segment]
				}
			}
		}
	}

	return result
}

func (c *Collection) Forget(key string) {

}

func (c *Collection) ToHashMap() *HashMap {
	return c.All()
}

func (c *Collection) ToJson(option int) (string, error) {
	return JsonEncode(c.items)
}
func (c *Collection) ToString() string {
	strJson, _ := c.ToJson(0)
	return strJson
}

func (c *Collection) Count() int {
	return len(*c.items)
}

func (c *Collection) Unserialize(serialized string) *HashMap {

	return c.items
}
