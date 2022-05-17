package object

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MergeStringMap(t *testing.T) {

	base := &StringMap{
		"key1": "123",
		"key2": "xxx",
		"key3": "",
		"key4": "101112",

	}

	toMap := &StringMap{
		"key1": "",
		"key2": "456",
		"key3": "789",

	}

	toMap = MergeStringMap(toMap, base)

	assert.EqualValues(t, &StringMap{
		"key1": "123",
		"key2": "456",
		"key3": "789",
		"key4": "101112",
	}, toMap)

}

