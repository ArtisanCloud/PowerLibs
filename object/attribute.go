package object

import (
	"errors"
	"fmt"
	"strings"
)

type Attribute struct {
	Attributes HashMap
}

func NewAttribute(attributes *HashMap) *Attribute {
	return &Attribute{
		Attributes: *attributes,
	}
}

func (attr *Attribute) SetAttributes(attributes *HashMap) *Attribute {
	attr.Attributes = *attributes

	return attr

}

func (attr *Attribute) SetAttribute(name string, value interface{}) *Attribute {

	segments := strings.Split(name, ".")
	newItem := attr.Attributes

	var segment string
	for len(segments) > 1 {
		segment, segments = segments[0], segments[1:]
		if newItem[segment] == nil {
			newItem[segment] = map[string]interface{}{}
		}
		newItem = newItem[segment].(map[string]interface{})
	}

	// transform the hashobject to string
	switch value.(type) {
	case map[string]interface{}:
	case HashMap:
	case *HashMap:
		value, _ = JsonEncode(value)
		break
	default:
	}

	newItem[segments[0]] = value

	return attr
}

func (attr *Attribute) IsRequired(attribute string) bool {

	requiredAttributes := attr.GetRequired()
	has := ContainsString(requiredAttributes, attribute)
	return has
}

func (attr *Attribute) GetRequired() []string {
	if attr.Attributes["required"] != nil {
		return attr.Attributes["required"].([]string)
	} else {
		return []string{}
	}
}

func (attr *Attribute) GetAttributes() *HashMap {
	return &attr.Attributes
}

func (attr *Attribute) GetAttribute(name string, defaultValue interface{}) interface{} {

	var result interface{}

	hashedObject := attr.Attributes

	if name == "" {
		return &hashedObject
	}

	if hashedObject[name] != nil {
		return hashedObject[name]
	} else {
		result = defaultValue
	}

	segments := strings.Split(name, ".")
	if len(segments) > 1 {
		for _, segment := range segments {
			if hashedObject[segment] == nil {
				return defaultValue
			} else {
				switch hashedObject[segment].(type) {
				case HashMap:
					hashedObject = hashedObject[segment].(HashMap)
				case map[string]interface{}:
					hashedObject = hashedObject[segment].(map[string]interface{})
				default:
					return hashedObject[segment]
				}
			}
		}
	}

	return result

}

func (attr *Attribute) Get(attribute string, defaultValue interface{}) interface{} {
	return attr.GetAttribute(attribute, defaultValue)
}

func (attr *Attribute) Has(key string) bool {
	return attr.Attributes[key] != nil
}

func (attr *Attribute) GetString(attribute string, defaultValue string) string {
	strResult := attr.Get(attribute, defaultValue).(string)
	if strResult == "" {
		strResult = defaultValue
	}
	return strResult
}

func (attr *Attribute) Merge(attributes *HashMap) *Attribute {

	MergeHashMap(&attr.Attributes, attributes)

	return attr
}

func (attr *Attribute) CheckRequiredAttributes() error {

	requiredAttributes := attr.GetRequired()
	for _, attribute := range requiredAttributes {
		if attr.GetAttribute(attribute, nil) == nil {
			return errors.New(fmt.Sprintf("\"%s\" cannot be empty.", attribute))
		}
	}
	return nil
}
