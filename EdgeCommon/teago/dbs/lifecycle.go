package dbs

import "sync"

var readyCallbacks []func()
var readyMu sync.Mutex
var isReady bool

// NotifyReady triggers registered callbacks.
func NotifyReady() {
	readyMu.Lock()
	isReady = true
	callbacks := append([]func(){}, readyCallbacks...)
	readyCallbacks = nil
	readyMu.Unlock()
	for _, cb := range callbacks {
		cb()
	}
}

// OnReadyDone registers a callback for ready event.
func OnReadyDone(cb func()) {
	if cb == nil {
		return
	}
	readyMu.Lock()
	if isReady {
		readyMu.Unlock()
		cb()
		return
	}
	readyCallbacks = append(readyCallbacks, cb)
	readyMu.Unlock()
}

// OnReady is an alias for OnReadyDone.
func OnReady(cb func()) {
	OnReadyDone(cb)
}
