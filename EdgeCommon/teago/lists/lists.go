package lists

import (
	"reflect"
	"sort"
)

// ContainsString reports whether value exists in the slice.
func ContainsString(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

// ContainsInt reports whether value exists in the slice.
func ContainsInt(list []int, value int) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

// ContainsInt64 reports whether value exists in the int64 slice.
func ContainsInt64(list []int64, value int64) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

// Contains performs a generic scan over a slice.
func Contains(list any, value any) bool {
	rv := reflect.ValueOf(list)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return false
	}
	for i := 0; i < rv.Len(); i++ {
		if reflect.DeepEqual(rv.Index(i).Interface(), value) {
			return true
		}
	}
	return false
}

// Sort wraps sort.Slice to match the TeaGo signature.
func Sort(slice any, less func(i, j int) bool) {
	sort.Slice(slice, less)
}

// Reverse reverses slice elements in place.
func Reverse(slice any) {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return
	}
	left, right := 0, rv.Len()-1
	for left < right {
		li := rv.Index(left)
		ri := rv.Index(right)
		tmp := li.Interface()
		li.Set(ri)
		ri.Set(reflect.ValueOf(tmp))
		left++
		right--
	}
}
