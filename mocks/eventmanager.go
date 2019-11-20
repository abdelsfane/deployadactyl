package mocks

import (
	I "github.com/compozed/deployadactyl/interfaces"
)

// EventManager handmade mock for tests.
type EventManager struct {
	AddHandlerCall struct {
		Received struct {
			Handler   I.Handler
			EventType string
		}
		Returns struct {
			Error error
		}
	}
	EmitCall struct {
		TimesCalled int
		Received    struct {
			Events []I.Event
		}
		Returns struct {
			Error []error
		}
	}
	EmitEventCall struct {
		TimesCalled int
		Received    struct {
			Events []I.IEvent
		}
		Returns struct {
			Error []error
		}
	}
}

// AddHandler mock method.
func (e *EventManager) AddHandler(handler I.Handler, eventType string) error {
	e.AddHandlerCall.Received.Handler = handler
	e.AddHandlerCall.Received.EventType = eventType

	return e.AddHandlerCall.Returns.Error
}

// Emit mock method.
func (e *EventManager) Emit(event I.Event) error {
	defer func() { e.EmitCall.TimesCalled++ }()

	e.EmitCall.Received.Events = append(e.EmitCall.Received.Events, event)

	if len(e.EmitCall.Returns.Error) > e.EmitCall.TimesCalled {
		return e.EmitCall.Returns.Error[e.EmitCall.TimesCalled]
	}

	return nil
}

func (e *EventManager) EmitEvent(event I.IEvent) error {
	defer func() { e.EmitEventCall.TimesCalled++ }()

	e.EmitEventCall.Received.Events = append(e.EmitEventCall.Received.Events, event)

	if len(e.EmitEventCall.Returns.Error) > e.EmitEventCall.TimesCalled {
		return e.EmitEventCall.Returns.Error[e.EmitEventCall.TimesCalled]
	}
	return nil

}

func (e *EventManager) AddBinding(binding I.Binding) {}
