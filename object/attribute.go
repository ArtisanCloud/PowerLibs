package object

import (
	"errors"
	"fmt"
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
	attr.Attributes[name] = value
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
	if attr.Attributes[name] != nil {
		return attr.Attributes[name]
	} else {
		return defaultValue
	}
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
