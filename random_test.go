package random

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_UniformInt64(t *testing.T) {
	number, err := UniformInt64(1, 1000)
	assert.Nil(t, err)

	assert.True(t, number > 0)
	assert.True(t, number < 1001)
}

func Test_UniformFloat64(t *testing.T) {
	number, errRandomUniformFloat64 := UniformFloat64()

	assert.Nil(t, errRandomUniformFloat64)
	assert.True(t, number >= 0)
	assert.True(t, number < 1)
}

func Test_TransformToExponential(t *testing.T) {
	number, errRandomUniformFloat64 := TransformToExponential(10, 0.1)
	assert.Nil(t, errRandomUniformFloat64)
	assert.Equal(t, 0.000123409, number)

	number, errRandomUniformFloat64 = TransformToExponential(10, 0.5)
	assert.Nil(t, errRandomUniformFloat64)
	assert.Equal(t, 0.006737946, number)

	number, errRandomUniformFloat64 = TransformToExponential(10, 0.9)
	assert.Nil(t, errRandomUniformFloat64)
	assert.Equal(t, 0.367879441, number)

	number, errRandomUniformFloat64 = TransformToExponential(10, 0.99)
	assert.Nil(t, errRandomUniformFloat64)
	assert.Equal(t, 0.904837418, number)
}

func Test_Truncate(t *testing.T) {
	number := Truncate(0.123456789, 1)
	assert.Equal(t, 0.1, number)

	number = Truncate(0.123456789, 2)
	assert.Equal(t, 0.12, number)

	number = Truncate(0.123456789, 5)
	assert.Equal(t, 0.12345, number)

	number = Truncate(0.123456789, 8)
	assert.Equal(t, 0.12345678, number)
}

func Test_Round(t *testing.T) {
	number := Round(0.123456789, 1)
	assert.Equal(t, 0.1, number)

	number = Round(0.123456789, 2)
	assert.Equal(t, 0.12, number)

	number = Round(0.123456789, 5)
	assert.Equal(t, 0.12346, number)

	number = Round(0.123456789, 8)
	assert.Equal(t, 0.12345679, number)
}
