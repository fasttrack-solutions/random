package random

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

// UniformInt64 generates an int64 in the range (min, max) using a uniform distribution
func UniformInt64(min int64, max int64) (int64, error) {
	num := max - min + 1
	res, err := rand.Int(rand.Reader, big.NewInt(int64(num)))
	if err != nil {
		return 0, err
	}

	return res.Int64() + min, nil
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

// TransformToExponential transforms a float64 between (0,1) to (0,1) using a linear to exponential transformation
// The aggression parameter must be larger than 5 and less than 100
// The progression parameter must be between 0 and 1
// i.e. TransformToExponential(10, 0.5) ->
func TransformToExponential(aggression float64, progression float64) (float64, error) {
	if aggression < 5 || aggression > 100 {
		return 0, fmt.Errorf("aggression must be between 5 and 100")
	} else if progression <= 0 {
		return 0, nil
	} else if progression >= 1 {
		return 1, nil
	}

	expProgression := math.Pow(math.E, aggression*(progression-1))
	truncated := Truncate(expProgression, 9)

	if truncated == 0 {
		return 0.000000001, nil
	}

	return truncated, nil
}

// Truncate chops off decimals after precision.
// i.e. Truncate(1.2345678, 5) returns 1.23456 instead of a rounded 1.23457
func Truncate(val float64, precision int) float64 {
	multiplier := math.Pow(10, float64(precision))
	return math.Floor(val*multiplier) / multiplier
}

// Round function to round a float to a specific number of decimal places
func Round(val float64, precision int) float64 {
	multiplier := math.Pow(10, float64(precision))
	return math.Round(val*multiplier) / multiplier
}
