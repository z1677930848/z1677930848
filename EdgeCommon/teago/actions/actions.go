package actions

import "net/http"

// ActionObject mirrors the minimal fields used in EdgeCommon to avoid pulling TeaGo.
type ActionObject struct {
	Request *http.Request
}

// ActionWrapper is the minimal interface consumed by callers.
type ActionWrapper interface {
	Object() *ActionObject
}

// NewActionObject is a helper to build an ActionObject from a request.
func NewActionObject(req *http.Request) *ActionObject {
	return &ActionObject{Request: req}
}
