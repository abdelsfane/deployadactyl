package delete

import (
	"io"
	"reflect"

	"github.com/compozed/deployadactyl/eventmanager"
	"github.com/compozed/deployadactyl/interfaces"
	"github.com/compozed/deployadactyl/structs"
	"github.com/go-errors/errors"
)

type eventBinding struct {
	etype   reflect.Type
	handler func(event interface{}) error
}

func (s eventBinding) Accepts(event interface{}) bool {
	return reflect.TypeOf(event) == s.etype
}

func (b eventBinding) Emit(event interface{}) error {
	return b.handler(event)
}

type DeleteFailureEvent struct {
	CFContext     interfaces.CFContext
	Data          map[string]interface{}
	Authorization interfaces.Authorization
	Environment   structs.Environment
	Error         error
	Response      io.ReadWriter
	Log           interfaces.DeploymentLogger
}

func (e DeleteFailureEvent) Name() string {
	return "DeleteFailureEvent"
}

func NewDeleteFailureEventBinding(handler func(event DeleteFailureEvent) error) interfaces.Binding {
	return eventBinding{
		etype: reflect.TypeOf(DeleteFailureEvent{}),
		handler: func(gevent interface{}) error {
			event, ok := gevent.(DeleteFailureEvent)
			if ok {
				return handler(event)
			} else {
				return eventmanager.InvalidEventType{errors.New("invalid event type")}
			}
		},
	}
}

type DeleteSuccessEvent struct {
	CFContext     interfaces.CFContext
	Data          map[string]interface{}
	Authorization interfaces.Authorization
	Environment   structs.Environment
	Response      io.ReadWriter
	Log           interfaces.DeploymentLogger
}

func (e DeleteSuccessEvent) Name() string {
	return "DeleteSuccessEvent"
}

func NewDeleteSuccessEventBinding(handler func(event DeleteSuccessEvent) error) interfaces.Binding {
	return eventBinding{
		etype: reflect.TypeOf(DeleteSuccessEvent{}),
		handler: func(gevent interface{}) error {
			event, ok := gevent.(DeleteSuccessEvent)
			if ok {
				return handler(event)
			} else {
				return eventmanager.InvalidEventType{errors.New("invalid event type")}
			}
		},
	}
}

type DeleteStartedEvent struct {
	CFContext     interfaces.CFContext
	Data          map[string]interface{}
	Environment   structs.Environment
	Authorization interfaces.Authorization
	Response      io.ReadWriter
	Log           interfaces.DeploymentLogger
}

func (e DeleteStartedEvent) Name() string {
	return "DeleteStartedEvent"
}

func NewDeleteStartedEventBinding(handler func(event DeleteStartedEvent) error) interfaces.Binding {
	return eventBinding{
		etype: reflect.TypeOf(DeleteStartedEvent{}),
		handler: func(gevent interface{}) error {
			event, ok := gevent.(DeleteStartedEvent)
			if ok {
				return handler(event)
			} else {
				return eventmanager.InvalidEventType{errors.New("invalid event type")}
			}
		},
	}
}

type DeleteFinishedEvent struct {
	CFContext     interfaces.CFContext
	Data          map[string]interface{}
	Authorization interfaces.Authorization
	Environment   structs.Environment
	Response      io.ReadWriter
	Log           interfaces.DeploymentLogger
}

func (e DeleteFinishedEvent) Name() string {
	return "DeleteFinishedEvent"
}

func NewDeleteFinishedEventBinding(handler func(event DeleteFinishedEvent) error) interfaces.Binding {
	return eventBinding{
		etype: reflect.TypeOf(DeleteFinishedEvent{}),
		handler: func(gevent interface{}) error {
			event, ok := gevent.(DeleteFinishedEvent)
			if ok {
				return handler(event)
			} else {
				return eventmanager.InvalidEventType{errors.New("invalid event type")}
			}
		},
	}
}
