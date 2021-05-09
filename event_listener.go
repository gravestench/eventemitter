package eventemitter

type EventListener struct {
	fn   func(...interface{})
	once bool
}
