package interfaces

// Handler interface.
type Handler interface {
	OnEvent(event Event) error
}
