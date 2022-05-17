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
