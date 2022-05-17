package object

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MergeHashMap(t *testing.T) {

	base := &HashMap{
		"key1": 123,
		"key2": "456",
		"key3": "",
		"key4": nil,
		"key5": &StringMap{},
		"key6": &HashMap{},
	}

	toMap := &HashMap{
		"key2": "",
		"key3": "123",
		"key4": &HashMap{},
		"key5": nil,
		"key6": &StringMap{},
	}

	toMap = MergeHashMap(toMap, base)

	assert.EqualValues(t, &HashMap{
		"key1": 123,
		"key2": "456",
		"key3": "123",
		"key4": &HashMap{},
		"key5": &StringMap{},
		"key6": &StringMap{},
	}, toMap)

}

func Test_ReplaceHashMapRecursive(t *testing.T) {

	base := &HashMap{
		"key1": 123,
		"key2": "456",
		"key3": "789",
		"key4": nil,
		"key5": map[string]int{},
		"key6": &map[string]float32{},
	}

	base2 := &HashMap{
		"key1": 456,
		"key2": "base456",
		"key3": "",
		"key4": nil,
		"key5": &StringMap{},
	}

	toMap := &HashMap{
		"key2": "",
		"key3": "123",
		"key4": &HashMap{},
		"key5": nil,
		"key6": &StringMap{},
	}

	toMap = ReplaceHashMapRecursive(toMap, base, base2)

	assert.EqualValues(t, &HashMap{
		"key1": 456,
		"key2": "base456",
		"key3": "",
		"key4": nil,
		"key5": &StringMap{},
		"key6": &map[string]float32{},
	}, toMap)

}
