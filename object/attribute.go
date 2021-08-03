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

func (attr *Attribute) IsRequired(attributes interface{}) bool {

	has, _ := InHash(attributes, attr.GetRequired().(*HashMap))
	return has
}

func (attr *Attribute) GetRequired() interface{} {
	if attr.Attributes["required"] != nil {
		return attr.Attributes["required"]
	} else {
		return HashMap{}
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

func (attr *Attribute) Merge(attributes *HashMap) *Attribute {

	MergeHashMap(&attr.Attributes, attributes)

	return attr
}

func (attr *Attribute) CheckRequiredAttributes() error {
	requiredAttributes := attr.GetRequired().(*HashMap)
	for attribute, _ := range *requiredAttributes {
		if attr.GetAttribute(attribute, nil) == nil {
			return errors.New(fmt.Sprintf("\"%s\" cannot be empty.", attribute))
		}
	}
	return nil
}
