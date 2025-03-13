package random

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func Test_UniformInt64(t *testing.T) {
	number, err := UniformInt64(10, 50)
	assert.Nil(t, err)

	assert.True(t, number >= int64(10))
	assert.True(t, number <= int64(50))

	number, err = UniformInt64(-1, 0)
	assert.EqualError(t, err, "min must be larger than or equal to 0")

	number, err = UniformInt64(0, -1)
	assert.EqualError(t, err, "max must be larger than or equal to 0")

	number, err = UniformInt64(0, 0)
	assert.Equal(t, 0, 0)

	number, err = UniformInt64(math.MaxInt32, 1)
	assert.EqualError(t, err, "min must be less than 2,147,483,647")

	number, err = UniformInt64(0, math.MaxInt32)
	assert.EqualError(t, err, "max must be less than 2,147,483,647")

	number, err = UniformInt64(2, 1)
	assert.EqualError(t, err, "min must be less than max")

	number, err = UniformInt64(432, 432)
	assert.Equal(t, number, int64(432))
}

func Test_UniformFloat64(t *testing.T) {
	number, err := UniformFloat64()
	assert.Nil(t, err)

	assert.True(t, number >= 0)
	assert.True(t, number < 1)
}

func Test_Truncate(t *testing.T) {
	number := Truncate(0.123456789, 1)
	assert.Equal(t, 0.1, number)

	number = Truncate(0.123456789, 2)
	assert.Equal(t, 0.12, number)

	number = Truncate(0.123456789, 5)
	assert.Equal(t, 0.12345, number)

	number = Truncate(0.123456789, 9)
	assert.Equal(t, 0.123456789, number)

	number = Truncate(0.1234567891, 9)
	assert.Equal(t, 0.123456789, number)
}
