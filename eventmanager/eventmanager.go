// Package eventmanager emits events.
package eventmanager

import (
	"fmt"

	I "github.com/compozed/deployadactyl/interfaces"
	"github.com/go-errors/errors"
)

type Err interface {
	Error() string
}

type InvalidEventType struct {
	Err
}

type EventManagerConstructor func(log I.DeploymentLogger, bindings []I.Binding) I.EventManager

// EventManager has handlers for each registered event type.
type EventManager struct {
	Bindings []I.Binding
	Log      I.Logger
}

type legacyEventBinding struct {
	etype   string
	handler I.Handler
}

func (b legacyEventBinding) Accepts(event interface{}) bool {
	levent, ok := event.(I.Event)
	if !ok {
		return false
	}
	return levent.Type == b.etype
}

func (b legacyEventBinding) Emit(event interface{}) error {
	levent, ok := event.(I.Event)
	if !ok {
		return InvalidEventType{Err: errors.New("invalid event type")}
	}
	return b.handler.OnEvent(levent)
}

func NewEventManager(log I.DeploymentLogger, bindings []I.Binding) I.EventManager {
	if bindings == nil {
		bindings = make([]I.Binding, 0)
	}
	return &EventManager{
		Log:      log,
		Bindings: bindings,
	}
}

// AddHandler takes a handler and eventType and returns an error if a handler is not provided.
func (e *EventManager) AddHandler(handler I.Handler, eventType string) error {
	if handler == nil {
		return InvalidArgumentError{}
	}
	e.Bindings = append(e.Bindings, legacyEventBinding{
		etype:   eventType,
		handler: handler,
	})
	e.Log.Debugf("handler for [%s] event added successfully", eventType)
	return nil
}

// Emit emits an event.
func (e *EventManager) Emit(event I.Event) error {
	return e.EmitEvent(event)
}

func (e *EventManager) AddBinding(binding I.Binding) {
	e.Bindings = append(e.Bindings, binding)
}

func (e EventManager) EmitEvent(event I.IEvent) (err error) {
	defer func() {
		if r := recover(); r != nil {
			e.Log.Errorf("recovered from panic")
			err = errors.New(fmt.Sprintf("Recovered from a panic: %v", r))
		}
	}()
	for _, binding := range e.Bindings {
		if binding.Accepts(event) {
			err = binding.Emit(event)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
