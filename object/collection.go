package object

import "strings"

type Collection struct {
	items HashMap
}

func NewCollection(items *HashMap) *Collection {
	return &Collection{
		items: *items,
	}
}

func (c *Collection) All() *HashMap {
	return &c.items
}

func (c *Collection) Only(keys []string) (result *HashMap) {

	result = &HashMap{}

	for key, value := range c.items {
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

}

// Get an item from an hashMap using "dot" notation.
func (c *Collection) Get(key string, defaultValue interface{}) interface{} {

	var result interface{}

	hashedObject := c.items

	if key == "" {
		return &hashedObject
	}

	if hashedObject[key] != nil {
		return hashedObject[key]
	}else{
		result = defaultValue
	}

	segments := strings.Split(key, ".")
	if len(segments) > 1 {
		for _, segment := range segments {
			if hashedObject[segment] == nil {
				return defaultValue
			} else {
				result = hashedObject[segment]
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

func (c *Collection) ToJson(option int) string {
	return ""
}
func (c *Collection) ToString() string {
	return c.ToJson(0)
}

func (c *Collection) Count() int {
	return len(c.items)
}

func (c *Collection) Unserialize(serialized string) HashMap {

	return c.items
}
