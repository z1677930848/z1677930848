package types

import (
	"fmt"
	"strconv"
)

// String converts common values to string.
func String(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprint(v)
	}
}

// Int converts to int with a best-effort strategy.
func Int(value any) int {
	return int(Int64(value))
}

// Int32 converts to int32.
func Int32(value any) int32 {
	return int32(Int64(value))
}

// Uint16 converts to uint16.
func Uint16(value any) uint16 {
	return uint16(Int64(value))
}

// Uint64 converts to uint64.
func Uint64(value any) uint64 {
	if v, ok := value.(uint64); ok {
		return v
	}
	return uint64(Int64(value))
}

// Int8 converts to int8.
func Int8(value any) int8 {
	return int8(Int64(value))
}

// Int64 converts to int64.
func Int64(value any) int64 {
	switch v := value.(type) {
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	case uint:
		return int64(v)
	case uint8:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	case float32:
		return int64(v)
	case float64:
		return int64(v)
	case bool:
		if v {
			return 1
		}
		return 0
	case string:
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i
		}
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return int64(f)
		}
		return 0
	default:
		return 0
	}
}

// Bool converts to bool.
func Bool(value any) bool {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		if v == "" {
			return false
		}
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i != 0
		}
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f != 0
		}
		return false
	default:
		return Int64(value) != 0
	}
}

// Float64 converts to float64.
func Float64(value any) float64 {
	switch v := value.(type) {
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case float32:
		return float64(v)
	case float64:
		return v
	case bool:
		if v {
			return 1
		}
		return 0
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
		return 0
	default:
		return 0
	}
}

// Float32 converts to float32.
func Float32(value any) float32 {
	return float32(Float64(value))
}
