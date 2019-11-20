package interfaces

type Event struct {
	Type  string
	Data  interface{}
	Error error
}

func (e Event) Name() string {
	return e.Type
}

type Binding interface {
	Accepts(event interface{}) bool
	Emit(event interface{}) error
}

// EventManager interface.
type EventManager interface {
	AddHandler(handler Handler, eventType string) error
	Emit(event Event) error
	EmitEvent(event IEvent) error
	AddBinding(binding Binding)
}

type IEvent interface {
	Name() string
}
