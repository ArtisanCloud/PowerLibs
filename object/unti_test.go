package object

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ConvertToCentUnit(t *testing.T) {

	var money int
	money = ConvertToCentUnit(1.23)
	assert.Equal(t, 123, money)

	money = ConvertToCentUnit(123)
	assert.Equal(t, 12300, money)

	money = ConvertToCentUnit(1.23456)
	assert.Equal(t, 123, money)

	money = ConvertToCentUnit(0.123456)
	assert.Equal(t, 12, money)

	money = ConvertToCentUnit(0.163456)
	assert.Equal(t, 16, money)

	money = ConvertToCentUnit(0.166456)
	assert.Equal(t, 16, money)
}

func Test_ConvertToYuanUnit(t *testing.T) {

	var money float64
	money = ConvertToYuanUnit(123)
	assert.Equal(t, 1.23, money)

	money = ConvertToYuanUnit(12500)
	assert.Equal(t, 125.0, money)

	money = ConvertToYuanUnit(163456)
	assert.Equal(t, 1634.56, money)

}
