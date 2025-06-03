package random

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
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

// DeterministicRandom creates deterministic random numbers using a seed.
// The same seed, sequence number and probabilities generate the same outcome.
func DeterministicRandom(seedHex string, sequence int32, probabilities []float64) (int32, error) {
	// Validate input
	if len(seedHex) != 64 {
		return 0, errors.New("seedHex must be 64 bytes")
	} else if sequence < 0 {
		return 0, errors.New("sequence must be larger than than or equal to 0")
	} else if len(probabilities) == 0 {
		return 0, errors.New("probabilities must not be empty")
	}

	// Decode the seed
	seed, err := hex.DecodeString(seedHex)
	if err != nil {
		return 0, fmt.Errorf("invalid seed hex: %w", err)
	} else if len(seed) != 32 {
		return 0, errors.New("seed must decode to exactly 32 bytes")
	}

	// Validate and sum probabilities
	sum := 0.0
	for _, p := range probabilities {
		if p < 0 || p > 1 {
			return 0, fmt.Errorf("invalid input %v; valid range 0 <= p <= 1", p)
		}
		sum += p
	}

	const epsilon = 1e-12 // allow for minor float faults
	if math.Abs(sum-1.0) > epsilon {
		return 0, fmt.Errorf("sum of probabilities %v; must be exactly 1.0", sum)
	}

	// Build cumulative thresholds
	thresholds := make([]uint64, len(probabilities))
	cumulative := 0.0
	for i, p := range probabilities {
		cumulative += p
		if i == len(probabilities)-1 {
			thresholds[i] = math.MaxUint64 // ensure full coverage
		} else {
			thresholds[i] = uint64(cumulative * math.Pow(2, 64))
		}
	}

	// Compute random number
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(sequence))

	h := sha256.New()
	h.Write(seed)
	h.Write(buf[:])
	hash := h.Sum(nil)
	x := binary.BigEndian.Uint64(hash[:8])

	// Find the selected index
	for i, t := range thresholds {
		if x < t {
			if i < math.MinInt32 || i > math.MaxInt32 {
				return 0, fmt.Errorf("threshold index out of range for Int32")
			}
			return int32(i), nil
		}
	}

	// Should never happen if sum == 1.0
	return 0, errors.New("unexpected: no prize selected despite sum == 1.0")
}
