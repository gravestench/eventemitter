package eventemitter

import (
	"sync"
	"unsafe"
)

// EventEmitter is a struct that implements the IEventEmitter interface.
type EventEmitter struct {
	mu            sync.Mutex
	listeners     map[string][]func(...any)
	onceListeners map[string][]func(...any)
}

// New creates a new EventEmitter instance.
func New() *EventEmitter {
	return &EventEmitter{
		listeners:     make(map[string][]func(...any)),
		onceListeners: make(map[string][]func(...any)),
	}
}

// Emit emits an event to all registered listeners.
func (e *EventEmitter) Emit(event string, args ...any) *sync.WaitGroup {
	e.mu.Lock()
	defer e.mu.Unlock()

	var wg sync.WaitGroup

	listeners := e.listeners[event]
	funcs := make([]func(args ...any), 0)
	for _, listener := range listeners {
		wg.Add(1)
		l := listener // keep in scope
		funcs = append(funcs, func(args ...any) {
			l(args...)
			wg.Done()
		})
	}

	// Call once listeners and remove them after execution
	onceListeners := e.onceListeners[event]
	onceFuncs := make([]func(args ...any), 0)
	for _, listener := range onceListeners {
		wg.Add(1)
		l := listener // keep in scope
		onceFuncs = append(onceFuncs, func(args ...any) {
			l(args...)
			wg.Done()
		})
	}

	for _, fn := range funcs {
		go fn(args...)
	}

	for _, fn := range onceFuncs {
		go fn(args...)
	}

	delete(e.onceListeners, event)

	return &wg
}

// On registers a listener for a specific event.
func (e *EventEmitter) On(event string, fn func(...any)) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.listeners[event] = append(e.listeners[event], fn)
}

// Off removes a specific listener for a specific event.
func (e *EventEmitter) Off(event string, fn func(...any)) {
	e.mu.Lock()
	defer e.mu.Unlock()

	pointerOf := func(fn func(...any)) uintptr {
		return *(*uintptr)(unsafe.Pointer(&fn))
	}

	for idx, listener := range e.listeners[event] {
		if pointerOf(listener) == pointerOf(fn) {
			e.listeners[event] = append(e.listeners[event][:idx], e.listeners[event][idx+1:]...)
			break
		}
	}

	if len(e.listeners[event]) == 0 {
		delete(e.listeners, event)
	}
}

// Once registers a one-time listener for a specific event.
func (e *EventEmitter) Once(event string, fn func(...any)) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.onceListeners[event] = append(e.onceListeners[event], fn)
}

// RemoveAllListeners removes all listeners for specific event.
func (e *EventEmitter) RemoveAllListeners(events ...string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, event := range events {
		delete(e.listeners, event)
		delete(e.onceListeners, event)
	}
}
