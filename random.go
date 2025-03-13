package random

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/big"
)

// UniformInt64 generates an int64 in the range (min, max) using a uniform distribution
func UniformInt64(min int32, max int32) (int64, error) {
	if min < 0 {
		return 0, errors.New("min must be larger than or equal to 0")
	} else if max < 0 {
		return 0, errors.New("max must be larger than or equal to 0")
	} else if min > math.MaxInt32-1 {
		return 0, errors.New("min must be less than 2,147,483,647")
	} else if max > math.MaxInt32-1 {
		return 0, errors.New("max must be less than 2,147,483,647")
	} else if max < min {
		return 0, errors.New("min must be less than max")
	}

	if min == max {
		return int64(min), nil
	}

	num := max - min + 1
	res, err := rand.Int(rand.Reader, big.NewInt(int64(num)))
	if err != nil {
		return 0, err
	}

	return res.Int64() + int64(min), nil
}

// UniformFloat64 generates a float64 in the range (0, 1] using a uniform distribution
func UniformFloat64() (float64, error) {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return 0, fmt.Errorf("failed to generate secure random number: %s", err.Error())
	}
	r := float64(binary.LittleEndian.Uint64(b[:])) / (1 << 64)

	return Truncate(r, 9), nil
}

// Truncate chops off decimals after precision.
// i.e. Truncate(1.2345678, 5) returns 1.23456 instead of a rounded 1.23457
func Truncate(val float64, precision int) float64 {
	multiplier := math.Pow(10, float64(precision))
	return math.Floor(val*multiplier) / multiplier
}
