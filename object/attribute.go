package object

type Attribute struct {
	Attributes HashMap
}

func NewAttribute(attributes *HashMap) *Attribute {
	return &Attribute{
		Attributes: *attributes,
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

func (attr *Attribute) SetAttribute(name string, value interface{}) *Attribute {
	attr.Attributes[name] = value
	return attr
}

func (attr *Attribute) Merge(attributes *HashMap) *Attribute {

	MergeHashMap(&attr.Attributes, attributes)

	return attr
}
