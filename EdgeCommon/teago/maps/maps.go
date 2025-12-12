package maps

import (
	"encoding/json"
	"reflect"

	"github.com/iwind/TeaGo/types"
)

// Map is a light alias around map[string]any with helpers similar to TeaGo.
type Map map[string]any

// NewMap copies an existing map into a fresh Map.
func NewMap(source any) Map {
	switch v := source.(type) {
	case Map:
		return copyMap(v)
	case map[string]any:
		return copyMap(Map(v))
	case map[string]string:
		dst := Map{}
		for k, val := range v {
			dst[k] = val
		}
		return dst
	default:
		return Map{}
	}
}

func copyMap(src Map) Map {
	dst := Map{}
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// Has reports whether the key exists.
func (m Map) Has(key string) bool {
	_, ok := m[key]
	return ok
}

// GetString returns the value as string.
func (m Map) GetString(key string) string {
	return types.String(m[key])
}

// GetInt returns the value as int.
func (m Map) GetInt(key string) int {
	return types.Int(m[key])
}

// GetInt32 returns the value as int32.
func (m Map) GetInt32(key string) int32 {
	return types.Int32(m[key])
}

// GetInt64 returns the value as int64.
func (m Map) GetInt64(key string) int64 {
	return types.Int64(m[key])
}

// GetUint64 returns the value as uint64.
func (m Map) GetUint64(key string) uint64 {
	return types.Uint64(m[key])
}

// GetBool returns the value as bool.
func (m Map) GetBool(key string) bool {
	return types.Bool(m[key])
}

// GetBytes returns the value as []byte when possible.
func (m Map) GetBytes(key string) []byte {
	value := m[key]
	switch v := value.(type) {
	case []byte:
		return v
	case string:
		return []byte(v)
	default:
		return nil
	}
}

// GetFloat64 returns the value as float64.
func (m Map) GetFloat64(key string) float64 {
	return types.Float64(m[key])
}

// GetFloat32 returns the value as float32.
func (m Map) GetFloat32(key string) float32 {
	return float32(types.Float64(m[key]))
}

// Get returns raw value.
func (m Map) Get(key string) any {
	return m[key]
}

// FindCol returns a value by key.
func (m Map) FindCol(key string) any {
	return m[key]
}

// GetSlice returns the value as []string when possible.
func (m Map) GetSlice(key string) []string {
	value := m[key]
	switch v := value.(type) {
	case []string:
		return v
	case []any:
		result := make([]string, 0, len(v))
		for _, item := range v {
			result = append(result, types.String(item))
		}
		return result
	default:
		rv := reflect.ValueOf(value)
		if rv.IsValid() && (rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array) {
			result := make([]string, 0, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				result = append(result, types.String(rv.Index(i).Interface()))
			}
			return result
		}
		return nil
	}
}

// FindStringCol returns string value or default.
func (m Map) FindStringCol(key string) string {
	return types.String(m[key])
}

// FindInt64Col returns int64 value or default 0.
func (m Map) FindInt64Col(key string) int64 {
	return types.Int64(m[key])
}

// GetMap returns the value as Map when possible.
func (m Map) GetMap(key string) Map {
	value := m[key]
	switch v := value.(type) {
	case Map:
		return v
	case map[string]any:
		return Map(v)
	default:
		return nil
	}
}

// AsJSON marshals map to JSON bytes.
func (m Map) AsJSON() []byte {
	b, _ := json.Marshal(m)
	return b
}

// AsPrettyJSON marshals map to indented JSON bytes.
func (m Map) AsPrettyJSON() []byte {
	b, _ := json.MarshalIndent(m, "", "  ")
	return b
}
