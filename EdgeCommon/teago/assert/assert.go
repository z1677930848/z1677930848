package assert

import (
	"encoding/json"
	"reflect"
	"testing"
)

type Assertion struct {
	t     *testing.T
	quiet bool
}

// NewAssertion creates a lightweight test helper compatible with TeaGo's assert API.
func NewAssertion(t *testing.T) *Assertion {
	return &Assertion{t: t}
}

func (a *Assertion) helper() {
	if a.t != nil {
		a.t.Helper()
	}
}

func (a *Assertion) fail(msg string, args ...any) {
	a.helper()
	if a.t != nil {
		a.t.Fatalf(msg, args...)
	}
}

// IsTrue fails the test if value is not true.
func (a *Assertion) IsTrue(value bool) {
	if !value {
		a.fail("expected true but got false")
	}
}

// IsFalse fails the test if value is not false.
func (a *Assertion) IsFalse(value bool) {
	if value {
		a.fail("expected false but got true")
	}
}

// IsNil fails if the value is not nil.
func (a *Assertion) IsNil(value any) {
	if !isNil(value) {
		a.fail("expected nil but got %#v", value)
	}
}

// IsNotNil fails if the value is nil.
func (a *Assertion) IsNotNil(value any) {
	if isNil(value) {
		a.fail("expected non-nil value")
	}
}

// Quiet marks the assertion helper as quiet (only logs via t.Log).
func (a *Assertion) Quiet() *Assertion {
	a.quiet = true
	return a
}

// Log forwards logs to testing.T when available.
func (a *Assertion) Log(args ...any) {
	if a.t == nil {
		return
	}
	a.t.Log(args...)
}

// LogJSON pretty-prints value as JSON for debugging.
func (a *Assertion) LogJSON(value any) {
	if a.t == nil {
		return
	}
	b, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		a.t.Log(value)
		return
	}
	a.t.Log(string(b))
}

func isNil(value any) bool {
	if value == nil {
		return true
	}
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return val.IsNil()
	default:
		return false
	}
}
