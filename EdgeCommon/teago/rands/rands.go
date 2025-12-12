package rands

import (
	"crypto/rand"
	"encoding/hex"
	mrand "math/rand"
	"time"
)

func init() {
	mrand.Seed(time.Now().UnixNano())
}

// Int returns a random integer in [min, max].
func Int(min int, max int) int {
	if max <= min {
		return min
	}
	return mrand.Intn(max-min+1) + min
}

// Int64 returns a random int64 in [min, max].
func Int64(bounds ...int64) int64 {
	switch len(bounds) {
	case 0:
		return mrand.Int63()
	case 1:
		max := bounds[0]
		if max <= 0 {
			return 0
		}
		return mrand.Int63n(max)
	default:
		min, max := bounds[0], bounds[1]
		if max <= min {
			return min
		}
		return mrand.Int63n(max-min+1) + min
	}
}

// HexString returns a random hexadecimal string with the given length.
func HexString(length int) string {
	if length <= 0 {
		return ""
	}
	buf := make([]byte, (length+1)/2)
	_, _ = rand.Read(buf)
	hexed := hex.EncodeToString(buf)
	return hexed[:length]
}

// String returns a random alphanumeric string with given length.
func String(length int) string {
	if length <= 0 {
		return ""
	}
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[mrand.Intn(len(letters))]
	}
	return string(b)
}
